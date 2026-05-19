package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"unicode"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/pkg/metrics"
)

var tracer = otel.Tracer("user-auth-service/usecase")

// authUsecase — конкретная реализация AuthUsecase.
// Зависит исключительно от domain-интерфейсов (Dependency Inversion).
type authUsecase struct {
	userRepo  domain.UserRepository
	cache     domain.CacheRepository
	publisher domain.MessagePublisher
	email     domain.EmailService
	tokens    domain.TokenManager
	logger    *zap.Logger
}

// NewAuthUsecase — конструктор с внедрением зависимостей
func NewAuthUsecase(
	userRepo domain.UserRepository,
	cache domain.CacheRepository,
	publisher domain.MessagePublisher,
	email domain.EmailService,
	tokens domain.TokenManager,
	logger *zap.Logger,
) AuthUsecase {
	return &authUsecase{
		userRepo:  userRepo,
		cache:     cache,
		publisher: publisher,
		email:     email,
		tokens:    tokens,
		logger:    logger,
	}
}

// ─────────────────────────────────────────
//  Register
// ─────────────────────────────────────────

func (uc *authUsecase) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	// Трейсинг: создаём дочерний span
	ctx, span := tracer.Start(ctx, "AuthUsecase.Register")
	defer span.End()

	span.SetAttributes(attribute.String("user.email", input.Email))

	// Метрика: инкрементируем счётчик попыток регистрации
	metrics.RegistrationAttempts.Inc()

	// 1. Валидация пароля
	if err := validatePassword(input.Password); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// 2. Проверяем, что email не занят
	existing, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil && err != domain.ErrUserNotFound {
		span.RecordError(err)
		return nil, fmt.Errorf("register: check email: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// 3. Хэшируем пароль через bcrypt (cost=12)
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("register: hash password: %w", err)
	}

	// 4. Формируем сущность пользователя
	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: string(hash),
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Phone:        input.Phone,
		Role:         domain.RoleClient, // новый пользователь — всегда клиент
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 5. Сохраняем в БД
	if err := uc.userRepo.Create(ctx, user); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("register: create user: %w", err)
	}

	// 6. Генерируем токен верификации email
	verifyToken, err := generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("register: generate verify token: %w", err)
	}
	evToken := &domain.EmailVerificationToken{
		UserID:    user.ID,
		Token:     verifyToken,
		ExpiresAt: now.Add(24 * time.Hour),
	}
	if err := uc.userRepo.SaveEmailVerificationToken(ctx, evToken); err != nil {
		return nil, fmt.Errorf("register: save verify token: %w", err)
	}

	// 7. Отправляем письмо (асинхронно — не блокируем ответ)
	go func() {
		bgCtx := context.Background()
		if err := uc.email.SendVerificationEmail(bgCtx, user.Email, user.FirstName, verifyToken); err != nil {
			uc.logger.Error("failed to send verification email",
				zap.String("user_id", user.ID.String()),
				zap.Error(err),
			)
		}
	}()

	// 8. Публикуем событие user.registered в NATS
	payload := domain.UserRegisteredPayload{
		UserID:     user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       user.Role,
		OccurredAt: now,
	}
	if err := uc.publisher.Publish(ctx, domain.EventUserRegistered, payload); err != nil {
		// Ошибка публикации не должна откатывать регистрацию — логируем
		uc.logger.Warn("failed to publish user.registered event",
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
	}

	metrics.RegistrationSuccess.Inc()
	uc.logger.Info("user registered",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)

	span.SetAttributes(attribute.String("user.id", user.ID.String()))
	return user, nil
}

// ─────────────────────────────────────────
//  Login
// ─────────────────────────────────────────

func (uc *authUsecase) Login(ctx context.Context, email, password string) (*domain.Token, *domain.User, error) {
	ctx, span := tracer.Start(ctx, "AuthUsecase.Login")
	defer span.End()
	span.SetAttributes(attribute.String("user.email", email))

	metrics.LoginAttempts.Inc()

	// 1. Ищем пользователя
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, nil, domain.ErrInvalidCredentials // не раскрываем причину
		}
		span.RecordError(err)
		return nil, nil, fmt.Errorf("login: get user: %w", err)
	}

	// 2. Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		metrics.LoginFailure.Inc()
		return nil, nil, domain.ErrInvalidCredentials
	}

	// 3. Проверяем верификацию email
	if !user.IsVerified {
		return nil, nil, domain.ErrEmailNotVerified
	}

	// 4. Генерируем JWT пару
	tokenPair, err := uc.tokens.GenerateTokenPair(user)
	if err != nil {
		span.RecordError(err)
		return nil, nil, fmt.Errorf("login: generate tokens: %w", err)
	}

	// 5. Сохраняем refresh-токен в Redis для последующей ротации
	if err := uc.cache.SetRefreshToken(ctx, user.ID, tokenPair.RefreshToken, tokenPair.ExpiresIn*10); err != nil {
		uc.logger.Warn("failed to cache refresh token", zap.String("user_id", user.ID.String()), zap.Error(err))
	}

	metrics.LoginSuccess.Inc()
	uc.logger.Info("user logged in", zap.String("user_id", user.ID.String()))
	span.SetAttributes(attribute.String("user.id", user.ID.String()))

	return tokenPair, user, nil
}

// ─────────────────────────────────────────
//  Logout
// ─────────────────────────────────────────

func (uc *authUsecase) Logout(ctx context.Context, accessToken string) error {
	ctx, span := tracer.Start(ctx, "AuthUsecase.Logout")
	defer span.End()

	// Парсим токен, чтобы получить jti и время истечения
	claims, err := uc.tokens.ParseAccessToken(accessToken)
	if err != nil {
		return domain.ErrTokenInvalid
	}

	// Добавляем jti в чёрный список Redis (TTL = оставшееся время жизни токена)
	ttl := int64(30 * 60) // fallback: 30 минут
	if err := uc.cache.BlacklistToken(ctx, claims.TokenID, ttl); err != nil {
		span.RecordError(err)
		return fmt.Errorf("logout: blacklist token: %w", err)
	}

	// Удаляем refresh-токен из Redis
	if err := uc.cache.DeleteRefreshToken(ctx, claims.UserID); err != nil {
		uc.logger.Warn("failed to delete refresh token on logout",
			zap.String("user_id", claims.UserID.String()),
			zap.Error(err),
		)
	}

	uc.logger.Info("user logged out", zap.String("user_id", claims.UserID.String()))
	return nil
}

// ─────────────────────────────────────────
//  RefreshToken
// ─────────────────────────────────────────

func (uc *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	ctx, span := tracer.Start(ctx, "AuthUsecase.RefreshToken")
	defer span.End()

	claims, err := uc.tokens.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	// Проверяем, что токен ещё хранится в Redis (не был инвалидирован)
	stored, err := uc.cache.GetRefreshToken(ctx, claims.UserID)
	if err != nil || stored != refreshToken {
		return nil, domain.ErrTokenInvalid
	}

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	newPair, err := uc.tokens.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("refresh: generate tokens: %w", err)
	}

	// Ротация: заменяем старый refresh на новый
	if err := uc.cache.SetRefreshToken(ctx, user.ID, newPair.RefreshToken, newPair.ExpiresIn*10); err != nil {
		uc.logger.Warn("failed to rotate refresh token", zap.Error(err))
	}

	return newPair, nil
}

// ─────────────────────────────────────────
//  VerifyEmail
// ─────────────────────────────────────────

func (uc *authUsecase) VerifyEmail(ctx context.Context, token string) error {
	ctx, span := tracer.Start(ctx, "AuthUsecase.VerifyEmail")
	defer span.End()

	evToken, err := uc.userRepo.GetEmailVerificationToken(ctx, token)
	if err != nil {
		return domain.ErrTokenInvalid
	}
	if time.Now().UTC().After(evToken.ExpiresAt) {
		return domain.ErrTokenExpired
	}

	user, err := uc.userRepo.GetByID(ctx, evToken.UserID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	user.IsVerified = true
	user.UpdatedAt = time.Now().UTC()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		span.RecordError(err)
		return fmt.Errorf("verify email: update user: %w", err)
	}

	// Удаляем использованный токен
	_ = uc.userRepo.DeleteEmailVerificationToken(ctx, token)

	// Инвалидируем кэш профиля
	_ = uc.cache.InvalidateUserProfile(ctx, user.ID)

	// Публикуем событие
	_ = uc.publisher.Publish(ctx, domain.EventEmailVerified, domain.EmailVerifiedPayload{
		UserID:     user.ID,
		Email:      user.Email,
		OccurredAt: time.Now().UTC(),
	})

	uc.logger.Info("email verified", zap.String("user_id", user.ID.String()))
	return nil
}

// ─────────────────────────────────────────
//  RequestPasswordReset
// ─────────────────────────────────────────

func (uc *authUsecase) RequestPasswordReset(ctx context.Context, email string) error {
	ctx, span := tracer.Start(ctx, "AuthUsecase.RequestPasswordReset")
	defer span.End()

	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Не сообщаем, существует ли email — защита от перебора
		return nil
	}

	resetToken, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("password reset: generate token: %w", err)
	}

	prToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: time.Now().UTC().Add(1 * time.Hour),
	}
	if err := uc.userRepo.SavePasswordResetToken(ctx, prToken); err != nil {
		return fmt.Errorf("password reset: save token: %w", err)
	}

	go func() {
		bgCtx := context.Background()
		if err := uc.email.SendPasswordResetEmail(bgCtx, user.Email, user.FirstName, resetToken); err != nil {
			uc.logger.Error("failed to send password reset email",
				zap.String("user_id", user.ID.String()),
				zap.Error(err),
			)
		}
	}()

	uc.logger.Info("password reset requested", zap.String("user_id", user.ID.String()))
	return nil
}

// ─────────────────────────────────────────
//  ConfirmPasswordReset
// ─────────────────────────────────────────

func (uc *authUsecase) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	ctx, span := tracer.Start(ctx, "AuthUsecase.ConfirmPasswordReset")
	defer span.End()

	prToken, err := uc.userRepo.GetPasswordResetToken(ctx, token)
	if err != nil {
		return domain.ErrTokenInvalid
	}
	if time.Now().UTC().After(prToken.ExpiresAt) {
		return domain.ErrTokenExpired
	}

	if err := validatePassword(newPassword); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("confirm reset: hash password: %w", err)
	}

	user, err := uc.userRepo.GetByID(ctx, prToken.UserID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	user.PasswordHash = string(hash)
	user.UpdatedAt = time.Now().UTC()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		span.RecordError(err)
		return fmt.Errorf("confirm reset: update user: %w", err)
	}

	_ = uc.userRepo.DeletePasswordResetToken(ctx, token)
	_ = uc.cache.InvalidateUserProfile(ctx, user.ID)

	uc.logger.Info("password reset confirmed", zap.String("user_id", user.ID.String()))
	return nil
}

// ─────────────────────────────────────────
//  Helpers
// ─────────────────────────────────────────

// generateSecureToken — 32 байта крипто-случайности в hex
func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// validatePassword — минимальные требования к паролю
func validatePassword(password string) error {
	if len(password) < 8 {
		return domain.ErrWeakPassword
	}
	var hasUpper, hasDigit bool
	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	if !hasUpper || !hasDigit {
		return domain.ErrWeakPassword
	}
	return nil
}

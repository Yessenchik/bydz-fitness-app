package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	// Укажи здесь правильный импорт твоего сгенерированного proto-кода
	pb "github.com/Yessenchik/bydz-fitness-app/user-auth-service/gen/userauth/v1"
)

func main() {
	// 1. Подключаемся к нашему gRPC микросервису
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC service: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserAuthServiceClient(conn)

	// 2. Создаем HTTP роутер
	r := gin.Default()

	// Настраиваем CORS (разрешаем запросы от React)
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"}, // Порт Vite
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	// 3. Эндпоинт Регистрации (переводит REST -> gRPC)
	r.POST("/api/auth/register", func(c *gin.Context) {
		var req pb.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request format"})
			return
		}

		res, err := client.Register(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	// 4. Эндпоинт Логина
	r.POST("/api/auth/login", func(c *gin.Context) {
		var req pb.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request format"})
			return
		}

		res, err := client.Login(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Неверный email или пароль. Либо почта не подтверждена."})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	// 5. Защищенный эндпоинт Профиля
	r.GET("/api/profile", func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "user_id is required"})
			return
		}

		// Берем токен из HTTP заголовка React'а и кладем в gRPC метадату
		authHeader := c.GetHeader("Authorization")
		ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", authHeader)

		res, err := client.GetProfile(ctx, &pb.GetProfileRequest{UserId: userID})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized or token expired"})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	// --- НОВЫЙ БЛОК: Эндпоинт подтверждения Email ---
	r.GET("/auth/verify-email", func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.String(http.StatusBadRequest, "Токен не найден в ссылке")
			return
		}

		// Вызываем gRPC метод VerifyEmail
		_, err := client.VerifyEmail(context.Background(), &pb.VerifyEmailRequest{
			VerificationToken: token,
		})

		if err != nil {
			// Если токен истек или неверный
			c.String(http.StatusBadRequest, "Ошибка подтверждения: ссылка устарела или недействительна")
			return
		}

		// Если всё успешно, перенаправляем пользователя на фронтенд (на страницу логина)
		// Добавляем ?verified=true, чтобы фронтенд мог показать сообщение "Почта подтверждена!"
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/login?verified=true")
	})
	// --------------------------------------------------

	// 6. Запускаем Gateway на порту 8080
	log.Println("✅ API Gateway running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

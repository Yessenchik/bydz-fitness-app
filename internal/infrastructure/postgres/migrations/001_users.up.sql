-- Расширение для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Тип для ролей
CREATE TYPE user_role AS ENUM ('client', 'trainer', 'admin');

-- Основная таблица пользователей
CREATE TABLE IF NOT EXISTS users (
                                     id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         VARCHAR(255) NOT NULL,
    password_hash TEXT         NOT NULL,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    phone         VARCHAR(20)  NOT NULL DEFAULT '',
    role          user_role    NOT NULL DEFAULT 'client',
    is_verified   BOOLEAN      NOT NULL DEFAULT FALSE,
    is_deleted    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_unique UNIQUE (email)
    );

-- Индексы
CREATE INDEX idx_users_email      ON users (email)      WHERE is_deleted = false;
CREATE INDEX idx_users_role       ON users (role)        WHERE is_deleted = false;
CREATE INDEX idx_users_created_at ON users (created_at DESC);

-- Токены сброса пароля
CREATE TABLE IF NOT EXISTS password_reset_tokens (
                                                     user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,

    CONSTRAINT prt_user_id_unique UNIQUE (user_id),
    CONSTRAINT prt_token_unique   UNIQUE (token)
    );

-- Токены верификации email
CREATE TABLE IF NOT EXISTS email_verification_tokens (
                                                         user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,

    CONSTRAINT evt_user_id_unique UNIQUE (user_id),
    CONSTRAINT evt_token_unique   UNIQUE (token)
    );

-- Audit log для критичных операций
CREATE TABLE IF NOT EXISTS user_audit_log (
                                              id          BIGSERIAL    PRIMARY KEY,
                                              user_id     UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action      VARCHAR(50)  NOT NULL,
    occurred_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_audit_user_id ON user_audit_log (user_id);
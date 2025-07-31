-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     email TEXT NOT NULL UNIQUE,
                                     password_hash TEXT NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

-- Индекс по email
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);

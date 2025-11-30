-- Инициализация базы данных
-- Этот файл выполняется АВТОМАТИЧЕСКИ при первом запуске PostgreSQL контейнера

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_moderator BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    email VARCHAR(255) UNIQUE,
    full_name VARCHAR(255),
    phone VARCHAR(50),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы кредитных продуктов (услуг)
-- ВАЖНО: image_url - TEXT (не VARCHAR) для поддержки data URI
CREATE TABLE IF NOT EXISTS credits (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    icon VARCHAR(500),
    image_url TEXT,
    type VARCHAR(100) NOT NULL,
    sum_from INTEGER NOT NULL,
    sum_to INTEGER NOT NULL,
    percent NUMERIC(5, 2) NOT NULL,
    term_from INTEGER NOT NULL,
    term_to INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы заявок (оценка кредитоспособности)
CREATE TABLE IF NOT EXISTS applications (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    formed_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    creator_id INTEGER NOT NULL,
    moderator_id INTEGER,
    total_amount NUMERIC(15, 2) DEFAULT 0,
    full_name VARCHAR(255),
    income NUMERIC(15, 2),
    obligations NUMERIC(15, 2),
    credit_score INTEGER DEFAULT 0,
    scoring_result VARCHAR(50),
    max_credit_amount NUMERIC(15, 2) DEFAULT 0,
    rejection_reason TEXT,
    dti NUMERIC(5, 2) DEFAULT 0
);

-- Создание таблицы связи заявок и услуг (M-M)
CREATE TABLE IF NOT EXISTS application_products (
    id SERIAL PRIMARY KEY,
    application_id INTEGER NOT NULL,
    credit_id INTEGER NOT NULL,
    requested_amount NUMERIC(15, 2) DEFAULT 0,
    requested_term_days INTEGER DEFAULT 0,
    monthly_payment NUMERIC(15, 2) DEFAULT 0,
    interest_rate NUMERIC(5, 2) DEFAULT 0,
    UNIQUE(application_id, credit_id),
    FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (credit_id) REFERENCES credits(id) ON DELETE CASCADE
);

-- Добавляем внешние ключи для applications
ALTER TABLE applications 
    DROP CONSTRAINT IF EXISTS fk_applications_creator,
    ADD CONSTRAINT fk_applications_creator FOREIGN KEY (creator_id) REFERENCES users(id);

ALTER TABLE applications 
    DROP CONSTRAINT IF EXISTS fk_applications_moderator,
    ADD CONSTRAINT fk_applications_moderator FOREIGN KEY (moderator_id) REFERENCES users(id);

-- Создаем расширение для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаем функцию для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Применяем триггер к таблице users
DROP TRIGGER IF EXISTS set_users_updated_at ON users;
CREATE TRIGGER set_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Применяем триггер к таблице credits
DROP TRIGGER IF EXISTS set_credits_updated_at ON credits;
CREATE TRIGGER set_credits_updated_at
BEFORE UPDATE ON credits
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

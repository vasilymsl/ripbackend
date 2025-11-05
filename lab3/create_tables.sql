-- Создание базовых таблиц для Lab3
-- Выполнить: docker exec -i lab3-postgres-1 psql -U myuser -d credits < create_tables.sql

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
CREATE TABLE IF NOT EXISTS credits (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    icon VARCHAR(500),
    type VARCHAR(100) NOT NULL,
    sum_from INTEGER NOT NULL,
    sum_to INTEGER NOT NULL,
    percent NUMERIC(5, 2) NOT NULL,
    term_from INTEGER NOT NULL,
    term_to INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы заявок
CREATE TABLE IF NOT EXISTS applications (
    id VARCHAR(50) PRIMARY KEY,
    user_id INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    formed_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    creator_id INTEGER NOT NULL,
    moderator_id INTEGER,
    total_amount NUMERIC(15, 2) DEFAULT 0,
    full_name VARCHAR(255),
    income NUMERIC(15, 2),
    obligations NUMERIC(15, 2)
);

-- Создание таблицы связи заявок и услуг (M-M)
CREATE TABLE IF NOT EXISTS application_products (
    application_id VARCHAR(50) NOT NULL,
    credit_id INTEGER NOT NULL,
    quantity INTEGER DEFAULT 1,
    order_num INTEGER DEFAULT 0,
    PRIMARY KEY (application_id, credit_id),
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

ALTER TABLE applications 
    DROP CONSTRAINT IF EXISTS fk_applications_user,
    ADD CONSTRAINT fk_applications_user FOREIGN KEY (user_id) REFERENCES users(id);

-- Создаем расширение для UUID (если понадобится)
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

-- Вставка тестовых пользователей
INSERT INTO users (username, password, is_moderator, created_at) VALUES
('testuser', '$2a$10$8K1p/a0dL3.I4ZL8bQ8Uve5Ls4lZx8o6mQZ2bWYqhFLlKLfPn3XqS', true, NOW())
ON CONFLICT (username) DO NOTHING;

-- Вставка тестовых услуг
INSERT INTO credits (title, icon, type, sum_from, sum_to, percent, term_from, term_to) VALUES
('Потребительский кредит', 'percent', 'Потребительский кредит', 10000, 500000, 15.9, 3, 60),
('Кредит наличными', 'wallet', 'Потребительский кредит', 50000, 3000000, 14.5, 6, 60),
('Кредит под залог', 'bank', 'Ипотека', 500000, 30000000, 12.5, 12, 180),
('Ипотека "Льготная"', 'home', 'Ипотека', 500000, 15000000, 7.5, 36, 360),
('Ипотека "Молодая семья"', 'family', 'Ипотека', 300000, 8000000, 6.9, 24, 360),
('Автокредит', 'car', 'Автокредит', 300000, 5000000, 16.5, 12, 60)
ON CONFLICT DO NOTHING;

-- Вывод информации
SELECT 'База данных инициализирована успешно!' AS message;
SELECT COUNT(*) AS users_count FROM users;
SELECT COUNT(*) AS credits_count FROM credits;


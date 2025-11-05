-- Миграция для лабораторной работы №3
-- Обновление схемы базы данных для REST API

-- Создание расширения для uuid (если не существует)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Обновление таблицы users
-- Добавляем поля для полноценной работы с пользователями
ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(255) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS full_name VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Обновление таблицы credits (услуги)
-- Удаляем старое поле image_url, так как будем хранить только имя файла
-- (полный URL будет формироваться динамически через MinIO)
ALTER TABLE credits ADD COLUMN IF NOT EXISTS image_name VARCHAR(255);
ALTER TABLE credits ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE credits ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Обновление таблицы applications
-- Поля уже должны быть корректными, проверим их наличие
ALTER TABLE applications ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'draft';
ALTER TABLE applications ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS formed_at TIMESTAMP;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS creator_id INTEGER;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS moderator_id INTEGER;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS full_name VARCHAR(255);
ALTER TABLE applications ADD COLUMN IF NOT EXISTS income VARCHAR(100);
ALTER TABLE applications ADD COLUMN IF NOT EXISTS obligations VARCHAR(100);
ALTER TABLE applications ADD COLUMN IF NOT EXISTS total_amount DECIMAL(15,2) DEFAULT 0;

-- Обновление таблицы application_products (многие-ко-многим)
-- Добавляем дополнительные поля для работы
ALTER TABLE application_products ADD COLUMN IF NOT EXISTS quantity INTEGER DEFAULT 1;
ALTER TABLE application_products ADD COLUMN IF NOT EXISTS order_num INTEGER DEFAULT 0;

-- Создание индексов для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_applications_status ON applications(status);
CREATE INDEX IF NOT EXISTS idx_applications_creator ON applications(creator_id);
CREATE INDEX IF NOT EXISTS idx_applications_formed_at ON applications(formed_at);
CREATE INDEX IF NOT EXISTS idx_credits_status ON credits(status);
CREATE INDEX IF NOT EXISTS idx_application_products_app_id ON application_products(application_id);

-- Создание функции для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создание триггеров для автоматического обновления updated_at
DROP TRIGGER IF EXISTS update_credits_updated_at ON credits;
CREATE TRIGGER update_credits_updated_at
    BEFORE UPDATE ON credits
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Комментарии к таблицам и столбцам для документации
COMMENT ON TABLE credits IS 'Таблица кредитных продуктов (услуг)';
COMMENT ON TABLE applications IS 'Таблица заявок на кредиты';
COMMENT ON TABLE application_products IS 'Связь многие-ко-многим между заявками и продуктами';
COMMENT ON TABLE users IS 'Таблица пользователей системы';

COMMENT ON COLUMN applications.status IS 'Статус: draft, deleted, formed, completed, rejected';
COMMENT ON COLUMN applications.total_amount IS 'Итоговая сумма, рассчитываемая при завершении заявки';
COMMENT ON COLUMN credits.status IS 'Статус: active, deleted';


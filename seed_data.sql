-- Создание пользователей
INSERT INTO users (username, password, is_moderator) VALUES
('testuser', 'password123', false),
('moderator', 'modpass123', true);

-- Заполнение таблицы credits (кредитные продукты)
-- Поля: title, icon, image_url, type, sum_from, sum_to, percent, term_from, term_to
INSERT INTO credits (title, icon, image_url, type, sum_from, sum_to, percent, term_from, term_to) VALUES
('Потребительский кредит', 'percent', 'https://via.placeholder.com/400x300/1e3a5f/ffffff?text=Potrebitelskiy+Kredit', 'consumer', 50000, 1000000, 17.90, 6, 36),
('Кредит наличными', 'wallet', 'https://via.placeholder.com/400x300/2c5aa0/ffffff?text=Kredit+Nalichnymi', 'cash', 100000, 3000000, 15.90, 12, 60),
('Кредит под залог', 'bank', 'https://via.placeholder.com/400x300/3a6fb0/ffffff?text=Kredit+Pod+Zalog', 'secured', 500000, 30000000, 12.50, 12, 180),
('Кредитная карта', 'card', 'https://via.placeholder.com/400x300/4a7fc1/ffffff?text=Kreditnaya+Karta', 'card', 30000, 600000, 19.90, 1, 60),
('Автокредит', 'car', 'https://via.placeholder.com/400x300/5a8fd1/ffffff?text=Avtokredit', 'auto', 200000, 5000000, 16.50, 12, 60),
('Ипотека', 'home', 'https://via.placeholder.com/400x300/6a9fe2/ffffff?text=Ipoteka', 'mortgage', 1000000, 50000000, 8.50, 60, 360);

-- Заполнение таблицы applications (заявки)
TRUNCATE TABLE applications RESTART IDENTITY CASCADE;
INSERT INTO applications (id, status, created_at, creator_id, full_name, income, obligations) VALUES
(1, 'formed', NOW(), 1, 'Иванов Иван Иванович', 120000, 25000),
(2, 'formed', NOW(), 1, 'Петров Петр Петрович', 85000, 15000),
(3, 'completed', NOW(), 1, 'Сидорова Анна Михайловна', 150000, 30000);

-- Сброс sequence для applications после вставки данных с явными ID
SELECT setval('applications_id_seq', (SELECT MAX(id) FROM applications));

-- Заполнение таблицы application_products (связь заявок и продуктов)
-- Новые поля: requested_amount, requested_term_days, monthly_payment (число), interest_rate
INSERT INTO application_products (application_id, credit_id, requested_amount, requested_term_days, monthly_payment, interest_rate) VALUES
(1, 1, 1000000, 1080, 35000, 17.9),  -- Потребительский кредит на 36 мес (1080 дней)
(1, 4, 300000, 60, 8500, 19.9),      -- Кредитная карта
(2, 2, 1500000, 1800, 18000, 15.9),  -- Кредит наличными на 60 мес (1800 дней)
(2, 5, 2500000, 1800, 22000, 16.5),  -- Автокредит на 60 мес
(3, 3, 5000000, 5475, 42000, 12.5),  -- Кредит под залог на 15 лет (5475 дней)
(3, 4, 200000, 60, 5000, 19.9);      -- Кредитная карта


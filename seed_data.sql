-- Создание пользователей
INSERT INTO users (username, password, is_moderator) VALUES
('testuser', 'password123', false),
('moderator', 'modpass123', true);

-- Заполнение таблицы credits (кредитные продукты)
INSERT INTO credits (title, icon, image_url, rate, term, amount, description, status) VALUES
('Потребительский кредит', 'percent', 'http://localhost:9000/images/%20credit/1.jpg', 'Ставка: от 17.9% годовых', 'Срок: до 36 мес', 'Сумма: до 1 000 000 ₽', 'Пользовательский кредит — универсальный продукт, который можно оформить на любые личные нужды без залога и поручителей.', 'active'),
('Кредит наличными', 'wallet', 'http://localhost:9000/images/%20credit/2.jpg', 'Ставка: от 15.9% годовых', 'Срок: до 60 мес', 'Сумма: до 3 000 000 ₽', 'Кредит наличными — быстрое решение для получения денежных средств на карту или наличными.', 'active'),
('Кредит под залог на любые цели', 'bank', 'http://localhost:9000/images/%20credit/3.jpg', 'Ставка: от 12.5% годовых', 'Срок: до 15 лет', 'Сумма: до 30 000 000 ₽', 'Кредит под залог недвижимости — выгодное решение для получения крупной суммы на длительный срок.', 'active'),
('Кредитная карта', 'card', 'http://localhost:9000/images/%20credit/4.jpg', 'Ставка: от 19.9% годовых', 'Льготный период: до 60 дней', 'Лимит: до 600 000 ₽', 'Кредитная карта — удобный инструмент для ежедневных покупок и непредвиденных расходов.', 'active'),
('Автокредит', 'car', 'http://localhost:9000/images/%20credit/5.webp', 'Ставка: от 16.5% годовых', 'Срок: до 60 мес', 'Сумма: до 5 000 000 ₽', 'Автокредит — специальное предложение для покупки нового или подержанного автомобиля.', 'active');

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


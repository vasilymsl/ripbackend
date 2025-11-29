-- Создание пользователей
INSERT INTO users (username, password, is_moderator) VALUES
('testuser', 'password123', false),
('moderator', 'modpass123', true);

-- Заполнение таблицы credits (кредитные продукты)
-- Поля: title, icon, image_url, type, sum_from, sum_to, percent, term_from, term_to
-- Используем data URI со встроенными SVG (работают без интернета)
INSERT INTO credits (title, icon, image_url, type, sum_from, sum_to, percent, term_from, term_to) VALUES
('Потребительский кредит', 'percent', 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MDAiIGhlaWdodD0iMzAwIiB2aWV3Qm94PSIwIDAgNDAwIDMwMCI+CiAgPHJlY3Qgd2lkdGg9IjQwMCIgaGVpZ2h0PSIzMDAiIGZpbGw9IiMxZTNhNWYiLz4KICA8Y2lyY2xlIGN4PSIyMDAiIGN5PSIxMjAiIHI9IjUwIiBmaWxsPSJ3aGl0ZSIgb3BhY2l0eT0iMC4yIi8+CiAgPHRleHQgeD0iMjAwIiB5PSIxODAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIyMiIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGZvbnQtd2VpZ2h0PSJib2xkIj7Qn9C+0YLRgNC10LHQuNGC0LXQu9GM0YHQutC40Lkg0LrRgNC10LTQuNGCPC90ZXh0PgogIDx0ZXh0IHg9IjIwMCIgeT0iMjIwIiBmb250LWZhbWlseT0iQXJpYWwsIHNhbnMtc2VyaWYiIGZvbnQtc2l6ZT0iMTYiIGZpbGw9IndoaXRlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBvcGFjaXR5PSIwLjgiPtCa0YDQtdC00LjRgtC90YvQuSDQv9GA0L7QtNGD0LrRgjwvdGV4dD4KPC9zdmc+', 'consumer', 50000, 1000000, 17.90, 6, 36),
('Кредит наличными', 'wallet', 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MDAiIGhlaWdodD0iMzAwIiB2aWV3Qm94PSIwIDAgNDAwIDMwMCI+CiAgPHJlY3Qgd2lkdGg9IjQwMCIgaGVpZ2h0PSIzMDAiIGZpbGw9IiMyYzVhYTAiLz4KICA8Y2lyY2xlIGN4PSIyMDAiIGN5PSIxMjAiIHI9IjUwIiBmaWxsPSJ3aGl0ZSIgb3BhY2l0eT0iMC4yIi8+CiAgPHRleHQgeD0iMjAwIiB5PSIxODAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIyMiIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGZvbnQtd2VpZ2h0PSJib2xkIj7QmtGA0LXQtNC40YIg0L3QsNC70LjRh9C90YvQvNC4PC90ZXh0PgogIDx0ZXh0IHg9IjIwMCIgeT0iMjIwIiBmb250LWZhbWlseT0iQXJpYWwsIHNhbnMtc2VyaWYiIGZvbnQtc2l6ZT0iMTYiIGZpbGw9IndoaXRlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBvcGFjaXR5PSIwLjgiPtCa0YDQtdC00LjRgtC90YvQuSDQv9GA0L7QtNGD0LrRgjwvdGV4dD4KPC9zdmc+', 'cash', 100000, 3000000, 15.90, 12, 60),
('Кредит под залог', 'bank', 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MDAiIGhlaWdodD0iMzAwIiB2aWV3Qm94PSIwIDAgNDAwIDMwMCI+CiAgPHJlY3Qgd2lkdGg9IjQwMCIgaGVpZ2h0PSIzMDAiIGZpbGw9IiMzYTZmYjAiLz4KICA8Y2lyY2xlIGN4PSIyMDAiIGN5PSIxMjAiIHI9IjUwIiBmaWxsPSJ3aGl0ZSIgb3BhY2l0eT0iMC4yIi8+CiAgPHRleHQgeD0iMjAwIiB5PSIxODAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIyMiIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGZvbnQtd2VpZ2h0PSJib2xkIj7QmtGA0LXQtNC40YIg0L/QvtC0INC30LDQu9C+0LM8L3RleHQ+CiAgPHRleHQgeD0iMjAwIiB5PSIyMjAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIxNiIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiIG9wYWNpdHk9IjAuOCI+0JrRgNC10LTQuNGC0L3Ri9C5INC/0YDQvtC00YPQutGCPC90ZXh0Pgo8L3N2Zz4=', 'secured', 500000, 30000000, 12.50, 12, 180),
('Кредитная карта', 'card', 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MDAiIGhlaWdodD0iMzAwIiB2aWV3Qm94PSIwIDAgNDAwIDMwMCI+CiAgPHJlY3Qgd2lkdGg9IjQwMCIgaGVpZ2h0PSIzMDAiIGZpbGw9IiM0YTdmYzEiLz4KICA8Y2lyY2xlIGN4PSIyMDAiIGN5PSIxMjAiIHI9IjUwIiBmaWxsPSJ3aGl0ZSIgb3BhY2l0eT0iMC4yIi8+CiAgPHRleHQgeD0iMjAwIiB5PSIxODAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIyMiIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGZvbnQtd2VpZ2h0PSJib2xkIj7QmtGA0LXQtNC40YLQvdCw0Y8g0LrQsNGA0YLQsDwvdGV4dD4KICA8dGV4dCB4PSIyMDAiIHk9IjIyMCIgZm9udC1mYW1pbHk9IkFyaWFsLCBzYW5zLXNlcmlmIiBmb250LXNpemU9IjE2IiBmaWxsPSJ3aGl0ZSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgb3BhY2l0eT0iMC44Ij7QmtGA0LXQtNC40YLQvdGL0Lkg0L/RgNC+0LTRg9C60YI8L3RleHQ+Cjwvc3ZnPg==', 'card', 30000, 600000, 19.90, 1, 60),
('Автокредит', 'car', 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MDAiIGhlaWdodD0iMzAwIiB2aWV3Qm94PSIwIDAgNDAwIDMwMCI+CiAgPHJlY3Qgd2lkdGg9IjQwMCIgaGVpZ2h0PSIzMDAiIGZpbGw9IiM1YThmZDEiLz4KICA8Y2lyY2xlIGN4PSIyMDAiIGN5PSIxMjAiIHI9IjUwIiBmaWxsPSJ3aGl0ZSIgb3BhY2l0eT0iMC4yIi8+CiAgPHRleHQgeD0iMjAwIiB5PSIxODAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIyMiIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGZvbnQtd2VpZ2h0PSJib2xkIj7QkNCy0YLQvtC60YDQtdC00LjRgjwvdGV4dD4KICA8dGV4dCB4PSIyMDAiIHk9IjIyMCIgZm9udC1mYW1pbHk9IkFyaWFsLCBzYW5zLXNlcmlmIiBmb250LXNpemU9IjE2IiBmaWxsPSJ3aGl0ZSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgb3BhY2l0eT0iMC44Ij7QmtGA0LXQtNC40YLQvdGL0Lkg0L/RgNC+0LTRg9C60YI8L3RleHQ+Cjwvc3ZnPg==', 'auto', 200000, 5000000, 16.50, 12, 60),
('Ипотека', 'home', 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MDAiIGhlaWdodD0iMzAwIiB2aWV3Qm94PSIwIDAgNDAwIDMwMCI+CiAgPHJlY3Qgd2lkdGg9IjQwMCIgaGVpZ2h0PSIzMDAiIGZpbGw9IiM2YTlmZTIiLz4KICA8Y2lyY2xlIGN4PSIyMDAiIGN5PSIxMjAiIHI9IjUwIiBmaWxsPSJ3aGl0ZSIgb3BhY2l0eT0iMC4yIi8+CiAgPHRleHQgeD0iMjAwIiB5PSIxODAiIGZvbnQtZmFtaWx5PSJBcmlhbCwgc2Fucy1zZXJpZiIgZm9udC1zaXplPSIyMiIgZmlsbD0id2hpdGUiIHRleHQtYW5jaG9yPSJtaWRkbGUiIGZvbnQtd2VpZ2h0PSJib2xkIj7QmNC/0L7RgtC10LrQsDwvdGV4dD4KICA8dGV4dCB4PSIyMDAiIHk9IjIyMCIgZm9udC1mYW1pbHk9IkFyaWFsLCBzYW5zLXNlcmlmIiBmb250LXNpemU9IjE2IiBmaWxsPSJ3aGl0ZSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgb3BhY2l0eT0iMC44Ij7QmtGA0LXQtNC40YLQvdGL0Lkg0L/RgNC+0LTRg9C60YI8L3RleHQ+Cjwvc3ZnPg==', 'mortgage', 1000000, 50000000, 8.50, 60, 360);

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


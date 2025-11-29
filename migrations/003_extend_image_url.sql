-- Расширяем поле image_url для поддержки data URI (встроенных SVG)
ALTER TABLE credits ALTER COLUMN image_url TYPE TEXT;


-- Миграция для добавления полей оценки кредитоспособности
-- Тема: Оценка кредитоспособности заемщика по скоринговой модели

-- Добавляем поля для результата скоринга
ALTER TABLE applications ADD COLUMN IF NOT EXISTS credit_score INTEGER DEFAULT 0;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS scoring_result VARCHAR(50) DEFAULT NULL;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS max_credit_amount DECIMAL(15,2) DEFAULT 0;
ALTER TABLE applications ADD COLUMN IF NOT EXISTS rejection_reason TEXT DEFAULT NULL;

-- Добавляем индекс для быстрого поиска по результату скоринга
CREATE INDEX IF NOT EXISTS idx_applications_scoring_result ON applications(scoring_result);
CREATE INDEX IF NOT EXISTS idx_applications_credit_score ON applications(credit_score DESC);

-- Комментарии для документации
COMMENT ON COLUMN applications.credit_score IS 'Кредитный скор заемщика (0-1000). Рассчитывается автоматически при формировании заявки';
COMMENT ON COLUMN applications.scoring_result IS 'Результат оценки: approved (одобрено), rejected (отклонено), pending (ожидает оценки)';
COMMENT ON COLUMN applications.max_credit_amount IS 'Максимальная рекомендованная сумма кредита на основе дохода и обязательств';
COMMENT ON COLUMN applications.rejection_reason IS 'Причина отклонения заявки (если scoring_result = rejected)';


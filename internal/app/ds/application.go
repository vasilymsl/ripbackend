package ds

import "time"

const (
	ApplicationStatusDraft     = "draft"     // черновик
	ApplicationStatusDeleted   = "deleted"   // удалён
	ApplicationStatusFormed    = "formed"    // сформирован
	ApplicationStatusCompleted = "completed" // завершён
	ApplicationStatusRejected  = "rejected"  // отклонён
)

type Application struct {
	ID          int        `gorm:"primaryKey;autoIncrement"` // Числовой ID с автоинкрементом
	Status      string     `gorm:"not null;default:'draft'"` // Обязательное поле
	CreatedAt   time.Time  `gorm:"not null;autoCreateTime"`  // Обязательное поле
	CreatorID   int        `gorm:"not null"`                 // Обязательное поле
	Creator     User       `gorm:"foreignKey:CreatorID"`
	FormedAt    *time.Time `gorm:"default:null"` // Дата формирования (nullable)
	CompletedAt *time.Time `gorm:"default:null"` // Дата завершения (nullable)
	ModeratorID *int       `gorm:"default:null"` // Модератор (nullable)
	Moderator   *User      `gorm:"foreignKey:ModeratorID"`

	// Поля по предметной области (обязательные при формировании заявки)
	FullName    string  `gorm:"column:full_name;not null"` // ФИО обязательно
	Income      float64 `gorm:"not null"`                  // Доход обязательно (число!)
	Obligations float64 `gorm:"not null"`                  // Обязательства обязательно (число!)
	TotalAmount float64 `gorm:"default:0"`                 // Рассчитываемое поле при завершении

	// Поля скоринговой модели оценки кредитоспособности
	CreditScore      int     `gorm:"column:credit_score;default:0"`       // Кредитный скор (0-1000)
	ScoringResult    string  `gorm:"column:scoring_result"`               // approved/rejected/pending
	MaxCreditAmount  float64 `gorm:"column:max_credit_amount;default:0"`  // Макс. рекомендованная сумма
	RejectionReason  string  `gorm:"column:rejection_reason"`             // Причина отклонения

	Products []ApplicationProduct `gorm:"foreignKey:ApplicationID"`
}

// TableName переопределяет имя таблицы для GORM
func (Application) TableName() string {
	return "applications"
}

type ApplicationProduct struct {
	ID                int     `gorm:"primaryKey;autoIncrement"`
	ApplicationID     int     `gorm:"column:application_id;not null;uniqueIndex:idx_app_credit"` // Составной уникальный ключ
	CreditID          int     `gorm:"column:credit_id;not null;uniqueIndex:idx_app_credit"`      // Составной уникальный ключ
	RequestedAmount   float64 `gorm:"column:requested_amount"`                                   // Запрашиваемая сумма
	RequestedTermDays int     `gorm:"column:requested_term_days"`                                // Срок кредита в днях
	MonthlyPayment    float64 `gorm:"column:monthly_payment"`                                    // Ежемесячный платеж (число!)
	InterestRate      float64 `gorm:"column:interest_rate"`                                      // Процентная ставка
	Order             Order   `gorm:"foreignKey:CreditID"`
}

// TableName переопределяет имя таблицы для GORM
func (ApplicationProduct) TableName() string {
	return "application_products"
}

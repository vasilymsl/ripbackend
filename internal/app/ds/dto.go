package ds

import "time"

// ============================================
// DTO для Orders (Кредитные продукты / Услуги)
// ============================================

// OrderCreateRequest - запрос на создание услуги (без изображения)
type OrderCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Icon        string `json:"icon"`
	Rate        string `json:"rate"`
	Term        string `json:"term"`
	Amount      string `json:"amount"`
	Feature1    string `json:"feature1"`
	Feature2    string `json:"feature2"`
	Feature3    string `json:"feature3"`
	Description string `json:"description"`
}

// OrderUpdateRequest - запрос на обновление услуги
type OrderUpdateRequest struct {
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	Rate        string `json:"rate"`
	Term        string `json:"term"`
	Amount      string `json:"amount"`
	Feature1    string `json:"feature1"`
	Feature2    string `json:"feature2"`
	Feature3    string `json:"feature3"`
	Description string `json:"description"`
}

// OrderResponse - ответ с данными услуги
type OrderResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Icon        string    `json:"icon"`
	ImageURL    string    `json:"image_url"`
	Rate        string    `json:"rate"`
	Term        string    `json:"term"`
	Amount      string    `json:"amount"`
	Feature1    string    `json:"feature1"`
	Feature2    string    `json:"feature2"`
	Feature3    string    `json:"feature3"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrderListResponse - ответ со списком услуг
type OrderListResponse struct {
	Total  int             `json:"total"`
	Orders []OrderResponse `json:"orders"`
}

// ============================================
// DTO для Applications (Заявки)
// ============================================

// ApplicationUpdateRequest - запрос на изменение полей заявки
type ApplicationUpdateRequest struct {
	FullName    string  `json:"full_name"`
	Income      float64 `json:"income"`
	Obligations float64 `json:"obligations"`
}

// ApplicationResponse - ответ с данными заявки
type ApplicationResponse struct {
	ID           int                     `json:"id"`
	Status       string                  `json:"status"`
	CreatedAt    time.Time               `json:"created_at"`
	FormedAt     *time.Time              `json:"formed_at,omitempty"`
	CompletedAt  *time.Time              `json:"completed_at,omitempty"`
	CreatorLogin string                  `json:"creator_login"`
	CreatorID    int                     `json:"creator_id"`
	Moderator    *string                 `json:"moderator_login,omitempty"`
	ModeratorID  *int                    `json:"moderator_id,omitempty"`
	FullName     string                  `json:"full_name"`
	Income       float64                 `json:"income"`
	Obligations  float64                 `json:"obligations"`
	TotalAmount  float64                 `json:"total_amount"`
	
	// Результаты оценки кредитоспособности
	CreditScore     int     `json:"credit_score"`                // Кредитный скор (0-1000)
	ScoringResult   string  `json:"scoring_result,omitempty"`    // approved/rejected/pending
	MaxCreditAmount float64 `json:"max_credit_amount"`           // Макс. рекомендованная сумма
	RejectionReason string  `json:"rejection_reason,omitempty"`  // Причина отклонения
	
	Products     []ApplicationProductDTO `json:"products,omitempty"`
}

// ApplicationListResponse - ответ со списком заявок
type ApplicationListResponse struct {
	Total        int                   `json:"total"`
	Applications []ApplicationResponse `json:"applications"`
}

// ApplicationBasketResponse - ответ с информацией о корзине
type ApplicationBasketResponse struct {
	ApplicationID int `json:"application_id"`
	ProductCount  int `json:"product_count"`
}

// ApplicationFormRequest - запрос на формирование заявки
type ApplicationFormRequest struct {
	// Пустая структура, т.к. формирование происходит автоматически
}

// ApplicationCompleteRequest - запрос на завершение/отклонение заявки
type ApplicationCompleteRequest struct {
	Action string `json:"action" binding:"required,oneof=complete reject"` // "complete" или "reject"
}

// ============================================
// DTO для ApplicationProducts (М-М связь)
// ============================================

// ApplicationProductDTO - данные о продукте в заявке
type ApplicationProductDTO struct {
	ID                int     `json:"id"`
	CreditID          int     `json:"credit_id"`
	CreditTitle       string  `json:"credit_title"`
	ImageURL          string  `json:"image_url"`
	RequestedAmount   float64 `json:"requested_amount"`
	RequestedTermDays int     `json:"requested_term_days"`
	MonthlyPayment    float64 `json:"monthly_payment"`
	InterestRate      float64 `json:"interest_rate"`
}

// AddToApplicationRequest - запрос на добавление продукта в заявку
type AddToApplicationRequest struct {
	CreditID int `json:"credit_id" binding:"required"`
}

// UpdateApplicationProductRequest - запрос на изменение М-М связи
type UpdateApplicationProductRequest struct {
	RequestedAmount   *float64 `json:"requested_amount,omitempty"`
	RequestedTermDays *int     `json:"requested_term_days,omitempty"`
	MonthlyPayment    *float64 `json:"monthly_payment,omitempty"`
	InterestRate      *float64 `json:"interest_rate,omitempty"`
}

// ============================================
// DTO для Users (Пользователи)
// ============================================

// UserRegisterRequest - запрос на регистрацию пользователя
type UserRegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=6"`
	Email       string `json:"email" binding:"required,email"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	IsModerator bool   `json:"is_moderator"`
}

// UserLoginRequest - запрос на аутентификацию
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest - запрос на обновление профиля пользователя
type UserUpdateRequest struct {
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	FullName string `json:"full_name,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}

// UserResponse - ответ с данными пользователя
type UserResponse struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	Phone       string    `json:"phone"`
	IsModerator bool      `json:"is_moderator"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserLoginResponse - ответ на успешную аутентификацию
type UserLoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"` // Для упрощения - просто ID пользователя в виде строки
}

// ============================================
// Общие структуры ответов
// ============================================

// ErrorResponse - структура для ответов с ошибкой
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SuccessResponse - структура для успешных ответов
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

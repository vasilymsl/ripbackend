package ds

import (
	"fmt"
	"time"
)

// Константы статусов для Order (для совместимости с lab2)
const (
	OrderStatusActive  = "active"
	OrderStatusDeleted = "deleted"
)

type Order struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"column:title;not null" json:"title"`
	Icon        string    `gorm:"column:icon" json:"icon"`
	Type        string    `gorm:"column:type;not null" json:"type"`
	SumFrom     int       `gorm:"column:sum_from;not null" json:"sum_from"`
	SumTo       int       `gorm:"column:sum_to;not null" json:"sum_to"`
	Percent     float64   `gorm:"column:percent;not null" json:"percent"`
	TermFrom    int       `gorm:"column:term_from;not null" json:"term_from"`
	TermTo      int       `gorm:"column:term_to;not null" json:"term_to"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	
	// Поля для обратной совместимости с lab2 (computed fields или дополнительные)
	ImageURL string `gorm:"-" json:"-"`          // Генерируется из Icon
	Rate     string `gorm:"-" json:"-"`          // Генерируется из Percent
	Term     string `gorm:"-" json:"-"`          // Генерируется из TermFrom/TermTo
	Amount   string `gorm:"-" json:"-"`          // Генерируется из SumFrom/SumTo
	Feature1 string `gorm:"column:feature1" json:"feature1,omitempty"` // Дополнительные характеристики (опционально)
	Feature2 string `gorm:"column:feature2" json:"feature2,omitempty"`
	Feature3 string `gorm:"column:feature3" json:"feature3,omitempty"`
	Status   string `gorm:"column:status;default:'active'" json:"status"` // Статус (active/deleted)
}

// TableName переопределяет имя таблицы для GORM
func (Order) TableName() string {
	return "credits"
}

// PopulateCompatibilityFields заполняет поля для совместимости с lab2
func (o *Order) PopulateCompatibilityFields() {
	// ImageURL = Icon (в lab2 они одинаковые)
	o.ImageURL = o.Icon
	
	// Rate - процентная ставка
	o.Rate = fmt.Sprintf("Ставка от %.2f%%", o.Percent)
	
	// Term - срок кредита
	if o.TermFrom == o.TermTo {
		o.Term = fmt.Sprintf("Срок: %d мес", o.TermFrom)
	} else {
		o.Term = fmt.Sprintf("Срок: от %d до %d мес", o.TermFrom, o.TermTo)
	}
	
	// Amount - сумма кредита
	if o.SumFrom == o.SumTo {
		o.Amount = fmt.Sprintf("Сумма: %d ₽", o.SumFrom)
	} else {
		o.Amount = fmt.Sprintf("Сумма: от %d до %d ₽", o.SumFrom, o.SumTo)
	}
}

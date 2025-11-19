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
	ImageURL    string    `gorm:"column:image_url" json:"image_url"`  // URL картинки в MinIO
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

// getSourcesDescription возвращает описание источников данных для услуги
func (o *Order) getSourcesDescription() string {
	switch o.Title {
	case "Базовый скоринг заемщика":
		return "Источники: КИ, соц.-дем. данные, поведенческие паттерны"
	case "Анализ транзакций":
		return "Доходы/расходы, регулярность, финансовая устойчивость"
	case "Залоговый скоринг":
		return "Ликвидность объекта в прогнозе PD"
	case "Ипотечный скоринг":
		return "Анализ заемщика + объекта"
	case "Скоринг самозанятых":
		return "Анализ операций по кошелькам/платёжкам"
	default:
		return ""
	}
}

// PopulateCompatibilityFields заполняет поля для совместимости с lab2
// Адаптировано для платформы оценки кредитоспособности
func (o *Order) PopulateCompatibilityFields() {
	// Для главной страницы (карточки) - полные описания
	// Rate - точность модели (используем Percent как точность)
	if o.Percent > 0 {
		o.Rate = fmt.Sprintf("Точность: до %.0f%%", o.Percent)
	} else {
		o.Rate = "Риск-модель: заемщик + залог"
	}
	
	// Term - время расчета или срок внедрения
	if o.TermFrom == o.TermTo && o.TermFrom < 10 {
		// Для быстрых расчетов (в секундах)
		if o.TermFrom == 2 {
			o.Term = "Время расчёта: < 2 сек"
		} else if o.TermFrom >= 1 && o.TermFrom <= 3 {
			o.Term = fmt.Sprintf("Время: %d–3 сек", o.TermFrom)
		} else {
			o.Term = fmt.Sprintf("Срок внедрения: до %d дней", o.TermFrom)
		}
	} else if o.TermFrom == 120 {
		// Для ипотечного скоринга
		o.Term = "120+ факторов"
	} else if o.TermFrom == 0 {
		// Для скоринга самозанятых
		o.Term = "Доходы без справок"
	} else {
		o.Term = fmt.Sprintf("Срок внедрения: до %d дней", o.TermTo)
	}
	
	// Amount - источники данных или особенности
	o.Amount = o.getSourcesDescription()
}

// PopulateForTable заполняет поля для отображения в таблице заявки (короткие значения)
func (o *Order) PopulateForTable() {
	// Для таблицы - короткие значения
	if o.Percent > 0 {
		o.Rate = fmt.Sprintf("до %.0f%%", o.Percent)
	} else {
		o.Rate = "—"
	}
	
	if o.TermFrom == o.TermTo && o.TermFrom < 10 {
		if o.TermFrom == 2 {
			o.Term = "< 2 сек"
		} else if o.TermFrom >= 1 && o.TermFrom <= 3 {
			o.Term = fmt.Sprintf("%d–3 сек", o.TermFrom)
		} else {
			o.Term = fmt.Sprintf("до %d дней", o.TermFrom)
		}
	} else if o.TermFrom == 120 {
		o.Term = "до 120 дней"
	} else if o.TermFrom == 0 {
		o.Term = "мгновенно"
	} else {
		o.Term = fmt.Sprintf("до %d дней", o.TermTo)
	}
	
	o.Amount = ""
}

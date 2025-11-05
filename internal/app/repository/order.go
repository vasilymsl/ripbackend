package repository

import (
	"errors"
	"fmt"
	"lab1/internal/app/ds"

	"gorm.io/gorm"
)

// ============================================
// Методы репозитория для Orders (Услуги)
// ============================================

// GetOrders возвращает список всех услуг (для совместимости с lab2)
func (r *Repository) GetOrders() ([]ds.Order, error) {
	var orders []ds.Order
	err := r.db.Order("id ASC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}
	return orders, nil
}

// GetOrdersWithFilter возвращает список всех услуг с фильтрацией
func (r *Repository) GetOrdersWithFilter(titleFilter string) ([]ds.Order, error) {
	var orders []ds.Order
	query := r.db

	// Фильтрация по названию, если передан параметр
	if titleFilter != "" {
		query = query.Where("title ILIKE ?", "%"+titleFilter+"%")
	}

	err := query.Order("id ASC").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrdersByTitle возвращает список услуг по названию (для совместимости с lab2)
func (r *Repository) GetOrdersByTitle(title string) ([]ds.Order, error) {
	var orders []ds.Order
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Order("id ASC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrder возвращает услугу по ID (для совместимости с lab2)
func (r *Repository) GetOrder(id int) (ds.Order, error) {
	order := ds.Order{}
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return ds.Order{}, err
	}
	return order, nil
}

// GetOrderByID возвращает услугу по ID
func (r *Repository) GetOrderByID(id int) (ds.Order, error) {
	var order ds.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Order{}, fmt.Errorf("order with id %d not found", id)
		}
		return ds.Order{}, err
	}
	return order, nil
}

// CreateOrder создает новую услугу
func (r *Repository) CreateOrder(order *ds.Order) error {
	return r.db.Create(order).Error
}

// UpdateOrder обновляет существующую услугу
func (r *Repository) UpdateOrder(order *ds.Order) error {
	// Проверяем, что услуга существует
	_, err := r.GetOrderByID(order.ID)
	if err != nil {
		return err
	}

	// Обновляем только разрешенные поля
	updates := map[string]interface{}{
		"title":     order.Title,
		"icon":      order.Icon,
		"type":      order.Type,
		"sum_from":  order.SumFrom,
		"sum_to":    order.SumTo,
		"percent":   order.Percent,
		"term_from": order.TermFrom,
		"term_to":   order.TermTo,
	}

	return r.db.Model(&ds.Order{}).Where("id = ?", order.ID).Updates(updates).Error
}

// UpdateOrderImage обновляет изображение услуги
func (r *Repository) UpdateOrderImage(orderID int, imageURL string) error {
	return r.db.Model(&ds.Order{}).Where("id = ?", orderID).Update("icon", imageURL).Error
}

// DeleteOrder удаляет услугу
func (r *Repository) DeleteOrder(id int) error {
	// Проверяем существование услуги
	_, err := r.GetOrderByID(id)
	if err != nil {
		return err
	}

	return r.db.Where("id = ?", id).Delete(&ds.Order{}).Error
}

// GetOrdersByIDs возвращает список услуг по массиву ID
func (r *Repository) GetOrdersByIDs(ids []int) ([]ds.Order, error) {
	var orders []ds.Order
	err := r.db.Where("id IN ?", ids).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

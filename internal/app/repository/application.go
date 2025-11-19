package repository

import (
	"errors"
	"fmt"
	"lab1/internal/app/ds"
	"time"

	"gorm.io/gorm"
)

// ============================================
// Методы репозитория для Applications (Заявки)
// ============================================

// GetDraftApplication получает заявку-черновик для пользователя
func (r *Repository) GetDraftApplication(userID int) (*ds.Application, error) {
	var app ds.Application
	err := r.db.Preload("Products.Order").
		Where("creator_id = ? AND status = ?", userID, ds.ApplicationStatusDraft).
		First(&app).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Черновика нет
		}
		return nil, err
	}

	return &app, nil
}

// CreateDraftApplication создает новую заявку-черновик
func (r *Repository) CreateDraftApplication(userID int) (*ds.Application, error) {
	app := ds.Application{
		Status:      ds.ApplicationStatusDraft,
		CreatorID:   userID,
		FullName:    "",  // Будет заполнено пользователем
		Income:      0,   // Будет заполнено пользователем
		Obligations: 0,   // Будет заполнено пользователем
		TotalAmount: 0,
	}

	if err := r.db.Create(&app).Error; err != nil {
		return nil, fmt.Errorf("ошибка создания заявки: %w", err)
	}

	return &app, nil
}

// GetApplications возвращает список заявок с фильтрацией
// Исключает черновики и удаленные заявки
func (r *Repository) GetApplications(statusFilter string, dateFrom, dateTo *time.Time) ([]ds.Application, error) {
	var applications []ds.Application

	query := r.db.Preload("Products.Order").
		Preload("Creator").
		Preload("Moderator").
		Where("status != ? AND status != ?", ds.ApplicationStatusDraft, ds.ApplicationStatusDeleted)

	// Фильтрация по статусу
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	// Фильтрация по диапазону дат формирования
	if dateFrom != nil {
		query = query.Where("formed_at >= ?", dateFrom)
	}
	if dateTo != nil {
		query = query.Where("formed_at <= ?", dateTo)
	}

	err := query.Order("created_at DESC").Find(&applications).Error
	if err != nil {
		return nil, err
	}

	return applications, nil
}

// GetApplicationByID возвращает заявку по ID (с продуктами)
func (r *Repository) GetApplicationByID(id int) (*ds.Application, error) {
	var app ds.Application
	err := r.db.Preload("Products.Order").
		Preload("Creator").
		Preload("Moderator").
		Where("id = ? AND status != ?", id, ds.ApplicationStatusDeleted).
		First(&app).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("application not found")
		}
		return nil, err
	}

	return &app, nil
}

// GetApplication возвращает заявку по ID (для совместимости с lab2, возвращает не указатель)
func (r *Repository) GetApplication(id int) (ds.Application, error) {
	appPtr, err := r.GetApplicationByID(id)
	if err != nil {
		return ds.Application{}, err
	}
	if appPtr == nil {
		return ds.Application{}, fmt.Errorf("application not found")
	}
	return *appPtr, nil
}

// UpdateApplication обновляет поля заявки
func (r *Repository) UpdateApplication(app *ds.Application) error {
	// Проверяем, что заявка существует и в статусе черновик
	existing, err := r.GetApplicationByID(app.ID)
	if err != nil {
		return err
	}

	if existing.Status != ds.ApplicationStatusDraft {
		return fmt.Errorf("can only update draft applications")
	}

	// Обновляем только бизнес-поля
	updates := map[string]interface{}{
		"full_name":   app.FullName,
		"income":      app.Income,
		"obligations": app.Obligations,
	}

	return r.db.Model(&ds.Application{}).Where("id = ?", app.ID).Updates(updates).Error
}

// FormApplication формирует заявку (меняет статус на formed)
// При формировании выполняется оценка кредитоспособности по скоринговой модели
func (r *Repository) FormApplication(appID int, userID int) error {
	// Получаем заявку
	app, err := r.GetApplicationByID(appID)
	if err != nil {
		return err
	}

	// Проверки
	if app.CreatorID != userID {
		return fmt.Errorf("only creator can form the application")
	}
	if app.Status != ds.ApplicationStatusDraft {
		return fmt.Errorf("can only form draft applications")
	}

	// Проверяем обязательные поля
	if app.FullName == "" || app.Income == 0 || app.Obligations == 0 {
		return fmt.Errorf("required fields are missing: full_name, income, obligations")
	}

	// Проверяем наличие продуктов
	if len(app.Products) == 0 {
		return fmt.Errorf("application must have at least one product")
	}

	// Рассчитываем суммарный ежемесячный платеж по всем кредитам
	totalMonthlyPayment := float64(0)
	for _, product := range app.Products {
		totalMonthlyPayment += product.MonthlyPayment
	}

	// Выполняем оценку кредитоспособности по скоринговой модели
	scoringResult := ds.CalculateCreditworthiness(app.Income, app.Obligations, totalMonthlyPayment)

	// Формируем заявку с результатами скоринга
	now := time.Now()
	return r.db.Model(&ds.Application{}).Where("id = ?", appID).Updates(map[string]interface{}{
		"status":            ds.ApplicationStatusFormed,
		"formed_at":         now,
		"credit_score":      scoringResult.CreditScore,
		"scoring_result":    scoringResult.Result,
		"max_credit_amount": scoringResult.MaxCreditAmount,
		"rejection_reason":  scoringResult.RejectionReason,
		"total_amount":      totalMonthlyPayment,
	}).Error
}

// CompleteApplication завершает заявку (вычисляет итоговую сумму)
func (r *Repository) CompleteApplication(appID int, moderatorID int) error {
	// Получаем заявку
	app, err := r.GetApplicationByID(appID)
	if err != nil {
		return err
	}

	// Проверки
	if app.Status != ds.ApplicationStatusFormed {
		return fmt.Errorf("can only complete formed applications")
	}

	// Вычисляем итоговую сумму
	totalAmount := float64(0)
	for _, product := range app.Products {
		totalAmount += product.MonthlyPayment // Используем ежемесячный платеж
	}

	// Завершаем заявку
	now := time.Now()
	return r.db.Model(&ds.Application{}).Where("id = ?", appID).Updates(map[string]interface{}{
		"status":       ds.ApplicationStatusCompleted,
		"completed_at": now,
		"moderator_id": moderatorID,
		"total_amount": totalAmount,
	}).Error
}

// RejectApplication отклоняет заявку
func (r *Repository) RejectApplication(appID int, moderatorID int) error {
	// Получаем заявку
	app, err := r.GetApplicationByID(appID)
	if err != nil {
		return err
	}

	// Проверки
	if app.Status != ds.ApplicationStatusFormed {
		return fmt.Errorf("can only reject formed applications")
	}

	// Отклоняем заявку
	now := time.Now()
	return r.db.Model(&ds.Application{}).Where("id = ?", appID).Updates(map[string]interface{}{
		"status":       ds.ApplicationStatusRejected,
		"completed_at": now,
		"moderator_id": moderatorID,
	}).Error
}

// DeleteApplication логически удаляет заявку
func (r *Repository) DeleteApplication(appID int, userID int) error {
	// Получаем заявку
	app, err := r.GetApplicationByID(appID)
	if err != nil {
		return err
	}

	// Проверки
	if app.CreatorID != userID {
		return fmt.Errorf("only creator can delete the application")
	}
	if app.Status != ds.ApplicationStatusDraft {
		return fmt.Errorf("can only delete draft applications")
	}

	// Удаляем
	return r.db.Model(&ds.Application{}).Where("id = ?", appID).Update("status", ds.ApplicationStatusDeleted).Error
}

// ============================================
// Методы для ApplicationProducts (М-М связь)
// ============================================

// AddCreditToApplication добавляет услугу в заявку через ORM (для совместимости с lab2)
func (r *Repository) AddCreditToApplication(appID int, creditID int) error {
	// Проверяем существование заявки
	var app ds.Application
	if err := r.db.Where("id = ? AND status = ?", appID, ds.ApplicationStatusDraft).First(&app).Error; err != nil {
		return fmt.Errorf("заявка не найдена или не в статусе черновик")
	}

	// Проверяем существование кредита
	var credit ds.Order
	if err := r.db.Where("id = ?", creditID).First(&credit).Error; err != nil {
		return fmt.Errorf("кредит не найден")
	}

	// Проверяем, не добавлен ли уже этот кредит в заявку
	var existing ds.ApplicationProduct
	err := r.db.Where("application_id = ? AND credit_id = ?", appID, creditID).First(&existing).Error
	if err == nil {
		// Кредит уже есть в заявке - это нормально, просто возвращаем успех
		return nil
	}

	// Добавляем продукт
	appProduct := ds.ApplicationProduct{
		ApplicationID: appID,
		CreditID:      creditID,
	}

	return r.db.Create(&appProduct).Error
}

// AddProductToApplication добавляет продукт в заявку
func (r *Repository) AddProductToApplication(appID int, creditID int, userID int) error {
	// Проверяем заявку
	app, err := r.GetApplicationByID(appID)
	if err != nil {
		return err
	}

	if app.CreatorID != userID {
		return fmt.Errorf("access denied")
	}
	if app.Status != ds.ApplicationStatusDraft {
		return fmt.Errorf("can only add products to draft applications")
	}

	// Проверяем продукт
	_, err = r.GetOrderByID(creditID)
	if err != nil {
		return fmt.Errorf("credit not found")
	}

	// Проверяем, не добавлен ли уже
	var existing ds.ApplicationProduct
	err = r.db.Where("application_id = ? AND credit_id = ?", appID, creditID).First(&existing).Error
	if err == nil {
		// Уже добавлен - просто возвращаем успех
		return nil
	}

	// Добавляем новый продукт
	appProduct := ds.ApplicationProduct{
		ApplicationID:  appID,
		CreditID:       creditID,
		MonthlyPayment: 0,
	}

	return r.db.Create(&appProduct).Error
}

// RemoveProductFromApplication удаляет продукт из заявки
func (r *Repository) RemoveProductFromApplication(appID int, creditID int, userID int) error {
	// Проверяем заявку
	app, err := r.GetApplicationByID(appID)
	if err != nil {
		return err
	}

	if app.CreatorID != userID {
		return fmt.Errorf("access denied")
	}
	if app.Status != ds.ApplicationStatusDraft {
		return fmt.Errorf("can only remove products from draft applications")
	}

	// Удаляем продукт
	result := r.db.Where("application_id = ? AND credit_id = ?", appID, creditID).
		Delete(&ds.ApplicationProduct{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found in application")
	}

	return nil
}

// UpdateApplicationProduct обновляет М-М связь
func (r *Repository) UpdateApplicationProduct(appID int, creditID int, requestedAmount *float64, requestedTermDays *int, monthlyPayment *float64, interestRate *float64, userID int) error {
	// Проверяем заявку
	app, err := r.GetApplicationByID(appID)
	if err != nil {
		return err
	}

	if app.CreatorID != userID {
		return fmt.Errorf("access denied")
	}
	if app.Status != ds.ApplicationStatusDraft {
		return fmt.Errorf("can only update products in draft applications")
	}

	// Формируем updates
	updates := make(map[string]interface{})
	if requestedAmount != nil {
		updates["requested_amount"] = *requestedAmount
	}
	if requestedTermDays != nil {
		updates["requested_term_days"] = *requestedTermDays
	}
	if monthlyPayment != nil {
		updates["monthly_payment"] = *monthlyPayment
	}
	if interestRate != nil {
		updates["interest_rate"] = *interestRate
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	result := r.db.Model(&ds.ApplicationProduct{}).
		Where("application_id = ? AND credit_id = ?", appID, creditID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found in application")
	}

	return nil
}

// CountProductsInApplication возвращает количество продуктов в заявке
func (r *Repository) CountProductsInApplication(appID int) (int, error) {
	var count int64
	err := r.db.Model(&ds.ApplicationProduct{}).Where("application_id = ?", appID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

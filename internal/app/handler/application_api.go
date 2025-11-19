package handler

import (
	"fmt"
	"lab1/internal/app/ds"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ============================================
// API Handlers для Applications (Заявки)
// ============================================

// GetApplicationBasketAPI - GET /api/applications/basket - получить информацию о корзине
func (h *Handler) GetApplicationBasketAPI(c *gin.Context) {
	// Получаем текущего пользователя
	userID := h.getCurrentUserID(c)

	// Получаем черновик заявки
	app, err := h.Repository.GetDraftApplication(userID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Если черновика нет, создаем его
	if app == nil {
		app, err = h.Repository.CreateDraftApplication(userID)
		if err != nil {
			h.errorHandler(c, http.StatusInternalServerError, err)
			return
		}
	}

	// Подсчитываем количество продуктов
	count, err := h.Repository.CountProductsInApplication(app.ID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	response := ds.ApplicationBasketResponse{
		ApplicationID: app.ID,
		ProductCount:  count,
	}

	c.JSON(http.StatusOK, response)
}

// GetApplicationsAPI - GET /api/applications - получить список заявок с фильтрацией
func (h *Handler) GetApplicationsAPI(c *gin.Context) {
	// Получаем текущего пользователя из JWT
	userID := h.getCurrentUserID(c)
	
	// Проверяем является ли пользователь модератором
	isModerator := false
	if claims, ok := c.Get("claims"); ok {
		if jwtClaims, ok := claims.(*ds.JWTClaims); ok {
			isModerator = jwtClaims.IsModerator
		}
	}

	// Получаем параметры фильтрации
	statusFilter := c.Query("status")
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	var dateFrom, dateTo *time.Time

	// Парсим даты, если они переданы
	if dateFromStr != "" {
		t, err := time.Parse("2006-01-02", dateFromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ds.ErrorResponse{
				Status:  "fail",
				Message: "invalid date_from format, use YYYY-MM-DD",
			})
			return
		}
		dateFrom = &t
	}

	if dateToStr != "" {
		t, err := time.Parse("2006-01-02", dateToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ds.ErrorResponse{
				Status:  "fail",
				Message: "invalid date_to format, use YYYY-MM-DD",
			})
			return
		}
		// Устанавливаем конец дня
		endOfDay := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		dateTo = &endOfDay
	}

	// Получаем список заявок
	applications, err := h.Repository.GetApplications(statusFilter, dateFrom, dateTo)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Фильтруем заявки: обычный пользователь видит только свои, модератор - все
	var filteredApplications []ds.Application
	if isModerator {
		// Модератор видит все заявки
		filteredApplications = applications
	} else {
		// Обычный пользователь видит только свои заявки
		for _, app := range applications {
			if app.CreatorID == userID {
				filteredApplications = append(filteredApplications, app)
			}
		}
	}

	// Формируем ответ
	var appResponses []ds.ApplicationResponse
	for _, app := range filteredApplications {
		appResp := h.buildApplicationResponse(&app)
		appResponses = append(appResponses, appResp)
	}

	response := ds.ApplicationListResponse{
		Total:        len(appResponses),
		Applications: appResponses,
	}

	c.JSON(http.StatusOK, response)
}

// GetApplicationAPI - GET /api/applications/:id - получить одну заявку с продуктами
func (h *Handler) GetApplicationAPI(c *gin.Context) {
	// Получаем ID из URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid application id",
		})
		return
	}

	// Получаем заявку
	app, err := h.Repository.GetApplicationByID(id)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	// Формируем ответ с продуктами
	response := h.buildApplicationResponse(app)

	// Добавляем детальную информацию о продуктах
	var products []ds.ApplicationProductDTO
	for _, p := range app.Products {
		imageURL := ""
		if p.Order.ImageURL != "" {
			imageURL = h.MinioService.GetFileURL("credits", p.Order.ImageURL)
		}

		products = append(products, ds.ApplicationProductDTO{
			ID:                p.ID,
			CreditID:          p.CreditID,
			CreditTitle:       p.Order.Title,
			ImageURL:          imageURL,
			RequestedAmount:   p.RequestedAmount,
			RequestedTermDays: p.RequestedTermDays,
			MonthlyPayment:    p.MonthlyPayment,
			InterestRate:      p.InterestRate,
		})
	}
	response.Products = products

	c.JSON(http.StatusOK, response)
}

// UpdateApplicationAPI - PUT /api/applications/:id - обновить поля заявки
func (h *Handler) UpdateApplicationAPI(c *gin.Context) {
	// Получаем текущего пользователя
	userID := h.getCurrentUserID(c)

	// Получаем ID заявки
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid application id",
		})
		return
	}

	var req ds.ApplicationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Получаем заявку для проверки прав доступа
	app, err := h.Repository.GetApplicationByID(id)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	if app.CreatorID != userID {
		c.JSON(http.StatusForbidden, ds.ErrorResponse{
			Status:  "fail",
			Message: "access denied",
		})
		return
	}

	// Обновляем поля
	app.FullName = req.FullName
	app.Income = req.Income
	app.Obligations = req.Obligations

	if err := h.Repository.UpdateApplication(app); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("Application updated: ID=%d", id)

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Application updated successfully",
		Data:    h.buildApplicationResponse(app),
	})
}

// FormApplicationAPI - PUT /api/applications/:id/form - сформировать заявку
func (h *Handler) FormApplicationAPI(c *gin.Context) {
	// Получаем текущего пользователя
	userID := h.getCurrentUserID(c)

	// Получаем ID заявки
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid application id",
		})
		return
	}

	// Формируем заявку
	if err := h.Repository.FormApplication(id, userID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("Application formed: ID=%d", id)

	// Получаем обновленную заявку
	app, _ := h.Repository.GetApplicationByID(id)

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Application formed successfully",
		Data:    h.buildApplicationResponse(app),
	})
}

// CompleteApplicationAPI - PUT /api/applications/:id/complete - завершить/отклонить заявку модератором
func (h *Handler) CompleteApplicationAPI(c *gin.Context) {
	// Получаем текущего пользователя
	userID := h.getCurrentUserID(c)

	// Проверяем, что пользователь - модератор
	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	if !user.IsModerator {
		c.JSON(http.StatusForbidden, ds.ErrorResponse{
			Status:  "fail",
			Message: "only moderator can complete or reject applications",
		})
		return
	}

	// Получаем ID заявки
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid application id",
		})
		return
	}

	var req ds.ApplicationCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Выполняем действие
	if req.Action == "complete" {
		if err := h.Repository.CompleteApplication(id, userID); err != nil {
			h.errorHandler(c, http.StatusBadRequest, err)
			return
		}
		logrus.Infof("Application completed: ID=%d", id)
	} else if req.Action == "reject" {
		if err := h.Repository.RejectApplication(id, userID); err != nil {
			h.errorHandler(c, http.StatusBadRequest, err)
			return
		}
		logrus.Infof("Application rejected: ID=%d", id)
	}

	// Получаем обновленную заявку
	app, _ := h.Repository.GetApplicationByID(id)

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: fmt.Sprintf("Application %sd successfully", req.Action),
		Data:    h.buildApplicationResponse(app),
	})
}

// DeleteApplicationAPI - DELETE /api/applications/:id - удалить заявку
func (h *Handler) DeleteApplicationAPI(c *gin.Context) {
	// Получаем текущего пользователя
	userID := h.getCurrentUserID(c)

	// Получаем ID заявки
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid application id",
		})
		return
	}

	// Удаляем заявку
	if err := h.Repository.DeleteApplication(id, userID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("Application deleted: ID=%d", id)

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Application deleted successfully",
	})
}

// buildApplicationResponse формирует ответ с данными заявки
func (h *Handler) buildApplicationResponse(app *ds.Application) ds.ApplicationResponse {
	creatorLogin := ""
	if app.Creator.Username != "" {
		creatorLogin = app.Creator.Username
	}

	var moderatorLogin *string
	if app.Moderator != nil && app.Moderator.Username != "" {
		login := app.Moderator.Username
		moderatorLogin = &login
	}

	return ds.ApplicationResponse{
		ID:           app.ID,
		Status:       app.Status,
		CreatedAt:    app.CreatedAt,
		FormedAt:     app.FormedAt,
		CompletedAt:  app.CompletedAt,
		CreatorLogin: creatorLogin,
		CreatorID:    app.CreatorID,
		Moderator:    moderatorLogin,
		ModeratorID:  app.ModeratorID,
		FullName:     app.FullName,
		Income:       app.Income,
		Obligations:  app.Obligations,
		TotalAmount:  app.TotalAmount,
		
		// Результаты оценки кредитоспособности
		CreditScore:     app.CreditScore,
		ScoringResult:   app.ScoringResult,
		MaxCreditAmount: app.MaxCreditAmount,
		RejectionReason: app.RejectionReason,
	}
}

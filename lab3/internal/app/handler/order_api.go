package handler

import (
	"fmt"
	"lab1/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ============================================
// API Handlers для Orders (Услуги/Кредиты)
// ============================================

// GetOrdersAPI - GET /api/credits - получить список услуг с фильтрацией
func (h *Handler) GetOrdersAPI(c *gin.Context) {
	// Получаем параметр фильтрации по названию
	titleFilter := c.Query("title")

	// Получаем список услуг из репозитория
	orders, err := h.Repository.GetOrdersWithFilter(titleFilter)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Формируем ответ с URL изображений
	var orderResponses []ds.OrderResponse
	for _, order := range orders {
		imageURL := ""
		if order.ImageURL != "" {
			imageURL = h.MinioService.GetFileURL("credits", order.ImageURL)
		}

		orderResponses = append(orderResponses, ds.OrderResponse{
			ID:          order.ID,
			Title:       order.Title,
			Icon:        order.Icon,
			ImageURL:    imageURL,
			Rate:        order.Rate,
			Term:        order.Term,
			Amount:      order.Amount,
			Feature1:    order.Feature1,
			Feature2:    order.Feature2,
			Feature3:    order.Feature3,
			Description: order.Description,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		})
	}

	response := ds.OrderListResponse{
		Total:  len(orderResponses),
		Orders: orderResponses,
	}

	c.JSON(http.StatusOK, response)
}

// GetOrderAPI - GET /api/credits/:id - получить одну услугу
func (h *Handler) GetOrderAPI(c *gin.Context) {
	// Получаем ID из URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// Получаем услугу из репозитория
	order, err := h.Repository.GetOrderByID(id)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	// Формируем URL изображения
	imageURL := ""
	if order.ImageURL != "" {
		imageURL = h.MinioService.GetFileURL("credits", order.ImageURL)
	}

	response := ds.OrderResponse{
		ID:          order.ID,
		Title:       order.Title,
		Icon:        order.Icon,
		ImageURL:    imageURL,
		Rate:        order.Rate,
		Term:        order.Term,
		Amount:      order.Amount,
		Feature1:    order.Feature1,
		Feature2:    order.Feature2,
		Feature3:    order.Feature3,
		Description: order.Description,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// CreateOrderAPI - POST /api/credits - создать новую услугу (без изображения)
func (h *Handler) CreateOrderAPI(c *gin.Context) {
	var req ds.OrderCreateRequest

	// Парсим JSON запрос
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Создаем объект Order
	order := ds.Order{
		Title:       req.Title,
		Icon:        req.Icon,
		Rate:        req.Rate,
		Term:        req.Term,
		Amount:      req.Amount,
		Feature1:    req.Feature1,
		Feature2:    req.Feature2,
		Feature3:    req.Feature3,
		Description: req.Description,
		Status:      ds.OrderStatusActive,
	}

	// Сохраняем в БД
	if err := h.Repository.CreateOrder(&order); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("Order created successfully: ID=%d", order.ID)

	c.JSON(http.StatusCreated, ds.SuccessResponse{
		Status:  "success",
		Message: "Order created successfully",
		Data: ds.OrderResponse{
			ID:          order.ID,
			Title:       order.Title,
			Icon:        order.Icon,
			ImageURL:    "",
			Rate:        order.Rate,
			Term:        order.Term,
			Amount:      order.Amount,
			Feature1:    order.Feature1,
			Feature2:    order.Feature2,
			Feature3:    order.Feature3,
			Description: order.Description,
			Status:      order.Status,
		},
	})
}

// UpdateOrderAPI - PUT /api/credits/:id - обновить услугу
func (h *Handler) UpdateOrderAPI(c *gin.Context) {
	// Получаем ID из URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid id parameter",
		})
		return
	}

	var req ds.OrderUpdateRequest

	// Парсим JSON запрос
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Получаем существующую услугу
	existing, err := h.Repository.GetOrderByID(id)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	// Обновляем поля (только если они переданы)
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Icon != "" {
		existing.Icon = req.Icon
	}
	if req.Rate != "" {
		existing.Rate = req.Rate
	}
	if req.Term != "" {
		existing.Term = req.Term
	}
	if req.Amount != "" {
		existing.Amount = req.Amount
	}
	if req.Feature1 != "" {
		existing.Feature1 = req.Feature1
	}
	if req.Feature2 != "" {
		existing.Feature2 = req.Feature2
	}
	if req.Feature3 != "" {
		existing.Feature3 = req.Feature3
	}
	if req.Description != "" {
		existing.Description = req.Description
	}

	// Сохраняем изменения
	if err := h.Repository.UpdateOrder(&existing); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("Order updated successfully: ID=%d", id)

	// Формируем URL изображения
	imageURL := ""
	if existing.ImageURL != "" {
		imageURL = h.MinioService.GetFileURL("credits", existing.ImageURL)
	}

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Order updated successfully",
		Data: ds.OrderResponse{
			ID:          existing.ID,
			Title:       existing.Title,
			Icon:        existing.Icon,
			ImageURL:    imageURL,
			Rate:        existing.Rate,
			Term:        existing.Term,
			Amount:      existing.Amount,
			Feature1:    existing.Feature1,
			Feature2:    existing.Feature2,
			Feature3:    existing.Feature3,
			Description: existing.Description,
			Status:      existing.Status,
		},
	})
}

// DeleteOrderAPI - DELETE /api/credits/:id - удалить услугу
func (h *Handler) DeleteOrderAPI(c *gin.Context) {
	// Получаем ID из URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid id parameter",
		})
		return
	}

	// Получаем услугу для получения имени изображения
	order, err := h.Repository.GetOrderByID(id)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	// Удаляем изображение из MinIO, если оно есть
	if order.ImageURL != "" {
		if err := h.MinioService.DeleteFile("credits", order.ImageURL); err != nil {
			logrus.Warnf("Failed to delete image from MinIO: %v", err)
		}
	}

	// Логически удаляем услугу
	if err := h.Repository.DeleteOrder(id); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("Order deleted successfully: ID=%d", id)

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Order deleted successfully",
	})
}

// UploadOrderImageAPI - POST /api/credits/:id/image - загрузить изображение для услуги
func (h *Handler) UploadOrderImageAPI(c *gin.Context) {
	// Получаем ID из URL
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid id parameter",
		})
		return
	}

	// Проверяем существование услуги
	order, err := h.Repository.GetOrderByID(id)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	// Получаем файл из запроса
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "image file is required",
		})
		return
	}

	// Проверяем размер файла (макс 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "file size must not exceed 10MB",
		})
		return
	}

	// Удаляем старое изображение, если оно есть
	if order.ImageURL != "" {
		if err := h.MinioService.DeleteFile("credits", order.ImageURL); err != nil {
			logrus.Warnf("Failed to delete old image: %v", err)
		}
	}

	// Загружаем новое изображение в MinIO
	fileName, err := h.MinioService.UploadFile(file, "credits")
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	// Обновляем запись в БД
	if err := h.Repository.UpdateOrderImage(id, fileName); err != nil {
		// Откатываем загрузку файла
		_ = h.MinioService.DeleteFile("credits", fileName)
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	imageURL := h.MinioService.GetFileURL("credits", fileName)

	logrus.Infof("Order image uploaded successfully: ID=%d, file=%s", id, fileName)

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Image uploaded successfully",
		Data: gin.H{
			"image_url": imageURL,
		},
	})
}

// AddOrderToApplicationAPI - POST /api/credits/:id/add-to-application - добавить услугу в заявку
func (h *Handler) AddOrderToApplicationAPI(c *gin.Context) {
	// Получаем текущего пользователя (в реальности из сессии/токена)
	userID := getCurrentUserID(c)

	// Получаем ID услуги из URL
	idParam := c.Param("id")
	creditID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid id parameter",
		})
		return
	}

	// Получаем или создаем черновик заявки
	app, err := h.Repository.GetDraftApplication(userID)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	if app == nil {
		// Создаем новый черновик
		app, err = h.Repository.CreateDraftApplication(userID)
		if err != nil {
			h.errorHandler(c, http.StatusInternalServerError, err)
			return
		}
	}

	// Добавляем продукт в заявку
	if err := h.Repository.AddProductToApplication(app.ID, creditID, userID); err != nil {
		if err.Error() == "credit not found" {
			c.JSON(http.StatusNotFound, ds.ErrorResponse{
				Status:  "fail",
				Message: "Credit not found",
			})
			return
		}
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("Credit added to application: appID=%s, creditID=%d", app.ID, creditID)

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Credit added to application successfully",
		Data: gin.H{
			"application_id": app.ID,
		},
	})
}

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
// API Handlers для ApplicationProducts (М-М связь)
// ============================================

// RemoveProductFromApplicationAPI - DELETE /api/applications/:app_id/products/:credit_id - удалить продукт из заявки
func (h *Handler) RemoveProductFromApplicationAPI(c *gin.Context) {
	// Получаем текущего пользователя
	userID := getCurrentUserID(c)

	// Получаем ID заявки и продукта
	appIDStr := c.Param("app_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid application id",
		})
		return
	}

	creditIDParam := c.Param("credit_id")
	creditID, err := strconv.Atoi(creditIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid credit_id parameter",
		})
		return
	}

	// Удаляем продукт из заявки
	if err := h.Repository.RemoveProductFromApplication(appID, creditID, userID); err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, ds.ErrorResponse{
				Status:  "fail",
				Message: "access denied",
			})
			return
		}
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("Product removed from application: appID=%d, creditID=%d", appID, creditID)

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Product removed from application successfully",
	})
}

// UpdateApplicationProductAPI - PUT /api/applications/:app_id/products/:credit_id - обновить М-М связь
func (h *Handler) UpdateApplicationProductAPI(c *gin.Context) {
	// Получаем текущего пользователя
	userID := getCurrentUserID(c)

	// Получаем ID заявки и продукта
	appIDStr := c.Param("app_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid application id",
		})
		return
	}

	creditIDParam := c.Param("credit_id")
	creditID, err := strconv.Atoi(creditIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid credit_id parameter",
		})
		return
	}

	var req ds.UpdateApplicationProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Обновляем связь
	if err := h.Repository.UpdateApplicationProduct(appID, creditID, req.RequestedAmount, req.RequestedTermDays, req.MonthlyPayment, req.InterestRate, userID); err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, ds.ErrorResponse{
				Status:  "fail",
				Message: "access denied",
			})
			return
		}
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("ApplicationProduct updated: appID=%d, creditID=%d", appID, creditID)

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Product updated successfully",
	})
}

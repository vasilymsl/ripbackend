package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const DefaultUserID = 1 // Тестовый пользователь (в реальном приложении из сессии/JWT)

// AddToApplication добавляет услугу в заявку (POST через ORM)
func (h *Handler) AddToApplication(ctx *gin.Context) {
	creditIDStr := ctx.PostForm("credit_id")
	creditID, err := strconv.Atoi(creditIDStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Получаем или создаем черновик заявки
	draft, err := h.Repository.GetDraftApplication(DefaultUserID)
	if err != nil {
		// Создаем новую заявку в статусе черновик
		draft, err = h.Repository.CreateDraftApplication(DefaultUserID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		logrus.Infof("Создана новая заявка в статусе черновик: %d", draft.ID)
	}

	// Добавляем услугу в заявку
	if err := h.Repository.AddCreditToApplication(draft.ID, creditID); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("Услуга %d добавлена в заявку %d", creditID, draft.ID)

	// Редирект на страницу корзины
	ctx.Redirect(http.StatusFound, "/creditapplicationbasket")
}

// DeleteApplication логически удаляет заявку (POST через SQL UPDATE)
func (h *Handler) DeleteApplication(ctx *gin.Context) {
	appIDStr := ctx.PostForm("application_id")
	if appIDStr == "" {
		h.errorHandler(ctx, http.StatusBadRequest, nil)
		return
	}

	appID, parseErr := strconv.Atoi(appIDStr)
	if parseErr != nil {
		h.errorHandler(ctx, http.StatusBadRequest, parseErr)
		return
	}

	// Удаление через SQL UPDATE (не ORM!)
	// Получаем текущего пользователя (упрощенно)
	userID := DefaultUserID
	if err := h.Repository.DeleteApplication(appID, userID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("Заявка %d логически удалена (статус изменен на deleted)", appID)

	// Остаемся в корзине
	ctx.Redirect(http.StatusFound, "/creditapplicationbasket")
}

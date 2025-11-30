package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetApplications(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заявки из урла (то есть из /scoringapplicationbasket/:id)

	if idStr != "" {
		// Если указан конкретный ID
		id, parseErr := strconv.Atoi(idStr)
		if parseErr != nil {
			ctx.HTML(http.StatusBadRequest, "applications.html", gin.H{
				"error": "Неверный ID заявки",
			})
			return
		}
		appPtr, err := h.Repository.GetApplicationByID(id)
		if err != nil {
			logrus.Error(err)
			ctx.HTML(http.StatusNotFound, "applications.html", gin.H{
				"error": "Заявка не найдена или удалена",
			})
			return
		}
		// Заполняем поля для таблицы (короткие значения)
		for i := range appPtr.Products {
		appPtr.Products[i].Order.PopulateForTable()
		// Формируем полный URL изображения из MinIO
		if appPtr.Products[i].Order.ImageURL != "" {
			appPtr.Products[i].Order.ImageURL = h.MinioService.GetFileURL("scoring", appPtr.Products[i].Order.ImageURL)
		}
		}
		ctx.HTML(http.StatusOK, "applications.html", gin.H{
			"application": *appPtr,
		})
		return
	}

	// Получаем черновик текущего пользователя
	draft, err := h.Repository.GetDraftApplication(DefaultUserID)
	if err != nil || draft == nil {
		// Нет черновика - показываем пустую корзину
		ctx.HTML(http.StatusOK, "applications.html", gin.H{
			"message": "Корзина пуста. Добавьте услугу!",
		})
		return
	}

	// Заполняем поля для таблицы (короткие значения)
	for i := range draft.Products {
	draft.Products[i].Order.PopulateForTable()
	// Формируем полный URL изображения из MinIO
	if draft.Products[i].Order.ImageURL != "" {
		draft.Products[i].Order.ImageURL = h.MinioService.GetFileURL("scoring", draft.Products[i].Order.ImageURL)
	}
	}

	ctx.HTML(http.StatusOK, "applications.html", gin.H{
		"application": *draft,
	})
}

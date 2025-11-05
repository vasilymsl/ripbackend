package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetApplications(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заявки из урла (то есть из /creditapplicationbasket/:id)

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

	ctx.HTML(http.StatusOK, "applications.html", gin.H{
		"application": *draft,
	})
}

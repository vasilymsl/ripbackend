package handler

import (
	"net/http"
	"strconv"
	"time"

	"lab1/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetOrders(ctx *gin.Context) {
	var orders []ds.Order
	var err error

	searchQuery := ctx.Query("query") // получаем значение из нашего поля
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		orders, err = h.Repository.GetOrders()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		orders, err = h.Repository.GetOrdersByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	// Заполняем поля совместимости для шаблонов
	for i := range orders {
		orders[i].PopulateCompatibilityFields()
		// Формируем полный URL изображения из MinIO (бакет scoring)
		if orders[i].ImageURL != "" {
			orders[i].ImageURL = h.MinioService.GetFileURL("scoring", orders[i].ImageURL)
		}
	}

	// Получаем количество товаров в корзине
	cartCount := 0
	draft, errDraft := h.Repository.GetDraftApplication(DefaultUserID)
	if errDraft == nil && draft != nil {
		cartCount = len(draft.Products)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":       time.Now().Format("15:04:05"),
		"orders":     orders,
		"query":      searchQuery, // передаем введенный запрос обратно на страницу
		"cart_count": cartCount,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /scoringmodel/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	// Заполняем поля совместимости для шаблона
	order.PopulateCompatibilityFields()
	
	// Формируем полный URL изображения из MinIO
	if order.ImageURL != "" {
		order.ImageURL = h.MinioService.GetFileURL("scoring", order.ImageURL)
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order": order,
	})
}

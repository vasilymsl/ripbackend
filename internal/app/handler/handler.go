package handler

import (
	"lab1/internal/app/repository"
	"lab1/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository   *repository.Repository
	MinioService *service.MinioService
}

func NewHandler(r *repository.Repository, m *service.MinioService) *Handler {
	return &Handler{
		Repository:   r,
		MinioService: m,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	// Старые маршруты для обратной совместимости (если нужны)
	router.GET("/", h.GetOrders)
	router.GET("/credit-services", h.GetOrders)
	router.GET("/credit/:id", h.GetOrder)
	router.GET("/creditapplicationbasket/:id", h.GetApplications)
	router.GET("/creditapplicationbasket", h.GetApplications)
	router.POST("/add-to-application", h.AddToApplication)
	router.POST("/delete-application", h.DeleteApplication)

	// API маршруты с префиксом /api
	api := router.Group("/api")
	{
		// ============================================
		// Домен Orders (Услуги/Кредиты) - 7 методов
		// ============================================
		credits := api.Group("/credits")
		{
			credits.GET("", h.GetOrdersAPI)                                     // 1. GET список с фильтрацией
			credits.GET("/:id", h.GetOrderAPI)                                  // 2. GET одна запись
			credits.POST("", h.CreateOrderAPI)                                  // 3. POST добавление (без изображения)
			credits.PUT("/:id", h.UpdateOrderAPI)                               // 4. PUT изменение
			credits.DELETE("/:id", h.DeleteOrderAPI)                            // 5. DELETE удаление
			credits.POST("/:id/image", h.UploadOrderImageAPI)                   // 6. POST добавление изображения
			credits.POST("/:id/add-to-application", h.AddOrderToApplicationAPI) // 7. POST добавления в заявку
		}

		// ============================================
		// Домен Applications (Заявки) - 7 методов
		// ============================================
		applications := api.Group("/applications")
		{
			applications.GET("/basket", h.GetApplicationBasketAPI)      // 1. GET иконки корзины
			applications.GET("", h.GetApplicationsAPI)                  // 2. GET список с фильтрацией
			applications.GET("/:id", h.GetApplicationAPI)               // 3. GET одна запись
			applications.PUT("/:id/form", h.FormApplicationAPI)         // 5. PUT сформировать
			applications.PUT("/:id/complete", h.CompleteApplicationAPI) // 6. PUT завершить/отклонить
			applications.PUT("/:id", h.UpdateApplicationAPI)            // 4. PUT изменения полей
			applications.DELETE("/:id", h.DeleteApplicationAPI)         // 7. DELETE удаление
		}

		// ============================================
		// Домен ApplicationProducts (М-М связь) - 2 метода
		// Отдельная группа для избежания конфликта роутов
		// ============================================
		appProducts := api.Group("/application-products")
		{
			appProducts.DELETE("/:app_id/:credit_id", h.RemoveProductFromApplicationAPI) // DELETE из заявки
			appProducts.PUT("/:app_id/:credit_id", h.UpdateApplicationProductAPI)        // PUT изменение М-М
		}

		// ============================================
		// Домен Users (Пользователи) - 5 методов
		// ============================================
		users := api.Group("/users")
		{
			users.POST("/register", h.RegisterUserAPI)    // 1. POST регистрация
			users.POST("/login", h.LoginUserAPI)          // 2. POST аутентификация
			users.POST("/logout", h.LogoutUserAPI)        // 3. POST деавторизация
			users.GET("/profile", h.GetUserProfileAPI)    // 4. GET профиля
			users.PUT("/profile", h.UpdateUserProfileAPI) // 5. PUT профиля
		}
	}

	logrus.Info("API routes registered successfully")
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("/Users/vasilymaslovsky/Documents/lab3/templates/*")
	router.Static("/static", "/Users/vasilymaslovsky/Documents/lab3/resources")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

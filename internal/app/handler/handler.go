package handler

import (
	"html/template"
	"lab1/internal/app/config"
	"lab1/internal/app/middleware"
	"lab1/internal/app/redis"
	"lab1/internal/app/repository"
	"lab1/internal/app/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository   *repository.Repository
	MinioService *service.MinioService
	Config       *config.Config
	RedisClient  *redis.Client
	AuthMW       *middleware.AuthMiddleware
}

func NewHandler(r *repository.Repository, m *service.MinioService, cfg *config.Config, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository:   r,
		MinioService: m,
		Config:       cfg,
		RedisClient:  redisClient,
		AuthMW:       middleware.NewAuthMiddleware(cfg.JWT, redisClient),
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	// Старые маршруты для обратной совместимости (если нужны)
	router.GET("/", h.GetOrders)
	router.GET("/scoring-services", h.GetOrders)
	router.GET("/scoringmodel/:id", h.GetOrder)
	router.GET("/scoringapplicationbasket/:id", h.GetApplications)
	router.GET("/scoringapplicationbasket", h.GetApplications)
	router.POST("/add-to-application", h.AuthMW.OptionalAuth(), h.AddToApplication)
	router.POST("/delete-application", h.AuthMW.RequireAuth(), h.DeleteApplication)

	// API маршруты с префиксом /api
	api := router.Group("/api")
	{
		// ============================================
		// Домен Users (Пользователи) - 5 методов
		// Публичные эндпоинты (без авторизации)
		// ============================================
		users := api.Group("/users")
		{
			users.POST("/register", h.RegisterUserAPI) // 1. POST регистрация
			users.POST("/login", h.LoginUserAPI)       // 2. POST аутентификация

			// Защищенные эндпоинты (требуют авторизации)
			authenticated := users.Group("")
			authenticated.Use(h.AuthMW.RequireAuth())
			{
				authenticated.POST("/logout", h.LogoutUserAPI)        // 3. POST деавторизация
				authenticated.GET("/profile", h.GetUserProfileAPI)    // 4. GET профиля
				authenticated.PUT("/profile", h.UpdateUserProfileAPI) // 5. PUT профиля
			}
		}

		// ============================================
		// Домен Orders (Услуги/Кредиты) - 7 методов
		// Гости: только чтение (GET)
		// Пользователи: чтение + добавление в заявку
		// Модераторы: все операции
		// ============================================
		credits := api.Group("/credits")
		{
			// Публичные эндпоинты (доступны всем)
			credits.GET("", h.GetOrdersAPI)   // 1. GET список с фильтрацией
			credits.GET("/:id", h.GetOrderAPI) // 2. GET одна запись

			// Эндпоинты для авторизованных пользователей
			creditsAuth := credits.Group("")
			creditsAuth.Use(h.AuthMW.RequireAuth())
			{
				creditsAuth.POST("/:id/add-to-application", h.AddOrderToApplicationAPI) // 7. POST добавления в заявку
			}

			// Эндпоинты только для модераторов
			creditsModerator := credits.Group("")
			creditsModerator.Use(h.AuthMW.RequireAuth(), h.AuthMW.RequireModerator())
			{
				creditsModerator.POST("", h.CreateOrderAPI)                // 3. POST добавление
				creditsModerator.PUT("/:id", h.UpdateOrderAPI)             // 4. PUT изменение
				creditsModerator.DELETE("/:id", h.DeleteOrderAPI)          // 5. DELETE удаление
				creditsModerator.POST("/:id/image", h.UploadOrderImageAPI) // 6. POST добавление изображения
			}
		}

		// ============================================
		// Домен Applications (Заявки) - 7 методов
		// Гости: нет доступа (401/403)
		// Пользователи: только свои заявки
		// Модераторы: все заявки + завершение
		// ============================================
		applications := api.Group("/applications")
		applications.Use(h.AuthMW.RequireAuth()) // Все эндпоинты требуют авторизации
		{
			applications.GET("/basket", h.GetApplicationBasketAPI) // 1. GET иконки корзины
			applications.GET("", h.GetApplicationsAPI)             // 2. GET список (с фильтрацией по создателю)
			applications.GET("/:id", h.GetApplicationAPI)          // 3. GET одна запись
			applications.PUT("/:id/form", h.FormApplicationAPI)    // 5. PUT сформировать
			applications.PUT("/:id", h.UpdateApplicationAPI)       // 4. PUT изменения полей
			applications.DELETE("/:id", h.DeleteApplicationAPI)    // 7. DELETE удаление

			// Только для модераторов
			applicationsModerator := applications.Group("")
			applicationsModerator.Use(h.AuthMW.RequireModerator())
			{
				applicationsModerator.PUT("/:id/complete", h.CompleteApplicationAPI) // 6. PUT завершить/отклонить
			}
		}

		// ============================================
		// Домен ApplicationProducts (М-М связь) - 2 метода
		// Требует авторизации
		// ============================================
		appProducts := api.Group("/application-products")
		appProducts.Use(h.AuthMW.RequireAuth())
		{
			appProducts.DELETE("/:app_id/:credit_id", h.RemoveProductFromApplicationAPI) // DELETE из заявки
			appProducts.PUT("/:app_id/:credit_id", h.UpdateApplicationProductAPI)        // PUT изменение М-М
		}
	}

	logrus.Info("API routes registered successfully")
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	// Создаем кастомные функции для шаблонов
	router.SetFuncMap(template.FuncMap{
		"add": func(a, b float64) float64 {
			return a + b
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
	})
	
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "resources")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

package handler

import (
	"fmt"
	"lab1/internal/app/auth"
	"lab1/internal/app/ds"
	"lab1/internal/app/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ============================================
// API Handlers для Users (Пользователи)
// ============================================

// RegisterUserAPI - POST /api/users/register - регистрация нового пользователя
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ds.UserRegisterRequest true "Данные для регистрации"
// @Success 201 {object} ds.SuccessResponse "Пользователь успешно зарегистрирован"
// @Failure 400 {object} ds.ErrorResponse "Некорректные данные"
// @Failure 409 {object} ds.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} ds.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/register [post]
func (h *Handler) RegisterUserAPI(c *gin.Context) {
	var req ds.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Валидация данных
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "username is required",
		})
		return
	}

	if req.Password == "" || len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: "password must be at least 6 characters",
		})
		return
	}

	// Создаем пользователя
	user := ds.User{
		Username:    req.Username,
		Password:    req.Password,
		Email:       req.Email,
		FullName:    req.FullName,
		Phone:       req.Phone,
		IsModerator: req.IsModerator,
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, ds.ErrorResponse{
				Status:  "fail",
				Message: err.Error(),
			})
			return
		}
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("User registered successfully: username=%s, id=%d", user.Username, user.ID)

	// Формируем ответ (без пароля)
	response := ds.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Phone:       user.Phone,
		IsModerator: user.IsModerator,
		CreatedAt:   user.CreatedAt,
	}

	c.JSON(http.StatusCreated, ds.SuccessResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data:    response,
	})
}

// LoginUserAPI - POST /api/users/login - аутентификация пользователя
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя в систему и возвращает JWT токен
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ds.UserLoginRequest true "Учетные данные"
// @Success 200 {object} object{status=string,message=string,data=object{user=ds.UserResponse,access_token=string,token_type=string,expires_in=int}} "Успешная аутентификация"
// @Failure 400 {object} ds.ErrorResponse "Некорректные данные"
// @Failure 401 {object} ds.ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} ds.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/login [post]
func (h *Handler) LoginUserAPI(c *gin.Context) {
	var req ds.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Аутентифицируем пользователя
	user, err := h.Repository.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
			Status:  "fail",
			Message: "invalid credentials",
		})
		return
	}

	// Генерируем JWT токен
	tokenString, err := auth.GenerateJWT(user, h.Config.JWT)
	if err != nil {
		logrus.Errorf("Failed to generate JWT: %v", err)
		c.JSON(http.StatusInternalServerError, ds.ErrorResponse{
			Status:  "error",
			Message: "failed to generate token",
		})
		return
	}

	logrus.Infof("User logged in: username=%s, id=%d, is_moderator=%v", user.Username, user.ID, user.IsModerator)

	// Сохраняем сессию в Redis
	if err := h.RedisClient.SaveSession(c.Request.Context(), tokenString, user.ID, h.Config.JWT.ExpiresIn); err != nil {
		logrus.Errorf("Failed to save session to Redis: %v", err)
	} else {
		logrus.Infof("Session saved to Redis for user %d", user.ID)
	}

	// Также сохраняем токен в cookie для удобства (опционально)
	c.SetCookie("auth_token", tokenString, int(h.Config.JWT.ExpiresIn.Seconds()), "/", "", false, true)

	// Формируем ответ
	expiresInSeconds := int64(h.Config.JWT.ExpiresIn.Seconds())
	
	response := ds.UserLoginResponse{
		User: ds.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			Phone:       user.Phone,
			IsModerator: user.IsModerator,
			CreatedAt:   user.CreatedAt,
		},
		Token: tokenString,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data": gin.H{
			"user":         response.User,
			"access_token": tokenString,
			"token_type":   "Bearer",
			"expires_in":   expiresInSeconds,
		},
	})
}

// LogoutUserAPI - POST /api/users/logout - деавторизация пользователя
// @Summary Выход из системы
// @Description Завершает сессию пользователя и добавляет токен в blacklist
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ds.SuccessResponse "Успешный выход"
// @Failure 401 {object} ds.ErrorResponse "Требуется авторизация"
// @Router /users/logout [post]
func (h *Handler) LogoutUserAPI(c *gin.Context) {
	// Получаем токен из заголовка
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Получаем claims для определения времени жизни
		claims, err := auth.ParseJWT(tokenString, h.Config.JWT)
		if err == nil {
			// Добавляем токен в blacklist до момента его истечения
			ttl := auth.GetTokenExpiration(claims)
			if ttl > 0 {
				err = h.RedisClient.WriteJWTToBlacklist(c.Request.Context(), tokenString, ttl)
				if err != nil {
					logrus.Errorf("Failed to add JWT to blacklist: %v", err)
				}
			}
		}
	}

	// Удаляем cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	logrus.Info("User logged out")

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Logout successful",
	})
}

// GetUserProfileAPI - GET /api/users/profile - получить профиль текущего пользователя
// @Summary Получить профиль пользователя
// @Description Возвращает профиль текущего авторизованного пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ds.UserResponse "Профиль пользователя"
// @Failure 401 {object} ds.ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ds.ErrorResponse "Пользователь не найден"
// @Router /users/profile [get]
func (h *Handler) GetUserProfileAPI(c *gin.Context) {
	// Получаем ID пользователя из JWT токена (middleware уже проверил авторизацию)
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
			Status:  "fail",
			Message: "user not authenticated",
		})
		return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	// Формируем ответ
	response := ds.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Phone:       user.Phone,
		IsModerator: user.IsModerator,
		CreatedAt:   user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUserProfileAPI - PUT /api/users/profile - обновить профиль пользователя
// @Summary Обновить профиль пользователя
// @Description Обновляет данные профиля текущего пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ds.UserUpdateRequest true "Данные для обновления"
// @Success 200 {object} ds.SuccessResponse "Профиль успешно обновлен"
// @Failure 400 {object} ds.ErrorResponse "Некорректные данные"
// @Failure 401 {object} ds.ErrorResponse "Требуется авторизация"
// @Failure 404 {object} ds.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ds.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/profile [put]
func (h *Handler) UpdateUserProfileAPI(c *gin.Context) {
	// Получаем ID пользователя из JWT токена
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
			Status:  "fail",
			Message: "user not authenticated",
		})
		return
	}

	var req ds.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Получаем пользователя
	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	// Обновляем поля
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Password != "" {
		user.Password = req.Password
	}

	// Сохраняем изменения
	if err := h.Repository.UpdateUser(user); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("User profile updated: ID=%d", userID)

	// Формируем ответ
	response := ds.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Phone:       user.Phone,
		IsModerator: user.IsModerator,
		CreatedAt:   user.CreatedAt,
	}

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Profile updated successfully",
		Data:    response,
	})
}

// ============================================
// Вспомогательные функции
// ============================================

// getCurrentUserID возвращает ID текущего пользователя из контекста
func (h *Handler) getCurrentUserID(c *gin.Context) int {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return 0
		}
	return userID
}

// getCurrentUser возвращает текущего пользователя
func (h *Handler) getCurrentUser(c *gin.Context) (*ds.User, error) {
	userID := h.getCurrentUserID(c)
	if userID == 0 {
		return nil, fmt.Errorf("user not authenticated")
	}
	return h.Repository.GetUserByID(userID)
}

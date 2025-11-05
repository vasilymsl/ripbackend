package handler

import (
	"fmt"
	"lab1/internal/app/ds"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ============================================
// API Handlers для Users (Пользователи)
// ============================================

// RegisterUserAPI - POST /api/users/register - регистрация нового пользователя
func (h *Handler) RegisterUserAPI(c *gin.Context) {
	var req ds.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ds.ErrorResponse{
			Status:  "fail",
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	// Создаем пользователя
	user := ds.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		if err.Error() == fmt.Sprintf("user with username '%s' already exists", req.Username) {
			c.JSON(http.StatusConflict, ds.ErrorResponse{
				Status:  "fail",
				Message: err.Error(),
			})
			return
		}
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	logrus.Infof("User registered successfully: username=%s", user.Username)

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

	logrus.Infof("User logged in: username=%s", user.Username)

	// Создаем "токен" (упрощенный вариант - просто ID пользователя)
	token := fmt.Sprintf("user_%d", user.ID)

	// Сохраняем в сессии (упрощенный вариант через cookie)
	c.SetCookie("auth_token", token, 86400, "/", "", false, true)

	// Формируем ответ
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
		Token: token,
	}

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    response,
	})
}

// LogoutUserAPI - POST /api/users/logout - деавторизация пользователя
func (h *Handler) LogoutUserAPI(c *gin.Context) {
	// Удаляем cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	logrus.Info("User logged out")

	c.JSON(http.StatusOK, ds.SuccessResponse{
		Status:  "success",
		Message: "Logout successful",
	})
}

// GetUserProfileAPI - GET /api/users/profile - получить профиль текущего пользователя
func (h *Handler) GetUserProfileAPI(c *gin.Context) {
	// Получаем текущего пользователя
	userID := getCurrentUserID(c)

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
func (h *Handler) UpdateUserProfileAPI(c *gin.Context) {
	// Получаем текущего пользователя
	userID := getCurrentUserID(c)

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
// Упрощенная версия - для демонстрации используется фиксированный пользователь
func getCurrentUserID(c *gin.Context) int {
	// В реальном приложении здесь нужно получать ID из JWT токена или сессии
	// Для лабораторной работы используем фиксированного пользователя (как требуется в задании)

	// Пытаемся получить из cookie
	token, err := c.Cookie("auth_token")
	if err == nil && token != "" {
		// Парсим токен вида "user_123"
		var userID int
		_, err := fmt.Sscanf(token, "user_%d", &userID)
		if err == nil && userID > 0 {
			return userID
		}
	}

	// Если нет авторизации, используем дефолтного пользователя ID=1 (как в задании)
	return 1
}

// getCurrentUser возвращает текущего пользователя
func (h *Handler) getCurrentUser(c *gin.Context) (*ds.User, error) {
	userID := getCurrentUserID(c)
	return h.Repository.GetUserByID(userID)
}

// requireAuth middleware для проверки авторизации
func (h *Handler) requireAuth(c *gin.Context) {
	userID := getCurrentUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ds.ErrorResponse{
			Status:  "fail",
			Message: "authentication required",
		})
		c.Abort()
		return
	}
	c.Next()
}

// requireModerator middleware для проверки прав модератора
func (h *Handler) requireModerator(c *gin.Context) {
	user, err := h.getCurrentUser(c)
	if err != nil || !user.IsModerator {
		c.JSON(http.StatusForbidden, ds.ErrorResponse{
			Status:  "fail",
			Message: "moderator access required",
		})
		c.Abort()
		return
	}
	c.Next()
}

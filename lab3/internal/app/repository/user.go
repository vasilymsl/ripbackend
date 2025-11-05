package repository

import (
	"errors"
	"fmt"
	"lab1/internal/app/ds"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ============================================
// Методы репозитория для Users (Пользователи)
// ============================================

// CreateUser создает нового пользователя
func (r *Repository) CreateUser(user *ds.User) error {
	// Проверяем, не существует ли уже пользователь с таким username
	existing, _ := r.GetUserByUsername(user.Username)
	if existing != nil {
		return fmt.Errorf("user with username '%s' already exists", user.Username)
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Создаем пользователя
	return r.db.Create(user).Error
}

// GetUserByID возвращает пользователя по ID
func (r *Repository) GetUserByID(id int) (*ds.User, error) {
	var user ds.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername возвращает пользователя по username
func (r *Repository) GetUserByUsername(username string) (*ds.User, error) {
	var user ds.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser обновляет данные пользователя
func (r *Repository) UpdateUser(user *ds.User) error {
	// Проверяем существование пользователя
	existing, err := r.GetUserByID(user.ID)
	if err != nil {
		return err
	}

	updates := make(map[string]interface{})

	// Обновляем только те поля, которые переданы
	if user.Email != "" && user.Email != existing.Email {
		updates["email"] = user.Email
	}
	if user.FullName != "" && user.FullName != existing.FullName {
		updates["full_name"] = user.FullName
	}
	if user.Phone != "" && user.Phone != existing.Phone {
		updates["phone"] = user.Phone
	}

	// Если передан новый пароль, хешируем его
	if user.Password != "" && user.Password != existing.Password {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		updates["password"] = string(hashedPassword)
	}

	if len(updates) == 0 {
		return nil // Нечего обновлять
	}

	return r.db.Model(&ds.User{}).Where("id = ?", user.ID).Updates(updates).Error
}

// AuthenticateUser проверяет учетные данные пользователя
func (r *Repository) AuthenticateUser(username, password string) (*ds.User, error) {
	user, err := r.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// GetUsers возвращает список всех пользователей
func (r *Repository) GetUsers() ([]ds.User, error) {
	var users []ds.User
	err := r.db.Order("id ASC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

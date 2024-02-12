package models

import "time"

type UserID string

func (id UserID) String() string {
	return string(id)
}

// User Пользователь
type User struct {
	ID        UserID    `json:"-"`                                                                                              // ID пользователя
	Username  string    `json:"login" validate:"required" message:"required:{field} is required" label:"login"`                 // Логин
	Password  string    `json:"password,omitempty" validate:"required" message:"required:{field} is required" label:"password"` // Пароль
	CreatedAt time.Time `json:"-"`                                                                                              // Дата регистрации
}

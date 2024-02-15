package models

import "time"

type UserID string

func (id UserID) String() string {
	return string(id)
}

// User Пользователь
type User struct {
	ID        UserID    // ID пользователя
	Username  string    // Логин
	Password  string    // Пароль
	CreatedAt time.Time // Дата регистрации
}

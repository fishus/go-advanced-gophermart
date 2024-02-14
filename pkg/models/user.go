package models

import "time"

type UserID string

func (id UserID) String() string {
	return string(id)
}

// User Пользователь
type User struct {
	ID        UserID    `db:"id"`         // ID пользователя
	Username  string    `db:"username"`   // Логин
	Password  string    `db:"-" `         // Пароль
	CreatedAt time.Time `db:"created_at"` // Дата регистрации
}

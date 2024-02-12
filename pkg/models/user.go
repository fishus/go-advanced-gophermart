package models

import "time"

type UserID string

func (id UserID) String() string {
	return string(id)
}

// User Пользователь
type User struct {
	ID        UserID    `db:"id" json:"-"`                                                                                             // ID пользователя
	Username  string    `db:"username" json:"login" validate:"required" message:"required:{field} is required" label:"login"`          // Логин
	Password  string    `db:"-" json:"password,omitempty" validate:"required" message:"required:{field} is required" label:"password"` // Пароль
	CreatedAt time.Time `db:"created_at" json:"-"`                                                                                     // Дата регистрации
}

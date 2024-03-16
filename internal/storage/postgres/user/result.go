package user

import (
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

// UserResult Пользователь
type UserResult struct {
	ID        models.UserID `db:"id"`         // ID пользователя
	Username  string        `db:"username"`   // Логин
	Password  string        `db:"-"`          // Пароль
	CreatedAt time.Time     `db:"created_at"` // Дата регистрации
}

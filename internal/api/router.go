package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "github.com/fishus/go-advanced-gophermart/internal/api/middleware"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

func Router(s *server) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(mw.Decompress)
	r.Use(middleware.Compress(9, "application/json"))
	r.Use(middleware.RequestLogger(&logger.LogFormatter{}))

	r.Post("/api/user/register", s.userRegister) // Регистрация пользователя
	r.Post("/api/user/login", s.userLogin)       // Аутентификация пользователя
	r.Post("/api/user/orders", s.orderAdd)       // Загрузка номера заказа для расчёта
	r.Get("/api/user/orders", s.ordersList)      // Список загруженных номеров заказов
	r.Get("/api/user/balance", s.userBalance)    // Получение баланса пользователя

	return r
}

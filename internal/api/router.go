package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "github.com/fishus/go-advanced-gophermart/internal/api/middleware"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
)

func (s *server) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(mw.Decompress)
	r.Use(middleware.Compress(9, "application/json"))
	r.Use(middleware.RequestLogger(&logger.LogFormatter{}))

	r.Post("/api/user/register", s.user.Register)            // Регистрация пользователя
	r.Post("/api/user/login", s.user.Login)                  // Аутентификация пользователя
	r.Post("/api/user/orders", s.order.Add)                  // Загрузка номера заказа для расчёта
	r.Get("/api/user/orders", s.order.List)                  // Список загруженных номеров заказов
	r.Get("/api/user/balance", s.loyalty.Balance)            // Получение баланса пользователя
	r.Post("/api/user/balance/withdraw", s.loyalty.Withdraw) // Запрос на списание средств
	r.Get("/api/user/withdrawals", s.loyalty.Withdrawals)    // Информации о выводе средств

	return r
}

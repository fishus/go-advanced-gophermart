package order

import (
	"strconv"

	"github.com/theplant/luhn"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
)

// Проверка номера заказа на корректность с помощью алгоритма Луна
func validateNumLuhn(num string) error {
	i, err := strconv.Atoi(num)
	if err != nil {
		return serviceErr.ErrOrderWrongNum
	}
	if !luhn.Valid(i) {
		return serviceErr.ErrOrderWrongNum
	}
	return nil
}

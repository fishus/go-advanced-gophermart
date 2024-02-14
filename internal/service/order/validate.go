package order

import (
	"strconv"

	"github.com/gookit/validate"
	"github.com/theplant/luhn"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
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

func validateOrder(order models.Order) error {
	v := validate.Struct(&order)

	v.AddRule("UserID", "required")
	v.AddRule("Num", "required")

	v.AddMessages(map[string]string{
		"required": "the {field} is required",
	})

	v.AddTranslates(map[string]string{
		"UserID": "user id",
		"Num":    "order number",
	})

	if !v.Validate() {
		return serviceErr.NewValidationError(v.Errors)
	}

	return nil
}

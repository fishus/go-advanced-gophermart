package user

import (
	"github.com/gookit/validate"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func validateUser(user models.User) error {
	v := validate.Struct(&user)

	v.AddRule("Username", "required")
	v.AddRule("Password", "required")

	v.AddMessages(map[string]string{
		"required": "the {field} is required",
	})

	v.AddTranslates(map[string]string{
		"Username": "login",
		"Password": "password",
	})

	if !v.Validate() {
		return serviceErr.NewValidationError(v.Errors)
	}

	return nil
}

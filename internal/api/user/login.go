package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	apiCommon "github.com/fishus/go-advanced-gophermart/internal/api/common"
	"github.com/fishus/go-advanced-gophermart/internal/logger"
	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
)

// Login Аутентификация пользователя
func (a *api) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type reqData struct {
		Username string `json:"login"`              // Логин
		Password string `json:"password,omitempty"` // Пароль
	}

	var data reqData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error())
		return
	}

	user := models.User{
		Username: data.Username,
		Password: data.Password,
	}

	userID, err := a.service.User().Login(r.Context(), user)
	if err != nil {
		var validErr *serviceErr.ValidationError
		if errors.As(err, &validErr) {
			apiCommon.JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, serviceErr.ErrUserNotFound) {
			apiCommon.JSONError(w, "Wrong login or password", http.StatusUnauthorized)
			return
		}
		apiCommon.JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	// аутентификация пользователя
	token, err := a.service.User().BuildToken(userID)
	if err != nil {
		apiCommon.JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
	}
	w.Header().Set("Authorization", ("Bearer " + token))

	w.WriteHeader(http.StatusOK)
}

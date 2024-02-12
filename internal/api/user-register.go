package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fishus/go-advanced-gophermart/internal/logger"
	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *server) userRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error())
		return
	}

	user, err := s.service.User().Register(r.Context(), user)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			JSONError(w, "User already exists", http.StatusConflict)
			return
		}
		var validErr *serviceErr.ValidationError
		if errors.As(err, &validErr) {
			JSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
		JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	// автоматическая аутентификация пользователя.
	token, err := s.service.User().BuildToken(user)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
	}
	w.Header().Set("Authorization", token)

	w.WriteHeader(http.StatusOK)
}

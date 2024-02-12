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

// userLogin Аутентификация пользователя
func (s *server) userLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		logger.Log.Debug(err.Error())
		return
	}

	userID, err := s.service.User().Login(r.Context(), user)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			JSONError(w, "Wrong login or password", http.StatusUnauthorized)
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

	// аутентификация пользователя
	token, err := s.service.User().BuildToken(userID)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		logger.Log.Error(err.Error())
	}
	w.Header().Set("Authorization", token)

	w.WriteHeader(http.StatusOK)
}

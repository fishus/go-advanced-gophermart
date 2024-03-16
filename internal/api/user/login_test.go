package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	sMocks "github.com/fishus/go-advanced-gophermart/internal/service/mocks"
)

func (ts *APITestSuite) TestLogin() {
	ctx := context.Background()

	url := "/login"

	type reqData struct {
		Username string `json:"login"`              // Логин
		Password string `json:"password,omitempty"` // Пароль
	}

	tests := []struct {
		name       string
		data       reqData
		loginErr   error
		respStatus int
	}{
		{
			name: "Positive case",
			data: reqData{
				Username: "testuser",
				Password: "12345",
			},
			loginErr:   nil,
			respStatus: http.StatusOK,
		},
		{
			name: "Required fields",
			data: reqData{
				Username: "",
				Password: "",
			},
			loginErr:   serviceErr.NewValidationError(serviceErr.NewValidationError(errors.New("required fields error"))),
			respStatus: http.StatusBadRequest,
		},
		{
			name: "User not found",
			data: reqData{
				Username: "testuser",
				Password: "12345",
			},
			loginErr:   serviceErr.ErrUserNotFound,
			respStatus: http.StatusUnauthorized,
		},
	}
	for _, tc := range tests {
		ts.Run(tc.name, func() {
			token := "VALID-JWT-TOKEN"
			userID := models.UserID(uuid.New().String())
			user := models.User{
				Username: tc.data.Username,
				Password: tc.data.Password,
			}

			sUser := sMocks.NewUserer(ts.T())

			mockUserLogin := sUser.EXPECT().Login(mock.AnythingOfType("*context.valueCtx"), user)
			if tc.loginErr == nil {
				mockUserLogin.Return(userID, nil)
			} else {
				mockUserLogin.Return("", tc.loginErr)
			}

			if tc.respStatus == http.StatusOK {
				sUser.EXPECT().BuildToken(userID).Return(token, nil)
			}

			ts.setServiceUser(sUser)

			body, err := json.Marshal(tc.data)
			ts.Require().NoError(err)
			resp, err := ts.client.R().
				SetContext(ctx).
				SetHeader("Content-Type", "application/json; charset=utf-8").
				SetBody(body).
				Post(url)
			ts.Require().NoError(err)
			ts.Equal(tc.respStatus, resp.StatusCode())

			if resp.StatusCode() == http.StatusOK {
				auth := resp.Header().Get("Authorization")
				ts.Equal("Bearer "+token, auth)
			}
		})
	}
}

package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	sMocks "github.com/fishus/go-advanced-gophermart/internal/service/mocks"
	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

func (ts *APITestSuite) TestUserBalance() {
	ctx := context.Background()

	url := "/api/user/balance"

	type LoyaltyBalanceResult struct {
		UserID    models.UserID `json:"-"`         // ID пользователя
		Current   float64       `json:"current"`   // Текущий баланс
		Accrued   float64       `json:"-"`         // Начислено за всё время
		Withdrawn float64       `json:"withdrawn"` // Списано за всё время
	}

	userID := models.UserID(uuid.New().String())

	tests := []struct {
		name       string
		auth       string
		want       models.LoyaltyBalance
		respStatus int
	}{
		{
			name: "Positive case",
			auth: "VALID-JWT-TOKEN",
			want: models.LoyaltyBalance{
				UserID:    userID,
				Current:   530.865, // FIXME 654.321 - 123.456
				Accrued:   654.321,
				Withdrawn: 123.456,
			},
			respStatus: http.StatusOK,
		},
		{
			name: "No balance",
			auth: "VALID-JWT-TOKEN",
			want: models.LoyaltyBalance{
				UserID: userID,
			},
			respStatus: http.StatusOK,
		},
		{
			name:       "Invalid auth token",
			auth:       "INVALID-JWT-TOKEN",
			respStatus: http.StatusUnauthorized,
		},
		{
			name:       "No auth token",
			auth:       "",
			respStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		ts.Run(tc.name, func() {
			sUser := sMocks.NewUserer(ts.T())

			var authToken string
			if tc.auth != "" {
				authToken = "Bearer " + tc.auth
			}
			mockUserCheckAuthorizationHeader := sUser.EXPECT().CheckAuthorizationHeader(authToken)
			if tc.auth == "VALID-JWT-TOKEN" {
				mockUserCheckAuthorizationHeader.Return(&uService.JWTClaims{UserID: userID}, nil)

				sUser.EXPECT().LoyaltyUserBalance(mock.AnythingOfType("*context.valueCtx"), userID).Return(tc.want, nil)
			} else {
				mockUserCheckAuthorizationHeader.Return(nil, uService.ErrInvalidToken)
			}

			ts.setService(nil, sUser)

			req := ts.client.R().SetContext(ctx)

			if tc.auth != "" {
				req.SetHeader("Authorization", "Bearer "+tc.auth)
			}

			resp, err := req.Get(url)
			ts.Require().NoError(err)
			ts.Equal(tc.respStatus, resp.StatusCode())

			if tc.respStatus == http.StatusOK {
				jsonBody, err := json.Marshal(LoyaltyBalanceResult(tc.want))
				ts.Require().NoError(err)
				ts.JSONEq(string(jsonBody), string(resp.Body()))
			}
		})
	}
}

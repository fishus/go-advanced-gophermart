package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/app/config"
	sMocks "github.com/fishus/go-advanced-gophermart/internal/service/mocks"
	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

func (ts *APITestSuite) TestUserBalance() {
	ctx := context.Background()

	url := "/api/user/balance"

	type want struct {
		Current   float64 `json:"current"`   // Текущий баланс
		Withdrawn float64 `json:"withdrawn"` // Списано за всё время
	}

	userID := models.UserID(uuid.New().String())

	tests := []struct {
		name       string
		auth       string
		data       models.LoyaltyBalance
		respStatus int
	}{
		{
			name: "Positive case",
			auth: "VALID-JWT-TOKEN",
			data: models.LoyaltyBalance{
				UserID:    userID,
				Accrued:   decimal.NewFromFloatWithExponent(654.321, -config.DecimalExponent),
				Withdrawn: decimal.NewFromFloatWithExponent(123.456, -config.DecimalExponent),
			},
			respStatus: http.StatusOK,
		},
		{
			name: "No balance",
			auth: "VALID-JWT-TOKEN",
			data: models.LoyaltyBalance{
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
			tc.data.Current = tc.data.Accrued.Sub(tc.data.Withdrawn)

			sUser := sMocks.NewUserer(ts.T())

			var authToken string
			if tc.auth != "" {
				authToken = "Bearer " + tc.auth
			}
			mockUserCheckAuthorizationHeader := sUser.EXPECT().CheckAuthorizationHeader(authToken)
			if tc.auth == "VALID-JWT-TOKEN" {
				mockUserCheckAuthorizationHeader.Return(&uService.JWTClaims{UserID: userID}, nil)

				sUser.EXPECT().LoyaltyUserBalance(mock.AnythingOfType("*context.valueCtx"), userID).Return(tc.data, nil)
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
				wantData := want{
					Current: func() float64 {
						f, _ := tc.data.Current.Float64()
						return f
					}(),
					Withdrawn: func() float64 {
						f, _ := tc.data.Withdrawn.Float64()
						return f
					}(),
				}
				jsonBody, err := json.Marshal(wantData)
				ts.Require().NoError(err)
				ts.JSONEq(string(jsonBody), string(resp.Body()))
			}
		})
	}
}

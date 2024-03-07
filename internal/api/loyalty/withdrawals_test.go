package loyalty

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/app/config"
	sMocks "github.com/fishus/go-advanced-gophermart/internal/service/mocks"
	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

func (ts *APITestSuite) TestWithdrawals() {
	ctx := context.Background()

	url := "/withdrawals"

	type want struct {
		OrderNum    string    `json:"order"`        // Номер заказа
		Withdrawal  float64   `json:"sum"`          // Списание
		ProcessedAt time.Time `json:"processed_at"` // Дата зачисления или списания
	}

	userID := models.UserID(uuid.New().String())

	tests := []struct {
		name       string
		auth       string
		data       []models.LoyaltyHistory
		respStatus int
	}{
		{
			name: "Positive case",
			auth: "VALID-JWT-TOKEN",
			data: []models.LoyaltyHistory{
				{
					UserID:      userID,
					OrderNum:    "6825296715",
					Accrual:     decimal.NewFromFloatWithExponent(654.321, -config.DecimalExponent),
					Withdrawal:  decimal.NewFromFloat(0),
					ProcessedAt: time.Now().UTC().Round(time.Minute),
				},
				{
					UserID:      userID,
					OrderNum:    "5347676263",
					Accrual:     decimal.NewFromFloat(0),
					Withdrawal:  decimal.NewFromFloatWithExponent(123.456, -config.DecimalExponent),
					ProcessedAt: time.Now().UTC().Round(time.Minute),
				},
			},
			respStatus: http.StatusOK,
		},
		{
			name:       "No withdrawals",
			auth:       "VALID-JWT-TOKEN",
			data:       nil,
			respStatus: http.StatusNoContent,
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
			sLoyalty := sMocks.NewLoyaltier(ts.T())

			var authToken string
			if tc.auth != "" {
				authToken = "Bearer " + tc.auth
			}
			mockUserCheckAuthorizationHeader := sUser.EXPECT().CheckAuthorizationHeader(authToken)
			if tc.auth == "VALID-JWT-TOKEN" {
				mockUserCheckAuthorizationHeader.Return(&uService.JWTClaims{UserID: userID}, nil)

				sLoyalty.EXPECT().UserWithdrawals(mock.AnythingOfType("*context.valueCtx"), userID).Return(tc.data, nil)
			} else {
				mockUserCheckAuthorizationHeader.Return(nil, uService.ErrInvalidToken)
			}

			ts.setServiceUser(sUser)
			ts.setServiceLoyalty(sLoyalty)

			req := ts.client.R().SetContext(ctx)

			if tc.auth != "" {
				req.SetHeader("Authorization", "Bearer "+tc.auth)
			}

			resp, err := req.Get(url)
			ts.Require().NoError(err)
			ts.Equal(tc.respStatus, resp.StatusCode())

			if tc.respStatus == http.StatusOK {
				var wantList []want
				for _, h := range tc.data {
					w := want{
						OrderNum: h.OrderNum,
						Withdrawal: func() float64 {
							f, _ := h.Withdrawal.Float64()
							return f
						}(),
						ProcessedAt: h.ProcessedAt,
					}
					wantList = append(wantList, w)
				}
				jsonBody, err := json.Marshal(wantList)
				ts.Require().NoError(err)
				ts.JSONEq(string(jsonBody), string(resp.Body()))
			}
		})
	}
}

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
	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	sMocks "github.com/fishus/go-advanced-gophermart/internal/service/mocks"
	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

func (ts *APITestSuite) TestUserWithdraw() {
	ctx := context.Background()

	url := "/api/user/balance/withdraw"

	userID := models.UserID(uuid.New().String())

	type reqData struct {
		Num string          `json:"order"` // Номер заказа
		Sum decimal.Decimal `json:"sum"`   // Сумма баллов к списанию в счёт оплаты
	}

	tests := []struct {
		name       string
		auth       string
		data       reqData
		wErr       error
		respStatus int
	}{
		{
			name: "Positive case",
			auth: "VALID-JWT-TOKEN",
			data: reqData{
				Num: "8542143048",
				Sum: decimal.NewFromFloatWithExponent(123.456, -config.DecimalExponent),
			},
			wErr:       nil,
			respStatus: http.StatusOK,
		},
		{
			name: "Low balance",
			auth: "VALID-JWT-TOKEN",
			data: reqData{
				Num: "8542143048",
				Sum: decimal.NewFromFloatWithExponent(123.456, -config.DecimalExponent),
			},
			wErr:       serviceErr.ErrLowBalance,
			respStatus: http.StatusPaymentRequired,
		},
		{
			name: "Incorrect data",
			auth: "VALID-JWT-TOKEN",
			data: reqData{
				Num: "8542143048",
				Sum: decimal.NewFromInt(-100).Round(config.DecimalExponent),
			},
			wErr:       serviceErr.ErrIncorrectData,
			respStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid number",
			auth: "VALID-JWT-TOKEN",
			data: reqData{
				Num: "55555",
				Sum: decimal.NewFromFloat(0).Round(config.DecimalExponent),
			},
			wErr:       serviceErr.ErrOrderWrongNum,
			respStatus: http.StatusUnprocessableEntity,
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
			sOrder := sMocks.NewOrderer(ts.T())

			var authToken string
			if tc.auth != "" {
				authToken = "Bearer " + tc.auth
			}
			mockUserCheckAuthorizationHeader := sUser.EXPECT().CheckAuthorizationHeader(authToken)
			if tc.auth == "VALID-JWT-TOKEN" {
				mockUserCheckAuthorizationHeader.Return(&uService.JWTClaims{UserID: userID}, nil)

				sUser.EXPECT().LoyaltyAddWithdraw(mock.AnythingOfType("*context.valueCtx"), userID, tc.data.Num, tc.data.Sum).Return(tc.wErr)
			} else {
				mockUserCheckAuthorizationHeader.Return(nil, uService.ErrInvalidToken)
			}

			ts.setService(sOrder, sUser)

			body, err := json.Marshal(tc.data)
			ts.Require().NoError(err)

			req := ts.client.R().
				SetContext(ctx).
				SetHeader("Content-Type", "application/json; charset=utf-8").
				SetBody(body)

			if tc.auth != "" {
				req.SetHeader("Authorization", "Bearer "+tc.auth)
			}

			resp, err := req.Post(url)
			ts.Require().NoError(err)
			ts.Equal(tc.respStatus, resp.StatusCode())
		})
	}
}

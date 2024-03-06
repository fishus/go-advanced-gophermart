package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	sMocks "github.com/fishus/go-advanced-gophermart/internal/service/mocks"
	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

func (ts *APITestSuite) TestOrdersList() {
	ctx := context.Background()

	url := "/api/user/orders"

	type want struct {
		Num        string             `json:"number"`            // Номер заказа
		Accrual    float64            `json:"accrual,omitempty"` // Начислено баллов лояльности
		Status     models.OrderStatus `json:"status"`            // Статус заказа
		UploadedAt time.Time          `json:"uploaded_at"`       // Дата и время добавления заказа
	}

	userID := models.UserID(uuid.New().String())

	tests := []struct {
		name       string
		auth       string
		data       []models.Order
		respStatus int
	}{
		{
			name: "Positive case",
			auth: "VALID-JWT-TOKEN",
			data: []models.Order{
				{
					Num:        "6825296715",
					Accrual:    decimal.NewFromFloatWithExponent(123.456, -5),
					Status:     models.OrderStatusProcessed,
					UploadedAt: time.Now().UTC().Round(time.Minute),
				},
			},
			respStatus: http.StatusOK,
		},
		{
			name:       "No orders",
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
			sOrder := sMocks.NewOrderer(ts.T())

			var authToken string
			if tc.auth != "" {
				authToken = "Bearer " + tc.auth
			}
			mockUserCheckAuthorizationHeader := sUser.EXPECT().CheckAuthorizationHeader(authToken)
			if tc.auth == "VALID-JWT-TOKEN" {
				mockUserCheckAuthorizationHeader.Return(&uService.JWTClaims{UserID: userID}, nil)
				sOrder.EXPECT().ListByUser(mock.AnythingOfType("*context.valueCtx"), userID).Return(tc.data, nil)
			} else {
				mockUserCheckAuthorizationHeader.Return(nil, uService.ErrInvalidToken)
			}

			ts.setService(sOrder, sUser)

			req := ts.client.R().SetContext(ctx)

			if tc.auth != "" {
				req.SetHeader("Authorization", "Bearer "+tc.auth)
			}

			resp, err := req.Get(url)
			ts.Require().NoError(err)
			ts.Equal(tc.respStatus, resp.StatusCode())

			if tc.respStatus == http.StatusOK {
				var wantList []want
				for _, o := range tc.data {
					w := want{
						Num: o.Num,
						Accrual: func() float64 {
							f, _ := o.Accrual.Float64()
							return f
						}(),
						Status:     o.Status,
						UploadedAt: o.UploadedAt,
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

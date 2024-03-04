package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	sMocks "github.com/fishus/go-advanced-gophermart/internal/service/mocks"
	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

func (ts *APITestSuite) TestOrdersList() {
	ctx := context.Background()

	url := "/api/user/orders"

	type OrderResult struct {
		ID         models.OrderID     `json:"-"`                 // ID заказа
		UserID     models.UserID      `json:"-"`                 // ID пользователя
		Num        string             `json:"number"`            // Номер заказа
		Accrual    float64            `json:"accrual,omitempty"` // Начислено баллов лояльности // FIXME
		Status     models.OrderStatus `json:"status"`            // Статус заказа
		UploadedAt time.Time          `json:"uploaded_at"`       // Дата и время добавления заказа
		UpdatedAt  time.Time          `json:"-"`                 // Дата и время обновления статуса заказа
	}

	userID := models.UserID(uuid.New().String())

	tests := []struct {
		name       string
		auth       string
		want       []OrderResult
		respStatus int
	}{
		{
			name: "Positive case",
			auth: "VALID-JWT-TOKEN",
			want: []OrderResult{
				{
					Num:        "6825296715",
					Accrual:    123.456,
					Status:     models.OrderStatusProcessed,
					UploadedAt: time.Now().UTC().Round(time.Minute),
				},
			},
			respStatus: http.StatusOK,
		},
		{
			name:       "No orders",
			auth:       "VALID-JWT-TOKEN",
			want:       nil,
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
				var oList []models.Order
				for _, o := range tc.want {
					oList = append(oList, models.Order(o))
				}
				sOrder.EXPECT().ListByUser(mock.AnythingOfType("*context.valueCtx"), userID).Return(oList, nil)
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
				jsonBody, err := json.Marshal(tc.want)
				ts.Require().NoError(err)
				ts.JSONEq(string(jsonBody), string(resp.Body()))
			}
		})
	}
}

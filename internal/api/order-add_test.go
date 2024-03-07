package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	sMocks "github.com/fishus/go-advanced-gophermart/internal/service/mocks"
	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
)

func (ts *APITestSuite) TestOrderAdd() {
	ctx := context.Background()

	url := "/api/user/orders"

	tests := []struct {
		name       string
		num        string
		auth       string
		orderErr   error
		respStatus int
	}{
		{
			name:       "Positive case",
			num:        "8163091187",
			auth:       "VALID-JWT-TOKEN",
			orderErr:   nil,
			respStatus: http.StatusAccepted,
		},
		{
			name:       "Already exists",
			num:        "8163091187",
			auth:       "VALID-JWT-TOKEN",
			orderErr:   serviceErr.ErrOrderAlreadyExists,
			respStatus: http.StatusOK,
		},
		{
			name:       "Invalid auth token",
			num:        "8163091187",
			auth:       "INVALID-JWT-TOKEN",
			orderErr:   nil,
			respStatus: http.StatusUnauthorized,
		},
		{
			name:       "No auth token",
			num:        "8163091187",
			auth:       "",
			orderErr:   nil,
			respStatus: http.StatusUnauthorized,
		},
		{
			name:       "No order number",
			num:        "",
			auth:       "VALID-JWT-TOKEN",
			orderErr:   serviceErr.NewValidationError(serviceErr.NewValidationError(errors.New("required fields error"))),
			respStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid order number (Luhn)",
			num:        "55555",
			auth:       "VALID-JWT-TOKEN",
			orderErr:   serviceErr.ErrOrderWrongNum,
			respStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "Wrong order owner",
			num:        "8163091187",
			auth:       "VALID-JWT-TOKEN",
			orderErr:   serviceErr.ErrOrderWrongOwner,
			respStatus: http.StatusConflict,
		},
	}
	for _, tc := range tests {
		ts.Run(tc.name, func() {
			orderID := models.OrderID(uuid.New().String())
			userID := models.UserID(uuid.New().String())

			wantOrder := models.Order{
				ID:      orderID,
				UserID:  userID,
				Num:     tc.num,
				Accrual: decimal.NewFromFloat(0),
				Status:  models.OrderStatusNew,
			}

			sUser := sMocks.NewUserer(ts.T())
			sOrder := sMocks.NewOrderer(ts.T())

			var authToken string
			if tc.auth != "" {
				authToken = "Bearer " + tc.auth
			}
			mockUserCheckAuthorizationHeader := sUser.EXPECT().CheckAuthorizationHeader(authToken)
			if tc.auth == "VALID-JWT-TOKEN" {
				mockUserCheckAuthorizationHeader.Return(&uService.JWTClaims{UserID: userID}, nil)

				mockOrderAdd := sOrder.EXPECT().Add(mock.AnythingOfType("*context.valueCtx"), userID, tc.num)
				if tc.orderErr == nil {
					mockOrderAdd.Return(orderID, nil)

					sOrder.EXPECT().GetByID(mock.AnythingOfType("context.backgroundCtx"), orderID).Return(wantOrder, nil)

					ts.loyalty.EXPECT().AddNewOrder(mock.AnythingOfType("context.backgroundCtx"), wantOrder)
				} else {
					mockOrderAdd.Return("", tc.orderErr)
				}
			} else {
				mockUserCheckAuthorizationHeader.Return(nil, uService.ErrInvalidToken)
			}

			ts.setService(sOrder, sUser)

			req := ts.client.R().
				SetContext(ctx).
				SetHeader("Content-Type", "text/plain; charset=utf-8").
				SetBody(tc.num)

			if tc.auth != "" {
				req.SetHeader("Authorization", "Bearer "+tc.auth)
			}

			resp, err := req.Post(url)
			ts.Require().NoError(err)
			ts.Equal(tc.respStatus, resp.StatusCode())
		})
	}
}

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

func (ts *APITestSuite) TestUserWithdrawals() {
	ctx := context.Background()

	url := "/api/user/withdrawals"

	type LoyaltyHistoryResult struct {
		UserID      models.UserID `json:"-"`            // ID пользователя
		OrderNum    string        `json:"order"`        // Номер заказа
		Accrual     float64       `json:"-"`            // Начисление
		Withdrawal  float64       `json:"sum"`          // Списание
		ProcessedAt time.Time     `json:"processed_at"` // Дата зачисления или списания
	}

	userID := models.UserID(uuid.New().String())

	tests := []struct {
		name       string
		auth       string
		want       []models.LoyaltyHistory
		respStatus int
	}{
		{
			name: "Positive case",
			auth: "VALID-JWT-TOKEN",
			want: []models.LoyaltyHistory{
				{
					UserID:      userID,
					OrderNum:    "6825296715",
					Accrual:     654.321,
					Withdrawal:  0,
					ProcessedAt: time.Now().UTC().Round(time.Minute),
				},
				{
					UserID:      userID,
					OrderNum:    "5347676263",
					Accrual:     0,
					Withdrawal:  123.456,
					ProcessedAt: time.Now().UTC().Round(time.Minute),
				},
			},
			respStatus: http.StatusOK,
		},
		{
			name:       "No withdrawals",
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

			var authToken string
			if tc.auth != "" {
				authToken = "Bearer " + tc.auth
			}
			mockUserCheckAuthorizationHeader := sUser.EXPECT().CheckAuthorizationHeader(authToken)
			if tc.auth == "VALID-JWT-TOKEN" {
				mockUserCheckAuthorizationHeader.Return(&uService.JWTClaims{UserID: userID}, nil)

				sUser.EXPECT().LoyaltyUserWithdrawals(mock.AnythingOfType("*context.valueCtx"), userID).Return(tc.want, nil)
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
				var hList []LoyaltyHistoryResult
				for _, h := range tc.want {
					hList = append(hList, LoyaltyHistoryResult(h))
				}
				jsonBody, err := json.Marshal(hList)
				ts.Require().NoError(err)
				ts.JSONEq(string(jsonBody), string(resp.Body()))
			}
		})
	}
}

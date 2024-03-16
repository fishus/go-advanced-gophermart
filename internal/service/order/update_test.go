package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	stMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

func (ts *OrderServiceTestSuite) TestUpdateStatus() {
	ctx := context.Background()

	orderID := models.OrderID(uuid.New().String())

	testCases := []struct {
		name    string
		status  models.OrderStatus
		wantErr bool
	}{
		{
			"Status Processing",
			models.OrderStatusProcessing,
			false,
		},
		{
			"Status Invalid",
			models.OrderStatusInvalid,
			false,
		},
		{
			"Status Processed",
			models.OrderStatusProcessed,
			false,
		},
		{
			"Status New",
			models.OrderStatusNew,
			false,
		},
		{
			"Wrong status",
			"test1234",
			true,
		},
		{
			"Undefined status",
			models.OrderStatusUndefined,
			true,
		},
	}

	for _, tc := range testCases {
		ts.Run(tc.name, func() {
			stOrder := stMocks.NewOrderer(ts.T())
			switch tc.name {
			case "Wrong status",
				"Undefined status":
			default:
				stOrder.EXPECT().UpdateStatus(ctx, orderID, tc.status).Return(nil)
			}
			ts.setStorage(stOrder, nil, nil)

			err := ts.service.UpdateStatus(ctx, orderID, tc.status)
			if tc.wantErr {
				ts.Error(err)
			} else {
				ts.storage.AssertExpectations(ts.T())
				ts.NoError(err)
			}
		})
	}
}

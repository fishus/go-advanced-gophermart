package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (ts *PostgresTestSuite) TestOrderResetProcessingStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	bUsername := make([]byte, 10)
	_, err := rand.Read(bUsername)
	ts.Require().NoError(err)

	userData := models.User{
		Username: hex.EncodeToString(bUsername),
		Password: hex.EncodeToString(bUsername),
	}
	userID, err := ts.storage.UserAdd(ctx, userData)
	ts.Require().NoError(err)

	orderData := make([]models.Order, 3)
	orderData[0] = models.Order{
		Num:    "2731332660",
		Status: models.OrderStatusNew,
	}
	orderData[1] = models.Order{
		Num:    "6147497876",
		Status: models.OrderStatusProcessing,
	}
	orderData[2] = models.Order{
		Num:    "8853861980",
		Status: models.OrderStatusProcessed,
	}
	for i := 0; i < len(orderData); i++ {
		orderData[i].UserID = userID
		orderData[i].UploadedAt = time.Now().UTC().Round(1 * time.Second)
		orderID, err := ts.storage.OrderAdd(ctx, orderData[i])
		ts.Require().NoError(err)
		orderData[i].ID = orderID
	}
	orderData[1].Status = models.OrderStatusNew

	err = ts.storage.OrderResetProcessingStatus(ctx)
	ts.NoError(err)

	for _, want := range orderData {
		order, err := ts.storage.OrderByID(ctx, want.ID)
		ts.NoError(err)
		order.UploadedAt = order.UploadedAt.UTC().Round(1 * time.Second)
		ts.EqualValues(want, order)
	}

}

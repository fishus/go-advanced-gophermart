package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"slices"
	"time"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
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
		orderData[i].UploadedAt = time.Now().UTC().Round(5 * time.Second)
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
		order.UploadedAt = order.UploadedAt.UTC().Round(5 * time.Second)
		ts.EqualValues(want, order)
	}

}

func (ts *PostgresTestSuite) TestOrderMoveToProcessing() {
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

	orderData := make([]models.Order, 4)
	orderData[0] = models.Order{
		Num:    "2731332660",
		Status: models.OrderStatusProcessing,
	}
	orderData[1] = models.Order{
		Num:    "6147497876",
		Status: models.OrderStatusNew,
	}
	orderData[2] = models.Order{
		Num:    "8853861980",
		Status: models.OrderStatusProcessed,
	}
	orderData[3] = models.Order{
		Num:    "7819289724",
		Status: models.OrderStatusNew,
	}
	for i := 0; i < len(orderData); i++ {
		orderData[i].UserID = userID
		orderData[i].UploadedAt = time.Now().UTC().Round(5 * time.Second)
		orderID, err := ts.storage.OrderAdd(ctx, orderData[i])
		ts.Require().NoError(err)
		orderData[i].ID = orderID
	}
	orderData[1].Status = models.OrderStatusProcessing
	orderData[3].Status = models.OrderStatusProcessing

	orders, err := ts.storage.OrderMoveToProcessing(ctx, 10)
	ts.NoError(err)
	ts.Equal(2, len(orders))

	for _, order := range orders {
		i := slices.IndexFunc(orderData, func(o models.Order) bool {
			return o.Num == order.Num
		})
		order.UploadedAt = order.UploadedAt.UTC().Round(5 * time.Second)
		ts.EqualValues(orderData[i], order)
	}
}

func (ts *PostgresTestSuite) TestOrderSetStatus() {
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

	orderData := make([]models.Order, 4)
	orderData[0] = models.Order{
		Num:    "2731332660",
		Status: models.OrderStatusProcessing,
	}
	orderData[1] = models.Order{
		Num:    "6147497876",
		Status: models.OrderStatusNew,
	}
	orderData[2] = models.Order{
		Num:    "8853861980",
		Status: models.OrderStatusProcessing,
	}
	orderData[3] = models.Order{
		Num:    "7819289724",
		Status: models.OrderStatusNew,
	}
	for i := 0; i < len(orderData); i++ {
		orderData[i].UserID = userID
		orderData[i].UploadedAt = time.Now().UTC().Round(5 * time.Second)
		orderID, err := ts.storage.OrderAdd(ctx, orderData[i])
		ts.Require().NoError(err)
		orderData[i].ID = orderID
	}
	orderData[1].Status = models.OrderStatusProcessing
	orderData[3].Status = models.OrderStatusProcessing
	orderData[0].Status = models.OrderStatusProcessed
	orderData[2].Status = models.OrderStatusProcessed

	err = ts.storage.OrderSetStatus(ctx, []models.OrderID{orderData[1].ID, orderData[3].ID}, models.OrderStatusProcessing)
	ts.NoError(err)

	orders, err := ts.storage.OrdersByFilter(ctx, 10, store.WithOrderIDList([]models.OrderID{orderData[1].ID, orderData[3].ID}...))
	ts.NoError(err)
	for _, order := range orders {
		order.UploadedAt = order.UploadedAt.UTC().Round(5 * time.Second)
		i := slices.IndexFunc(orderData, func(o models.Order) bool {
			return o.Num == order.Num
		})
		ts.EqualValues(orderData[i], order)
	}

	err = ts.storage.OrderSetStatus(ctx, []models.OrderID{orderData[0].ID, orderData[2].ID}, models.OrderStatusProcessed)
	ts.NoError(err)

	orders, err = ts.storage.OrdersByFilter(ctx, 10, store.WithOrderIDList([]models.OrderID{orderData[0].ID, orderData[2].ID}...))
	ts.NoError(err)
	for _, order := range orders {
		order.UploadedAt = order.UploadedAt.UTC().Round(5 * time.Second)
		i := slices.IndexFunc(orderData, func(o models.Order) bool {
			return o.Num == order.Num
		})
		ts.EqualValues(orderData[i], order)
	}
}

package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"math"
	mrand "math/rand"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestOrdersByFilter() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	bUsername := make([]byte, 10)
	_, err := rand.Read(bUsername)
	ts.Require().NoError(err)

	userData := models.User{
		ID:        models.UserID(uuid.New().String()),
		Username:  hex.EncodeToString(bUsername),
		Password:  hex.EncodeToString(bUsername),
		CreatedAt: time.Now().UTC().Round(1 * time.Second),
	}

	_, err = ts.pool.Exec(ctx, `INSERT INTO users (id, username, password, created_at) VALUES (@id, @username, @password, @created_at);`,
		pgx.NamedArgs{
			"id":         userData.ID,
			"username":   userData.Username,
			"password":   userData.Password,
			"created_at": userData.CreatedAt,
		})
	ts.NoError(err)

	ratio := math.Pow(10, float64(5))
	orderNums := []string{"4581347475", "4941161632", "6860728705"}
	orderStatuses := []models.OrderStatus{models.OrderStatusNew, models.OrderStatusProcessing, models.OrderStatusProcessing}
	orderData := make([]models.Order, len(orderNums))
	for i, orderNum := range orderNums {
		orderData[i] = models.Order{
			ID:         models.OrderID(uuid.New().String()),
			UserID:     userData.ID,
			Num:        orderNum,
			Accrual:    (math.Round(mrand.Float64()*ratio) / ratio),
			Status:     orderStatuses[i],
			UploadedAt: time.Now().UTC().Round(1 * time.Second),
		}
		_, err = ts.pool.Exec(ctx, `INSERT INTO orders (id, user_id, num, accrual, status, uploaded_at) VALUES (@id, @userID, @num, @accrual, @status, @uploadedAt);`,
			pgx.NamedArgs{
				"id":         orderData[i].ID,
				"userID":     orderData[i].UserID,
				"num":        orderData[i].Num,
				"accrual":    orderData[i].Accrual,
				"status":     orderData[i].Status,
				"uploadedAt": orderData[i].UploadedAt,
			})
		ts.NoError(err)
	}

	ts.Run("WithOrderNum", func() {
		orders, err := ts.storage.OrdersByFilter(ctx, 10, store.WithOrderNum(orderData[0].Num))
		ts.NoError(err)
		ts.Equal(orderData[0].Num, orders[0].Num)
		ts.EqualValues(orderData[0], orders[0])
	})

	ts.Run("WithOrderUserID", func() {
		limit := 2
		orders, err := ts.storage.OrdersByFilter(ctx, limit, store.WithOrderUserID(userData.ID))
		ts.NoError(err)
		ts.Equal(limit, len(orders))
		for _, order := range orders {
			i := slices.Index(orderNums, order.Num)
			ts.Equal(userData.ID, order.UserID)
			ts.EqualValues(orderData[i], order)
		}
	})

	ts.Run("WithOrderStatus", func() {
		orders, err := ts.storage.OrdersByFilter(ctx, 10, store.WithOrderStatus(models.OrderStatusNew))
		ts.NoError(err)
		for _, order := range orders {
			i := slices.Index(orderNums, order.Num)
			ts.Equal(models.OrderStatusNew, order.Status)
			ts.EqualValues(orderData[i], order)
		}
	})

	ts.Run("WithOrderStatuses", func() {
		orders, err := ts.storage.OrdersByFilter(ctx, 10, store.WithOrderStatuses(models.OrderStatusProcessing))
		ts.NoError(err)
		for _, order := range orders {
			i := slices.Index(orderNums, order.Num)
			ts.Equal(models.OrderStatusProcessing, order.Status)
			ts.EqualValues(orderData[i], order)
		}
	})
}

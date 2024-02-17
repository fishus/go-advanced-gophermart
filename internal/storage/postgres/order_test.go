package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *PostgresTestSuite) TestOrderByID() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	bUsername := make([]byte, 10)
	_, err := rand.Read(bUsername)
	ts.Require().NoError(err)

	userData := &models.User{
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

	orderData := models.Order{
		ID:         models.OrderID(uuid.New().String()),
		UserID:     userData.ID,
		Num:        "8020122696",
		Accrual:    0,
		Status:     models.OrderStatusNew,
		UploadedAt: time.Now().UTC().Round(1 * time.Second),
	}

	_, err = ts.pool.Exec(ctx, `INSERT INTO orders (id, user_id, num, accrual, status, uploaded_at) VALUES (@id, @userID, @num, @accrual, @status, @uploadedAt);`,
		pgx.NamedArgs{
			"id":         orderData.ID,
			"userID":     orderData.UserID,
			"num":        orderData.Num,
			"accrual":    orderData.Accrual,
			"status":     orderData.Status,
			"uploadedAt": orderData.UploadedAt,
		})
	ts.NoError(err)

	order, err := ts.storage.OrderByID(ctx, orderData.ID)
	ts.NoError(err)
	ts.EqualValues(orderData, order)
}

func (ts *PostgresTestSuite) TestOrderByFilter() {
	ctx, cancel := context.WithTimeout(context.Background(), ts.cfg.QueryTimeout)
	defer cancel()

	orderNums := []string{"5431720977", "5882492415"}
	userData := make([]models.User, len(orderNums))
	orderData := make([]models.Order, len(orderNums))
	for i, orderNum := range orderNums {
		bUsername := make([]byte, 10)
		_, err := rand.Read(bUsername)
		ts.Require().NoError(err)

		userData[i] = models.User{
			ID:        models.UserID(uuid.New().String()),
			Username:  hex.EncodeToString(bUsername),
			Password:  hex.EncodeToString(bUsername),
			CreatedAt: time.Now().UTC().Round(1 * time.Second),
		}

		_, err = ts.pool.Exec(ctx, `INSERT INTO users (id, username, password, created_at) VALUES (@id, @username, @password, @created_at);`,
			pgx.NamedArgs{
				"id":         userData[i].ID,
				"username":   userData[i].Username,
				"password":   userData[i].Password,
				"created_at": userData[i].CreatedAt,
			})
		ts.NoError(err)

		orderData[i] = models.Order{
			ID:         models.OrderID(uuid.New().String()),
			UserID:     userData[i].ID,
			Num:        orderNum,
			Accrual:    0,
			Status:     models.OrderStatusNew,
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
		order, err := ts.storage.OrderByFilter(ctx, store.WithOrderNum(orderData[0].Num))
		ts.NoError(err)
		ts.EqualValues(orderData[0], order)
	})

	ts.Run("WithOrderUserID", func() {
		order, err := ts.storage.OrderByFilter(ctx, store.WithOrderUserID(orderData[1].UserID))
		ts.NoError(err)
		ts.EqualValues(orderData[1], order)
	})
}

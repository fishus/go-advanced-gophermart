package accrual

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	sMocks "github.com/fishus/go-advanced-gophermart/internal/service/mocks"
)

func (ts *LoyaltyTestSuite) TestAddNewOrders() {
	ctx := context.Background()

	userID := models.UserID(uuid.New().String())
	orderNum := "1853241857"

	ts.Run("Positive case", func() {
		sOrder := sMocks.NewOrderer(ts.T())
		list := []models.Order{
			{
				ID:         models.OrderID(uuid.New().String()),
				UserID:     userID,
				Num:        orderNum,
				Accrual:    0,
				Status:     models.OrderStatusNew,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
		}
		mockCall := sOrder.EXPECT().ListNew(ctx).Return(list, nil)
		defer mockCall.Unset()
		ts.setService(sOrder, nil)

		ts.daemon.addNewOrders(ctx)

		ctxTC, cancel := context.WithTimeout(ctx, (1 * time.Second))
		defer cancel()
		select {
		case <-ctxTC.Done():
			ts.Fail(ctxTC.Err().Error())
		case order := <-ts.daemon.chOrders:
			ts.EqualValues(list[0], order)
		}
		sOrder.AssertExpectations(ts.T())
	})

	ts.Run("No new orders", func() {
		sOrder := sMocks.NewOrderer(ts.T())
		mockCall := sOrder.EXPECT().ListNew(ctx).Return(nil, nil)
		defer mockCall.Unset()
		ts.setService(sOrder, nil)

		ts.daemon.addNewOrders(ctx)

		ctxTC, cancel := context.WithTimeout(ctx, (1 * time.Second))
		defer cancel()
		select {
		case <-ctxTC.Done():
			ts.Error(ctxTC.Err())
		case <-ts.daemon.chOrders:
			ts.Fail("Unexpected order")
		}
		sOrder.AssertExpectations(ts.T())
	})
}

func (ts *LoyaltyTestSuite) TestAddProcessingOrders() {
	ctx := context.Background()

	//ts.T().SkipNow()

	userID := models.UserID(uuid.New().String())
	orderNum := "5347676263"

	ts.Run("Positive case", func() {
		sOrder := sMocks.NewOrderer(ts.T())
		list := []models.Order{
			{
				ID:         models.OrderID(uuid.New().String()),
				UserID:     userID,
				Num:        orderNum,
				Accrual:    0,
				Status:     models.OrderStatusProcessing,
				UploadedAt: time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			},
		}
		mockCall := sOrder.EXPECT().ListProcessing(ctx, 1).Return(list, nil)
		defer mockCall.Unset()
		ts.setService(sOrder, nil)

		ts.daemon.addProcessingOrders(ctx)

		ctxT, cancel := context.WithTimeout(ctx, (1500 * time.Millisecond))
		defer cancel()

		select {
		case <-ctxT.Done():
			ts.Fail(ctxT.Err().Error())
		case order := <-ts.daemon.chOrders:
			ts.EqualValues(list[0], order)
		}
		close(ts.daemon.chShutdown)
		ts.daemon.wg.Wait()
		ts.daemon.chShutdown = make(chan struct{})
		sOrder.AssertExpectations(ts.T())
	})

	ts.Run("No new orders", func() {
		sOrder := sMocks.NewOrderer(ts.T())
		mockCall := sOrder.EXPECT().ListProcessing(ctx, 1).Return(nil, nil)
		defer mockCall.Unset()
		ts.setService(sOrder, nil)

		ts.daemon.addProcessingOrders(ctx)

		ctxT, cancel := context.WithTimeout(ctx, (1500 * time.Millisecond))
		defer cancel()

		select {
		case <-ctxT.Done():
			ts.Error(ctxT.Err())
		case <-ts.daemon.chOrders:
			ts.Fail("Unexpected order")
		}
		close(ts.daemon.chShutdown)
		ts.daemon.wg.Wait()
		ts.daemon.chShutdown = make(chan struct{})
		sOrder.AssertExpectations(ts.T())
	})
}

func (ts *LoyaltyTestSuite) TestAddNewOrder() {
	ctx := context.Background()

	ts.Run("Positive case", func() {
		wantOrder := models.Order{
			ID:         models.OrderID(uuid.New().String()),
			UserID:     models.UserID(uuid.New().String()),
			Num:        "9400781309",
			Accrual:    0,
			Status:     models.OrderStatusNew,
			UploadedAt: time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}

		ts.daemon.AddNewOrder(ctx, wantOrder)

		ctxTC, cancel := context.WithTimeout(ctx, (1 * time.Second))
		defer cancel()
		select {
		case <-ctxTC.Done():
			ts.Fail(ctxTC.Err().Error())
		case order := <-ts.daemon.chOrders:
			ts.EqualValues(wantOrder, order)
		}
	})
}

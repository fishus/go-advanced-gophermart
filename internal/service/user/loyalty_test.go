package user

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (ts *UserServiceTestSuite) TestLoyaltyUserBalance() {
	ctx := context.Background()
	userID := models.UserID(uuid.New().String())

	ts.Run("Positive case", func() {
		wantBalance := models.LoyaltyBalance{
			UserID:    userID,
			Current:   227.482,
			Accrued:   752.113,
			Withdrawn: 524.631,
		}
		mockCall := ts.storage.On("LoyaltyBalanceByUser", ctx, userID).Return(wantBalance, nil)

		balance, err := ts.service.LoyaltyUserBalance(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantBalance, balance)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})

	ts.Run("Balance not found", func() {
		mockCall := ts.storage.On("LoyaltyBalanceByUser", ctx, userID).Return(models.LoyaltyBalance{}, store.ErrNotFound)

		wantBalance := models.LoyaltyBalance{
			UserID:    userID,
			Current:   0,
			Accrued:   0,
			Withdrawn: 0,
		}
		balance, err := ts.service.LoyaltyUserBalance(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantBalance, balance)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})
}

func (ts *UserServiceTestSuite) TestLoyaltyUserWithdrawals() {
	ctx := context.Background()
	userID := models.UserID(uuid.New().String())

	ts.Run("Positive case", func() {
		userHistory := make([]models.LoyaltyHistory, 2)
		userHistory[0] = models.LoyaltyHistory{
			UserID:      userID,
			OrderNum:    "5347676263",
			Accrual:     1123.456,
			Withdrawal:  0,
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		userHistory[1] = models.LoyaltyHistory{
			UserID:      userID,
			OrderNum:    "8163091187",
			Accrual:     0,
			Withdrawal:  654.321,
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		mockCall := ts.storage.On("LoyaltyHistoryByUser", ctx, userID).Return(userHistory, nil)

		wantWithdrawals := make([]models.LoyaltyHistory, 0)
		wantWithdrawals = append(wantWithdrawals, userHistory[1])

		withdrawals, err := ts.service.LoyaltyUserWithdrawals(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantWithdrawals, withdrawals)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})

	ts.Run("No withdrawals", func() {
		userHistory := make([]models.LoyaltyHistory, 2)
		userHistory[0] = models.LoyaltyHistory{
			UserID:      userID,
			OrderNum:    "5347676263",
			Accrual:     1123.456,
			Withdrawal:  0,
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		mockCall := ts.storage.On("LoyaltyHistoryByUser", ctx, userID).Return(userHistory, nil)

		wantWithdrawals := make([]models.LoyaltyHistory, 0)

		withdrawals, err := ts.service.LoyaltyUserWithdrawals(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantWithdrawals, withdrawals)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})

	ts.Run("History not found", func() {
		var emptyHistory []models.LoyaltyHistory
		mockCall := ts.storage.On("LoyaltyHistoryByUser", ctx, userID).Return(emptyHistory, store.ErrNotFound)

		withdrawals, err := ts.service.LoyaltyUserWithdrawals(ctx, userID)
		ts.NoError(err)
		ts.Nil(withdrawals)
		ts.Equal(len(withdrawals), 0)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})
}

func (ts *UserServiceTestSuite) TestLoyaltyAddWithdraw() {
	ctx := context.Background()
	userID := models.UserID(uuid.New().String())

	ts.Run("Positive case", func() {
		orderNum := "3903733214"
		withdraw := 659.784 // FIXME
		mockCall := ts.storage.On("LoyaltyAddWithdraw", ctx, userID, orderNum, withdraw).Return(nil)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.NoError(err)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})

	ts.Run("Invalid number", func() {
		orderNum := "126378912"
		withdraw := 659.784 // FIXME
		ts.storage.AssertNotCalled(ts.T(), "LoyaltyAddWithdraw", ctx, userID, orderNum, withdraw)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderWrongNum)
	})

	ts.Run("Incorrect withdraw amount", func() {
		orderNum := "3903733214"
		withdraw := -1.0 // FIXME
		ts.storage.AssertNotCalled(ts.T(), "LoyaltyAddWithdraw", ctx, userID, orderNum, withdraw)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrIncorrectData)
	})

	ts.Run("Low balance", func() {
		orderNum := "3903733214"
		withdraw := 659.784 // FIXME
		mockCall := ts.storage.On("LoyaltyAddWithdraw", ctx, userID, orderNum, withdraw).Return(store.ErrLowBalance)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrLowBalance)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})

	ts.Run("Balance not found", func() {
		orderNum := "3903733214"
		withdraw := 659.784 // FIXME
		mockCall := ts.storage.On("LoyaltyAddWithdraw", ctx, userID, orderNum, withdraw).Return(store.ErrNotFound)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrLowBalance)
		ts.storage.AssertExpectations(ts.T())
		mockCall.Unset()
	})
}

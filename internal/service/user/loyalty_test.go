package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

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
			Accrued:   decimal.NewFromFloatWithExponent(752.113, -5),
			Withdrawn: decimal.NewFromFloatWithExponent(524.631, -5),
		}
		wantBalance.Current = wantBalance.Accrued.Sub(wantBalance.Withdrawn)
		mockCall := ts.storage.EXPECT().LoyaltyBalanceByUser(ctx, userID).Return(wantBalance, nil)
		defer mockCall.Unset()

		balance, err := ts.service.LoyaltyUserBalance(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantBalance, balance)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Balance not found", func() {
		mockCall := ts.storage.EXPECT().LoyaltyBalanceByUser(ctx, userID).Return(models.LoyaltyBalance{}, store.ErrNotFound)
		defer mockCall.Unset()

		wantBalance := models.LoyaltyBalance{
			UserID: userID,
		}
		balance, err := ts.service.LoyaltyUserBalance(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantBalance, balance)
		ts.storage.AssertExpectations(ts.T())
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
			Accrual:     decimal.NewFromFloatWithExponent(1123.456, -5),
			Withdrawal:  decimal.NewFromFloat(0),
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		userHistory[1] = models.LoyaltyHistory{
			UserID:      userID,
			OrderNum:    "8163091187",
			Accrual:     decimal.NewFromFloat(0),
			Withdrawal:  decimal.NewFromFloatWithExponent(654.321, -5),
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		mockCall := ts.storage.EXPECT().LoyaltyHistoryByUser(ctx, userID).Return(userHistory, nil)
		defer mockCall.Unset()

		wantWithdrawals := make([]models.LoyaltyHistory, 0)
		wantWithdrawals = append(wantWithdrawals, userHistory[1])

		withdrawals, err := ts.service.LoyaltyUserWithdrawals(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantWithdrawals, withdrawals)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("No withdrawals", func() {
		userHistory := make([]models.LoyaltyHistory, 2)
		userHistory[0] = models.LoyaltyHistory{
			UserID:      userID,
			OrderNum:    "5347676263",
			Accrual:     decimal.NewFromFloatWithExponent(1123.456, -5),
			Withdrawal:  decimal.NewFromFloat(0),
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		mockCall := ts.storage.EXPECT().LoyaltyHistoryByUser(ctx, userID).Return(userHistory, nil)
		defer mockCall.Unset()

		wantWithdrawals := make([]models.LoyaltyHistory, 0)

		withdrawals, err := ts.service.LoyaltyUserWithdrawals(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantWithdrawals, withdrawals)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("History not found", func() {
		var emptyHistory []models.LoyaltyHistory
		mockCall := ts.storage.EXPECT().LoyaltyHistoryByUser(ctx, userID).Return(emptyHistory, store.ErrNotFound)
		defer mockCall.Unset()

		withdrawals, err := ts.service.LoyaltyUserWithdrawals(ctx, userID)
		ts.NoError(err)
		ts.Nil(withdrawals)
		ts.Equal(len(withdrawals), 0)
		ts.storage.AssertExpectations(ts.T())
	})
}

func (ts *UserServiceTestSuite) TestLoyaltyAddWithdraw() {
	ctx := context.Background()
	userID := models.UserID(uuid.New().String())

	ts.Run("Positive case", func() {
		orderNum := "3903733214"
		withdraw := decimal.NewFromFloatWithExponent(659.784, -5)
		mockCall := ts.storage.EXPECT().LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw).Return(nil)
		defer mockCall.Unset()

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.NoError(err)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Invalid number", func() {
		orderNum := "126378912"
		withdraw := decimal.NewFromFloatWithExponent(659.784, -5)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderWrongNum)
	})

	ts.Run("Incorrect withdraw amount", func() {
		orderNum := "3903733214"
		withdraw := decimal.NewFromFloatWithExponent(-1.0, -5)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrIncorrectData)
	})

	ts.Run("Low balance", func() {
		orderNum := "3903733214"
		withdraw := decimal.NewFromFloatWithExponent(659.784, -5)
		mockCall := ts.storage.EXPECT().LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw).Return(store.ErrLowBalance)
		defer mockCall.Unset()

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrLowBalance)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Balance not found", func() {
		orderNum := "3903733214"
		withdraw := decimal.NewFromFloatWithExponent(659.784, -5)
		mockCall := ts.storage.EXPECT().LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw).Return(store.ErrNotFound)
		defer mockCall.Unset()

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrLowBalance)
		ts.storage.AssertExpectations(ts.T())
	})
}

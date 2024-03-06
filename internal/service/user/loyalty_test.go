package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	"github.com/fishus/go-advanced-gophermart/internal/app/config"
	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
	stMocks "github.com/fishus/go-advanced-gophermart/internal/storage/mocks"
)

func (ts *UserServiceTestSuite) TestLoyaltyUserBalance() {
	ctx := context.Background()
	userID := models.UserID(uuid.New().String())

	ts.Run("Positive case", func() {
		wantBalance := models.LoyaltyBalance{
			UserID:    userID,
			Accrued:   decimal.NewFromFloatWithExponent(752.113, -config.DecimalExponent),
			Withdrawn: decimal.NewFromFloatWithExponent(524.631, -config.DecimalExponent),
		}
		wantBalance.Current = wantBalance.Accrued.Sub(wantBalance.Withdrawn)

		stLoyalty := stMocks.NewLoyaltier(ts.T())
		stLoyalty.EXPECT().BalanceByUser(ctx, userID).Return(wantBalance, nil)
		ts.setStorage(nil, nil, stLoyalty)

		balance, err := ts.service.LoyaltyUserBalance(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantBalance, balance)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Balance not found", func() {
		stLoyalty := stMocks.NewLoyaltier(ts.T())
		stLoyalty.EXPECT().BalanceByUser(ctx, userID).Return(models.LoyaltyBalance{}, store.ErrNotFound)
		ts.setStorage(nil, nil, stLoyalty)

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
			Accrual:     decimal.NewFromFloatWithExponent(1123.456, -config.DecimalExponent),
			Withdrawal:  decimal.NewFromFloat(0),
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}
		userHistory[1] = models.LoyaltyHistory{
			UserID:      userID,
			OrderNum:    "8163091187",
			Accrual:     decimal.NewFromFloat(0),
			Withdrawal:  decimal.NewFromFloatWithExponent(654.321, -config.DecimalExponent),
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}

		stLoyalty := stMocks.NewLoyaltier(ts.T())
		stLoyalty.EXPECT().HistoryByUser(ctx, userID).Return(userHistory, nil)
		ts.setStorage(nil, nil, stLoyalty)

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
			Accrual:     decimal.NewFromFloatWithExponent(1123.456, -config.DecimalExponent),
			Withdrawal:  decimal.NewFromFloat(0),
			ProcessedAt: time.Now().UTC().Round(time.Minute),
		}

		stLoyalty := stMocks.NewLoyaltier(ts.T())
		stLoyalty.EXPECT().HistoryByUser(ctx, userID).Return(userHistory, nil)
		ts.setStorage(nil, nil, stLoyalty)

		wantWithdrawals := make([]models.LoyaltyHistory, 0)

		withdrawals, err := ts.service.LoyaltyUserWithdrawals(ctx, userID)
		ts.NoError(err)
		ts.EqualValues(wantWithdrawals, withdrawals)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("History not found", func() {
		var emptyHistory []models.LoyaltyHistory

		stLoyalty := stMocks.NewLoyaltier(ts.T())
		stLoyalty.EXPECT().HistoryByUser(ctx, userID).Return(emptyHistory, store.ErrNotFound)
		ts.setStorage(nil, nil, stLoyalty)

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
		withdraw := decimal.NewFromFloatWithExponent(659.784, -config.DecimalExponent)

		stLoyalty := stMocks.NewLoyaltier(ts.T())
		stLoyalty.EXPECT().AddWithdraw(ctx, userID, orderNum, withdraw).Return(nil)
		ts.setStorage(nil, nil, stLoyalty)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.NoError(err)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Invalid number", func() {
		orderNum := "126378912"
		withdraw := decimal.NewFromFloatWithExponent(659.784, -config.DecimalExponent)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrOrderWrongNum)
	})

	ts.Run("Incorrect withdraw amount", func() {
		orderNum := "3903733214"
		withdraw := decimal.NewFromFloatWithExponent(-1.0, -config.DecimalExponent)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrIncorrectData)
	})

	ts.Run("Low balance", func() {
		orderNum := "3903733214"
		withdraw := decimal.NewFromFloatWithExponent(659.784, -config.DecimalExponent)

		stLoyalty := stMocks.NewLoyaltier(ts.T())
		stLoyalty.EXPECT().AddWithdraw(ctx, userID, orderNum, withdraw).Return(store.ErrLowBalance)
		ts.setStorage(nil, nil, stLoyalty)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrLowBalance)
		ts.storage.AssertExpectations(ts.T())
	})

	ts.Run("Balance not found", func() {
		orderNum := "3903733214"
		withdraw := decimal.NewFromFloatWithExponent(659.784, -config.DecimalExponent)

		stLoyalty := stMocks.NewLoyaltier(ts.T())
		stLoyalty.EXPECT().AddWithdraw(ctx, userID, orderNum, withdraw).Return(store.ErrNotFound)
		ts.setStorage(nil, nil, stLoyalty)

		err := ts.service.LoyaltyAddWithdraw(ctx, userID, orderNum, withdraw)
		ts.Error(err)
		ts.ErrorIs(err, serviceErr.ErrLowBalance)
		ts.storage.AssertExpectations(ts.T())
	})
}

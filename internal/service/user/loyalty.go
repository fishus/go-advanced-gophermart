package user

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	serviceErr "github.com/fishus/go-advanced-gophermart/internal/service/err"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *service) LoyaltyUserBalance(ctx context.Context, userID models.UserID) (balance models.LoyaltyBalance, err error) {
	balance, err = s.storage.Loyalty().BalanceByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			balance.UserID = userID
			err = nil
			return
		}
		return
	}
	return
}

func (s *service) LoyaltyUserWithdrawals(ctx context.Context, userID models.UserID) ([]models.LoyaltyHistory, error) {
	history, err := s.storage.Loyalty().HistoryByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	withdrawals := make([]models.LoyaltyHistory, 0, len(history))
	for _, h := range history {
		if h.Withdrawal.GreaterThan(decimal.NewFromFloat(0)) {
			withdrawals = append(withdrawals, h)
		}
	}

	return withdrawals, nil
}

func (s *service) LoyaltyAddWithdraw(ctx context.Context, userID models.UserID, orderNum string, withdraw decimal.Decimal) error {
	// Проверка номера заказа на корректность с помощью алгоритма Луна
	if err := s.order.ValidateNumLuhn(orderNum); err != nil {
		return serviceErr.ErrOrderWrongNum
	}

	if withdraw.LessThan(decimal.NewFromFloat(0)) {
		return serviceErr.ErrIncorrectData
	}

	err := s.storage.Loyalty().AddWithdraw(ctx, userID, orderNum, withdraw)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return serviceErr.ErrLowBalance

		case errors.Is(err, store.ErrLowBalance):
			return serviceErr.ErrLowBalance
		}
	}
	return err
}

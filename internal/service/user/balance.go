package user

import (
	"context"
	"errors"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

func (s *service) UserLoyaltyBalance(ctx context.Context, userID models.UserID) (balance models.LoyaltyBalance, err error) {
	balance, err = s.storage.LoyaltyBalanceByUser(ctx, userID)
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

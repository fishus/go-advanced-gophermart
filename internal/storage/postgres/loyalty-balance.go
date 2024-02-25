package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/fishus/go-advanced-gophermart/pkg/models"
)

func (s *storage) LoyaltyBalanceUpdate(ctx context.Context, tx pgx.Tx, balance models.LoyaltyBalance) error {
	ctxQuery, cancel := context.WithTimeout(ctx, s.cfg.QueryTimeout)
	defer cancel()

	_, err := tx.Exec(ctxQuery, `
WITH b AS (
	SELECT @accrued::numeric as accrued, @withdrawn::numeric as withdrawn
)
INSERT INTO loyalty_balance SELECT @userID::uuid AS user_id, (accrued - withdrawn) AS current, accrued, withdrawn FROM b
ON CONFLICT (user_id) DO UPDATE SET
accrued = loyalty_balance.accrued + EXCLUDED.accrued,
withdrawn = loyalty_balance.withdrawn + EXCLUDED.withdrawn,
current = ((loyalty_balance.accrued + EXCLUDED.accrued) - (loyalty_balance.withdrawn + EXCLUDED.withdrawn));
`,
		pgx.NamedArgs{
			"userID":    balance.UserID.String(),
			"accrued":   balance.Accrued,
			"withdrawn": balance.Withdrawn,
		})
	return err
}

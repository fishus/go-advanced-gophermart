package service

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	oService "github.com/fishus/go-advanced-gophermart/internal/service/order"
	uService "github.com/fishus/go-advanced-gophermart/internal/service/user"
	store "github.com/fishus/go-advanced-gophermart/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Userer  --with-expecter
type Userer interface {
	Register(context.Context, models.User) (models.UserID, error)
	Login(context.Context, models.User) (models.UserID, error)
	GetByID(context.Context, models.UserID) (models.User, error)
	BuildToken(models.UserID) (string, error)
	DecryptToken(tokenString string) (*uService.JWTClaims, error)
	CheckAuthorizationHeader(auth string) (*uService.JWTClaims, error)
	LoyaltyUserBalance(context.Context, models.UserID) (models.LoyaltyBalance, error)
	LoyaltyAddWithdraw(ctx context.Context, userID models.UserID, orderNum string, withdraw decimal.Decimal) error
	LoyaltyUserWithdrawals(context.Context, models.UserID) ([]models.LoyaltyHistory, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Orderer  --with-expecter
type Orderer interface {
	ValidateNumLuhn(num string) error
	GetByID(context.Context, models.OrderID) (models.Order, error)
	Add(ctx context.Context, userID models.UserID, orderNum string) (models.OrderID, error)
	ListNew(context.Context) ([]models.Order, error)
	ListProcessing(ctx context.Context, limit int) ([]models.Order, error)
	UpdateStatus(context.Context, models.OrderID, models.OrderStatus) error
	AddAccrual(ctx context.Context, id models.OrderID, accrual decimal.Decimal) error
	ListByUser(context.Context, models.UserID) ([]models.Order, error)
}

type service struct {
	cfg     *Config
	storage store.Storager
	user    Userer
	order   Orderer
}

func New(cfg *Config, s store.Storager) *service {
	order := oService.New(s)
	userCfg := &uService.Config{
		JWTExpires:   cfg.JWTExpires,
		JWTSecretKey: cfg.JWTSecretKey,
	}
	user := uService.New(userCfg, s).SetOrder(order)

	return &service{
		storage: s,
		cfg:     cfg,
		user:    user,
		order:   order,
	}
}

func (s *service) Storage() store.Storager {
	return s.storage
}

func (s *service) User() Userer {
	return s.user
}

func (s *service) Order() Orderer {
	return s.order
}

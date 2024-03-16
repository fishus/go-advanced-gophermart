package service

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/fishus/go-advanced-gophermart/pkg/models"

	lService "github.com/fishus/go-advanced-gophermart/internal/service/loyalty"
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
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Orderer  --with-expecter
type Orderer interface {
	ValidateNumLuhn(num string) error
	GetByID(context.Context, models.OrderID) (models.Order, error)
	Add(ctx context.Context, userID models.UserID, orderNum string) (models.OrderID, error)
	ListNew(context.Context) ([]models.Order, error)
	ListProcessing(ctx context.Context, limit int) ([]models.Order, error)
	ListByUser(context.Context, models.UserID) ([]models.Order, error)
	UpdateStatus(context.Context, models.OrderID, models.OrderStatus) error
	AddAccrual(ctx context.Context, id models.OrderID, accrual decimal.Decimal) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Loyaltier  --with-expecter
type Loyaltier interface {
	UserBalance(context.Context, models.UserID) (models.LoyaltyBalance, error)
	AddWithdraw(ctx context.Context, userID models.UserID, orderNum string, withdraw decimal.Decimal) error
	UserWithdrawals(context.Context, models.UserID) ([]models.LoyaltyHistory, error)
}

type service struct {
	cfg     *Config
	storage store.Storager
	user    Userer
	order   Orderer
	loyalty Loyaltier
}

func New(cfg *Config, s store.Storager) *service {
	order := oService.New(s)
	userCfg := &uService.Config{
		JWTExpires:   cfg.JWTExpires,
		JWTSecretKey: cfg.JWTSecretKey,
	}
	user := uService.New(userCfg, s)
	loyalty := lService.New(s).SetOrder(order)

	return &service{
		storage: s,
		cfg:     cfg,
		user:    user,
		order:   order,
		loyalty: loyalty,
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

func (s *service) Loyalty() Loyaltier {
	return s.loyalty
}

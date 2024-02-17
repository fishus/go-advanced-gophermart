package storage

import "github.com/fishus/go-advanced-gophermart/pkg/models"

type OrderFilters struct {
	UserID   models.UserID
	Num      string
	Statuses []models.OrderStatus
}

func (o OrderFilters) IsEmpty() bool {
	isEmpty := true
	if o.UserID != "" {
		isEmpty = false
	}
	if o.Num != "" {
		isEmpty = false
	}
	return isEmpty
}

type OrderFilter func(o *OrderFilters)

func WithOrderUserID(userID models.UserID) OrderFilter {
	return func(f *OrderFilters) {
		f.UserID = userID
	}
}

func WithOrderNum(num string) OrderFilter {
	return func(f *OrderFilters) {
		f.Num = num
	}
}

func WithOrderStatus(status models.OrderStatus) OrderFilter {
	return func(f *OrderFilters) {
		f.Statuses = append(f.Statuses, status)
	}
}

func WithOrderStatuses(statuses ...models.OrderStatus) OrderFilter {
	return func(f *OrderFilters) {
		f.Statuses = statuses
	}
}

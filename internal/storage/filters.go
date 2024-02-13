package storage

import "github.com/fishus/go-advanced-gophermart/pkg/models"

type OrderFilters struct {
	UserID models.UserID
	Num    string
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

func WithOrderNum(num string) OrderFilter {
	return func(f *OrderFilters) {
		f.Num = num
	}
}

func WithOrderUserID(userID models.UserID) OrderFilter {
	return func(f *OrderFilters) {
		f.UserID = userID
	}
}

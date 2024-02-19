package storage

import "github.com/fishus/go-advanced-gophermart/pkg/models"

type OrderByDirection string

func (d OrderByDirection) String() string {
	return string(d)
}

const (
	OrderByAsc  OrderByDirection = "ASC"
	OrderByDesc OrderByDirection = "DESC"
)

type OrderByField string

func (f OrderByField) String() string {
	return string(f)
}

const (
	OrderByUploadedAt OrderByField = "uploaded_at"
)

type OrderFilters struct {
	UserID   models.UserID
	Num      string
	Statuses []models.OrderStatus
	OrderBy  []struct {
		Field OrderByField
		Dir   OrderByDirection
	}
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

func WithOrderBy(field OrderByField, direction OrderByDirection) OrderFilter {
	return func(f *OrderFilters) {
		f.OrderBy = append(f.OrderBy, struct {
			Field OrderByField
			Dir   OrderByDirection
		}{Field: field, Dir: direction})
	}
}

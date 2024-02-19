package accrual

type Config struct {
	APIAddr        string // API service address host:post
	LimitNewOrders int    // Limit buffered channel of new orders
}

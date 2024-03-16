package accrual

import "time"

type Config struct {
	APIAddr        string // API service address host:post
	RequestTimeout time.Duration
	MaxAttempts    int
	WorkersNum     int
}

package postgres

import "time"

type Config struct {
	ConnString     string
	ConnectTimeout time.Duration
	QueryTimeout   time.Duration
}

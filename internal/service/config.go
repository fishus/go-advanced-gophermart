package service

import (
	"time"
)

type Config struct {
	JWTExpires   time.Duration
	JWTSecretKey string
}

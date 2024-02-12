package storage

import "errors"

var ErrAlreadyExists = errors.New("entity already exists")
var ErrNotFound = errors.New("entity not found")

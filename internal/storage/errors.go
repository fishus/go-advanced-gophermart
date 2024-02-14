package storage

import "errors"

var ErrIncorrectData = errors.New("incorrect input data")
var ErrAlreadyExists = errors.New("entity already exists")
var ErrNotFound = errors.New("entity not found")

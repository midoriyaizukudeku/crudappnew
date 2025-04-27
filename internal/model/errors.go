package model

import "errors"

var (
	ErrnoRecord = errors.New("model :NO RECORD FOUND")

	InvalidCredientials = errors.New("model: Invalid credientilas")

	DuplicateEmail = errors.New("models: duplicate email")
)

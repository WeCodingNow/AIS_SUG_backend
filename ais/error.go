package ais

import "errors"

var (
	ErrCathedraNotFound = errors.New("cathedra not found")
	ErrSemesterNotFound = errors.New("semester not found")
)

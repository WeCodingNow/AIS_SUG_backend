package ais

import "errors"

var (
	ErrCathedraNotFound    = errors.New("cathedra not found")
	ErrSemesterNotFound    = errors.New("semester not found")
	ErrContactTypeNotFound = errors.New("contact type not found")
	ErrContactNotFound     = errors.New("contact not found")
)

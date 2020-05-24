package aisauth

import "errors"

var (
	ErrDuplicateUser = errors.New("Student already has a user associated with it")
	ErrNoStudentUser = errors.New("Student has no user associated with it")
)

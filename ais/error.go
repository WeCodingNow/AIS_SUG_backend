package ais

import "errors"

var (
	ErrCathedraNotFound         = errors.New("cathedra not found")
	ErrDisciplineNotFound       = errors.New("discipline not found")
	ErrSemesterNotFound         = errors.New("semester not found")
	ErrContactTypeNotFound      = errors.New("contact type not found")
	ErrContactNotFound          = errors.New("contact not found")
	ErrResidenceNotFound        = errors.New("residence not found")
	ErrControlEventTypeNotFound = errors.New("control event type not found")
	ErrControlEventNotFound     = errors.New("control event not found")
)

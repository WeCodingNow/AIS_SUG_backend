package http

import (
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
)

type Handler struct {
	useCase ais.UseCase
}

func NewHandler(useCase ais.UseCase) *Handler {
	return &Handler{
		useCase,
	}
}

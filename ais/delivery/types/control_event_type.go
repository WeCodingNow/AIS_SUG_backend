package types

import "github.com/WeCodingNow/AIS_SUG_backend/models"

type JSONControlEventType struct {
	ID  int    `json:"id"`
	Def string `json:"def"`
}

func ToJsonControlEventType(controlEventType *models.ControlEventType) *JSONControlEventType {
	return &JSONControlEventType{
		ID:  controlEventType.ID,
		Def: controlEventType.Def,
	}
}

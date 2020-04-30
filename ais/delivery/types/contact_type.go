package types

import "github.com/WeCodingNow/AIS_SUG_backend/models"

type JSONContactType struct {
	ID  int    `json:"id"`
	Def string `json:"def"`
}

func ToJsonContactType(contactType *models.ContactType) *JSONContactType {
	return &JSONContactType{
		ID:  contactType.ID,
		Def: contactType.Def,
	}
}

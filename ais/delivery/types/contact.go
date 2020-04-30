package types

import "github.com/WeCodingNow/AIS_SUG_backend/models"

type JSONContact struct {
	ID          int              `json:"id"`
	Def         string           `json:"def"`
	ContactType *JSONContactType `json:"type"`
}

func toJsonContact(contact *models.Contact) *JSONContact {
	return &JSONContact{
		ID:          contact.ID,
		Def:         contact.Def,
		ContactType: ToJsonContactType(contact.ContactType),
	}
}

type ContactJSONContact struct {
	*JSONContact
	Student *JSONStudent `json:"student"`
}

func ToContactJsonContact(contact *models.Contact) *ContactJSONContact {
	return &ContactJSONContact{
		JSONContact: toJsonContact(contact),
		Student:     toJsonStudent(contact.Student),
	}
}

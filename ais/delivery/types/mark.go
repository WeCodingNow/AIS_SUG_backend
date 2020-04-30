package types

import (
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONMark struct {
	ID    int       `json:"id"`
	Date  time.Time `json:"date"`
	Value int       `json:"value"`
}

func toJsonMark(mark *models.Mark) *JSONMark {
	return &JSONMark{
		ID:    mark.ID,
		Date:  mark.Date,
		Value: mark.Value,
	}
}

type MarkJSONStudent struct {
	*JSONStudent
	Group     *StudentJSONGroup `json:"group"`
	Contacts  []*JSONContact    `json:"contacts"`
	Residence *JSONResidence    `json:"residence"`
}

func toMarkJsonStudent(student *models.Student) *MarkJSONStudent {
	contactJSONs := make([]*JSONContact, 0, len(student.Contacts))
	for _, contact := range student.Contacts {
		contactJSONs = append(contactJSONs, toJsonContact(contact))
	}

	return &MarkJSONStudent{
		JSONStudent: toJsonStudent(student),
		Group:       toStudentJsonGroup(student.Group),
		Contacts:    contactJSONs,
		Residence:   ToJsonResidence(student.Residence),
	}
}

type MarkJSONMark struct {
	*JSONMark
	ControlEvent *JSONControlEvent `json:"control_event"`
	Student      *MarkJSONStudent  `json:"student"`
}

func ToMarkJsonMark(mark *models.Mark) *MarkJSONMark {
	return &MarkJSONMark{
		JSONMark:     toJsonMark(mark),
		ControlEvent: toJsonControlEvent(mark.ControlEvent),
		Student:      toMarkJsonStudent(mark.Student),
	}
}

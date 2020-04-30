package types

import (
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONGroup struct {
	ID     int `json:"id"`
	Number int `json:"number"`
}

func toJsonGroup(group *models.Group) *JSONGroup {
	return &JSONGroup{
		ID:     group.ID,
		Number: group.Number,
	}
}

type GroupJSONStudent struct {
	*JSONStudent
	Marks     []*studentJSONMark `json:"marks"`
	Contacts  []*JSONContact     `json:"contacts"`
	Residence *JSONResidence     `json:"residence"`
}

func toGroupJsonStudent(student *models.Student) *GroupJSONStudent {
	contactJSONs := make([]*JSONContact, 0, len(student.Contacts))
	for _, contact := range student.Contacts {
		contactJSONs = append(contactJSONs, toJsonContact(contact))
	}

	markJSONs := make([]*studentJSONMark, 0, len(student.Marks))
	for _, mark := range student.Marks {
		markJSONs = append(markJSONs, toStudentJSONMark(mark))
	}

	return &GroupJSONStudent{
		JSONStudent: toJsonStudent(student),
		Marks:       markJSONs,
		Contacts:    contactJSONs,
		Residence:   toJsonResidence(student.Residence),
	}
}

type GroupJSONGroup struct {
	*JSONGroup
	Cathedra *JSONCathedra       `json:"cathedra"`
	Students []*GroupJSONStudent `json:"students"`
}

func ToGroupJsonGroup(group *models.Group) *GroupJSONGroup {
	studentJSONs := make([]*GroupJSONStudent, 0, len(group.Students))

	for _, studentModel := range group.Students {
		studentJSONs = append(studentJSONs, toGroupJsonStudent(studentModel))
	}

	return &GroupJSONGroup{
		JSONGroup: toJsonGroup(group),
		Students:  studentJSONs,
		Cathedra:  toJsonCathedra(group.Cathedra),
	}
}

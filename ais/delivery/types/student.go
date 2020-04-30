package types

import (
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONStudent struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	SecondName string  `json:"second_name"`
	ThirdName  *string `json:"third_name"`
}

func toJsonStudent(student *models.Student) *JSONStudent {
	return &JSONStudent{
		ID:         student.ID,
		Name:       student.Name,
		SecondName: student.SecondName,
		ThirdName:  student.ThirdName,
	}
}

type StudentJSONGroup struct {
	*JSONGroup
	Cathedra *JSONCathedra `json:"cathedra"`
}

func toStudentJsonGroup(group *models.Group) *StudentJSONGroup {
	return &StudentJSONGroup{
		JSONGroup: toJsonGroup(group),
		Cathedra:  toJsonCathedra(group.Cathedra),
	}
}

type StudentJSONControlEvent struct {
	*JSONControlEvent
	Discipline *JSONDiscipline `json:"discipline"`
}

func toStudentJSONControlEvent(controlEvent *models.ControlEvent) *StudentJSONControlEvent {
	return &StudentJSONControlEvent{
		JSONControlEvent: toJsonControlEvent(controlEvent),
		Discipline:       ToJsonDiscipline(controlEvent.Discipline),
	}
}

type studentJSONMark struct {
	*JSONMark
	ControlEvent *StudentJSONControlEvent
}

func toStudentJSONMark(mark *models.Mark) *studentJSONMark {
	return &studentJSONMark{
		JSONMark:     toJsonMark(mark),
		ControlEvent: toStudentJSONControlEvent(mark.ControlEvent),
	}
}

type StudentJSONStudent struct {
	*JSONStudent
	Group     *StudentJSONGroup  `json:"group"`
	Marks     []*studentJSONMark `json:"marks"`
	Contacts  []*JSONContact     `json:"contacts"`
	Residence *JSONResidence     `json:"residence"`
}

func ToStudentJsonStudent(student *models.Student) *StudentJSONStudent {
	contactJSONs := make([]*JSONContact, 0, len(student.Contacts))
	for _, contact := range student.Contacts {
		contactJSONs = append(contactJSONs, toJsonContact(contact))
	}

	markJSONs := make([]*studentJSONMark, 0, len(student.Marks))
	for _, mark := range student.Marks {
		markJSONs = append(markJSONs, toStudentJSONMark(mark))
	}

	return &StudentJSONStudent{
		JSONStudent: toJsonStudent(student),
		Group:       toStudentJsonGroup(student.Group),
		Marks:       markJSONs,
		Contacts:    contactJSONs,
		Residence:   toJsonResidence(student.Residence),
	}
}

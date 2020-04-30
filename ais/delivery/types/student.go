package types

import (
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONStudent struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	SecondName string  `json:"second_name"`
	ThirdName  *string `json:"third_name"`
	// Group      *StudentShortGroup     `json:"group"`
	// Residence  *StudentShortResidence `json:"residence"`
	// Contacts   []*StudentShortContact `json:"contacts"`
	// Marks      []*StudentShortMark    `json:"marks"`
}

func ToJsonStudent(student *models.Student) *JSONStudent {
	// contactJSONs := make([]*StudentShortContact, 0, len(student.Contacts))

	// for _, contact := range student.Contacts {
	// 	contactJSONs = append(contactJSONs, toStudentShortContact(contact))
	// }

	return &JSONStudent{
		ID:         student.ID,
		Name:       student.Name,
		SecondName: student.SecondName,
		ThirdName:  student.ThirdName,
		// Group:      toStudentShortGroup(student.Group),
		// Residence:  toStudentShortResidence(student.Residence),
		// Contacts:   contactJSONs,
	}
}

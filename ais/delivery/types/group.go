package types

import (
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type JSONGroup struct {
	ID     int `json:"id"`
	Number int `json:"number"`
	// Cathedra *StudentShortCathedra `json:"cathedra"`
	// Students []*GroupShortStudent  `json:"students"`
}

func ToJsonGroup(group *models.Group) *JSONGroup {
	// studentJSONs := make([]*GroupShortStudent, 0, len(group.Students))

	// for _, studentModel := range group.Students {
	// 	log.Print(*studentModel)
	// 	studentJSONs = append(studentJSONs, toGroupShortStudent(studentModel))
	// }

	return &JSONGroup{
		ID: group.ID,
		// Cathedra: toStudentShortCathedra(group.Cathedra),
		// Students: studentJSONs,
		Number: group.Number,
	}
}

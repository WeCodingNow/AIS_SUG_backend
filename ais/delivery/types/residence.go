package types

import "github.com/WeCodingNow/AIS_SUG_backend/models"

type JSONResidence struct {
	ID        int    `json:"id"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Community bool   `json:"community"`
}

func toJsonResidence(residence *models.Residence) *JSONResidence {
	return &JSONResidence{
		ID:        residence.ID,
		Address:   residence.Address,
		City:      residence.City,
		Community: residence.Community,
	}
}

type ResidenceJSONStudent struct {
	*JSONStudent
	Group    *JSONGroup     `json:"group"`
	Contacts []*JSONContact `json:"contacts"`
}

func toResidenceJsonStudent(student *models.Student) *ResidenceJSONStudent {
	contactJSONs := make([]*JSONContact, 0, len(student.Contacts))

	for _, contact := range student.Contacts {
		contactJSONs = append(contactJSONs, toJsonContact(contact))
	}

	return &ResidenceJSONStudent{
		JSONStudent: toJsonStudent(student),
		Group:       toJsonGroup(student.Group),
		Contacts:    contactJSONs,
	}
}

type ResidenceJSONResidence struct {
	*JSONResidence
	Students []*ResidenceJSONStudent `json:"students"`
}

func ToResidenceJSONResidence(residence *models.Residence) *ResidenceJSONResidence {
	studentJSONs := make([]*ResidenceJSONStudent, 0, len(residence.Students))

	for _, student := range residence.Students {
		studentJSONs = append(studentJSONs, toResidenceJsonStudent(student))
	}

	return &ResidenceJSONResidence{
		JSONResidence: toJsonResidence(residence),
		Students:      studentJSONs,
	}
}

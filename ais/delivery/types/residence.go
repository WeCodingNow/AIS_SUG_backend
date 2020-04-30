package types

import "github.com/WeCodingNow/AIS_SUG_backend/models"

type JSONResidence struct {
	ID        int    `json:"id"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Community bool   `json:"community"`
	// Students  []string `json:"students"`
}

func ToJsonResidence(residence *models.Residence) *JSONResidence {
	return &JSONResidence{
		ID:        residence.ID,
		Address:   residence.Address,
		City:      residence.City,
		Community: residence.Community,
		// Students:  []string{"student1", "student2"},
	}
}

package types

import "github.com/WeCodingNow/AIS_SUG_backend/models"

type JSONCathedra struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	// Groups    []*CathedraShortGroup `json:"groups"`
}

func ToJsonCathedra(cathedra *models.Cathedra) *JSONCathedra {
	// groupJSONs := make([]*CathedraShortGroup, 0, len(cathedra.Groups))

	// for _, groupModel := range cathedra.Groups {
	// 	groupJSONs = append(groupJSONs, toJsonShortGroup(groupModel))
	// }

	return &JSONCathedra{
		ID:        cathedra.ID,
		Name:      cathedra.Name,
		ShortName: cathedra.ShortName,
		// Groups:    groupJSONs,
	}
}

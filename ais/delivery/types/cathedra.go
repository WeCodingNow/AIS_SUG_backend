package types

import "github.com/WeCodingNow/AIS_SUG_backend/models"

type JSONCathedra struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

func toJsonCathedra(cathedra *models.Cathedra) *JSONCathedra {
	return &JSONCathedra{
		ID:        cathedra.ID,
		Name:      cathedra.Name,
		ShortName: cathedra.ShortName,
	}
}

type CathedraJSONGroup struct {
	*JSONGroup
	Students []*GroupJSONStudent `json:"students"`
}

func toCathedraJsonGroup(group *models.Group) *CathedraJSONGroup {
	studentJSONs := make([]*GroupJSONStudent, 0, len(group.Students))

	for _, studentModel := range group.Students {
		studentJSONs = append(studentJSONs, toGroupJsonStudent(studentModel))
	}

	return &CathedraJSONGroup{
		JSONGroup: toJsonGroup(group),
		Students:  studentJSONs,
	}
}

type CathedraJSONCathedra struct {
	*JSONCathedra
	Groups []*CathedraJSONGroup
}

func ToCathedraJSONCathedra(cathedra *models.Cathedra) *CathedraJSONCathedra {
	groupJSONs := make([]*CathedraJSONGroup, 0, len(cathedra.Groups))

	for _, group := range cathedra.Groups {
		groupJSONs = append(groupJSONs, toCathedraJsonGroup(group))
	}

	return &CathedraJSONCathedra{
		JSONCathedra: toJsonCathedra(cathedra),
		Groups:       groupJSONs,
	}
}

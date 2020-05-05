package models

type Cathedra struct {
	ID        int
	Name      string
	ShortName string

	Groups []*Group
}

func ToJSONCathedra(c *Cathedra, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[CathedraT] = true

	retMap := JSONMap{
		"id":         c.ID,
		"name":       c.Name,
		"short_name": c.ShortName,
	}

	if _, ok := refs[GroupT]; !ok {
		groupJSONs := make([]JSONMap, 0, len(c.Groups))
		for _, group := range c.Groups {
			groupJSONs = append(groupJSONs, ToJSONGroup(group, refs))
		}

		retMap["groups"] = groupJSONs
	}

	return retMap
}

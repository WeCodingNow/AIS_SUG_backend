package models

type Cathedra struct {
	ID        int
	Name      string
	ShortName string

	Groups []*Group
}

func ToJSONCathedra(c *Cathedra, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":         c.ID,
		"name":       c.Name,
		"short_name": c.ShortName,
	}

	if filled, ok := refs[GroupT]; !(ok && filled) {
		groupJSONs := make([]JSONMap, 0, len(c.Groups))
		for _, group := range c.Groups {
			groupJSONs = append(groupJSONs, ToJSONGroup(group, withDontWant(refs, StudentT)))
		}
		retMap["groups"] = groupJSONs
	}

	return retMap
}

package models

type Role struct {
	ID  int
	Def string
}

func ToJSONRole(r *Role, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":  r.ID,
		"def": r.Def,
	}

	return retMap
}

package models

type ContactType struct {
	ID  int
	Def string
}

func ToJSONContactType(cet *ContactType, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":  cet.ID,
		"def": cet.Def,
	}

	return retMap
}

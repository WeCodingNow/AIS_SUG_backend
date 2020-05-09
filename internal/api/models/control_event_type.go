package models

type ControlEventType struct {
	ID  int
	Def string
}

func ToJSONControlEventType(cet *ControlEventType, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":  cet.ID,
		"def": cet.Def,
	}

	return retMap
}

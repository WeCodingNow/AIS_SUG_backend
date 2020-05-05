package models

type ControlEventType struct {
	ID  int
	Def string
}

func ToJSONControlEventType(cet *ControlEventType, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[ControlEventT] = true

	retMap := JSONMap{
		"id":  cet.ID,
		"def": cet.Def,
	}

	return retMap
}

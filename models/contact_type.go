package models

type ContactType struct {
	ID  int
	Def string
}

func ToJSONContactType(cet *ContactType, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[ContactTypeT] = true

	retMap := JSONMap{
		"id":  cet.ID,
		"def": cet.Def,
	}

	return retMap
}

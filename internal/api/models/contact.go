package models

type Contact struct {
	ID  int
	Def string

	*ContactType
	*Student
}

func ToJSONContact(c *Contact, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[ContactT] = true

	retMap := JSONMap{
		"id":  c.ID,
		"def": c.Def,
	}

	if _, ok := refs[ContactTypeT]; !ok {
		retMap["type"] = ToJSONContactType(c.ContactType, refs)
	}

	if _, ok := refs[StudentT]; !ok {
		retMap["student"] = ToJSONStudent(c.Student, refs)
	}

	return retMap
}

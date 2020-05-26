package models

type Contact struct {
	ID  int
	Def string

	*ContactType
	*Student
}

func ToJSONContact(c *Contact, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":  c.ID,
		"def": c.Def,
	}

	if filled, ok := refs[ContactTypeT]; !(ok && filled) {
		retMap["type"] = ToJSONContactType(c.ContactType, withDontWant(refs, ContactT))
	}

	if filled, ok := refs[StudentT]; !(ok && filled) {
		retMap["student"] = ToJSONStudent(c.Student, withDontWant(refs, ContactT, GroupT, MarkT))
	}

	return retMap
}

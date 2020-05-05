package models

type Student struct {
	ID         int
	Name       string
	SecondName string
	ThirdName  *string

	*Group
	*Residence
	Contacts []*Contact
	Marks    []*Mark
}

func ToJSONStudent(s *Student, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[StudentT] = true

	retMap := JSONMap{
		"id":          s.ID,
		"name":        s.Name,
		"second_name": s.SecondName,
		"third_name":  s.ThirdName,
	}

	if _, ok := refs[GroupT]; !ok {
		retMap["group"] = ToJSONGroup(s.Group, refs)
	}

	if _, ok := refs[ResidenceT]; !ok {
		retMap["residence"] = ToJSONResidence(s.Residence, refs)
	}

	if _, ok := refs[ContactT]; !ok {
		contactJSONs := make([]JSONMap, 0, len(s.Contacts))
		for _, contact := range s.Contacts {
			contactJSONs = append(contactJSONs, ToJSONContact(contact, refs))
		}
		retMap["contacts"] = contactJSONs
	}

	if _, ok := refs[MarkT]; !ok {
		markJSONs := make([]JSONMap, 0, len(s.Marks))
		for _, mark := range s.Marks {
			markJSONs = append(markJSONs, ToJSONMark(mark, refs))
		}
		retMap["marks"] = markJSONs
	}

	return retMap
}

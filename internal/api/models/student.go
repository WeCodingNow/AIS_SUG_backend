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
	Backlogs []*Backlog
}

func ToJSONStudent(s *Student, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":          s.ID,
		"name":        s.Name,
		"second_name": s.SecondName,
		"third_name":  s.ThirdName,
	}

	if filled, ok := refs[GroupT]; !(ok && filled) {
		retMap["group"] = ToJSONGroup(s.Group, withDontWant(refs, StudentT))
	}

	if filled, ok := refs[ResidenceT]; !(ok && filled) {
		retMap["residence"] = ToJSONResidence(s.Residence, withDontWant(refs, StudentT))
	}

	if filled, ok := refs[ContactT]; !(ok && filled) {
		contactJSONs := make([]JSONMap, 0, len(s.Contacts))
		for _, contact := range s.Contacts {
			contactJSONs = append(contactJSONs, ToJSONContact(contact, withDontWant(refs, StudentT)))
		}
		retMap["contacts"] = contactJSONs
	}

	if filled, ok := refs[MarkT]; !(ok && filled) {
		markJSONs := make([]JSONMap, 0, len(s.Marks))
		for _, mark := range s.Marks {
			markJSONs = append(markJSONs, ToJSONMark(mark, withDontWant(refs, StudentT)))
		}
		retMap["marks"] = markJSONs
	}

	if filled, ok := refs[BacklogT]; !(ok && filled) {
		backlogJSONs := make([]JSONMap, 0, len(s.Backlogs))
		for _, backlog := range s.Backlogs {
			backlogJSONs = append(backlogJSONs, ToJsonBacklog(backlog, withDontWant(refs, StudentT, DisciplineT)))
		}
		retMap["backlogs"] = backlogJSONs
	}

	return retMap
}

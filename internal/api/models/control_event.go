package models

import "time"

type ControlEvent struct {
	ID   int
	Date time.Time

	ControlEventType *ControlEventType
	Discipline       *Discipline
	Marks            []*Mark
	Semester         *Semester
}

func ToJSONControlEvent(ce *ControlEvent, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[ControlEventT] = true

	retMap := JSONMap{
		"id":   ce.ID,
		"date": ce.Date,
	}

	if _, ok := refs[ControlEventTypeT]; !ok {
		retMap["type"] = ToJSONControlEventType(ce.ControlEventType, refs)
	}

	if _, ok := refs[DisciplineT]; !ok {
		retMap["discipline"] = ToJSONDiscipline(ce.Discipline, refs)
	}

	if _, ok := refs[MarkT]; !ok {
		markJSONs := make([]JSONMap, 0, len(ce.Marks))
		for _, mark := range ce.Marks {
			markJSONs = append(markJSONs, ToJSONMark(mark, refs))
		}

		retMap["marks"] = markJSONs
	}

	if _, ok := refs[SemesterT]; !ok {
		retMap["semester"] = ToJSONSemester(ce.Semester, refs)
	}

	return retMap
}

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
	retMap := JSONMap{
		"id":   ce.ID,
		"date": ce.Date,
	}

	if filled, ok := refs[ControlEventTypeT]; !(ok && filled) {
		retMap["type"] = ToJSONControlEventType(ce.ControlEventType, withDontWant(refs, ControlEventT))
	}

	if filled, ok := refs[DisciplineT]; !(ok && filled) {
		retMap["discipline"] = ToJSONDiscipline(ce.Discipline, withDontWant(refs, ControlEventT))
	}

	if filled, ok := refs[MarkT]; !(ok && filled) {
		markJSONs := make([]JSONMap, 0, len(ce.Marks))
		for _, mark := range ce.Marks {
			markJSONs = append(markJSONs, ToJSONMark(mark, withDontWant(refs, ControlEventT, GroupT)))
		}
		retMap["marks"] = markJSONs
	}

	if filled, ok := refs[SemesterT]; !(ok && filled) {
		retMap["semester"] = ToJSONSemester(ce.Semester, withDontWant(refs, ControlEventT, GroupT))
	}

	return retMap
}

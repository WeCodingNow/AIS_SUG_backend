package models

import "time"

type Mark struct {
	ID    int
	Date  time.Time
	Value int

	*ControlEvent
	*Student
}

func ToJSONMark(m *Mark, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[MarkT] = true

	retMap := JSONMap{
		"id":    m.ID,
		"date":  m.Date,
		"value": m.Value,
	}

	if _, ok := refs[StudentT]; !ok {
		retMap["student"] = ToJSONStudent(m.Student, refs)
	}

	if _, ok := refs[ControlEventT]; !ok {
		retMap["control_event"] = ToJSONControlEvent(m.ControlEvent, refs)
	}

	return retMap
}

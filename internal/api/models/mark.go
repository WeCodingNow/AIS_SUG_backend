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
	retMap := JSONMap{
		"id":    m.ID,
		"date":  m.Date,
		"value": m.Value,
	}

	if filled, ok := refs[StudentT]; !(ok && filled) {
		retMap["student"] = ToJSONStudent(m.Student, withDontWant(refs, MarkT))
	}

	if filled, ok := refs[ControlEventT]; !(ok && filled) {
		retMap["control_event"] = ToJSONControlEvent(m.ControlEvent, withDontWant(refs, MarkT))
	}

	return retMap
}

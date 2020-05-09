package models

import "time"

type Semester struct {
	ID        int        `json:"id"`
	Number    int        `json:"number"`
	Beginning time.Time  `json:"beginning"`
	End       *time.Time `json:"end"`

	Groups        []*Group
	ControlEvents []*ControlEvent
}

func ToJSONSemester(s *Semester, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":        s.ID,
		"number":    s.Number,
		"beginning": s.Beginning,
		"end":       s.End,
	}

	if filled, ok := refs[GroupT]; !(ok && filled) {
		groupJSONs := make([]JSONMap, 0, len(s.Groups))
		for _, group := range s.Groups {
			groupJSONs = append(groupJSONs, ToJSONGroup(group, withDontWant(refs, SemesterT, MarkT)))
		}
		retMap["groups"] = groupJSONs
	}

	if filled, ok := refs[ControlEventT]; !(ok && filled) {
		controlEventJSONs := make([]JSONMap, 0, len(s.ControlEvents))
		for _, controlEvent := range s.ControlEvents {
			controlEventJSONs = append(controlEventJSONs, ToJSONControlEvent(controlEvent, withDontWant(refs, SemesterT, GroupT, ResidenceT, ContactT)))
		}
		retMap["control_events"] = controlEventJSONs
	}

	return retMap
}

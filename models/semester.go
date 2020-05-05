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
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[SemesterT] = true

	retMap := JSONMap{
		"id":        s.ID,
		"number":    s.Number,
		"beginning": s.Beginning,
		"end":       s.End,
	}

	if _, ok := refs[GroupT]; !ok {
		groupJSONs := make([]JSONMap, 0, len(s.Groups))
		for _, group := range s.Groups {
			groupJSONs = append(groupJSONs, ToJSONGroup(group, refs))
		}

		retMap["groups"] = groupJSONs
	}

	return retMap
}

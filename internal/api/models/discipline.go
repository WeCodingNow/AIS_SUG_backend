package models

type Discipline struct {
	ID    int
	Name  string
	Hours int

	ControlEvents []*ControlEvent
	Semesters     []*Semester
	Backlogs      []*Backlog
}

func ToJSONDiscipline(d *Discipline, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":    d.ID,
		"name":  d.Name,
		"hours": d.Hours,
	}

	if filled, ok := refs[ControlEventT]; !(ok && filled) {
		controlEventJSONs := make([]JSONMap, 0, len(d.ControlEvents))
		for _, controlEvent := range d.ControlEvents {
			controlEventJSONs = append(controlEventJSONs, ToJSONControlEvent(controlEvent, withDontWant(refs, DisciplineT)))
		}
		retMap["control_events"] = controlEventJSONs
	}

	if filled, ok := refs[SemesterT]; !(ok && filled) {
		semesterJSONs := make([]JSONMap, 0, len(d.Semesters))
		for _, semester := range d.Semesters {
			semesterJSONs = append(semesterJSONs, ToJSONSemester(semester, withDontWant(refs, DisciplineT, GroupT, ControlEventT)))
		}
		retMap["semesters"] = semesterJSONs
	}

	if filled, ok := refs[BacklogT]; !(ok && filled) {
		backlogJSONs := make([]JSONMap, 0, len(d.Backlogs))
		for _, backlog := range d.Backlogs {
			backlogJSONs = append(backlogJSONs, ToJsonBacklog(backlog, withDontWant(refs, DisciplineT, StudentT)))
		}
		retMap["backlogs"] = backlogJSONs
	}

	return retMap
}

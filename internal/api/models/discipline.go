package models

type Discipline struct {
	ID    int
	Name  string
	Hours int

	ControlEvents []*ControlEvent
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
			// controlEventJSONs = append(controlEventJSONs, ToJSONControlEvent(controlEvent, refs))
			controlEventJSONs = append(controlEventJSONs, ToJSONControlEvent(controlEvent, withDontWant(refs, DisciplineT)))
		}
		retMap["control_events"] = controlEventJSONs
	}

	return retMap
}

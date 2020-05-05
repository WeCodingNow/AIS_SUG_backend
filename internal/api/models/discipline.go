package models

type Discipline struct {
	ID    int
	Name  string
	Hours int

	ControlEvents []*ControlEvent
}

func ToJSONDiscipline(d *Discipline, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[DisciplineT] = true

	retMap := JSONMap{
		"id":    d.ID,
		"name":  d.Name,
		"hours": d.Hours,
	}

	if _, ok := refs[ControlEventT]; !ok {
		controlEventJSONs := make([]JSONMap, 0, len(d.ControlEvents))
		for _, controlEvent := range d.ControlEvents {
			controlEventJSONs = append(controlEventJSONs, ToJSONControlEvent(controlEvent, refs))
		}
		retMap["control_events"] = controlEventJSONs
	}

	return retMap
}

package models

type Residence struct {
	ID        int
	Address   string
	City      string
	Community bool

	Students []*Student
}

func ToJSONResidence(r *Residence, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":        r.ID,
		"address":   r.Address,
		"city":      r.City,
		"community": r.Community,
	}

	if filled, ok := refs[StudentT]; !(ok && filled) {
		studentJSONs := make([]JSONMap, 0, len(r.Students))
		for _, student := range r.Students {
			// studentJSONs = append(studentJSONs, ToJSONStudent(student, refs))
			studentJSONs = append(studentJSONs, ToJSONStudent(student, withDontWant(refs, ResidenceT)))
		}

		retMap["students"] = studentJSONs
	}

	return retMap
}

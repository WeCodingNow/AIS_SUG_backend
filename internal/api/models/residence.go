package models

type Residence struct {
	ID        int
	Address   string
	City      string
	Community bool

	Students []*Student
}

func ToJSONResidence(r *Residence, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[ResidenceT] = true

	retMap := JSONMap{
		"id":        r.ID,
		"address":   r.Address,
		"city":      r.City,
		"community": r.Community,
	}

	if _, ok := refs[StudentT]; !ok {
		studentJSONs := make([]JSONMap, 0, len(r.Students))
		for _, student := range r.Students {
			studentJSONs = append(studentJSONs, ToJSONStudent(student, refs))
		}

		retMap["students"] = studentJSONs
	}

	return retMap
}

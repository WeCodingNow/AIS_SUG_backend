package models

type Group struct {
	ID     int
	Number int

	Cathedra  *Cathedra
	Students  []*Student
	Semesters []*Semester
}

func ToJSONGroup(g *Group, refs JSONRefTable) JSONMap {
	if refs == nil {
		refs = make(JSONRefTable)
	}

	refs[GroupT] = true

	retMap := JSONMap{
		"id":     g.ID,
		"number": g.Number,
	}

	if _, ok := refs[CathedraT]; !ok {
		retMap["cathedra"] = ToJSONCathedra(g.Cathedra, refs)
	}

	refs[StudentT] = true
	if filled, ok := refs[SemesterT]; !ok || !filled {
		semesterJSONs := make([]JSONMap, 0, len(g.Semesters))
		for _, semester := range g.Semesters {
			semesterJSONs = append(semesterJSONs, ToJSONSemester(semester, refs))
		}

		retMap["semesters"] = semesterJSONs
	}
	//TODO: переделать это
	delete(refs, StudentT)

	if _, ok := refs[StudentT]; !ok {
		studentJSONs := make([]JSONMap, 0, len(g.Students))
		for _, student := range g.Students {
			studentJSONs = append(studentJSONs, ToJSONStudent(student, refs))
		}

		retMap["students"] = studentJSONs
	}

	return retMap
}

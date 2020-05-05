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

	if _, ok := refs[StudentT]; !ok {
		studentJSONs := make([]JSONMap, 0, len(g.Students))
		for _, student := range g.Students {
			studentJSONs = append(studentJSONs, ToJSONStudent(student, refs))
		}

		retMap["students"] = studentJSONs
	}

	if _, ok := refs[SemesterT]; !ok {
		semesterJSONs := make([]JSONMap, 0, len(g.Semesters))
		for _, semester := range g.Semesters {
			semesterJSONs = append(semesterJSONs, ToJSONSemester(semester, refs))
		}

		retMap["semesters"] = semesterJSONs
	}

	return retMap
}

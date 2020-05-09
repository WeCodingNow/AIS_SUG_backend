package models

type Group struct {
	ID     int
	Number int

	Cathedra  *Cathedra
	Students  []*Student
	Semesters []*Semester
}

func ToJSONGroup(g *Group, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":     g.ID,
		"number": g.Number,
	}

	if filled, ok := refs[CathedraT]; !(ok && filled) {
		retMap["cathedra"] = ToJSONCathedra(g.Cathedra, withDontWant(refs, GroupT))
	}

	if filled, ok := refs[SemesterT]; !(ok && filled) {
		semesterJSONs := make([]JSONMap, 0, len(g.Semesters))
		for _, semester := range g.Semesters {
			semesterJSONs = append(semesterJSONs, ToJSONSemester(semester, withDontWant(refs, GroupT)))
		}
		retMap["semesters"] = semesterJSONs
	}

	if filled, ok := refs[StudentT]; !(ok && filled) {
		studentJSONs := make([]JSONMap, 0, len(g.Students))
		for _, student := range g.Students {
			studentJSONs = append(studentJSONs, ToJSONStudent(student, withDontWant(refs, GroupT)))
		}
		retMap["students"] = studentJSONs
	}

	return retMap
}

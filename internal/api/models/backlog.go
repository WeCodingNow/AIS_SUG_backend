package models

// CREATE TABLE Задолженность(
//     id SERIAL,
//     описание varchar(200),
//     ликвидирована boolean NOT NULL DEFAULT false,
//     id_дисциплины int NOT NULL references Дисциплина(id) ON DELETE CASCADE,
// 	   id_студента int NOT NULL references Студент(id) ON DELETE CASCADE,
//     CONSTRAINT задолженность_pk PRIMARY KEY (id)
// );

type Backlog struct {
	ID          int
	Description string
	Done        bool

	Discipline *Discipline
	Student    *Student
}

func ToJsonBacklog(b *Backlog, refs JSONRefTable) JSONMap {
	retMap := JSONMap{
		"id":   b.ID,
		"desc": b.Description,
		"done": b.Done,
	}

	if filled, ok := refs[DisciplineT]; !(ok && filled) {
		retMap["discipline"] = ToJSONDiscipline(b.Discipline, withDontWant(refs, BacklogT, ControlEventT, SemesterT))
	}

	if filled, ok := refs[StudentT]; !(ok && filled) {
		retMap["student"] = ToJSONStudent(b.Student, withDontWant(refs, BacklogT, ControlEventT, GroupT, ResidenceT, MarkT))
	}

	return retMap
}

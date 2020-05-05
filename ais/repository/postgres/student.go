package postgres

import (
	"context"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

// CREATE TABLE Студент(
//     id SERIAL,
//     id_группы int NOT NULL references Группа(id) ON DELETE CASCADE,
//     имя varchar(50) NOT NULL,
//     фамилия varchar(50) NOT NULL,
//     отчество varchar(50),
//     CONSTRAINT студент_pk PRIMARY KEY (id)
// );

type repoStudent struct {
	ID         int
	Name       string
	SecondName string
	ThirdName  *string

	Group     *repoGroup
	Residence *repoResidence
	Marks     map[int]*repoMark
	Contacts  map[int]*repoContact

	model *models.Student
}

func NewRepoStudent() *repoStudent {
	return &repoStudent{
		Marks:    make(map[int]*repoMark),
		Contacts: make(map[int]*repoContact),
	}
}

func (s *repoStudent) Fill(scannable Scannable) {
	scannable.Scan(&s.ID, &s.Name, &s.SecondName, &s.ThirdName)
}

func (s repoStudent) GetID() int {
	return s.ID
}

const studentTable = "Студент"
const studentFields = "id,имя,фамилия,отчество"
const studentGroupFK = "id_группы"
const studentResidenceFK = "id_места_жительства"

func (c repoStudent) GetDescription() ModelDescription {
	return ModelDescription{
		Table:  studentTable,
		Fields: studentFields,
		Dependencies: []ModelDependency{
			{
				DependencyType:  ManyToOne,
				ForeignKeyField: studentGroupFK,
				ModelMaker:      func() RepoModel { return NewRepoGroup() },
			},
			{
				DependencyType:  ManyToOne,
				ForeignKeyField: studentResidenceFK,
				ModelMaker:      func() RepoModel { return NewRepoResidence() },
			},
			{
				DependencyType:     OneToMany,
				DepForeignKeyField: markStudentFK,
				ModelMaker:         func() RepoModel { return NewRepoMark() },
			},
			{
				DependencyType:     OneToMany,
				DepForeignKeyField: contactStudentFK,
				ModelMaker:         func() RepoModel { return NewRepoContact() },
			},
		},
	}
}

func (c *repoStudent) toModel() *models.Student {
	if c.model == nil {
		c.model = &models.Student{
			ID:         c.ID,
			Name:       c.Name,
			SecondName: c.SecondName,
		}

		if c.ThirdName != nil {
			c.model.ThirdName = c.ThirdName
		}

		c.model.Group = c.Group.toModel()
		c.model.Residence = c.Residence.toModel()

		marks := make([]*models.Mark, 0, len(c.Marks))
		for _, repoM := range c.Marks {
			marks = append(marks, repoM.toModel())
		}
		c.model.Marks = marks

		contacts := make([]*models.Contact, 0, len(c.Contacts))
		for _, contactM := range c.Contacts {
			contacts = append(contacts, contactM.toModel())
		}
		c.model.Contacts = contacts
	}

	return c.model
}

func (s *repoStudent) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoGroup:
		s.Group = dep
	case *repoResidence:
		s.Residence = dep
	case *repoMark:
		s.Marks[dep.ID] = dep
	case *repoContact:
		s.Contacts[dep.ID] = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetStudent(ctx context.Context, id int) (*models.Student, error) {
	student := NewRepoStudent()
	filler, err := MakeFiller(ctx, r.db, studentFields, studentTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrStudentNotFound
	}

	err = filler.Fill(student)

	return student.toModel(), nil
}

func (r *DBAisRepository) GetAllStudents(ctx context.Context) ([]*models.Student, error) {
	students := make([]*models.Student, 0)
	filler, err := MakeFiller(ctx, r.db, studentFields, studentTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoStudent := NewRepoStudent()
		err = filler.Fill(newRepoStudent)
		if err != nil {
			return nil, err
		}
		students = append(students, newRepoStudent.toModel())
	}

	return students, nil
}

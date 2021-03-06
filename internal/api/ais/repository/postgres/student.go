package postgres

import (
	"context"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
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
	Backlogs  map[int]*repoBacklog

	model *models.Student
}

func NewRepoStudent() *repoStudent {
	return &repoStudent{
		Marks:    make(map[int]*repoMark),
		Contacts: make(map[int]*repoContact),
		Backlogs: make(map[int]*repoBacklog),
	}
}

func (s *repoStudent) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&s.ID, &s.Name, &s.SecondName, &s.ThirdName)
}

func (s repoStudent) GetID() int {
	return s.ID
}

const studentTable = "Студент"
const studentFields = "id,имя,фамилия,отчество"
const studentGroupFK = "id_группы"
const studentResidenceFK = "id_места_жительства"

func (c repoStudent) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  studentTable,
		Fields: studentFields,
		Dependencies: []pgorm.ModelDependency{
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: studentGroupFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoGroup() },
			},
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: studentResidenceFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoResidence() },
			},
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: markStudentFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoMark() },
			},
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: backlogStudentFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoBacklog() },
			},
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: contactStudentFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoContact() },
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

		backlogs := make([]*models.Backlog, 0, len(c.Backlogs))
		for _, repoB := range c.Backlogs {
			backlogs = append(backlogs, repoB.toModel())
		}
		c.model.Backlogs = backlogs

		contacts := make([]*models.Contact, 0, len(c.Contacts))
		for _, repoC := range c.Contacts {
			contacts = append(contacts, repoC.toModel())
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
	case *repoBacklog:
		s.Backlogs[dep.ID] = dep
	case *repoContact:
		s.Contacts[dep.ID] = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetStudent(ctx context.Context, id int) (*models.Student, error) {
	student := NewRepoStudent()
	filler, err := pgorm.MakeFiller(ctx, r.db, studentFields, studentTable, &id)

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
	filler, err := pgorm.MakeFiller(ctx, r.db, studentFields, studentTable, nil)

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

// const insertQuery

func (r *DBAisRepository) CreateStudent(ctx context.Context, name, secondName string, thirdName *string, groupID, residenceID int) (*models.Student, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO Студент(имя,фамилия,отчество,id_группы, id_места_жительства) VALUES ( $1, $2, $3, $4, $5) RETURNING id`,
		name, secondName, thirdName, groupID, residenceID,
	)

	var newID int
	err := row.Scan(&newID)

	if err != nil {
		return nil, err
	}

	return r.GetStudent(ctx, newID)
}

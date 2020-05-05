package pgorm

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Scannable interface {
	Scan(...interface{}) error
}

type DependencyType int

const (
	OneToOne DependencyType = iota + 1
	OneToMany
	ManyToOne
	ManyToMany
)

type ModelDependency struct {
	DependencyType     DependencyType
	ModelMaker         func() RepoModel
	ForeignKeyField    string
	DepForeignKeyField string
	ManyToManyTable    string
}

type assignedModelDependency struct {
	ModelDependency
	CurModel RepoModel
	CurID    int
}

type ModelDescription struct {
	Table        string
	Fields       string
	Dependencies []ModelDependency
}

type RepoModel interface {
	GetDescription() ModelDescription
	// GetTable() string
	// GetFields() string
	// GetDeps() []ModelDependency
	GetID() int
	Fill(scannable Scannable)
	AcceptDep(dep interface{}) error
}

func (md ModelDependency) Assign(rm RepoModel) assignedModelDependency {
	return assignedModelDependency{
		ModelDependency: md,
		CurModel:        rm,
		CurID:           rm.GetID(),
	}
	// md.CurModel = rm
}

func prefixFields(table, fieldsString string) string {
	fields := strings.Split(fieldsString, ",")

	for i := range fields {
		fields[i] = fmt.Sprintf("%s.%s", table, fields[i])
	}

	return strings.Join(fields, ",")
}

type DepTable = map[string]map[int]RepoModel

type RepoModelFiller struct {
	db        *sql.DB
	rows      *sql.Rows
	ctx       context.Context
	depsTable DepTable
}

// type Filter = map[string]interface{}
type FilteredFields = []string

func makeFilterString(f FilteredFields) string {
	str := "WHERE "

	for i, v := range f {
		str += fmt.Sprintf("%s = $%d", v, i+1)
		if i < len(f)-1 {
			str += " AND "
		}
	}

	return str
}

func MakeFillerWithFilter(
	ctx context.Context, db *sql.DB,
	fields, table string, filter FilteredFields,
	values ...interface{},
) (*RepoModelFiller, error) {
	query := fmt.Sprintf("SELECT %s FROM %s %s", fields, table, makeFilterString(filter))
	log.Print(query)

	rows, err := db.QueryContext(ctx, query, values...)

	if err != nil {
		return nil, err
	}

	return &RepoModelFiller{
		db:        db,
		rows:      rows,
		ctx:       ctx,
		depsTable: make(DepTable),
	}, nil
}

func MakeFiller(ctx context.Context, db *sql.DB, fields, table string, id *int) (*RepoModelFiller, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", fields, table)

	var rows *sql.Rows
	var err error
	if id != nil {
		rows, err = db.QueryContext(ctx, fmt.Sprintf("%s WHERE id = $1", query), *id)
	} else {
		rows, err = db.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}

	return &RepoModelFiller{
		db:        db,
		rows:      rows,
		ctx:       ctx,
		depsTable: make(DepTable),
	}, nil
}

func (rmf *RepoModelFiller) Next() bool {
	return rmf.rows.Next()
}

func (rmf *RepoModelFiller) Fill(rm RepoModel) error {
	desc := rm.GetDescription()

	rm.Fill(rmf.rows)
	return rmf.resolveDeps(rm, desc.Table, rm.GetID())
}

func (rmf *RepoModelFiller) resolveDeps(rm RepoModel, table string, id int) error {
	if _, ok := rmf.depsTable[table]; !ok {
		rmf.depsTable[table] = make(map[int]RepoModel)
	}
	rmf.depsTable[table][id] = rm

	desc := rm.GetDescription()

	deps := desc.Dependencies
	depsQueue := make([]assignedModelDependency, 0, len(deps))
	for _, d := range deps {
		depsQueue = append(depsQueue, d.Assign(rm))
	}

	for len(depsQueue) > 0 {
		// log.Print("resolving a dep")

		dep := depsQueue[0]
		depsQueue = depsQueue[1:]

		desc := dep.CurModel.GetDescription()
		curModel := dep.CurModel
		table := desc.Table
		curID := curModel.GetID()

		depDesc := dep.ModelMaker().GetDescription()
		depTable, depFields := depDesc.Table, depDesc.Fields

		var joinQuery string

		switch dep.DependencyType {
		case OneToMany:
			// log.Printf("one to many dep %s.[]%s", table, depTable)
			depFK := dep.DepForeignKeyField
			joinQuery = fmt.Sprintf(
				"SELECT %s FROM %s JOIN %s ON %s.%s = %s.id WHERE %s.id = $1",
				prefixFields(depTable, depFields), depTable, table, depTable, depFK, table, table,
			)

		case ManyToOne:
			// log.Printf("many to one dep [%s].%s", table, depTable)
			depFK := dep.ForeignKeyField
			joinQuery = fmt.Sprintf(
				"SELECT %s FROM %s JOIN %s ON %s.%s = %s.id WHERE %s.id = $1",
				prefixFields(depTable, depFields), depTable, table, table, depFK, depTable, table,
			)

		case ManyToMany:
			// log.Println("many to many dep [%s].[%s]", table, depTable)

			leftKey := dep.ForeignKeyField
			rightKey := dep.DepForeignKeyField
			mtm := dep.ManyToManyTable

			joinQuery = fmt.Sprintf("SELECT %s FROM %s JOIN %s ON %s.%s = %s.id JOIN %s ON %s.%s = %s.id WHERE %s.id = $1",
				prefixFields(depTable, depFields), mtm, table, mtm, leftKey, table, depTable, mtm, rightKey, depTable, table,
			)
		default:
			panic("ne znau takoy dependency")
		}

		// log.Println(joinQuery)
		depRows, err := rmf.db.QueryContext(rmf.ctx, joinQuery, curID)

		if err != nil {
			log.Print(joinQuery)
			return err
		}

		for depRows.Next() {
			// log.Print("made a dependency model")
			newModel := dep.ModelMaker()
			newModel.Fill(depRows)
			if _, ok := rmf.depsTable[depTable]; !ok {
				rmf.depsTable[depTable] = make(map[int]RepoModel)
			}

			depID := newModel.GetID()

			if foundDep, ok := rmf.depsTable[depTable][depID]; !ok {
				// log.Print("saving new model")

				rmf.depsTable[depTable][depID] = newModel
				curModel.AcceptDep(newModel)
				newModel.AcceptDep(curModel)

				for _, newDep := range newModel.GetDescription().Dependencies {
					// newDep.Assign(newModel)
					depsQueue = append(depsQueue, newDep.Assign(newModel))
				}
			} else {
				// log.Println("dep ", foundDep, " - ", depTable, " was worked already!")
				// log.Println("discarding new model")

				curModel.AcceptDep(foundDep)
				foundDep.AcceptDep(curModel)
			}
		}
	}

	// log.Print("deps table in the end")
	// log.Print(rmf.depsTable)

	return nil
}

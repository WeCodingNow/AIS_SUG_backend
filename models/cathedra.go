package models

type Cathedra struct {
	ID        int
	Name      string
	ShortName string

	Groups []*Group
}

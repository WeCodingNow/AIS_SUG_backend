package postgres

type Residence struct {
	ID        int
	address   string
	city      string
	community bool
	Students  []*Student
}

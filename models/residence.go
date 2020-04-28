package models

type Residence struct {
	ID        int
	Address   string
	City      string
	Community bool
	Students  []*Student
}

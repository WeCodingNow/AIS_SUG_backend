package postgres

type Scannable interface {
	Scan(...interface{}) error
}

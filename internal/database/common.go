package database

type Scannable interface {
	Scan(dest ...interface{}) error
}

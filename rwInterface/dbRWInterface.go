package rwInterface

type DB interface {
	Open()
	Close()
	Query()
	QueryRow()
	Exec()
	Begin()
	Commit()
	Rollback()
	master()
	slave()
}

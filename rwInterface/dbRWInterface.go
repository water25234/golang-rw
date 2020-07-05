package rwinterface

import "database/sql"

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

type DBConfig struct {
	HOST     string
	POST     int
	DATABASE string
	USER     string
	PASSWORD string
}

type DBExecute struct {
	WRITE *sql.DB
	READ  *sql.DB
}

type DBExecuteMulti struct {
	WRITE *sql.DB
	READ  []*sql.DB
}

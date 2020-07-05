package rw

import (
	"database/sql"
	"fmt"
)

var tx *sql.Tx
var err error
var isTransaction bool = false
var IsCloseTransactionQueryToWrite bool = false

type DriverDBSource struct {
	WRITE DBConfig
	READ  DBConfig
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

type DriverConfig struct {
	DriverName     string
	DriverDBSource DriverDBSource
}

func Open(driverName string, driverDBsource DriverDBSource) (*DBExecute, error) {

	driverConfig := &DriverConfig{
		DriverName:     driverName,
		DriverDBSource: driverDBsource,
	}

	dBExecute := &DBExecute{
		WRITE: master(driverConfig),
		READ:  slave(driverConfig),
	}

	return dBExecute, nil
}

func (db *DBExecute) Close() {
	err = db.WRITE.Close()
	checkError(err)
	err = db.READ.Close()
	checkError(err)
}

func (db *DBExecute) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if isTransaction == true {
		println("Query Transaction Start")
		return db.WRITE.Query(query, args...)
	}
	println("Query Transaction silence")
	return db.READ.Query(query, args...)
}

func (db *DBExecute) QueryRow(query string, args ...interface{}) *sql.Row {
	if isTransaction == true {
		println("Query Transaction Start")
		return db.WRITE.QueryRow(query, args...)
	}
	println("Query Transaction silence")
	return db.READ.QueryRow(query, args...)
}

func (db *DBExecute) Exec(query string, args ...interface{}) (sql.Result, error) {
	if isTransaction == true {
		return db.WRITE.Exec(query, args...)
	}
	return db.WRITE.Exec(query, args...)
}

func (db *DBExecute) Begin() (*sql.Tx, error) {
	if IsCloseTransactionQueryToWrite == false {
		isTransaction = true
	}

	tx, err = db.WRITE.Begin()
	return tx, err
}

func Commit() error {
	err = tx.Commit()
	isTransaction = false
	return err
}

func Rollback() error {
	err = tx.Rollback()
	isTransaction = false
	return err
}

func master(driverConfig *DriverConfig) *sql.DB {
	dbMaster, err := sql.Open(driverConfig.DriverName, ConnectionString(driverConfig.DriverDBSource.WRITE))
	checkError(err)

	err = dbMaster.Ping()
	checkError(err)
	return dbMaster
}

func slave(driverConfig *DriverConfig) *sql.DB {
	dbSlave, err := sql.Open(driverConfig.DriverName, ConnectionString(driverConfig.DriverDBSource.READ))
	checkError(err)

	err = dbSlave.Ping()
	checkError(err)
	return dbSlave
}

func ConnectionString(config DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.HOST,
		config.POST,
		config.USER,
		config.PASSWORD,
		config.DATABASE,
	)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

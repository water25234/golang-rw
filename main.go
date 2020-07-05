package main

import (
	"database/sql"
	"fmt"

	"github.com/water25234/golang-rw/rw"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func main() {

	driverDBSource := rw.DriverDBSource{
		WRITE: rw.DBConfig{
			HOST:     "127.0.0.1",
			POST:     5432,
			DATABASE: "golangdb",
			USER:     "golang",
			PASSWORD: "golang",
		},
		READ: rw.DBConfig{
			HOST:     "127.0.0.1",
			POST:     5432,
			DATABASE: "golangdb",
			USER:     "golang",
			PASSWORD: "golang",
		},
	}

	// If you want to transaction query that the write for data, default false.
	// rw.IsCloseTransactionQueryToWrite = true

	db, err := rw.Open("postgres", driverDBSource)
	checkError(err)

	var id int
	var name string
	var quantity int

	err = db.QueryRow("SELECT * FROM inventory;").Scan(&id, &name, &quantity)
	checkError(err)
	fmt.Println("Data row = (%d, %s, %d)\n", id, name, quantity)

	sql := "INSERT INTO inventory (name, quantity) VALUES ($1, $2);"
	_, err = db.Exec(sql, "coconut", 300)
	checkError(err)
	fmt.Println("Inserted 1 rows of data")

	tx, err := db.Begin()
	checkError(err)

	err = db.QueryRow("SELECT * FROM inventory;").Scan(&id, &name, &quantity)
	checkError(err)
	fmt.Println("Data row = (%d, %s, %d)\n", id, name, quantity)

	fmt.Println(id)

	sql = "UPDATE inventory SET quantity = 2000 WHERE id = $1;"
	db.Exec(sql, id)

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		checkError(err)
	}

	fmt.Printf("Successful created connection to database")
}

func queryInventory() {

	sql_query := "SELECT * FROM inventory;"
	rows, err := db.Query(sql_query)
	checkError(err)

	var id int
	var name string
	var quantity int

	defer rows.Close()

	for rows.Next() {
		switch err := rows.Scan(&id, &name, &quantity); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned")
		case nil:
			fmt.Printf("Data row = (%d, %s, %d)\n", id, name, quantity)
		default:
			checkError(err)
		}
	}
}

func createInventory() {
	_, err = db.Exec("DROP TABLE IF EXISTS inventory;")
	checkError(err)
	fmt.Println("Finished dropping table (if existed)")

	_, err = db.Exec("CREATE TABLE inventory (id serial PRIMARY KEY, name VARCHAR(50), quantity INTEGER);")
	checkError(err)
	fmt.Println("Finished creating table")

	sql_statement := "INSERT INTO inventory (name, quantity) VALUES ($1, $2);"
	_, err = db.Exec(sql_statement, "banana", 150)
	checkError(err)
	_, err = db.Exec(sql_statement, "orange", 154)
	checkError(err)
	_, err = db.Exec(sql_statement, "apple", 100)
	checkError(err)
	fmt.Println("Inserted 3 rows of data")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

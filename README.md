#golang-rw

## Description
golang-rw is a library that execute DB write & read project.

## Install
```shell
$ go get github.com/water25234/golang-rw
```

## Work
```go
go run main.go

	// create table Inventory()
	createInventory()

	// add postgresql config driver parameters
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

	// connection DB & test ping 
	db, err := rw.Open("postgres", driverDBSource)
	checkError(err)
  
  	defer db.Close()
  
  
	// Query data
	var id int
	var name string
	var quantity int

	err = db.QueryRow("SELECT * FROM inventory;").Scan(&id, &name, &quantity)
	checkError(err)
	fmt.Println("Data row = (%d, %s, %d)\n", id, name, quantity)

	// insert data
	sql := "INSERT INTO inventory (name, quantity) VALUES ($1, $2);"
	_, err = db.Exec(sql, "coconut", 300)
	checkError(err)
	fmt.Println("Inserted 1 rows of data")

	// start Transaction
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

```

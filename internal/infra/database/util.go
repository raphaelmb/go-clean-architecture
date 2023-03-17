package database

import (
	"database/sql"
	"log"
)

func TryCreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	if err != nil {
		log.Fatal(err)
	}
}

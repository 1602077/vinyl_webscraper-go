package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	r "github.com/1602077/webscraper/pkg/records"
	_ "github.com/lib/pq"
)

type pgConfig struct {
	host, user, password, dbname string
	port                         int
}

var config = pgConfig{
	host:     "localhost",
	port:     5432,
	user:     os.Getenv("DBUSER"),
	password: os.Getenv("DBPASS"),
	dbname:   "records",
}

func main() {
	db := ConnectToDB(config)
	defer db.Close()

	result := QueryAllRows(db, "record")
	defer result.Close()
	Records := ReadQueryToRecord(result)
	Records.PrintRecords()
}

func ConnectToDB(dbConfig pgConfig) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.host, dbConfig.port, dbConfig.user, dbConfig.password, dbConfig.dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("err: opening connection to database '%s' failed.", dbConfig.dbname)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("err: ping to database '%s' failed.", dbConfig.dbname)
	}

	log.Printf("connection to database '%s' successfully opened.", dbConfig.dbname)

	return db
}

func QueryAllRows(db *sql.DB, table string) *sql.Rows {
	rows, err := db.Query("SELECT * FROM record") // FIXME: hardcorded table currently
	if err != nil {
		log.Fatal("err: quering all rows on database failed.")
	}
	return rows
}

func ReadQueryToRecord(rows *sql.Rows) *r.Records {
	var Records r.Records
	for rows.Next() {
		var id, art, alb string
		if err := rows.Scan(&id, &art, &alb); err != nil {
			break
		}
		Records = append(Records, r.NewRecord(art, alb, "", ""))
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("error: query row read failed")
	}
	return &Records
}

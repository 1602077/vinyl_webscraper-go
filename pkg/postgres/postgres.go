package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

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

/*
type database struct {
	db *sql.DB
	pgConfig
}

var db = database{
	db:       ConnectToDB(config),
	pgConfig: config,
}
*/

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

func ExecuteFromSQLFile(db *sql.DB, c pgConfig, filename string) {
	cmd := exec.Command("psql", "-U", c.user, "-h", c.host, "-d", c.dbname, "-a", "-f", filename)

	var out, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}
}

func QueryRecordAllRows(db *sql.DB) *sql.Rows {
	rows, err := db.Query("SELECT * FROM record")
	if err != nil {
		log.Fatal("err: quering all rows on records table failed.")
	}
	return rows
}

func ReadQueryToRecord(rows *sql.Rows) r.Records {
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
	return Records
}

func InsertIntoRecordTable(db *sql.DB, rec r.Records) {
	insertStatement := `
		INSERT INTO
			record (artist, album)
		VALUES
			($1, $2)
		RETURNING ID;`
	// var id int
	for _, rr := range rec {
		_, err := db.Exec(insertStatement, rr.GetArtist(), rr.GetAlbum())
		if err != nil {
			log.Printf("err: insert failed: %v", err)
		}
		// err := rows.Scan(&id)
		// if err != nil {
		// 	log.Printf("id retrival failed")
	}

	// log.Printf("new record inserted at id: %v", id)

}

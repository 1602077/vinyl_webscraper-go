// api methods for writing data to postgres db
package postgres

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	r "github.com/1602077/webscraper/pkg/records"
	_ "github.com/lib/pq"
)

type PgConfig struct {
	host, user, password, dbname string
	port                         int
}

func NewPgConfig(dbname string) *PgConfig {
	return &PgConfig{
		host:     "localhost",
		port:     5432,
		user:     os.Getenv("DBUSER"),
		password: os.Getenv("DBPASS"),
		dbname:   dbname,
	}
}

const DBNAME = "dev"

type PgInstance struct {
	config *PgConfig
	db     *sql.DB
}

func NewPostgresCli(dbname string) *PgInstance {
	return &PgInstance{
		config: NewPgConfig(dbname),
	}
}

// Opens a connection to database specified by config field
func (pg *PgInstance) Connect() *PgInstance {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pg.config.host, pg.config.port, pg.config.user, pg.config.password, pg.config.dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("err: opening connection to database '%s' failed.", pg.config.dbname)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("err: ping to database '%s' failed.", pg.config.dbname)
	}

	log.Printf("connection to database '%s' successfully opened.\n", pg.config.dbname)

	pg.db = db

	return pg
}

// Closes connection to database as specified by config field
func (pg *PgInstance) Close() {
	log.Printf("closing connection to database '%s'.", pg.config.dbname)
	pg.db.Close()
}

// Run command from '.sql' file on database.
func (pg *PgInstance) executeFromSQLFile(filename string) {
	cmd := exec.Command("psql", "-U", pg.config.user, "-h", pg.config.host, "-d", pg.config.dbname, "-a", "-f", filename)

	var out, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}
}

// Clears all data from tables
func (pg *PgInstance) wipe() *PgInstance {
	pg.executeFromSQLFile("../../data/wipeTables.sql")
	return pg
}

// Insert sample data into pg for testing
func (pg *PgInstance) insertTestData() *PgInstance {
	pg.executeFromSQLFile("../../data/testData.sql")
	return pg
}

// Runs "SELECT * FROM records"
func (pg *PgInstance) QueryRecordAllRows() *sql.Rows {
	rows, err := pg.db.Query("SELECT * FROM records;")
	if err != nil {
		log.Fatal("err: querying all rows on 'records' table failed.")
	}
	return rows
}

// Reads in the result of a db.Query(...) [*sql.Rows] to r.Records type
func ReadPartialQueryToRecord(rows *sql.Rows) r.Records {
	var Records r.Records
	for rows.Next() {
		var id, art, alb string
		if err := rows.Scan(&id, &art, &alb); err != nil {
			break
		}
		Records = append(Records, r.NewRecord(art, alb, "", 0))
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("error: query row read failed")
	}
	return Records
}

// Inserts a single record into 'record' table and returns it's id
func (pg *PgInstance) InsertRecordMaster(rec *r.Record) int {
	insertQuery := `
		INSERT INTO
			records (artist, album)
		VALUES
			($1, $2)
		RETURNING ID;`

	var id int
	if rows := pg.db.QueryRow(insertQuery, rec.GetArtist(), rec.GetAlbum()).Scan(&id); rows == sql.ErrNoRows {
		log.Print("err: insert into 'records' table failed.")
	}
	return id
}

// Retrieves the id of a record from 'record' table
func (pg *PgInstance) GetRecordID(rec *r.Record) (int, bool) {
	existsQuery := `
		SELECT id
		FROM records
		WHERE artist = $1 AND album = $2
		LIMIT 1;`

	var id int
	if rows := pg.db.QueryRow(existsQuery, rec.GetArtist(), rec.GetAlbum()).Scan(&id); rows == sql.ErrNoRows {
		return 0, false
	}
	return id, true
}

// Retrieves the id of a price row for a given record_id and date
func (pg *PgInstance) GetPriceID(recordID int, date time.Time) (int, bool) {
	existsQuery := `
		SELECT id
		FROM prices
		WHERE date = $1 AND record_id = $2
		LIMIT 1;`

	var id int
	if rows := pg.db.QueryRow(existsQuery, date, recordID).Scan(&id); rows == sql.ErrNoRows {
		return 0, false
	}
	return id, true
}

// Checks if record exists in 'record' table and adds if not & then inserts into pricing table
func (pg *PgInstance) InsertRecordAllTables(rec *r.Record) int {
	// check if exists in 'record table'
	recordID, ok := pg.GetRecordID(rec)
	if !ok {
		recordID = pg.InsertRecordMaster(rec)
	}

	today := time.Now()
	// today := time.Now().Format("2006-01-02")

	priceID, ok := pg.GetPriceID(recordID, today)
	if ok {
		// replace instead
		updateQuery := `
			UPDATE prices
			SET price = $1
			WHERE date = $2 AND record_id = $3
			RETURNING ID;`

		if rows := pg.db.QueryRow(updateQuery, rec.GetPrice(), today, recordID).Scan(&priceID); rows == sql.ErrNoRows {
			log.Printf("%s: no price currently stored  'prices' table.", rec.GetAlbum())
			return priceID
		}
		log.Printf("%s: updated in db.", rec.GetAlbum())
		return priceID
	}

	insertQuery := `
		INSERT INTO
			prices (date, price, record_id)
		VALUES
			($1, $2, $3)
		RETURNING ID;`

	if rows := pg.db.QueryRow(insertQuery, today, rec.GetPrice(), recordID).Scan(&priceID); rows == sql.ErrNoRows {
		log.Fatalf("%s: price insert into db[prices] failed.", rec.GetAlbum())
	}
	log.Printf("%s: written to db.", rec.GetAlbum())
	return priceID
}

// Runs "SELECT * FROM prices"
func (pg *PgInstance) QueryPriceAllRows() *sql.Rows {
	rows, err := pg.db.Query("SELECT * FROM prices;")
	if err != nil {
		log.Fatal("err: quering all rows on 'prices' table failed.")
	}
	return rows
}

// Gets most recent prices of all records
func (pg *PgInstance) GetCurrentRecordPrices() *sql.Rows {
	rows, err := pg.db.Query(`
		SELECT r.Artist, r.Album, p.MaxPrice
		FROM records r
		INNER JOIN (
			SELECT record_id, MAX(Date) as MaxDate, MAX(price) as MaxPrice
			FROM prices
			GROUP BY record_id
		) p ON p.record_id = r.id;`)

	if err != nil {
		log.Fatal("err: queries all records current prices failed.")
	}
	return rows
}

func ReadQueryToRecord(rows *sql.Rows) r.Records {
	var Records r.Records
	for rows.Next() {
		var art, alb string
		var price float64
		if err := rows.Scan(&art, &alb, &price); err != nil {
			break
		}
		Records = append(Records, r.NewRecord(art, alb, "", price))
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("error: query row read failed")
	}
	return Records
}

func (pg *PgInstance) PrintCurrentRecordPrices() {
	rows := pg.GetCurrentRecordPrices()
	rec := ReadQueryToRecord(rows)
	rec.PrintRecords()
}
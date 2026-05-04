package main

import (
	//"bufio"
	"database/sql"
	"fmt"

	//"log"
	//"os"
	//"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Album struct {
	id     int64
	title  string
	artist string
	price  float32
}

const USER = "root"
const PASSWD = "helloWorld" // put in your password here
const DATABASE = "recordings"
const CONNECTION = "tcp"
const HOST = "127.0.0.1"
const PORT = "3306"

func connectToADatabase() (bool, *sql.DB) {
	var db *sql.DB
	var err error

	query := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", USER, PASSWD, CONNECTION, HOST, PORT, DATABASE)

	// Get a database handle.
	db, err = sql.Open("mysql", query)
	if err != nil {
		fmt.Println("MySQL Open Error: ", err)
		return false, db
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("MySQL Ping Error: ", err)
		return false, db
	}

	fmt.Println("Connected to database: ", DATABASE)

	return true, db
}

func deleteARecord(db *sql.DB, title string) error {
	query := "DELETE FROM album WHERE title = ?"

	// Execute the query using db.Exec()
	result, err := db.Exec(query, title)

	if err != nil {
		return err
	}

	// Optional: Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Could not get rows affected: %v\n", err)
	}

	fmt.Printf("Deleted %d row(s)\n", rowsAffected)

	return fmt.Errorf("")
}

func displayAllRecords(db *sql.DB) error {
	var albums []Album
	var tmp Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return fmt.Errorf("displayAllRecords: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		err = rows.Scan(&tmp.id, &tmp.title, &tmp.artist, &tmp.price)
		if err != nil {
			return fmt.Errorf("displayAllRecords: %v", err)
		}
		albums = append(albums, tmp)
	}
	err = rows.Err()

	if err != nil {
		return fmt.Errorf("displayAllRecords: %v", err)
	} else {
		i := 1
		for _, value := range albums {
			fmt.Printf("%-3d %-65s %-25s $%6.2f\n", i, value.title, value.artist, value.price)
			i++
		}
	}

	return fmt.Errorf("")
}

func insertNewRecord(db *sql.DB) (int64, error) {
	var newAlbum Album

	newAlbum.artist = "Howard Shore"
	newAlbum.title = "Lord of the Ring: the Soundtrack"
	newAlbum.price = 21.99

	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", newAlbum.title, newAlbum.artist, newAlbum.price)
	if err != nil {
		return 0, fmt.Errorf("Error Inserting a New Album: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Error Inserting a New Album: %v", err)
	}
	return id, nil
}

// albumsByArtist queries for albums that have the specified artist name.
func searchAlbumsByArtist(db *sql.DB, name string) error {
	// An albums slice to hold data from returned rows.
	var albums []Album
	var tmp Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		err := rows.Scan(&tmp.id, &tmp.title, &tmp.artist, &tmp.price)

		if err != nil {
			return fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, tmp)
	}

	err = rows.Err()

	if err != nil {
		return fmt.Errorf("albumsByArtist %q: %v", name, err)
	} else {
		i := 1
		for _, value := range albums {
			fmt.Printf("%-3d %-65s %-25s $%6.2f\n", i, value.title, value.artist, value.price)
			i++
		}
	}

	return nil
}

func main() {
	var db *sql.DB
	var connected bool = false

	connected, db = connectToADatabase()

	if connected == false {
		return
	}

	displayAllRecords(db)

	id, err := insertNewRecord(db)

	if err != nil {
		fmt.Println("Error in Adding a New Record")
	} else {
		fmt.Println("Added a New Record with id = ", id)
	}

	displayAllRecords(db)

	err = deleteARecord(db, "Lord of the Ring: the Soundtrack")

	if err != nil {
		fmt.Println(err)
	}

	displayAllRecords(db)

	err = searchAlbumsByArtist(db, "Enya")
	if err != nil {
		fmt.Println(err)
	}

	db.Close()
}

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Advisee struct {
	id    string
	name  string
	dob   string
	email string
	phone string
	major string
	gpa   float64
	class string
}

const USER = "root"
const PASSWD = "helloWorld" // put in your password here
const DATABASE = "advising"
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

func displayAllAdvisees(db *sql.DB) error {
	var advisees []Advisee
	var tmp Advisee

	rows, err := db.Query(`
        SELECT p.id, p.name, p.dob, p.email, p.phone, s.major, s.gpa, s.class
        FROM person p
        JOIN student s ON p.id = s.id
    `)
	if err != nil {
		return fmt.Errorf("displayAllAdvisees: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		err = rows.Scan(&tmp.id, &tmp.name, &tmp.dob, &tmp.email, &tmp.phone, &tmp.major, &tmp.gpa, &tmp.class)
		if err != nil {
			return fmt.Errorf("displayAllPeople: %v", err)
		}
		advisees = append(advisees, tmp)
	}

	for _, value := range advisees {
		fmt.Printf("%-3s %-20s %-12s %-25s %-15s %-20s %-4.2f %-10s\n", value.id, value.name, value.dob, value.email, value.phone, value.major, value.gpa, value.class)
	}

	return nil
}

func insertNewAdvisee(db *sql.DB, id string, name string, dob string, email string, phone string, major string, gpa float64, class string) (string, error) {
	_, err := db.Exec("INSERT INTO person (id, name, dob, email, phone) VALUES (?, ?, ?, ?, ?)", id, name, dob, email, phone)
	if err != nil {
		return "", fmt.Errorf("Error Inserting a New Person: %v", err)
	}
	_, err = db.Exec("INSERT INTO student (id, major, gpa, class) VALUES (?, ?, ?, ?)", id, major, gpa, class)
	if err != nil {
		return "", fmt.Errorf("Error Inserting a New Student: %v", err)
	}
	return id, nil
}

func searchAdviseeByName(db *sql.DB, name string) error {
	var advisees []Advisee
	var tmp Advisee

	rows, err := db.Query(`
    SELECT p.id, p.name, p.dob, p.email, p.phone, s.major, s.gpa, s.class
    FROM person p
    JOIN student s ON p.id = s.id
    WHERE p.name = ?
`, name)
	if err != nil {
		return fmt.Errorf("AdviseeByName %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&tmp.id, &tmp.name, &tmp.dob, &tmp.email, &tmp.phone, &tmp.major, &tmp.gpa, &tmp.class)

		if err != nil {
			return fmt.Errorf("AdviseeByName %q: %v", name, err)
		}
		advisees = append(advisees, tmp)
	}

	err = rows.Err()

	for _, value := range advisees {
		fmt.Printf("%-3s %-20s %-12s %-25s %-15s %-20s %-4.2f %-10s\n", value.id, value.name, value.dob, value.email, value.phone, value.major, value.gpa, value.class)
	}

	return nil
}

func deleteAdvisee(db *sql.DB, name string) error {
	var id string

	err := db.QueryRow("SELECT id FROM person WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		fmt.Println("No advisee found with name:", name)
		return nil
	} else if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM student WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("Error deleting from student: %v", err)
	}
	// Then delete from person
	result, err := db.Exec("DELETE FROM person WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("Error deleting from person: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Deleted %d row(s) from person (and corresponding student)\n", rowsAffected)
	return nil
}

func updateAdvisee(db *sql.DB, name string) error {
	var id string

	// Get the id of the advisee
	err := db.QueryRow("SELECT id FROM person WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		fmt.Println("No advisee found with name:", name)
		return nil
	} else if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("What would you like to update?")
	fmt.Println("Options: name, dob, email, phone, major, gpa, class")
	choice := readLine(reader)

	switch choice {

	case "name":
		fmt.Println("Enter new name:")
		newValue := readLine(reader)
		_, err = db.Exec("UPDATE person SET name = ? WHERE id = ?", newValue, id)

	case "dob":
		fmt.Println("Enter new date of birth (mm/dd/yyyy):")
		newValue := readLine(reader)
		_, err = db.Exec("UPDATE person SET dob = ? WHERE id = ?", newValue, id)

	case "email":
		fmt.Println("Enter new email:")
		newValue := readLine(reader)
		_, err = db.Exec("UPDATE person SET email = ? WHERE id = ?", newValue, id)

	case "phone":
		fmt.Println("Enter new phone:")
		newValue := readLine(reader)
		_, err = db.Exec("UPDATE person SET phone = ? WHERE id = ?", newValue, id)

	case "major":
		fmt.Println("Enter new major:")
		newValue := readLine(reader)
		_, err = db.Exec("UPDATE student SET major = ? WHERE id = ?", newValue, id)

	case "gpa":
		fmt.Println("Enter new GPA:")
		var newGPA float64
		fmt.Scan(&newGPA)
		fmt.Scanln()
		_, err = db.Exec("UPDATE student SET gpa = ? WHERE id = ?", newGPA, id)

	case "class":
		fmt.Println("Enter new class:")
		newValue := readLine(reader)
		_, err = db.Exec("UPDATE student SET class = ? WHERE id = ?", newValue, id)

	default:
		fmt.Println("Invalid option.")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error updating advisee: %v", err)
	}

	fmt.Println("Advisee updated successfully.")
	return nil
}

func readLine(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	var db *sql.DB
	var connected bool = false
	applicationOpen := true
	var userResponse string
	reader := bufio.NewReader(os.Stdin)

	connected, db = connectToADatabase()

	if connected == false {
		return
	}

	//Open Application Loop
	for applicationOpen {
		fmt.Println("To display all of the advisees in the database, type 'display'\nTo insert a new advisee to the database, type 'insert'\nTo remove an advisee from the database, type 'remove'\nTo search for an advisee from the database, type 'search'\nTo update the information of an advisee, type 'update'\nTo close the application, type 'close'")
		userResponse = readLine(reader)
		var gpa float64 = 0

		switch userResponse {

		//Use Display fuction to display table of advisees
		case "display":
			displayAllAdvisees(db)

		//Ask user for all of the information for new advisee and then add the advisee to the database
		case "insert":
			fmt.Println("Enter advisee's id")
			var id = readLine(reader)
			fmt.Println("Enter advisee's name")
			var name = readLine(reader)
			fmt.Println("Enter advisee's date of birth (mm/dd/yyyy)")
			var dob = readLine(reader)
			fmt.Println("Enter advisee's email")
			var email = readLine(reader)
			fmt.Println("Enter advisee's phone")
			var phone = readLine(reader)
			fmt.Println("Enter advisee's major")
			var major = readLine(reader)
			fmt.Println("Enter student's GPA")
			fmt.Scan(&gpa)
			fmt.Scanln()
			fmt.Println("Enter advisee's class")
			var class = readLine(reader)

			insertNewAdvisee(db, id, name, dob, email, phone, major, gpa, class)

		//Stop the loop from continuing to end the program
		case "close":
			fmt.Println("Have a nice day")
			applicationOpen = false

		case "remove":
			fmt.Println("Enter the advisee's name to be deleted")
			var name = readLine(reader)
			deleteAdvisee(db, name)

		case "search":
			fmt.Println("Enter the advisee's name")
			var name = readLine(reader)
			searchAdviseeByName(db, name)

		case "update":
			fmt.Println("Enter the advisee's name")
			var name = readLine(reader)
			updateAdvisee(db, name)

		//Error handling for misinputs
		default:
			fmt.Println("That was not an option")
		}
	}

	db.Close()
}

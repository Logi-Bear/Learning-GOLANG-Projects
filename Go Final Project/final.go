package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const URL = "https://classlist.champlain.edu/api3/courses/semester/fall/type/all/filter/ug"

const USER = "root"
const PASSWD = // put in your password here
const DATABASE = "coursesList"
const CONNECTION = "tcp"
const HOST = "127.0.0.1"
const PORT = "3306"

type Course struct {
	ID              int    `json:"id"`
	Number          string `json:"number"`
	Credit          string `json:"credit"`
	OpenSeats       string `json:"openseats"`
	Days            string `json:"days"`
	Times           string `json:"times"`
	InstructorFname string `json:"instructor_fname"`
	InstructorLname string `json:"instructor_lname"`
	Description     string `json:"description"`
	Room            string `json:"room"`
	Subject         string `json:"subject"`
	CourseType      string `json:"type"`
	Prereq          string `json:"prereq"`
	Title           string `json:"title"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
}

type Semester struct {
	ID          string   `json:"id"`
	Year        int      `json:"year"`
	DisplayDate string   `json:"display_date"`
	CanRegister bool     `json:"can_register"`
	Type        string   `json:"type"`
	Children    []Course `json:"children"`
}

type CoursesAPIResponse struct {
	Identifier  string     `json:"identifier"`
	DataObjects []Semester `json:"items"`
}

var courseList []Course

func scrapeWeb(wg *sync.WaitGroup) {
	defer wg.Done()

	// Make an HTTP GET request
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalf("Error fetching API URL %s: %v", URL, err)
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Unmarshal the JSON data into the Post struct
	var data CoursesAPIResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	//print the structured data
	for _, item := range data.DataObjects {
		for _, course := range item.Children {
			courseList = append(courseList, course)
		}
	}
	fmt.Println("Web scraping complete. Data stored.")
}

func createDatabase(wg *sync.WaitGroup) {
	defer wg.Done()

	query := fmt.Sprintf("%s:%s@%s(%s:%s)/?parseTime=true", USER, PASSWD, CONNECTION, HOST, PORT)
	db, err := sql.Open("mysql", query)
	if err != nil {
		log.Fatalf("failed to connect to MySQL server: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", DATABASE))
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	query = fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", USER, PASSWD, CONNECTION, HOST, PORT, DATABASE)
	db, err = sql.Open("mysql", query)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS courses;")
	if err != nil {
		log.Fatalf("failed to drop existing courses table: %v", err)
	}

	createTableQuery := `
	CREATE TABLE courses (
		number VARCHAR(15) NOT NULL,
		credit VARCHAR(2) NOT NULL,
		openSeats VARCHAR(10) NOT NULL,
		days VARCHAR(15) NOT NULL,
		times VARCHAR(40) NOT NULL,
		instructorFname VARCHAR(50) NOT NULL,
		instructorLname VARCHAR(30) NOT NULL,
		description LONGTEXT NOT NULL,
		room VARCHAR(20) NOT NULL,
		subject VARCHAR(50) NOT NULL,
		courseType VARCHAR(20) NOT NULL,
		prereq LONGTEXT NOT NULL,
		title VARCHAR(100) NOT NULL,
		startDate VARCHAR(15) NOT NULL,
		endDate VARCHAR(15) NOT NULL
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("failed to create courses table: %v", err)
	}

	fmt.Println("Database and table created successfully.")
}

func insertCourses(db *sql.DB) {
	for _, course := range courseList {
		_, err := db.Exec(`INSERT INTO courses (number, credit, openSeats, days, times, instructorFname, instructorLname, description, room, 
		subject, courseType, prereq, title, startDate, endDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
			course.Number, course.Credit, course.OpenSeats, course.Days, course.Times, course.InstructorFname, course.InstructorLname,
			course.Description, course.Room, course.Subject, course.CourseType, course.Prereq, course.Title, course.StartDate, course.EndDate)

		if err != nil {
			log.Fatalf("failed to insert course %s: %v", course.Number, err)
		}
	}
	fmt.Println("Data loaded into the database successfully.")
}

func displayAllCourses() {
	for _, course := range courseList {
		fmt.Printf(
			"Title: %s\nNumber: %s\nDays: %s\nTimes: %s\nRoom: %s\nInstructor: %s %s\nOpen Seats: %v\n\n",
			course.Title,
			course.Number,
			course.Days,
			course.Times,
			course.Room,
			course.InstructorFname,
			course.InstructorLname,
			course.OpenSeats,
		)
	}
}

func searchForOnCampus(db *sql.DB) ([]Course, error) {
	rows, err := db.Query(`
		SELECT number, credit, openSeats, days, times,
		       instructorFname, instructorLname, description,
		       room, subject, courseType, prereq,
		       title, startDate, endDate
		FROM courses
		WHERE courseType = ?`, "Day/Evening")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Course

	for rows.Next() {
		var c Course

		err := rows.Scan(
			&c.Number,
			&c.Credit,
			&c.OpenSeats,
			&c.Days,
			&c.Times,
			&c.InstructorFname,
			&c.InstructorLname,
			&c.Description,
			&c.Room,
			&c.Subject,
			&c.CourseType,
			&c.Prereq,
			&c.Title,
			&c.StartDate,
			&c.EndDate,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func searchForCCO(db *sql.DB) ([]Course, error) {
	rows, err := db.Query(`
		SELECT number, credit, openSeats, days, times,
		       instructorFname, instructorLname, description,
		       room, subject, courseType, prereq,
		       title, startDate, endDate
		FROM courses
		WHERE courseType IN (?, ?)`, "Online", "Accelerated")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Course

	for rows.Next() {
		var c Course

		err := rows.Scan(
			&c.Number,
			&c.Credit,
			&c.OpenSeats,
			&c.Days,
			&c.Times,
			&c.InstructorFname,
			&c.InstructorLname,
			&c.Description,
			&c.Room,
			&c.Subject,
			&c.CourseType,
			&c.Prereq,
			&c.Title,
			&c.StartDate,
			&c.EndDate,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func searchForPrefix(db *sql.DB, prefix string) ([]Course, error) {
	rows, err := db.Query(`
		SELECT number, credit, openSeats, days, times,
		       instructorFname, instructorLname, description,
		       room, subject, courseType, prereq,
		       title, startDate, endDate
		FROM courses
		WHERE number LIKE ?`, prefix)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Course

	for rows.Next() {
		var c Course

		err := rows.Scan(
			&c.Number,
			&c.Credit,
			&c.OpenSeats,
			&c.Days,
			&c.Times,
			&c.InstructorFname,
			&c.InstructorLname,
			&c.Description,
			&c.Room,
			&c.Subject,
			&c.CourseType,
			&c.Prereq,
			&c.Title,
			&c.StartDate,
			&c.EndDate,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func searchForPrefixAndLevel(db *sql.DB, prefix string, level string) ([]Course, error) {
	rows, err := db.Query(`
		SELECT number, credit, openSeats, days, times,
		       instructorFname, instructorLname, description,
		       room, subject, courseType, prereq,
		       title, startDate, endDate
		FROM courses
		WHERE number LIKE ?
  		AND SUBSTRING(number, LOCATE(' ', number) + 1, 1) = ?`,
		prefix+"%", level)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Course

	for rows.Next() {
		var c Course

		err := rows.Scan(
			&c.Number,
			&c.Credit,
			&c.OpenSeats,
			&c.Days,
			&c.Times,
			&c.InstructorFname,
			&c.InstructorLname,
			&c.Description,
			&c.Room,
			&c.Subject,
			&c.CourseType,
			&c.Prereq,
			&c.Title,
			&c.StartDate,
			&c.EndDate,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func searchForLevel(db *sql.DB, level string) ([]Course, error) {
	rows, err := db.Query(`
		SELECT number, credit, openSeats, days, times,
		       instructorFname, instructorLname, description,
		       room, subject, courseType, prereq,
		       title, startDate, endDate
		FROM courses
		WHERE SUBSTRING(number, LOCATE(' ', number) + 1, 1) = ?`,
		level)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Course

	for rows.Next() {
		var c Course

		err := rows.Scan(
			&c.Number,
			&c.Credit,
			&c.OpenSeats,
			&c.Days,
			&c.Times,
			&c.InstructorFname,
			&c.InstructorLname,
			&c.Description,
			&c.Room,
			&c.Subject,
			&c.CourseType,
			&c.Prereq,
			&c.Title,
			&c.StartDate,
			&c.EndDate,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func searchForOpen(db *sql.DB) ([]Course, error) {
	rows, err := db.Query(`
		SELECT number, credit, openSeats, days, times,
		       instructorFname, instructorLname, description,
		       room, subject, courseType, prereq,
		       title, startDate, endDate
		FROM courses
		WHERE openSeats NOT LIKE '0'`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Course

	for rows.Next() {
		var c Course

		err := rows.Scan(
			&c.Number,
			&c.Credit,
			&c.OpenSeats,
			&c.Days,
			&c.Times,
			&c.InstructorFname,
			&c.InstructorLname,
			&c.Description,
			&c.Room,
			&c.Subject,
			&c.CourseType,
			&c.Prereq,
			&c.Title,
			&c.StartDate,
			&c.EndDate,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func searchForCourseByNumber(db *sql.DB, number string) (*Course, error) {
	var c Course
	err := db.QueryRow(`
    SELECT number, credit, openSeats, days, times,
           instructorFname, instructorLname, description,
           room, subject, courseType, prereq,
           title, startDate, endDate
    FROM courses
    WHERE number = ?`, number).Scan(
		&c.Number,
		&c.Credit,
		&c.OpenSeats,
		&c.Days,
		&c.Times,
		&c.InstructorFname,
		&c.InstructorLname,
		&c.Description,
		&c.Room,
		&c.Subject,
		&c.CourseType,
		&c.Prereq,
		&c.Title,
		&c.StartDate,
		&c.EndDate,
	)
	if err == sql.ErrNoRows {
		fmt.Println("No course found with that number")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &c, nil
}

func searchForCourseByTitle(db *sql.DB, title string) (*Course, error) {
	var c Course
	err := db.QueryRow(`
    SELECT number, credit, openSeats, days, times,
           instructorFname, instructorLname, description,
           room, subject, courseType, prereq,
           title, startDate, endDate
    FROM courses
    WHERE title = ?`, title).Scan(
		&c.Number,
		&c.Credit,
		&c.OpenSeats,
		&c.Days,
		&c.Times,
		&c.InstructorFname,
		&c.InstructorLname,
		&c.Description,
		&c.Room,
		&c.Subject,
		&c.CourseType,
		&c.Prereq,
		&c.Title,
		&c.StartDate,
		&c.EndDate,
	)
	if err == sql.ErrNoRows {
		fmt.Println("No course found with that title")
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &c, nil
}

func displayCourses(courses []Course) {
	for _, course := range courses {
		fmt.Printf(
			"Title: %s\nNumber: %s\nDays: %s\nTimes: %s\nRoom: %s\nInstructor: %s %s\nOpen Seats: %v\n\n",
			course.Title,
			course.Number,
			course.Days,
			course.Times,
			course.Room,
			course.InstructorFname,
			course.InstructorLname,
			course.OpenSeats,
		)
	}
}

func printCourse(course Course) {
	fmt.Printf(
		"----------------------------\n"+
			"Title: %s\n"+
			"Number: %s\n"+
			"Subject: %s\n"+
			"Type: %s\n"+
			"Credit: %s\n"+
			"Instructor: %s %s\n"+
			"Days/Times: %s / %s\n"+
			"Room: %s\n"+
			"Open Seats: %s\n"+
			"Prerequisites: %s\n"+
			"Description: %s\n"+
			"Start Date: %s\n"+
			"End Date: %s\n"+
			"----------------------------\n",
		course.Title,
		course.Number,
		course.Subject,
		course.CourseType,
		course.Credit,
		course.InstructorFname,
		course.InstructorLname,
		course.Days,
		course.Times,
		course.Room,
		course.OpenSeats,
		course.Prereq,
		course.Description,
		course.StartDate,
		course.EndDate,
	)
}

func readLine(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	var wg sync.WaitGroup
	applicationOpen := true
	var userResponse string
	reader := bufio.NewReader(os.Stdin)

	wg.Add(2)
	go scrapeWeb(&wg)
	go createDatabase(&wg)

	wg.Wait()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", USER, PASSWD, CONNECTION, HOST, PORT, DATABASE))
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	insertCourses(db)

	for applicationOpen {
		fmt.Println(`
================= COURSE Database =================

 1 | Display all courses
 2 | Display all on-campus courses
 3 | Display all CCO courses
 4 | Search by course prefix
 5 | Search by prefix and level
 6 | Search by course level
 7 | Search for open courses
 8 | Display course details
 9 | Close application

Enter your choice (1-9):`)

		userResponse = readLine(reader)

		switch userResponse {
		case "1":
			displayAllCourses()

		case "2":
			courses, err := searchForOnCampus(db)
			if err == nil {
				displayCourses(courses)
			}

		case "3":
			courses, err := searchForCCO(db)
			if err == nil {
				displayCourses(courses)
			}

		case "4":
			fmt.Println("What type of course are you looking for? (e.g., CSI, GPR, MTH)")
			prefix := readLine(reader)
			courses, err := searchForPrefix(db, prefix+"%")
			if err == nil {
				displayCourses(courses)
			}

		case "5":
			fmt.Println("What type of course are you looking for? (e.g., CSI, GPR, MTH)")
			prefix := readLine(reader)
			fmt.Println("What level of course? (Enter 1, 2, 3, 4, or 5)")
			level := readLine(reader)
			courses, err := searchForPrefixAndLevel(db, prefix+"%", level)
			if err == nil {
				displayCourses(courses)
			}

		case "6":
			fmt.Println("What level of course? (Enter 1, 2, 3, 4, or 5)")
			level := readLine(reader)
			courses, err := searchForLevel(db, level)
			if err == nil {
				displayCourses(courses)
			}

		case "7":
			courses, err := searchForOpen(db)
			if err == nil {
				displayCourses(courses)
			}

		case "8":
			fmt.Println("Are you using the course's number or title?")
			userChoice := readLine(reader)
			if userChoice == "number" {
				fmt.Println("Enter the course's number:")
				number := readLine(reader)
				course, err := searchForCourseByNumber(db, number)
				if err == nil {
					printCourse(*course)
				}
			} else {
				fmt.Println("Enter the course's title:")
				title := readLine(reader)
				course, err := searchForCourseByTitle(db, title)
				if err == nil && course != nil {
					printCourse(*course)
				}
			}

		case "9":
			fmt.Println("Closing Database")
			applicationOpen = false

		default:
			fmt.Println("That was not an option")
		}
	}
}

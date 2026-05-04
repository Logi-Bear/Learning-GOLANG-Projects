package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type AddressStruct struct {
	street string
	city   string
	zip    string
	state  string
}

type PersonStruct struct {
	name        string
	address     AddressStruct
	phoneNumber string
}

type StudentStruct struct {
	person PersonStruct
	id     string
	major  string
	gpa    float32
}

func display(student map[string]StudentStruct) {
	fmt.Println("Students:")
	for _, value := range student {
		fmt.Println("ID:", value.id)
		fmt.Println("Name:", value.person.name)
		fmt.Println("Major:", value.major)
		fmt.Println("GPA:", value.gpa)
		fmt.Println("Address:", value.person.address.street,
			value.person.address.city, value.person.address.state,
			value.person.address.zip)
		fmt.Println("Phone Number:", value.person.phoneNumber)
		fmt.Println()
	}
}

func readLine(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
    //Initialize Variables
	applicationOpen := true
	var userResponse string

	var students map[string]StudentStruct

	var userName string
	var userStreet string
	var userCity string
	var userZip string
	var userState string
	var userPN string
	var userID string
	var userMajor string
	var userGPA float32

	students = make(map[string]StudentStruct)

	reader := bufio.NewReader(os.Stdin)

    //Open Application Loop
	for applicationOpen {
		fmt.Println("To display the information of all students in the directory, type 'display'\nTo add a new student to the directory, type 'add'\nTo close the application, type 'close'")
		userResponse = readLine(reader)

		switch userResponse {

            //Use Display fuction to display map of students
            case "display":
                display(students)

            //Ask user for all of the information for new student and then add the student to the map
            case "add":
                fmt.Println("Enter student's name")
                userName = readLine(reader)
                fmt.Println("Enter student's street")
                userStreet = readLine(reader)
                fmt.Println("Enter student's city")
                userCity = readLine(reader)
                fmt.Println("Enter student's zip")
                userZip = readLine(reader)
                fmt.Println("Enter student's Phone Number")
                userPN = readLine(reader)
                fmt.Println("Enter student's ID")
                userID = readLine(reader)
                fmt.Println("Enter student's major")
                userMajor = readLine(reader)
                fmt.Println("Enter student's GPA")
                fmt.Scan(&userGPA)
                fmt.Scanln()

                students[userID] = StudentStruct{
                    person: PersonStruct{
                        name: userName,
                        address: AddressStruct{
                            street: userStreet,
                            city:   userCity,
                            zip:    userZip,
                            state:  userState,
                        },
                        phoneNumber: userPN,
                    },
                    id:    userID,
                    major: userMajor,
                    gpa:   userGPA,
                }
                fmt.Println()
                display(students)

            //Stop the loop from continuing to end the program
            case "close":
                fmt.Println("Have a nice day")
                applicationOpen = false

            //Error handling for misinputs
            default:
                fmt.Println("That was not an option")
        }
	}
}
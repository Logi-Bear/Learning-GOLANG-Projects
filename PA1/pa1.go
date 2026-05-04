/*
Author: Logan McCandless
Course: CSC-380
PA #1 – Fortune Telling Machine
Date: February 5, 2026
Description: Simulates a fortune teller machine that outputs age, zodiac sign, and a random fortune message.
*/

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

//Types

type date struct {
	month int
	day   int
	year  int
}

type zodiac struct {
	name       string
	startMonth int
	startDay   int
	endMonth   int
	endDay     int
}

//Constants

var DAYS_IN_MONTH = [12]int{
	31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31,
}

var ZODIACS = []zodiac{
	{"Capricorn", 12, 22, 1, 20},
	{"Aquarius", 1, 21, 2, 19},
	{"Pisces", 2, 20, 3, 20},
	{"Aries", 3, 21, 4, 20},
	{"Taurus", 4, 21, 5, 21},
	{"Gemini", 5, 22, 6, 21},
	{"Cancer", 6, 22, 7, 22},
	{"Leo", 7, 23, 8, 21},
	{"Virgo", 8, 22, 9, 23},
	{"Libra", 9, 24, 10, 23},
	{"Scorpio", 10, 24, 11, 22},
	{"Sagittarius", 11, 23, 12, 21},
}

const MIN_MESSAGES = 10

//Functions

func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	return year%4 == 0
}

func validateDate(d date) bool {
	// Year must be positive
	if d.year <= 0 {
		return false
	}
	// Month must be 1–12
	if d.month < 1 || d.month > 12 {
		return false
	}
	days := DAYS_IN_MONTH[d.month-1]

	// Handle leap year for February
	if d.month == 2 && isLeapYear(d.year) {
		days = 29
	}
	// Day must be valid for the given month
	if d.day < 1 || d.day > days {
		return false
	}
	return true
}

func calculateAge(todayDate date, birthDate date) int {
	age := todayDate.year - birthDate.year
	// If birthday hasn't happened yet this year, subtract 1
	if todayDate.month < birthDate.month ||
		(todayDate.month == birthDate.month && todayDate.day < birthDate.day) {
		age--
	}
	return age
}

func calculateZodiac(month int, day int) string {
	for i := 0; i < len(ZODIACS); i++ {
		z := ZODIACS[i]
		if (month == z.startMonth && day >= z.startDay) || (month == z.endMonth && day <= z.endDay) || (z.startMonth > z.endMonth && ((month == z.startMonth && day >= z.startDay) || (month == z.endMonth && day <= z.endDay))) {
			return z.name
		}
	}
	return "Unknown"
}

func loadMessages(fileName string) []string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file")
		return nil
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	if len(lines) == 0 {
		fmt.Println("No messages found")
		return nil
	}
	return lines
}

func readLine(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

//Main

func main() {
	//init variables
	applicationOpen := true
	var userName string
	var userAge int
	var dateOfBirth date
	var zodiacSign string
	var fortuneMessage string
	var choice string

	rand.Seed(time.Now().UnixNano())
	currentTime := time.Now()
	var todayDate date
	todayDate.day = currentTime.Day()
	todayDate.month = int(currentTime.Month())
	todayDate.year = currentTime.Year()

	messages := loadMessages("messages.txt")
	if len(messages) < MIN_MESSAGES {
		fmt.Printf("Error: You must have at least %d fortune messages.\n", MIN_MESSAGES)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the fortune teller machine!")

	for applicationOpen {
		//User input
		fmt.Print("Enter your name: ")
		userName = readLine(reader)

		fmt.Print("Enter your date of birth (MM/DD/YYYY): ")
		fmt.Scanf("%d/%d/%d", &dateOfBirth.month, &dateOfBirth.day, &dateOfBirth.year)

		//Calculate results
		if validateDate(dateOfBirth) {
			userAge = calculateAge(todayDate, dateOfBirth)
			zodiacSign = calculateZodiac(dateOfBirth.month, dateOfBirth.day)
			fortuneMessage = messages[rand.Intn(len(messages))]

			//Display Fortune
			fmt.Println(" Your Fortune ")
			fmt.Println("-------------------")
			fmt.Printf("Name: %s\n", userName)
			fmt.Printf("Age: %d\n", userAge)
			fmt.Printf("Zodiac Sign: %s\n", zodiacSign)
			fmt.Printf("Fortune: %s\n\n", fortuneMessage)

			//Continue or exit
			fmt.Print("Would you like another fortune? (y/n): ")
			fmt.Scan(&choice)
			fmt.Scanln()
			if choice != "y" {
				applicationOpen = false
				fmt.Println("Thank you for using the fortune teller machine. Goodbye!")
			}

		} else {
			fmt.Println("Invalid date. Please try again.")
			continue
		}
	}
}

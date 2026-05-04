package main

import "fmt"

func main() {
	//init variables
	var userYear int
	var isLeapYear bool
	//ask for input
	fmt.Println("Enter your year")
	//take input
	fmt.Scan(&userYear)
	//calculate if input is leap year
	if userYear < 0 {
		fmt.Println("Error! Invalid Year")
	} else {
		if userYear%4 == 0 {
			isLeapYear = true
			if userYear%100 == 0 && userYear%400 != 0 {
				isLeapYear = false
			}
		}
			//display results
		if isLeapYear {
			fmt.Println("Your year is a leap year")
	} else {
			fmt.Println("Your year is not a leap year")
		}
	}
}

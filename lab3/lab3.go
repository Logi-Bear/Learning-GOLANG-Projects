package main

import "math/rand"
import "fmt"

func main() {
//setup variables
var num1 int
var num2 int
var operation int
var playing bool = true
var userResponse string
var correctAnswer int
var symbol string
var userAnswer int
//start loop
for playing {
	fmt.Println("Welcome to number game. Type 'play' if you would like to play or type 'stop' if you would like to stop playing.")

	//take user input
	fmt.Scan(&userResponse)
	if userResponse == "play"{
		num1 = rand.Intn(11)
		num2 = rand.Intn(11)
		operation = rand.Intn(3)
		fmt.Println("Answer this equation:")
			switch operation {
			case 0:
				correctAnswer = num1 + num2
				symbol = "+"
			case 1:
				correctAnswer = num1 - num2
				symbol = "-"
			case 2:
				correctAnswer = num1 * num2
				symbol = "*"
			}
						
			fmt.Printf("What is %d %s %d?\n", num1, symbol, num2)
			fmt.Scan(&userAnswer)
			//calculate if answer is correct
			if userAnswer == correctAnswer {
				fmt.Println("Correct!")
			} else {
				fmt.Printf("Incorrect. The correct answer is %d\n", correctAnswer)

			}

	
	} else if userResponse == "stop" {
		fmt.Println("See you later!")
		playing = false
	} else {
		fmt.Println("That wasn't an option")
	}
}
}

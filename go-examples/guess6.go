package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randomInt(low int, high int) int {
	return rand.Intn(high-low+1) + low
}

func assessGuess(guess int, target int) bool {

	if guess == target {
		return true
	}
	if guess < target {
		fmt.Println("That's too low. Try again.")
	} else {
		fmt.Println("That's too high. Try again.")
	}
	return false
}

func promptForGuess(tries int, bound int) {
	if tries == 0 {
		fmt.Print("Enter a number: ")
	} else if tries == bound-1 {
		fmt.Print("This is your final guess. What's my number? ")
	} else {
		fmt.Print("What's your next guess? ")
	}
}

func playGame(number int, bound int) bool {

	success := false
	tries := 0

	for !success && tries < 6 {
		promptForGuess(tries, bound)
		var guess int
		fmt.Scanln(&guess)
		tries = tries + 1
		success = assessGuess(guess, number)
	}

	return success
}

func main() {

	rand.Seed(time.Now().UnixNano())
	number := randomInt(1, 100)

	fmt.Print("I've chosen a number from 1 to 100. ")
	fmt.Println("Try to guess what it is.")

	theyWon := playGame(number, 6)

	if theyWon {
		fmt.Print("Well done! ")
	} else {
		fmt.Println("Sorry, you are out of guesses...")
	}
	fmt.Println(number, "was the number I chose.")
}

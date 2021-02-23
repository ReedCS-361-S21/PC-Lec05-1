package main

import "fmt"

func cToF(degreesC int) int {
	degreesF := degreesC*9/5 + 32
	return degreesF
}

func main() {

	var c int
	fmt.Print("Enter a temperature in degrees celsius: ")
	fmt.Scanln(&c)
	fmt.Println("That's", cToF(c), "degrees fahrenheit.")

}

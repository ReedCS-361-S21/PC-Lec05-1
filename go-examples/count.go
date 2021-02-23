package main

import "fmt"

func main() {

	var top int
	fmt.Print("Enter the ending count: ")
	fmt.Scanln(&top)

	count := 0
	for count <= top {
		fmt.Println(count)
		count = count + 1
	}
	fmt.Println("Woo!")

}

package main

import "fmt"

func main() {
	var a, b int
	fmt.Scanln(&a, &b)
	fmt.Printf("(a + b): %v\n", (a + b))
}

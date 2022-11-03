package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 1 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "up":
		fmt.Println("Running migrations up")
	case "down":
		fmt.Println("Running migrations down")
	default:
		fmt.Println("Please provide a valid command")
	}

}
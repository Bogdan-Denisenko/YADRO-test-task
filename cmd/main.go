package main

import (
	"fmt"
	"os"
	"testTaskYADRO/internal/computerclub"
	"testTaskYADRO/internal/parser"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Requires one argument - file name")
		os.Exit(1)
	}

	filePath := os.Args[1]

	numTables, startTime, endTime, hourCost, events, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	computerClub := computerclub.NewComputerClub(numTables, startTime, endTime, hourCost)
	err = computerClub.Simulate(events)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

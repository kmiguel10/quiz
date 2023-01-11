package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	// once file is open, create a csv reader
	r := csv.NewReader(file)

	//parse the csv - read all lines
	lines, err := r.ReadAll()

	//see if there is an error
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)

	//keep track of correct answers
	correct := 0

	//print out the problems line by line for the users
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		var answer string
		fmt.Scanf("%s\n", &answer) //we use & because we need to be able to access the answer once we have it
		if answer == p.a {
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))

}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]), //trim the space - corner case
		}
	}
	return ret
}

// to make whatever problem input regardless of the source, uniform to our program
type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println((msg))
	os.Exit(1) //exit with status code and error message
}

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 10, "the time limit for the quiz in seconds")
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

	//create timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	//<-timer.C //wait for a message from this channel, so code is blocked until we get a message from this channel

	//keep track of correct answers
	correct := 0

problemloop:
	//print out the problems line by line for the users
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string) //create a channel to send answers

		//a go routine to listen to answers
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer) //we use & because we need to be able to access the answer once we have it
			answerCh <- answer
		}()
		//keep presenting problems until we receive a message from the channel which means that time is up
		//now we only have 2 options, which is receiving data from 2 channels.... if the timer is up or if there is answer to a question
		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh: //if there is an answer
			if answer == p.a {
				correct++
			}
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

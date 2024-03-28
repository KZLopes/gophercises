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
	csvFileName := flag.String("csv", "problems.csv", "path to a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("timer", 30, "Time limit for the quiz (in seconds)")
	flag.Parse()
	_ = timeLimit

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)

	records, _ := r.ReadAll()

	problems := parseRecords(records)

	score := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

problemLoop:
	for i, p := range problems {
		fmt.Printf("Question #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case answer := <-answerCh:
			if answer == p.a {
				score++
			}
		case <-timer.C:
			fmt.Println("\nYou ran out of time!")
			break problemLoop
		}
		if i == len(problems)-1 {
			fmt.Println("All questions answered!!!")
		}
	}
	fmt.Printf("You got %d out of %d questions correct!\n", score, len(problems))
}

type problem struct {
	q string
	a string
}

func parseRecords(records [][]string) []problem {
	ret := make([]problem, len(records))
	for i, record := range records {
		ret[i] = problem{
			q: record[0],
			a: strings.TrimSpace(record[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

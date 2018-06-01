package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	problemsFile string
	timeLimit    time.Duration
	randomize    bool
)

func main() {
	flag.StringVar(&problemsFile, "problems", "problems.csv", "define the path of a CSV list of in the format of `question,answer`")
	flag.DurationVar(&timeLimit, "timer", time.Duration(30*time.Second), "sets the time limit of the quiz")
	flag.BoolVar(&randomize, "randomize", false, "sets the quiz to be randomized each run")
	flag.Parse()

	p, err := filepath.Abs(problemsFile)
	if err != nil {
		panic(fmt.Errorf("failed to find file `%s`", problemsFile))
	}

	f, err := os.Open(p)
	if err != nil {
		panic(fmt.Errorf("failed to open file `%s`", problemsFile))
	}

	r := csv.NewReader(f)
	problems, err := r.ReadAll()
	if err != nil {
		panic(fmt.Errorf("failed to parse CSV file `%s`", problemsFile))
	}

	fmt.Printf("Press enter to start quiz...")
	fmt.Scanln()

	if randomize {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	correctAnswers := 0
	timer := time.After(timeLimit)
	doneChan := make(chan struct{})

	go func() {
		for _, problem := range problems {
			fmt.Printf("%s: ", problem[0])
			var answer string
			fmt.Scanln(&answer)

			if answer == strings.TrimSpace(problem[1]) {
				correctAnswers++
			}
		}
		<-doneChan
	}()

questionAnswerLoop:
	for {
		select {
		case <-timer:
			fmt.Println("\nout of time!")
			break questionAnswerLoop
		case <-doneChan:
			fmt.Println("done!")
			break questionAnswerLoop
		}
	}

	fmt.Printf("%d/%d correct\n", correctAnswers, len(problems))
}

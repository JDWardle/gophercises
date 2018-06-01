package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	problemsFile string
)

func main() {
	flag.StringVar(&problemsFile, "problems", "problems.csv", "define the path of a CSV list of in the format of `question,answer`")
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

	correctAnswers := 0

	for _, problem := range problems {
		fmt.Printf("%s: ", problem[0])
		var answer string
		fmt.Scanln(&answer)

		if answer == problem[1] {
			correctAnswers++
		}
	}

	fmt.Printf("%d/%d correct\n", correctAnswers, len(problems))
}

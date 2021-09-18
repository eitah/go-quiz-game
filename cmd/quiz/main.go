package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

func mainErr() error {
	quizFile := flag.String("in", "./problems.csv", "Path to a quiz file to load with data. In the csv format of question,answer, e.g. 1+2,3")
	flag.Parse()
	f, err := os.Open(*quizFile)
	if err != nil {
		return err
	}
	defer f.Close()

	parsed, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	var questions []string
	var answers []string
	for _, qna := range parsed {
		for index, item := range qna {
			if index == 0 {
				questions = append(questions, item)
			}
			if index == 1 {
				answers = append(answers, item)
			}
		}
	}

	buf := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the quiz game.")

	var countCorrect int
	var countWrong int
	for index, question := range questions {
		fmt.Printf("Here is your next question. %s\n", question)
		fmt.Print(">")
		sentence, err := buf.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Printf("You guessed: %s. The right answer was %s\n", strings.TrimSpace(sentence), answers[index])
		if strings.TrimSpace(sentence) == answers[index] {
			countCorrect++
			fmt.Println("Good job!")
		} else {
			countWrong++
			fmt.Println("Keep Trying...")
		}
	}

	fmt.Printf("Quiz complete! you scored %d/%d correct, which means you missed %d\n", countCorrect, len(questions), countWrong)

	return nil
}

package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func mainErr() error {
	timer := time.NewTimer(2 * time.Second)

	quit := make(chan bool)
	errc := make(chan error)
	done := make(chan error)
	defer close(quit)
	defer close(errc)
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-timer.C:
				fmt.Println("Sorry, time is up!")
				errc <- errors.New("timer finished")
			case <-errc:
				return
			}
		}
	}()

	for err := range errc {
		if err != nil {
			return err
		}
	}

	quizFile := flag.String("in", "./problems.csv", "Path to a quiz file to load with data. In the csv format of question,answer, e.g. 1+2,3")
	shouldShuffle := flag.Bool("shuffle", false, "Flag to indicate if you should shuffle quiz qs between runs.")
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

	if *shouldShuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(parsed), func(i, j int) { parsed[i], parsed[j] = parsed[j], parsed[i] })
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
		delimiter := "next"
		if index == 0 {
			delimiter = "first"
		}
		fmt.Printf("Here is your %s question. %s\n", delimiter, question)
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

	countDidNotAnswer := len(questions) - countCorrect - countWrong
	fmt.Printf("Quiz complete! you scored %d/%d correct, which means you got %d wrong and did not answer %d questions\n", countCorrect, len(questions), countWrong, countDidNotAnswer)

	return nil
}

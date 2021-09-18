package main

import (
	"bufio"
	"encoding/csv"
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
	quizFile := flag.String("in", "./problems.csv", "Path to a quiz file to load with data. In the csv format of question,answer, e.g. 1+2,3")
	shouldShuffle := flag.Bool("shuffle", false, "Flag to indicate if you should shuffle quiz qs between runs.")
	timeLimit := flag.Int("limit", 30, "teh time limit for the quiz in seconds")
	flag.Parse()

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

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

	questions, answers := makeQandA(parsed)

	buf := bufio.NewReader(os.Stdin)

	fmt.Printf("Welcome to the quiz game. Your questions are timed with limit %d seconds\n", *timeLimit)

	var countCorrect int
	var countWrong int

	// The problemloop here acts like a goto in ending the loop that we're in
problemLoop:
	for index, question := range questions {

		fmt.Printf("Here is your #%d question. %s\n", index+1, question)
		fmt.Print(">")

		answerCh := make(chan string)
		errCh := make(chan error)
		go func() {
			sentence, err := buf.ReadString('\n')
			if err != nil {
				errCh <- err
			}
			answerCh <- strings.TrimSpace(sentence)
		}()

		select {
		case <-timer.C:
			fmt.Println("Sorry time is up.")
			break problemLoop
		case err := <-errCh:
			if err != nil {
				return err
			}
		case answer := <-answerCh:
			fmt.Printf("You guessed: %s. The right answer was %s\n", answer, answers[index])
			if answer == answers[index] {
				countCorrect++
				fmt.Println("Good job!")
			} else {
				countWrong++
				fmt.Println("Keep Trying...")
			}
		}
	}

	countDidNotAnswer := len(questions) - countCorrect - countWrong
	fmt.Printf("Quiz complete! you scored %d/%d correct. You got %d wrong and did not answer %d questions\n", countCorrect, len(questions), countWrong, countDidNotAnswer)

	return nil
}

func makeQandA(parsed [][]string) ([]string, []string) {
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
	return questions, answers
}

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	// TimeToAnswer is the amount of time quiz taker has to answer, in seconds
	TimeToAnswer time.Duration = time.Duration(10)
)

// CSVFileReader takes a file name as input and returns a csv.Reader
func CSVFileReader(filename string) *csv.Reader {
	fr, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return csv.NewReader(strings.NewReader(string(fr)))
}

func getInput(ch chan<- string) {
	input := bufio.NewReader(os.Stdin)
	in, _ := input.ReadString('\n')
	in = strings.Replace(in, "\r\n", "", -1)
	in = strings.TrimSpace(in)
	ch <- in
	return
}
func main() {
	r := CSVFileReader("quiz.csv")
	r.TrimLeadingSpace = true
	qa, err := r.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	correct := 0
	guesses := make(map[string]string)
	input := make(chan string)
	guess := make(chan string)
	quiz := make(map[string]string)
	gameover := false
	for _, a := range qa {
		quiz[a[0]] = a[1]
	}

	for q, a := range quiz {
		fmt.Printf("What is the answer to %v?\n", q)
		go getInput(input)
		go func() {
			select {
			case i := <-input:
				guess <- i
				break
			case <-time.After(TimeToAnswer * time.Second):
				guess <- ""
				gameover = true
				break
			}
		}()
		if gameover {
			break
		}
		guesses[q] = <-guess
		if strings.ToLower(guesses[q]) == strings.ToLower(a) {
			correct = correct + 1
		}
	}
	fmt.Println("You answered", correct, "questions correctly out of", len(quiz))
}

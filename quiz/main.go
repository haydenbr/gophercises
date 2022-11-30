package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type QuizQuestion struct {
	Question  string
	Delimiter string
	Answer    string
	Response  string
}

type ResponseReader struct {
	Reader        *bufio.Reader
	Response      chan string
	ResponseError chan error
}

func NewResponseReader() *ResponseReader {
	responseReader := new(ResponseReader)

	responseReader.Reader = bufio.NewReader(os.Stdin)
	responseReader.Response = make(chan string)
	responseReader.ResponseError = make(chan error)

	return responseReader
}

func (r *ResponseReader) ReadResponse() {
	input, inputError := r.Reader.ReadString('\n')

	if inputError != nil {
		r.ResponseError <- inputError
	} else {
		r.Response <- strings.ToLower(strings.TrimSpace(input))
	}
}

func (q *QuizQuestion) IsCorrect() bool {
	return q.Answer == q.Response
}

func (q *QuizQuestion) DisplayQuestion() {
	fmt.Print(q.Question, q.Delimiter)
}

func main() {
	fileName := flag.String("csv", "problems-short.csv", "a csv file in the format of 'question,delimiter<optional>,answer'")
	limit := flag.Int64("limit", 30, "the time limit for the quiz, in seconds")
	flag.Parse()

	quizQuestions := parseCsv(*fileName)
	responseReader := NewResponseReader()
	limitTimer := time.NewTimer(time.Duration(*limit) * time.Second)

problemLoop:
	for _, q := range quizQuestions {
		q.DisplayQuestion()

		go responseReader.ReadResponse()

		select {
		case response := <-responseReader.Response:
			q.Response = response
		case responseError := <-responseReader.ResponseError:
			log.Fatalln(responseError)
		case <-limitTimer.C:
			fmt.Println("\nTime is up!")
			break problemLoop
		}
	}

	limitTimer.Stop()

	finalScore := computeScore(quizQuestions)

	fmt.Printf("Your score: %.2f%%\n", finalScore*100)
}

func parseCsv(fileName string) []*QuizQuestion {
	fileContents, readFileError := os.ReadFile(fileName)

	if readFileError != nil {
		log.Fatalln(readFileError)
	}

	quizQuestions := make([]*QuizQuestion, 0)
	csvReader := csv.NewReader(strings.NewReader(string(fileContents)))

	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln(err)
		}

		quizQuestions = append(quizQuestions, &QuizQuestion{
			Question:  record[0],
			Delimiter: record[1],
			Answer:    record[2],
		})
	}

	return quizQuestions
}

func computeScore(quizQuestions []*QuizQuestion) float64 {
	if len(quizQuestions) == 0 {
		return float64(0)
	}

	correctResponses := fold(quizQuestions, 0, func(total int, question *QuizQuestion, _ int) int {
		if question.IsCorrect() {
			return total + 1
		} else {
			return total
		}
	})

	return float64(correctResponses) / float64(len(quizQuestions))
}

func fold[T any, R any](s []T, initial R, f func(R, T, int) R) R {
	result := initial

	for i, v := range s {
		result = f(result, v, i)
	}

	return result
}

// source: https://golang.google.cn/src/os/file.go?s=20888:20930#L662
// copied here to walk through with a debugger to see how it works
func ReadFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var size int
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF

	// If a file claims a small size, read at least 512 bytes.
	// In particular, files in Linux's /proc claim size 0 but
	// then do not work right if read in small pieces,
	// so an initial read of 1 byte would not work correctly.
	if size < 512 {
		size = 512
	}

	data := make([]byte, size, size)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		// n, err := f.Read(data[len(data):cap(data)])
		_, err := f.Read(data)
		// data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}
	}
}

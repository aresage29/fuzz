package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// read wordlist
	dat, err := os.ReadFile("./wordlist.txt")
	if err != nil {
		fmt.Println("Some error: ", err)
	}

	words := strings.Split(string(dat), "\r\n")

	// create workers
	wordChan := make(chan string, 10)
	doneChan := make(chan bool, 400)
	for i := 1; i < 400; i++ {
		go worker(i, wordChan, doneChan)
	}

	// fill word channel with all the words we want to fuzz
	fmt.Println("Filling wordChan")
	for _, word := range words {
		wordChan <- word
	}

	// check that all the workers are done before ending
	for i := 1; i < 400; i++ {
		<-doneChan
	}
}

func worker(id int, wordChan chan string, doneChan chan bool) {
out:
	for {
		select {

		case url := <-wordChan:
			url = fmt.Sprintf("https://emile.space/%s", url)
			request(id, url)

		case <-time.After(3 * time.Second):
			fmt.Printf("worker %d couldn't get a new url after 3 seconds, quitting\n", id)
			break out

		}
	}
	doneChan <- true
}

func request(id int, url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Some error: ", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Some error: ", err)
	}

	if len(body) > 146 {
		fmt.Printf("[%d] req to url %s (%d)\n", id, url, len(body))
	}
	fmt.Printf("[%d] req to url %s (%d)\n", id, url, len(body))

}

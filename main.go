package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

const BaseUrl = "http://www.69shu.com/"

func printUsage() {
	fmt.Println("Usage: sixtynine <book-id> \"<book-name>\" <optional:kindle-email>")
	fmt.Println("  book-id       ID associated to book in 69shu.com")
	fmt.Println("  book-name     English name for the book")
	fmt.Println("  kindle-email  (optional) Kindle email to send ebook to")
}

func processBookUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tt := tokenizer.Next()
		switch {
		case tt == html.ErrorToken:
			// End of document.
			return
		case tt == html.StartTagToken:
			token := tokenizer.Token()
			fmt.Println(token.Data)
		}
	}
}

// Download and parse chapter URLs.
func processChapterUrl(url string, ch chan string, chFinished chan bool) {
}

func main() {
	args := os.Args[1:]

	numOfArgs := len(args)
	if numOfArgs < 2 || numOfArgs > 3 {
		printUsage()
		os.Exit(1)
	}

	bookId := args[0]
	bookName := args[1]
	url := BaseUrl + bookId + "/"
	fmt.Printf("Processing %s with name: \"%s\"\n", url, bookName)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	processBookUrl(url)

	if numOfArgs == 3 {
		kindleEmail := args[2]
		fmt.Printf("Sending \"%s\" to %s.\n", bookName, kindleEmail)
	}
}

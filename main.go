package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

const BaseUrl = "http://www.69shu.com/"

func printUsage() {
	fmt.Println("Usage: sixtynine <book-id> \"<book-name>\" <optional:kindle-email>")
	fmt.Println("  book-id       ID associated to book in 69shu.com")
	fmt.Println("  book-name     English name for the book")
	fmt.Println("  kindle-email  (optional) Kindle email to send ebook to")
}

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
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

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	fmt.Printf(buf.String())
}

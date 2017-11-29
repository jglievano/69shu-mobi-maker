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

func getHref(node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func nodeIsOfClass(node *html.Node, className string) bool {
	for _, attr := range node.Attr {
		if attr.Key == "class" {
			if attr.Val == className {
				return true
			}
		}
	}
	return false
}

// Gets FIRST child in parent node with class.
func getChildOfClass(parent *html.Node, tag string, className string) *html.Node {
	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			if nodeIsOfClass(c, className) {
				return c
			}
		}
	}
	return nil
}

// Gets FIRST child in parent node with tag.
func getChildWithTag(parent *html.Node, tag string) *html.Node {
	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			return c
		}
	}
	return nil
}

func parseBookHtmlNode(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "div" {
		if nodeIsOfClass(node, "mu_contain") {
			// Check if mu_h1 class element exists. Which means is not the complete chapter list.
			if getChildOfClass(node, "div", "mu_h1") == nil {
				list := getChildOfClass(node, "ul", "mulu_list")
				if list != nil {
					for li := list.FirstChild; li != nil; li = li.NextSibling {
						if li.Type == html.ElementNode {
							link := getChildWithTag(li, "a")
							for text := link.FirstChild; text != nil; text = text.NextSibling {
								if text.Type == html.TextNode {
									fmt.Println(text.Data)
								}
							}
							fmt.Println(getHref(link))
						}
					}
				}
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		parseBookHtmlNode(c)
	}
}

func processBookUrl(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	parseBookHtmlNode(doc)
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

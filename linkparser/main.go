package main

import (
	"fmt"
	"os"
	"unicode"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func isEmpty(s string) bool {
	isEmpty := true

	for _, r := range s {
		if !isEmpty {
			break
		}

		isEmpty = isEmpty && unicode.IsSpace(r)
	}

	return isEmpty
}

func getTextContent(node *html.Node) string {
	textContent := ""

	if node.FirstChild != nil {
		textContent += getTextContent(node.FirstChild)
	}

	if node.Type == html.TextNode && !isEmpty(node.Data) {
		textContent += node.Data
	}

	if node.NextSibling != nil {
		textContent += getTextContent(node.NextSibling)
	}

	return textContent
}

func walkHtmlTree(node *html.Node, links []Link) {

}

func main() {
	// fileBytes, readFileErr := ioutil.ReadFile("ex1.html")
	file, openFileErr := os.Open("ex1.html")

	if openFileErr != nil {
		panic(openFileErr)
	}

	rootNode, htmlParseErr := html.Parse(file)

	if htmlParseErr != nil {
		panic(htmlParseErr)
	}

	fmt.Println(getTextContent(rootNode))
}

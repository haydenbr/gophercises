package main

import (
	"fmt"
	"os"
	"strings"
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

func getHref(node *html.Node) string {
	href := ""

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}

	return href
}

func getTextContent(node *html.Node) string {
	textContent := ""

	if node.FirstChild != nil {
		textContent += getTextContent(node.FirstChild)
	}

	if node.Type == html.TextNode && !isEmpty(node.Data) {
		textContent += strings.TrimSpace(node.Data)
	}

	if node.NextSibling != nil {
		textContent += getTextContent(node.NextSibling)
	}

	return textContent
}

func walkHtmlTree(node *html.Node) []Link {
	links := make([]Link, 0, 10)

	if node.FirstChild != nil {
		links = append(links, walkHtmlTree(node.FirstChild)...)
	}

	if node.Type == html.ElementNode && node.Data == "a" {
		links = append(links, Link{Href: getHref(node), Text: getTextContent(node.FirstChild)})
	}

	if node.NextSibling != nil {
		links = append(links, walkHtmlTree(node.NextSibling)...)
	}

	return links
}

func main() {
	file, openFileErr := os.Open("ex4.html")

	if openFileErr != nil {
		panic(openFileErr)
	}

	rootNode, htmlParseErr := html.Parse(file)

	if htmlParseErr != nil {
		panic(htmlParseErr)
	}

	fmt.Println(walkHtmlTree(rootNode))
}

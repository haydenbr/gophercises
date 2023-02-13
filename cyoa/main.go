package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type StoryHandler struct {
	Story         *Story
	RootTemplate  *template.Template
	StoryTemplate *template.Template
}

const rootTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.intro.Title}}</title>
	</head>
	<body>
		<h1>{{.intro.Title}}</h1>
		<p>A choose your own adventure story.<p>
		<img src="https://go.dev/images/gophers/ladder.svg" />
		<a href="/arc/intro">Get Started!</a>
	</body>
<html>
`
const storyTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		{{range .Story}}
			<p>{{.}}</p>
		{{end}}
		<br/>
		{{range .Options}}
			<a href="/arc/{{.Arc}}" style="display: block">{{.Text}}</a>
		{{else}}
			<a href="/">return to start</a>
		{{end}}
	</body>
</html>`

type StoryArcOption struct {
	Text string
	Arc  string
}

type StoryArc struct {
	Title   string
	Story   []string
	Options []StoryArcOption
}

type Story map[string]StoryArc

func (h *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		h.RootTemplate.Execute(w, h.Story)
		return
	}

	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	isArcPath := pathParts != nil && len(pathParts) >= 2 && pathParts[0] == "arc"

	if isArcPath {
		requestedArc := pathParts[1]
		storyArc, hasStoryArc := (*h.Story)[requestedArc]

		if hasStoryArc {
			h.StoryTemplate.Execute(w, storyArc)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func getStory(fileName string) *Story {
	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	story := make(Story)
	json.Unmarshal(content, &story)

	return &story
}

func NewStoryHandler(fileName string, rootTemplateRaw string, storyTemplateRaw string) *StoryHandler {
	storyHandler := StoryHandler{}
	storyHandler.Story = getStory("story.json")

	rootTemplate, _ := template.New("root").Parse(rootTemplateRaw)
	storyTemplate, _ := template.New("story").Parse(storyTemplate)

	storyHandler.RootTemplate = rootTemplate
	storyHandler.StoryTemplate = storyTemplate

	return &storyHandler
}

func main() {
	storyHandler := NewStoryHandler("story.json", rootTemplate, storyTemplate)

	storyMux := http.NewServeMux()
	storyMux.Handle("/", storyHandler)
	fmt.Println("starting server on port 8080")
	http.ListenAndServe(":8080", storyMux)
}

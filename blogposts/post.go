package blogposts

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

type Post struct {
	Title       string
	Description string
	Date        time.Time
	Tags        map[string]struct{}
	Url         string
	Body        template.HTML
}

const (
	titleSeparator = "Title: "
	descSeparator  = "Description: "
	dateSeparator  = "Date: "
	tagsSeparator  = "Tags: "
)

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readLines := func(separator string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), separator)
	}

	title := readLines(titleSeparator)
	desc := readLines(descSeparator)
	date := readLines(dateSeparator)
	tagSlice := strings.Split(readLines(tagsSeparator), ", ")
	url := UrlCreator(title)
	body := template.HTML(readBody(scanner))

	const dateForm = "2006-01-02"
	parsedDate, err := time.Parse(dateForm, date)
	if err != nil {
		return Post{}, nil
	}

	tags := make(map[string]struct{})
	for _, tag := range tagSlice {
		tags[fmt.Sprintf(tag)] = struct{}{}
	}

	return Post{
		Title:       title,
		Description: desc,
		Date:        parsedDate,
		Tags:        tags,
		Url:         url,
		Body:        body,
	}, nil
}

func readBody(scanner *bufio.Scanner) []byte {
	scanner.Scan()
	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}

	newBuf := buf.Bytes()
	content := bytes.TrimSpace(blackfriday.Run(newBuf))
	return content
}

func UrlCreator(title string) string {
	title = strings.ToLower(strings.Replace(title, " ", "-", -1))
	reg, err := regexp.Compile("[^a-z0-9\\-]+")
	if err != nil {
		log.Fatal(err)
	}

	url := reg.ReplaceAllString(title, "")
	return "/blog/" + url
}

package blogposts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

type Post struct {
	Title       string
	Description string
	Date        time.Time
	Tags        []string
	Body        string
}

const (
	titleSeparator = "Title: "
	descSeparator  = "Description: "
	dateSeparator  = "Date: "
	tagsSeparator  = "Tags: "
)

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	readLines := func(tag string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tag)
	}

	title := readLines(titleSeparator)
	desc := readLines(descSeparator)
	date := readLines(dateSeparator)
	tags := strings.Split(readLines(tagsSeparator), ", ")
	body := strings.TrimSuffix(readBody(scanner), "\n")

	const dateForm = "2006-01-02"
	parsedDate, err := time.Parse(dateForm, date)
	if err != nil {
		return Post{}, nil
	}

	return Post{
		Title:       title,
		Description: desc,
		Date:        parsedDate,
		Tags:        tags,
		Body:        body,
	}, nil
}

func readBody(scanner *bufio.Scanner) string {
	scanner.Scan()
	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}

	newBuf := buf.String()
	content := blackfriday.Run([]byte(newBuf))
	return string(content)
}

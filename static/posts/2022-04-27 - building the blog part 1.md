Title: Building this site from scratch - part 1
Description: How i build this blog from scratch using Go.
Date: 2022-04-27
Tags: go
---
# Building this blog from scratch - part 1

In my quest to become a developer, I decided to rewrite my website from scratch. Before it was just a static page hosted on Github Pages. I had a [bash script](https://github.com/crdpa/bsg) that converted markdown files using [Pandoc](https://pandoc.org/) and generated the website for me, including an index of posts based on he filenames. It was hacky, but did the job.

It is time to take things to the next level and since I'm learning Go, which is perfect for web related stuff, I decided to write the site from the ground up.

## Where to start?

This was not my first time developing something. I have plenty of small projects that I really like ([Musyca](https://github.com/crdpa/musyca), [Kolekti](https://github.com/crdpa/kolekti) so it didn't seem like a daunting task, but building a website is definitely something new for me and documenting some of the steps and decicions I made here seems like a cool start for the first posts.

What I want is pretty basic: a website that works as a portfolio and a place to write about things that I find interesting.

## Getting the site up and running

The first step was setting up a server and... serve:

```go
// main.go
package main

// Request function to be used in http.HandleFunc, it will direct the root address
// to index.html e run the function that will return the front page posts.
// It will also redirect crdpa.net/blog to blog.html and check if it has a tag (?tag=)
// in the address.
func httpFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
		case "/", "/index.html":
		executeTemplate(w, "index.html", blogposts.FrontPage(posts))
		return
		case "/blog":
		tag = r.URL.Query().Get("tag")
		executeTemplate(w, "blog.html", blogposts.Archive(posts, tag))
	}
}

func executeTemplate(w http.ResponseWriter, templ string, content interface{}) {
	templates := template.Must(template.ParseGlob("./static/*.html"))
	err := templates.ExecuteTemplate(w, templ, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Here we serve the files we need
	stylesheets := http.FileServer(http.Dir("./static/css/"))
	http.Handle("/css/", http.StripPrefix("/css/", stylesheets))
	images := http.FileServer(http.Dir("./static/img/"))
	http.Handle("/img/", http.StripPrefix("/img/", images))

	// Running the request functions
	http.HandleFunc("/", httpFunc)
	http.HandleFunc("/blog", httpFunc)

	port := ":8000"
	log.Println("Server is running on port" + port)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
```

Done! Just run the app and type *localhost:8000* in Firefox.

## Blog posts

The blogposts will be markdown files that i will convert to HTML. I found this neat markdown processor called [Blackfriday](https://github.com/russross/blackfriday/tree/v2) and decided to give it a go.

The post structure is like this:

```go
type Post struct {
    Title       string
    Description string
    Date        time.Time
    Tags        []string
    Body        template.HTML
}
```

The field *Date* is time.Time so I can order posts by date in descending order and I can change the way the date is show in a [lot of ways](https://yourbasic.org/golang/format-parse-string-time-date-example/).

Markdown files does not have any metadata and I need something similar. My solution was to put the metadata first before the content, like this:

```
Title: Título do post
Description: Descrição do post
Date: 2006-01-02
Tags: go, website
---
Conteúdo do post.
```

The application will read the file and assign each line to it's expected Post field. In the end it will jump the "---" and read the content, convert to HTML and assign to Body.

Here is what i came up with:

```go
// post.go
package blogposts

const (
	titleSeparator = "Title: "
	descSeparator  = "Description: "
	dateSeparator  = "Date: "
	tagsSeparator  = "Tags: "
)

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	// Function to read the lines and trim the separator
	readLines := func(separator string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), separator)
	}

	title := readLines(titleSeparator)
	desc := readLines(descSeparator)
	date := readLines(dateSeparator)
	tags := strings.Split(readLines(tagsSeparator), ", ")
	body := template.HTML(readBody(scanner))

	// Parsing the date to time.Time
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

// function to read the content of the post and convert to HTML
func readBody(scanner *bufio.Scanner) []byte {
	scanner.Scan()
	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}

	newBuf := buf.Bytes()
	// Blackfriday doing it's thing
	// https://github.com/russross/blackfriday/tree/v2
	content := bytes.TrimSpace(blackfriday.Run(newBuf))
	return content
}
```
There are a lot of ways to read files from the filesystem in Go. I won't talk about this here because it will become a wall of text. but i made a function called NewPostsFromFS that returns a slice of Post structs ordered by date. Here is how i sorted:

```go
// blogposts.go
sort.Slice(posts, func(i, j int) bool {
	return posts[i].Date.After(posts[j].Date)
})
```

## The front page

For the front page i decided to show the last five blog posts. Remember the FrontPage function up there in the beginning of the post?

It receives the slice of Post structs and make a copy of the first posts.

```go
// blogposts.go
func FrontPage(posts []Post) []Post {
	numPosts := 5
	// show the all posts if the number of posts is less then 5
	if len(posts) < 5 {
		numPosts = len(posts)
	}
	fpPostList := make([]Post, numPosts)

	copy(fpPostList, posts[:numPosts])

	return fpPostList
}
```

This is all for the first part. In the next part I will explain how I implemented tag filtering in the blog archive page and some changes i had to make along the way.

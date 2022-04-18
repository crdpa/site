package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/crdpa/site/blogposts"
)

var (
	templates = template.Must(template.ParseGlob("./static/*.html"))
	posts     []blogposts.Post
)

func httpFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/", "/index.html":
		executeTemplate(w, "index.html", blogposts.FrontPage(posts))
		return
	}

	http.NotFound(w, r)
}

func executeTemplate(w http.ResponseWriter, templ string, content interface{}) {
	err := templates.ExecuteTemplate(w, templ, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	fsys := os.DirFS("./static/posts/")
	var err error
	posts, err = blogposts.NewPostsFromFS(fsys)
	if err != nil {
		log.Fatal(err)
	}

	stylesheets := http.FileServer(http.Dir("./static/css/"))
	http.Handle("/css/", http.StripPrefix("/css/", stylesheets))
	images := http.FileServer(http.Dir("./static/img/"))
	http.Handle("/img/", http.StripPrefix("/img/", images))

	http.HandleFunc("/", httpFunc)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

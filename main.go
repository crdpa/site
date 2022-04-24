package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/crdpa/site/blogposts"
)

var (
	posts   []blogposts.Post
	curPost blogposts.Post
	tag     string
)

func httpFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/", "/index.html":
		log.Println("i'm in case /")
		executeTemplate(w, "index.html", blogposts.FrontPage(posts))
		return
	case "/blog":
		log.Println("i'm in case /blog")
		tag = r.URL.Query().Get("tag")
		executeTemplate(w, "blog.html", blogposts.Archive(posts, tag))
		return
	case curPost.Url:
		log.Println("i'm in case curPost")
		executeTemplate(w, "post.html", curPost)
		return
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
	fsys := os.DirFS("./static/posts/")
	var err error
	posts, err = blogposts.NewPostsFromFS(fsys, tag)
	if err != nil {
		log.Fatal(err)
	}

	stylesheets := http.FileServer(http.Dir("./static/css/"))
	http.Handle("/css/", http.StripPrefix("/css/", stylesheets))
	images := http.FileServer(http.Dir("./static/img/"))
	http.Handle("/img/", http.StripPrefix("/img/", images))

	http.HandleFunc("/", httpFunc)
	http.HandleFunc("/blog", httpFunc)

	for _, post := range posts {
		curPost = post
		http.HandleFunc(post.Url, httpFunc)
	}

	port := ":8000"
	log.Println("Server is running on port" + port)

	log.Fatal(http.ListenAndServe(port, nil))
}

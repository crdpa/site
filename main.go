package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/crdpa/site/blogposts"
)

var (
	posts []blogposts.Post
	tag   string
)

func httpFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/", "/index.html":
		executeTemplate(w, "index.html", blogposts.FrontPage(posts))
		return
	case "/blog":
		tag = r.URL.Query().Get("tag")
		executeTemplate(w, "blog.html", blogposts.BlogArchive(posts, tag))
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

func makePostHandler(post blogposts.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		executeTemplate(w, "post.html", post)
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
	http.HandleFunc("/blog", httpFunc)

	for _, post := range posts {
		http.HandleFunc(post.Url, makePostHandler(post))
	}

	//port := "8000"
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Println("Server is running on port" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/crdpa/site/blogposts"
)

type blog struct {
	posts []blogposts.Post
	tag   string
}

func (bl *blog) httpFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/", "/index.html":
		executeTemplate(w, "index.html", blogposts.FrontPage(bl.posts))
		return
	case "/blog":
		bl.tag = r.URL.Query().Get("tag")
		executeTemplate(w, "blog.html", blogposts.BlogArchive(bl.posts, bl.tag))
		return
	default:
		http.NotFound(w, r)
		return
	}
}

func executeTemplate(w http.ResponseWriter, templ string, content interface{}) {
	templates := template.Must(template.ParseGlob("./ui/html/*.html"))
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
	var (
		posts []blogposts.Post
		tag   string
	)
	deploy := flag.Bool("deploy", false, "get environment $PORT")
	flag.Parse()

	bl := &blog{
		posts: posts,
		tag:   tag,
	}

	fsys := os.DirFS("./posts/")
	var err error
	bl.posts, err = blogposts.NewPostsFromFS(fsys)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", bl.httpFunc)
	mux.HandleFunc("/blog", bl.httpFunc)

	for _, post := range bl.posts {
		mux.HandleFunc(post.Url, makePostHandler(post))
	}

	port := "8000"
	if *deploy == true {
		port = os.Getenv("PORT")
		if port == "" {
			log.Fatal("$PORT must be set")
		}
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server is running on port " + port)
	log.Fatal(srv.ListenAndServe())
}

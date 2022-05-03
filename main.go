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
	default:
		http.NotFound(w, r)
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
	deploy := flag.Bool("deploy", false, "get environment $PORT")
	flag.Parse()

	fsys := os.DirFS("./static/posts/")
	var err error
	posts, err = blogposts.NewPostsFromFS(fsys)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	stylesheets := http.FileServer(http.Dir("./static/css/"))
	mux.Handle("/css/", http.StripPrefix("/css/", stylesheets))
	images := http.FileServer(http.Dir("./static/img/"))
	mux.Handle("/img/", http.StripPrefix("/img/", images))
	files := http.FileServer(http.Dir("./static/files/"))
	mux.Handle("/files/", http.StripPrefix("/files/", files))

	mux.HandleFunc("/", httpFunc)
	mux.HandleFunc("/blog", httpFunc)

	for _, post := range posts {
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

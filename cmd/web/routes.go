package main

import (
	"github.com/crdpa/site/blogposts"
	"html/template"
	"net/http"
)

func (bl *blog) httpFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/", "/index.html":
		bl.executeTemplate(w, "index.html", blogposts.FrontPage(bl.posts))
		return
	case "/blog":
		bl.tag = r.URL.Query().Get("tag")
		bl.executeTemplate(w, "blog.html", blogposts.BlogArchive(bl.posts, bl.tag))
		return
	default:
		http.NotFound(w, r)
		return
	}
}

func (bl *blog) executeTemplate(w http.ResponseWriter, templ string, content interface{}) {
	templates := template.Must(template.ParseGlob("./ui/html/*.html"))
	err := templates.ExecuteTemplate(w, templ, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (bl *blog) makePostHandler(post blogposts.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bl.executeTemplate(w, "post.html", post)
	}
}

func (bl *blog) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", bl.httpFunc)
	mux.HandleFunc("/blog", bl.httpFunc)

	for _, post := range bl.posts {
		mux.HandleFunc(post.Url, bl.makePostHandler(post))
	}

	return secureHeaders(mux)
}
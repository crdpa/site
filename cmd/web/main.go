package main

import (
	"flag"
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

	port := "8000"
	if *deploy == true {
		port = os.Getenv("PORT")
		if port == "" {
			log.Fatal("$PORT must be set")
		}
	}

	srv := &http.Server{
		Addr:         ":" + port,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      bl.routes(),
	}

	log.Println("Server is running on port " + port)
	log.Fatal(srv.ListenAndServe())
}

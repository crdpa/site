package blogposts_test

import (
	"reflect"
	"testing"
	"testing/fstest"
	"time"

	blogposts "github.com/crdpa/site/blogposts"
)

func assertPost(t *testing.T, got blogposts.Post, want blogposts.Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestNewBlogPosts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
Date: 2006-01-02
Tags: tdd, go
---
*Hello World*`
		secondBody = `Title: Post 2
Description: Description 2
Date: 2006-01-02
Tags: javascript, glue
---
Test Blog`
	)

	fs := fstest.MapFS{
		"hello-world1.md": {Data: []byte(firstBody)},
		"hello-world2.md": {Data: []byte(secondBody)},
	}

	posts, err := blogposts.NewPostsFromFS(fs)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != len(fs) {
		t.Errorf("got %d posts, want %d posts", len(posts), len(fs))
	}

	assertPost(t, posts[0], blogposts.Post{
		Title:       "Post 1",
		Description: "Description 1",
		Date:        time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		Tags:        []string{"tdd", "go"},
		Body:        `<p><em>Hello World</em></p>`,
	})
}

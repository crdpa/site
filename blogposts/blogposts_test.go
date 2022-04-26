package blogposts_test

import (
	"fmt"
	"reflect"
	"testing"
	"testing/fstest"
	"time"

	blogposts "github.com/crdpa/site/blogposts"
)

var (
	testPost1 = blogposts.Post{
		Title:       "Post 1",
		Description: "Description 1",
		Date:        time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		Tags: map[string]struct{}{
			"tdd": {},
			"go":  {},
		},
		Url:  "/blog/post-1",
		Body: `<p>Test 1</p>`,
	}
	testPost2 = blogposts.Post{
		Title:       "Post 2",
		Description: "Description 2",
		Date:        time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		Tags: map[string]struct{}{
			"javascript": {},
			"glue":       {},
		},
		Url:  "/blog/post-2",
		Body: `<p>Test 2</p>`,
	}
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
Test 1`
		secondBody = `Title: Post 2
Description: Description 2
Date: 2006-01-02
Tags: javascript, glue
---
Test 2`
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

	assertPost(t, posts[0], testPost1)

	assertPost(t, posts[1], testPost2)
}

func TestArchive(t *testing.T) {
	var tests = struct {
		post []blogposts.Post
		tag  string
		want []blogposts.Post
	}{
		post: []blogposts.Post{testPost1, testPost2},
		tag:  "tdd",
		want: []blogposts.Post{testPost1},
	}

	testname := fmt.Sprintf("%+v, %s, %+v", tests.post, tests.tag, tests.want)
	t.Run(testname, func(t *testing.T) {
		got := blogposts.Archive(tests.post, tests.tag)
		if !reflect.DeepEqual(got, tests.want) {
			t.Errorf("got %+v, want %+v", got, tests.want)
		}
	})
}

func TestUrlCreator(t *testing.T) {
	var tests = []struct {
		title string
		want  string
	}{
		{"Primeiro link 01!", "/blog/primeiro-link-01"},
		{"Segundo link - 2@", "/blog/segundo-link---2"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s, %s", tt.title, tt.want)
		t.Run(testname, func(t *testing.T) {
			ans := blogposts.UrlCreator(tt.title)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}

package blogposts

import (
	"io/fs"
	"sort"
)

func NewPostsFromFS(filesystem fs.FS, tag string) ([]Post, error) {
	dir, err := fs.ReadDir(filesystem, ".")
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, f := range dir {
		post, err := getPost(filesystem, f.Name())
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	// sort slice of posts by date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

func getPost(filesystem fs.FS, fileName string) (Post, error) {
	postFile, err := filesystem.Open(fileName)
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()

	return newPost(postFile)
}

func FrontPage(posts []Post) []Post {
	numPosts := 5
	if len(posts) < 5 {
		numPosts = len(posts)
	}
	fpPostList := make([]Post, numPosts)

	copy(fpPostList, posts[:numPosts])

	return fpPostList
}

func Archive(posts []Post, tag string) []Post {
	if tag == "" {
		return posts
	}

	var filterPosts []Post
	for _, v1 := range posts {
		for _, v2 := range v1.Tags {
			if v2 == tag {
				filterPosts = append(filterPosts, v1)
			}
		}
	}
	return filterPosts
}

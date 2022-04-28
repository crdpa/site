package blogposts

import (
	"io/fs"
	"sort"
)

type Archive struct {
	Posts []Post
	Tags  []string
}

func NewPostsFromFS(filesystem fs.FS) ([]Post, error) {
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

func tagList(posts []Post) []string {
	tagsMap := make(map[string]struct{})
	for _, post := range posts {
		for key := range post.Tags {
			tagsMap[key] = struct{}{}
		}
	}

	var tagsSlice []string
	for key := range tagsMap {
		tagsSlice = append(tagsSlice, key)
	}

	sort.Strings(tagsSlice)
	return tagsSlice
}

func BlogArchive(posts []Post, tag string) Archive {
	allTags := tagList(posts)

	if tag == "" {
		return Archive{
			Posts: posts,
			Tags:  allTags,
		}
	}

	var filterPosts []Post
	for _, post := range posts {
		_, has := post.Tags[tag]
		if !has {
			continue
		}

		filterPosts = append(filterPosts, post)
	}

	return Archive{
		Posts: filterPosts,
		Tags:  allTags,
	}
}

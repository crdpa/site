Title: Building this site from scratch - part 2
Description: Final steps and adjustments of the site.
Date: 2022-04-28
Tags: go
---
# Building this blog from scratch - part 2

Continuing from the last post, now I'll show some things I wanted to implement and the problems that arised with it.

## URLs

To implement the url for the posts, I decided to add a new field in the post struct:

```go
// post.go
type Post struct {
	Title       string
	Description string
	Date        time.Time
	Tags        map[string]struct{}
	Url         string
	Body        template.HTML
}
```

You can see here that the type of the *Tags* field is now a map. There is a valid reason for that and I'll get to it later.

Back to the URLs, here is the function to get a URL based on the title of the post:

```go
//post.go
func UrlCreator(title string) string {
	// replace spaces with "-"
	title = strings.ToLower(strings.Replace(title, " ", "-", -1))
	reg, err := regexp.Compile("[^a-z0-9\\-]+")
	if err != nil {
		log.Fatal(err)
	}

	url := reg.ReplaceAllString(title, "")
	// if the title has "-" with spaces around, it will become "---"
	// it's ugly, so let's fix that
	url = strings.Replace(url, "---", "-", -1)
	return "/blog/" + url
}
```

You can check the resulting URL in your address bar right now.

Now we need to implement the routes with the posts URLs. Let's make a new function and range over []Post:

```go
//main.go
func makePostHandler(post blogposts.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		executeTemplate(w, "post.html", post)
	}
}

func main() {
//...
	http.HandleFunc("/", httpFunc)
	http.HandleFunc("/blog", httpFunc)

	// range over the slice of Post
	for _, post := range posts {
		http.HandleFunc(post.Url, makePostHandler(post))
	}

```

All done!

## Tags

I wanted tags for my blog posts so the users could filter posts by subject. Not that this site will have a lot of posts, but why not? One more thing to learn.

First things first: why map[string]struct{} for tags?

The internal design of maps in Go is highly optimized for performance and memory management. An empty struct (struct{}) has no fields and cannot hold any pointers so it does not require memory to represent it. If your map will have thousands of entries it will need less memory. It is not my case, but this is a simple optimization and there is no reason to not use it.

Now we need to change the function *newPost* to use this new type:

```go
// post.go
func newPost(postFile io.Reader) (Post, error) {
	// ...
	// this is the new URL creator function
	url := UrlCreator(title)

	// put tags in a strings slice
	tagsSlice := strings.Split(readLines(tagsSeparator), ", ")

	tags := make(map[string]struct{})
	// range over the strings slice and add the tags to the map
	for _, tag := range tagSlice {
		tags[tag] = struct{}{}
	}
```

## Archive

For the blog archive page, i needed two values. The list of posts and the list of all tags available for filtering.

The problem is that the executeTemplate function can only receive a single value. I would have to choose between the slice of posts or the slice of tags, but there is a way around this.

It is time for a new struct:

```go
// blogposts.go
type Archive struct {
	Posts []Post
	Tags  []string
}
```

This struct has two fields: Post holds a slice of Post and Tags a slice of string. Now i can have single variable holding multiple data.

And here is the tagList and BlogArchive function that will generate the content for the blog.html page:

```go
// post.go

// create a slice containing all tags
func tagList(posts []Post) []string {
	tagsMap := make(map[string]struct{})
	// range over posts adding the tags to a map
	for _, post := range posts {
		for key := range post.Tags {
			tagsMap[key] = struct{}{}
		}
	}

	var tagsSlice []string
	// a map is unordered and we don't want the list of tags to appear in random
	// order every refresh of blog.html. Since there is no way to sort a map,
	// we need to pass the values to a slice for sorting
	for key := range tagsMap {
		tagsSlice = append(tagsSlice, key)
	}

	sort.Strings(tagsSlice)
	return tagsSlice
}

// the function returns the new Archive struct
func BlogArchive(posts []Post, tag string) Archive {
	// function to return a slice with the sorted tags
	allTags := tagList(posts)

	// if the user didn't select a tag for filtering, return all posts
	if tag == "" {
		return Archive{
			Posts: posts,
			Tags:  allTags,
		}
	}

	// if the user selects a particular tag, filter the posts
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
```

And that's all there is to it. Now we just have to adjust the HTML templates which is pretty straightforward.

Here is what we need to show everything in the blog archive page:

```html
<h1>Blog</h1>
<div class='tags'>
  <p>Tags:
  <!-- range over the Tags field of the Archive struct -->
  {{range .Tags}}
    <a href='/blog?tag={{.}}'>{{.}}</a>
  {{end}}
  </p>
</div>
<div class='archive'>
  <!-- range over the Posts field of the Archive struct -->
  {{range .Posts}}
    {{if ne .Title ""}}
      <div class='archive-item'>
        <p>{{.Date.Format "Jan 02, 2006"}} -
        <a href='{{.Url}}'>{{.Title}}</a></p>
        {{range $key, $value := .Tags}}<a href='/blog?tag={{$key}}'><p>{{$key}}</p></a>&nbsp{{end}}
      </div>
    {{end}}
  {{end}}
</div>
```

And now our blog is done.

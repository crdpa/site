Title: Date manipulation with CSS and Go
Description: How to display date formats in your website based on screen size
Date: 2022-05-11
Tags: go, css, tips
---
# Date manipulation with CSS and Go

There are various ways of formatting a time.Time value. The standard library gives you the ability to specify the pattern you want.

When you get the current time:

```go
fmt.Println(time.Now())
```

Here is the output:

```
Output: 2022-05-11 23:00:00 +0000 UTC m=+0.000000001
```

This is not very readable if you are building a website. Fortunately, you can do this:

```go
fmt.Println(time.Now().Format("January 2, 2006"))
```

```
Output: May 11, 2022
```

Or this:

```go
fmt.Println(time.Now().Format("2006-01-02"))
```

```
Output: 2022-05-11
```

Go uses the date "01/02 03:04:05PM '06 -0700" to format instead of the traditional "YYYY-MM-DD HH:MM:SS whatever". You can read more about the reason in the [documentation](https://pkg.go.dev/time#pkg-constants).

What is really cool is that you can format the date right in the HTML code of your website. By doing this, you can use CSS rules to display date in different formats when the screen size changes.

Maybe you want the full date when visiting the website in a big screen and the YYYY-MM-DD format when using a small screen so you don't lose information.

Here is a simple way of doing this:

```html
<div class='date-full'>
    {{.Date.Format "January 02, 2006"}}
</div>
<div class='date-small'>
    {{.Date.Format "2006-01-02"}}
</div>
```

The {{.Date.Format "2006-01-02"}} is the Go syntax to format the "Date" value which is of time.Time type. It is the [Format](https://pkg.go.dev/time#Time.Format) function we used in the Go code above.

Now we can use CSS rules to hide one and display the other based on screen size.

When the screen is big, we display the 'date-full' <div> and when the screen goes below a specified width, we reverse the condition. Here it is:

```css
.date-full {
    display: inline-block;
}

.date-small{
   display: none;
}

/* Rules for when the screen width is 800px or less */
@media only screen and (min-width: 0px) and (max-width: 800px) {
    .date-full {
        display: none;
    }

    .date-small{
        display: inline-block;
    }
}
```

And voila. You can write a mock-up site with hardcoded dates to see how it works without having to write any Go code, but the concept is pretty simple.

I use something similar on this website, but decided to hide the date when the screen is small.

Visit the [front page](https://crdpa.net) and resize the browser window. You should see the date disappearing when the screen goes below 800px width (I did this for the picture in the About section too).

I figured there is no need to have dates in the front page since it always shows the last five posts. Dates are more important in the blog archive.

I hope those tips will help you come up with cool ideas in your future projects.
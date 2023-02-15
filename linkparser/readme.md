# Exercise 4: HTML Link Parser

source: https://github.com/gophercises/link

given an html file, get all the anchor tags, strip out the href and text content into a data structure

input:
```html
<a href="/dog">
  <span>Something in a span</span>
  Text not in a span
  <b>Bold text!</b>
</a>
```

output:
```go
Link{
  Href: "/dog",
  Text: "Something in a span Text not in a span Bold text!",
}
```

- use this package https://pkg.go.dev/golang.org/x/net/html
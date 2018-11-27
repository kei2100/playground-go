package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		io.WriteString(w, `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
  </head>
  <form action="dict" method="post" enctype="application/x-www-form-urlencoded">
     text: <input type="text" name="text">
     <input type="submit" value="open dictionary">
  </form>
</html>
`)
	})

	http.HandleFunc("/dict", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		text := r.FormValue("text")
		u := url.URL{
			Scheme: "dict",
			Path:   html.UnescapeString(text),
		}
		w.Header().Set("Location", u.String())
		w.WriteHeader(302)
	})

	err := http.ListenAndServe(":8080", nil)
	log.Println(err)
}

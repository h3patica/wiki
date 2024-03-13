package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
)

const head = `<!DOCTYPE html>
<html lang="en">
<head>
<title>hdb</title>
<link rel="stylesheet" type="text/css" href="/css">
</head>
<header style="padding: 10px;border-bottom: 1px solid #fff">
<a href="/site"><img src="/img/brim2.png" height=75 width=75></a>
</header>
`

// Generates <nav> string for each page
func genNav(name string) string {
	var ret = `<nav><ol>`
	// Leveraging hdb.go to get child and parents of page
	dir, name := prune(name)
	parents := getAllParents(dir + name)
	children := getChildren(dir + name, 1)
	for _, x := range parents {
		ret += fmt.Sprintf("\n<li class='backlink'><a href='/site/%s'>%s</a>", x, end(x))
	}
	ret += "\n<li class='backlink'>" + end(name)
	if len(children) > 0 {
		ret += fmt.Sprintf("\n<li class='ghost'><a href='/site/%s'>%s</a><ol>", children[0], end(children[0]))
		for _, y := range children[1:] {
			ret += fmt.Sprintf("\n<li><a href='/site/%s'>%s</a>", y, end(y))
		}
		ret += "\n</ol>"
	}
	return ret + "\n</ol></nav>"
}

func handler(w http.ResponseWriter, r *http.Request) {
 	name := r.URL.Path[len("/site/"):]
	// serve default page (index) if attempting to go to /site/
	if name == "" { name = "index.html" }
	p, err := os.ReadFile("./site/" + name)
	if err != nil { fmt.Fprintf(w, "404")  } else {
		fmt.Fprintf(w, "%s%s\n<main>%s</main>", head, genNav("./site/"+name), p)
	}
}

func main() {
	http.HandleFunc("/css/",
		func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./css/style.css") })
	http.HandleFunc("/img/",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./img/" + r.URL.Path[len("/img/"):])
		})
	http.Handle("/", http.RedirectHandler("/site", http.StatusSeeOther))
	http.HandleFunc("/site/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

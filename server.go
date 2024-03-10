package main

import (
	"fmt"
	"os"
	"log"
	"strings"
	"net/http"
)

const head = `<!DOCTYPE html>
<html lang="en">
<head>
<title>hdb</title>
<link rel="stylesheet" type="text/css" href="/css">
</head>
<header style="padding: 10px;border-bottom: 1px solid #fff">
<h1><a href="/site">hdb</a></h1>
</header>
`

func findName(list []string, name string) int {
	for i, j := range list[0:len(list)-2] {
		fmt.Println("checking " + name + " against: " + j)
		if j == name && i == len(list)-3 {
			return i+1
		}
	}
	return -1
}

func getLinks(name string) []string {
	files, _ := os.ReadDir("./site/")
	var links []string
	for _, x := range files {
		fmt.Println("seeing if " + name + " links to " + x.Name())
		y := strings.Split(x.Name(), ".")
		idx := findName(y, name)
		if idx != -1 {
			fmt.Println("yes")
			links = append(links, y[idx])
		}
	}
	return links
}

func genNav(name string) string {
	var ret = `<nav><ol>`
	cats := strings.Split(name, ".")
	cats = cats[0:len(cats)-1]
	for i, x := range cats[0:len(cats)-1] {
		ret += "\n<li class='backlink'>" + `<a href="/site/` + strings.Join(cats[0:i+1], ".") + `.html">` + x + "</a>"
	}
	ret += "\n<li class='backlink'>" + cats[len(cats)-1]
	pruned_name := cats[len(cats)-1]
	links := getLinks(pruned_name)
	if len(links) > 0 {
		ret += "\n<li class='ghost'>" + `<a href="/site/` + strings.Join(cats, ".") + "." + links[0] + `.html">` + links[0] + "</a><ol>"
		for _, y := range links[1:] {
			ret += "\n<li style='opacity:0.5'>" + `<a href="/site/` + pruned_name + "." +  y + `.html">` + y + "</a>"
		}
		ret += "\n</ol>"
	}
	return ret + "\n</ol></nav>"
}

func handler(w http.ResponseWriter, r *http.Request) {
 	name := r.URL.Path[len("/site/"):]
	if name == "" { name = "index.html" }
	p, err := os.ReadFile("./site/" + name)
	if err != nil { fmt.Fprintf(w, "404")  } else {
		fmt.Fprintf(w, "%s%s\n<main>%s</main>", head, genNav(name), p)
	}
}

func main() {
	http.HandleFunc("/css/",
		func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./css/style.css") })
	http.Handle("/", http.RedirectHandler("/site", http.StatusSeeOther))
	http.HandleFunc("/site/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

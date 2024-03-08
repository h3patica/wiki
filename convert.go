package main

import (
    "os"
    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/html"
    "github.com/gomarkdown/markdown/parser"
)

func main() {
    if len(os.Args) != 3 {  panic("Must provide markdown file to parse to html and a title.") }
    fi, err := os.ReadFile(os.Args[1])
    if err != nil { panic(err) }
    var name = "html/" + os.Args[2] + ".html"
    err = os.WriteFile(name, mdToHTML(fi), 0644)
    if err != nil { panic(err) }
}

func mdToHTML(md []byte) []byte {
    // create markdown parser with extensions
    extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
    p := parser.NewWithExtensions(extensions)
    doc := p.Parse(md)

    // create HTML renderer with extensions
    htmlFlags := html.CommonFlags | html.HrefTargetBlank
    opts := html.RendererOptions{Flags: htmlFlags}
    renderer := html.NewRenderer(opts)

    return markdown.Render(doc, renderer)
}

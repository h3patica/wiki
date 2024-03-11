package main

import (
	"os"
	"slices"
	"strings"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

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

func removeExts(files *[]os.DirEntry, ext string) []string  {
	var ret []string
	for _, f := range *files {
		name, _ := strings.CutSuffix(f.Name(), ext)
		ret = append(ret, name)
	}
	return ret
}

func main() {
	mdDir, _ := os.ReadDir("./resources/")
	mdFiles := removeExts(&mdDir, ".md")
	htmlDir, _ := os.ReadDir("./site/")
	htmlFiles := removeExts(&htmlDir, ".html")
	for _, x := range htmlFiles {
		if !(slices.Contains(mdFiles, x)) {
			err := os.Remove("./site/"+x+".html")
			if err != nil { panic(err) }
		}
	}
	for _, y := range mdFiles {
		if !(slices.Contains(htmlFiles, y)) {
			md, err := os.ReadFile("./resources/"+y+".md")
			if err != nil { panic(err) }
			err = os.WriteFile("./site/"+y+".html", mdToHTML(md), 0600)
			if err != nil { panic(err)  }
		}
	}
}


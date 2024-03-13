package main

import (
	"os"
	"strings"
	"syscall"
	"fmt"
	"os/signal"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/bwmarrin/discordgo"
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	content := string(m.Content[:])
	if string(content[0]) == "+" || string(content[0]) == "-" {
		lines := strings.Split(content, "\n")
		filename := strings.Trim(lines[1], "\n")
		post := strings.Join(lines[2:], "\n")
		path := "./site/" + filename + ".html"
		if string(content[0]) == "-" {
			err := os.Remove(path)
			if err != nil { s.ChannelMessageSend(m.ChannelID, "Couldn't remove post")  }
			s.ChannelMessageSend(m.ChannelID, "Deleted post " + filename)
		}
		if string(content[0]) == "+" {
			err := os.WriteFile(path, mdToHTML([]byte(post)), 0600)
			if err != nil { s.ChannelMessageSend(m.ChannelID, "Couldn't create post")  }
			s.ChannelMessageSend(m.ChannelID, "Created post " + filename)
		}
	}
}

func main() {
	f, err := os.ReadFile("../token.txt")
	token := fmt.Sprintf("%s", f)
	token = strings.Trim(token, "\n ")
	if err != nil { panic(err) }
	dg, err := discordgo.New("Bot " + token)
	if err != nil { panic(err) }
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil { panic(err) }
	fmt.Println("Bot now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}


// Command firstbot is the smallest useful Starling bot.
//
// Run it:
//
//	set DISCORD_TOKEN=your-token-here     (Windows)
//	export DISCORD_TOKEN=your-token-here  (macOS/Linux)
//	go run ./examples/firstbot
//
// Then type `!ping` in any channel the bot can see.
package main

import (
	"log"
	"os"
	"strings"

	"github.com/razecrs/starling"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("set DISCORD_TOKEN first - get one from discord.com/developers/applications")
	}

	// MessageContent is privileged and must also be enabled in the developer
	// portal. Starling warns at startup if this handler cannot receive content.
	bot := starling.New(token,
		starling.WithIntents(
			starling.IntentGuilds|
				starling.IntentGuildMessages|
				starling.IntentMessageContent,
		),
		starling.WithStatus(starling.StatusOnline, starling.Playing("with starling")),
	)

	bot.Command("ping", func(m *starling.MessageCreate, args []string) {
		if _, err := m.Reply("pong"); err != nil {
			log.Printf("could not reply: %v", err)
		}
	})

	bot.Command("say", func(m *starling.MessageCreate, args []string) {
		if len(args) == 0 {
			m.Reply("say what? try `!say hello`")
			return
		}
		// Anything built from user input should suppress mentions, or someone
		// will make the bot ping @everyone for them.
		m.ReplyComplex(starling.SendData{
			Content:         strings.Join(args, " "),
			AllowedMentions: starling.NoMentions(),
		})
	})

	bot.On(func(r *starling.Ready) {
		log.Printf("online as %s", r.User.Tag())
	})

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}

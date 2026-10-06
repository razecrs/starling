// Command ping is a minimal starling bot: it replies "pong" to "ping".
//
// Run it with a bot token in the environment:
//
//	DISCORD_TOKEN=... go run ./examples/ping
package main

import (
	"github.com/razecrs/starling"
	"log"
	"os"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("set DISCORD_TOKEN")
	}

	// MessageContent is a privileged intent: switch it on for the bot in the
	// Discord developer portal, or m.Content arrives empty.
	bot := starling.New(token,
		starling.WithIntents(starling.IntentGuilds|starling.IntentGuildMessages|starling.IntentMessageContent),
		starling.WithStatus(starling.StatusOnline, starling.Playing("with starling")),
	)

	bot.On(func(r *starling.Ready) {
		log.Printf("online as %s", r.User.Tag())
	})

	bot.On(func(m *starling.MessageCreate) {
		// Without this a bot that echoes anything will talk to itself forever.
		if m.IsFromBot() {
			return
		}
		if m.Content == "ping" {
			if _, err := m.Reply("pong"); err != nil {
				log.Printf("replying: %v", err)
			}
		}
	})

	// Run blocks until Ctrl-C, reconnecting on its own if the connection drops.
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}

// Command voice joins an occupied voice channel and plays a local file. If
// the input is a video, FFmpeg uses its audio track.
//
//	VOICE_FILE=./song.mp3 DISCORD_TOKEN=... go run ./examples/voice
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/razecrs/starling"
)

func main() {
	loadDotEnv()
	token, path := os.Getenv("DISCORD_TOKEN"), os.Getenv("VOICE_FILE")
	if token == "" || path == "" {
		log.Fatal("set DISCORD_TOKEN and VOICE_FILE")
	}

	bot := starling.New(token,
		starling.WithIntents(starling.IntentGuilds|starling.IntentGuildVoiceStates),
		starling.WithStatus(starling.StatusOnline, starling.Listening("voice over starling")),
	)
	var once sync.Once
	bot.On(func(g *starling.GuildCreate) {
		channelID := requestedChannel()
		if channelID.IsZero() {
			channelID = occupiedChannel(bot, g)
		}
		if channelID.IsZero() {
			return
		}
		once.Do(func() { go play(bot, g, channelID, path) })
	})

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}

func play(bot *starling.Client, guild *starling.GuildCreate, channelID starling.Snowflake, path string) {
	status := func(format string, args ...any) {
		message := fmt.Sprintf(format, args...)
		log.Print(message)
		for _, channel := range guild.Channels {
			if channel.Type == starling.ChannelGuildText && strings.EqualFold(channel.Name, "dev-chat") {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				_, _ = bot.Send(ctx, channel.ID, message)
				cancel()
				break
			}
		}
	}

	joinCtx, cancelJoin := context.WithTimeout(context.Background(), 15*time.Second)
	connection, err := bot.ConnectVoice(joinCtx, guild.ID, channelID)
	cancelJoin()
	if err != nil {
		status("voice test failed to connect: %v", err)
		_ = bot.Close()
		return
	}
	defer connection.Close(context.Background())

	playCtx, cancelPlay := playbackContext()
	defer cancelPlay()
	provider, err := connection.PlayFile(playCtx, path)
	if err != nil {
		status("voice test failed to start FFmpeg: %v", err)
		connection.Close(context.Background())
		_ = bot.Close()
		return
	}
	status("voice test playing %q in <#%s>", path, channelID)
	if err = provider.Wait(); err != nil {
		status("voice test stopped with an error: %v", err)
	} else {
		status("voice test finished")
	}
	provider.Close()
	connection.Close(context.Background())
	_ = bot.Close()
}

func occupiedChannel(bot *starling.Client, guild *starling.GuildCreate) starling.Snowflake {
	self := bot.Self()
	for _, state := range guild.VoiceStates {
		if !state.ChannelID.IsZero() && (self == nil || state.UserID != self.ID) {
			return state.ChannelID
		}
	}
	return 0
}

func requestedChannel() starling.Snowflake {
	value := os.Getenv("VOICE_CHANNEL_ID")
	if value == "" {
		return 0
	}
	id, err := starling.ParseSnowflake(value)
	if err != nil {
		log.Fatalf("VOICE_CHANNEL_ID: %v", err)
	}
	return id
}

func playbackContext() (context.Context, context.CancelFunc) {
	seconds, _ := strconv.Atoi(os.Getenv("VOICE_SECONDS"))
	if seconds > 0 {
		return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	}
	return context.WithCancel(context.Background())
}

func loadDotEnv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if ok && os.Getenv(strings.TrimSpace(key)) == "" {
			_ = os.Setenv(strings.TrimSpace(key), strings.Trim(strings.TrimSpace(value), `"'`))
		}
	}
}

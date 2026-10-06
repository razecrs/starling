package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"

	"github.com/razecrs/starling"
)

func main() {
	logs := starling.NewStarlog("my bot",
		starling.StarlogColor(true),
		starling.StarlogSystemMedia(),
		starling.StarlogWithPets(
			starling.NewStarPet(starling.StarPetNova),
			starling.NewStarPet(starling.StarPetComet),
			starling.NewStarPet(starling.StarPetNebula),
		))
	bot := starling.New(os.Getenv("DISCORD_TOKEN"),
		starling.WithStarlog(logs))

	bot.On(func(ready *starling.Ready) {
		logs.Info("ready", "user", ready.User.Tag())
	})

	bot.Command("build", func(message *starling.MessageCreate, _ []string) {
		logs.Build("running go test")
		cmd := exec.CommandContext(context.Background(), "go", "test", "./...")
		stdout := logs.BuildWriter("go test")
		stderr := logs.Writer(slog.LevelWarn, "go test")
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			_ = stdout.Close()
			_ = stderr.Close()
			logs.Error("build failed", "err", err)
			return
		}
		_ = stdout.Close()
		_ = stderr.Close()
		logs.Info("build passed")
	})

	if err := bot.Run(); err != nil {
		logs.Error("bot stopped", "err", err)
	}
}

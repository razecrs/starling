# Messages and REST

## Common message operations

```go
message, err := bot.Send(ctx, channelID, "hello")
message, err = bot.EditMessage(ctx, channelID, message.ID,
	starling.SendData{Content: "hello again"})
err = bot.React(ctx, channelID, message.ID, starling.Emoji{Name: "👍"})
err = bot.DeleteMessage(ctx, channelID, message.ID, "cleanup")
```

Event convenience methods (`Reply`, `Send`, `React`, and `Delete`) use the event's client and IDs. Client methods accept a context and are the better fit when the operation belongs to a larger request lifecycle.

## Rich messages and mentions

```go
_, err := bot.SendComplex(ctx, channelID, starling.SendData{
	Content: "deployment ready",
	Embeds: []starling.Embed{{
		Title: "Build 42",
		Color: 0x5865F2,
	}},
	Components: []starling.Component{
		starling.ActionRow(
			starling.Button(starling.ButtonPrimary, "Deploy", "deploy:42"),
			starling.LinkButton("Logs", logsURL),
		),
	},
	AllowedMentions: starling.NoMentions(),
})
```

Discord parses mentions by default. Use `NoMentions` for echoed user input, logs, external content, and templates unless the bot intentionally mentions people or roles.

Components V2 use `ComponentsMessage`, `Container`, `Section`, `TextDisplay`, `Thumbnail`, and related builders. Polls start with `NewPoll`; set duration and layout fields before sending.

## Uploads

```go
file, err := starling.FileFromPath("report.pdf")
if err != nil {
	log.Fatal(err)
}
file.Description = "Build report"

_, err = bot.SendFiles(ctx, channelID,
	starling.SendData{Content: "report", AllowedMentions: starling.NoMentions()},
	file)
```

`FileFromPath` streams from disk while Starling builds the multipart request; it does not retain an open file. `FileFromBytes` is convenient for generated data already in memory. Give every file an extension so Discord can infer previews and content type.

## History and pagination

`Messages` accepts `MessagesQuery` for before/after/around pagination. Discord page limits still apply. Use a context deadline and advance with the oldest/newest returned ID; do not retry the same page indefinitely when a channel changes during traversal.

## Errors and retries

Starling handles Discord rate-limit buckets, global rate limits, 429 responses, and bounded retries for transient server failures. Application code should classify the final error:

```go
if _, err := bot.Send(ctx, channelID, "hello"); err != nil {
	switch {
	case starling.IsUnauthorized(err):
		log.Fatal("the token was rejected")
	case starling.IsForbidden(err):
		log.Printf("missing permission: %v", err)
	case starling.IsNotFound(err):
		log.Printf("channel was removed")
	default:
		log.Printf("send failed: %v", err)
	}
}
```

`APIError` preserves the HTTP status, Discord error code, field validation failures, and a bounded raw body. Error strings sanitise untrusted control characters. Do not blindly retry permission, validation, or not-found errors.

## Audit-log reasons

Administrative helpers that can appear in Discord's audit log accept a `reason` argument. Starling encodes it into `X-Audit-Log-Reason`. Keep it concise and do not put credentials or private user data into it.

## Raw REST

Use the escape hatch for a new Discord route while keeping Starling's HTTP transport and limiter:

```go
var out MyResponse
err := bot.Request(ctx, starling.RESTRequest{
	Method: http.MethodGet,
	Path:   "/guilds/" + guildID.String() + "/new-resource",
	Route:  "GET /guilds/" + guildID.String() + "/new-resource",
}, &out)
```

`Route` is the rate-limit bucket identity, not merely a display label. Preserve major parameters such as guild/channel/webhook IDs and replace non-major IDs consistently. `RequestRaw` returns bounded bytes for non-JSON or not-yet-modelled responses. Prefer typed helpers when one exists; the generated [REST surface](rest-surface.md) is the index.

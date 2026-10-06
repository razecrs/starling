# Moving from DiscordGo

Starling does not try to rename every Discord concept. Most structs and IDs should feel familiar; the main differences are lifecycle, handler registration, safe state snapshots, interactions, and current voice support.

## Client and handlers

DiscordGo commonly starts with `discordgo.New`, `AddHandler`, and `Open`. Starling uses one blocking lifecycle:

```go
bot := starling.New(token, starling.WithIntents(intents))
bot.On(func(m *starling.MessageCreate) {
	// ...
})
log.Fatal(bot.Run())
```

`Run` reconnects and shuts down on process signals. `Listen` is the removable equivalent of a handler registration; its returned function is idempotent. `Once` handles one event.

## IDs

Starling uses the numeric `Snowflake` type rather than plain strings. Use `ParseSnowflake` at input boundaries and `.String()` for URLs, logs, and external storage. Keeping IDs numeric makes shard routing and timestamps direct operations.

## State

```go
member, ok := bot.State.Member(guildID, userID)
permissions, err := bot.State.ChannelPermissions(channelID, userID)
```

Starling returns deep snapshots. Mutating a returned member, role slice, embed, or component cannot mutate the live cache. This costs more than returning an internal pointer and avoids application/cache data races.

The default cache includes guilds, channels, threads, members, users, roles, emoji, stickers, thread members, voice states, presences, and bounded messages. Configure it with `DefaultStateConfig` and `WithStateCache`; use manual mode only when the application must own mutation timing.

## Messages and REST

```go
message, err := bot.Send(ctx, channelID, "hello")
_, err = bot.SendComplex(ctx, channelID, starling.SendData{Embeds: embeds})
```

Every client REST method accepts a context. Administrative methods take the audit-log reason on the same call. `Request`/`RequestRaw` are the equivalent escape hatch for a route without a typed helper.

Errors are `*APIError`, with helpers for common HTTP classifications. Starling's limiter handles route buckets, global limits, 429s, and bounded transient retries internally.

## Commands and interactions

Prefix commands are built in through `Command`. Slash definitions and handlers are registered together with `Slash`, then published with `SyncCommands`. Responses, defer/edit/followup, components, modals, autocomplete, and HTTP delivery are methods on `InteractionCreate` or `Client`; no second interaction router is required.

## Voice

Starling voice uses gateway v8, current AEAD transport modes, and DAVE. `ConnectVoice` completes the handshake, `PlayFile` uses FFmpeg, `Play` accepts a custom `OpusProvider`, and `Receive` exposes incoming Opus packets.

Do not carry over assumptions from DiscordGo's legacy voice transport. Joining waits for main-gateway voice events, so call it outside the event loop or from a goroutine.

## Raw access

`OnRaw`, `Request`, `RequestRaw`, `SendGateway`, `SendGuildGateway`, and `HandleGatewayFrame` preserve the low-level path. Moving to Starling does not require giving up new Discord fields or custom transport/replay work while waiting for a typed wrapper.

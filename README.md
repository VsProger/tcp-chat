# TCP Chat Server

A multi-room chat server written in Go over raw TCP sockets, using only the standard library. Any line-oriented TCP client — `ncat`, `telnet`, `nc` — works as the client.

## Overview

The point of the project was to build a concurrent network server from the socket up, without an HTTP framework, a WebSocket library or a message broker in between. That leaves the parts that usually stay hidden: accepting connections, giving each one its own goroutine, sharing mutable room state across those goroutines without racing, and detecting when a peer disappears.

The server supports several named rooms, room ownership with moderation rights, and broadcast delivery to everyone in a room.

## Technical Approach

**Connection handling.** The accept loop stays on the main goroutine and hands each accepted `net.Conn` to its own goroutine, so a slow or silent client cannot block anyone else. Each connection goroutine owns its read loop and closes the socket via `defer` on any exit path, including read errors from an abrupt disconnect.

**Shared state.** Rooms are the shared mutable state: a client list and a kicked-user set per room. Each `ChatRoom` carries its own `sync.Mutex`, so lock contention is scoped per room rather than serialised through one global lock, and membership changes cannot interleave with a broadcast walking the same slice.

**Message dispatch.** Incoming lines route through a `MessageHandler` interface. `HandleMessage` splits commands from chat text, and the command handler dispatches on the verb. Keeping dispatch behind an interface means the transport layer in `main.go` never learns what the protocol does.

**Moderation.** Room creators may kick and ban. A ban is recorded in the room's `KickedUsers` set and checked on join, so a banned user cannot immediately reconnect into the same room.

## Technologies

Go · `net` · `sync` · `bufio` — standard library only, no external dependencies.

## Commands

| Command | Action |
|---|---|
| `/help` | list available commands |
| `/create <name>` | create a room and join it |
| `/join <name>` | join an existing room |
| `/users` | list users in the current room |
| `/kick <user>` | remove a user from the room (creator only) |
| `/ban <user>` | remove and block rejoining (creator only) |
| `/logout` | leave the room and disconnect |

Joins and departures are announced to everyone in the room.

## How to Run

```bash
git clone https://github.com/VsProger/tcp-chat.git
cd tcp-chat

go run ./cmd/server/          # listens on :8989 by default
go run ./cmd/server/ 8080     # or pass a port
```

Connect from any number of terminals:

```bash
ncat localhost 8989     # or: nc localhost 8989 / telnet localhost 8989
```

The server greets you, asks for a username, and from there `/help` lists what is available.

## Project Structure

```
├── cmd/server/main.go        accept loop, per-connection goroutines
├── internal/core/
│   ├── client.go                 connection wrapper: greeting, line reads
│   └── chatroom.go               room state guarded by a per-room mutex
├── internal/chat/chat.go     room registry and initialisation
├── internal/handlers/        command dispatch, broadcast, moderation
├── internal/utils/           helpers
└── greetMessage.txt          banner shown on connect
```

## Future Improvements

- **Replace the per-room mutex with a coordinating goroutine** owning room state behind a channel, which removes a whole class of lock-ordering mistakes as features grow.
- **Add write timeouts and a bounded send buffer per client.** A client that stops reading currently blocks the broadcasting goroutine; `SetWriteDeadline` plus a buffered outbound channel would isolate it.
- **Persist history**, so a joining user sees recent messages instead of an empty room.
- **Add tests.** Command parsing and room membership transitions are testable with `net.Pipe` and need no real sockets.

## Authors

Azat Abdirashituly · Baurzhan Saliyev · Dias Imakanov

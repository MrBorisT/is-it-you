# Is It You?

Small real-time multiplayer game prototype with a Godot client and Go WebSocket server.

## Goal

Players hide in a crowd and try to reach the finish line without revealing themselves.

## Current status

- Local Godot toy prototype
- Go WebSocket connection spike
- Server sends welcome message
- Godot client connects to server

## Planned next step

Server-authoritative multiplayer movement.

## Run server

```bash
cd server
go run .
```

## Run client

Open client/ in Godot and run the network test scene.
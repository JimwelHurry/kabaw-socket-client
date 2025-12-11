# Kabaw WebSocket React Client

This is a React + TypeScript single-page app that connects to the Kabaw WebSocket server at `ws://localhost:8080/ws`. It handles connection state, auto-reconnect, message streaming, and sending messages to a channel.

## Features
- Connects to `ws://localhost:8080/ws?username=...&channel=...` (overridable via `VITE_WS_URL`)
- Live feed with message metadata, user ID display, and channel badges
- Auto-reconnect with backoff when the server restarts
- Send messages with button or Enter key
- Connection URL preview + error surface for malformed payloads

## Prerequisites
- Node.js 18+ and npm
- Go 1.21+ to run the provided WebSocket backend

## Running the backend
From the repository root:
```bash
go run main.go
# WebSocket endpoint: ws://localhost:8080/ws
```

## Running the React client (dev)
```bash
cd client
npm install          # first run only
npm run dev -- --host
```
Open the URL printed by Vite (defaults to http://localhost:5173).

## Build for production
```bash
cd client
npm run build
npm run preview -- --host
```

## Configuration
- `VITE_WS_URL` (optional): override the WebSocket base URL. Example: `VITE_WS_URL=ws://localhost:8080/ws`.
- Channel and username can be set from the UI before hitting “Connect”.

## How to demo
1. Start the Go server (`go run main.go`).
2. Start the React dev server (`npm run dev -- --host`).
3. Click “Connect” (the UI auto-connects on load with a random username). Messages from the simulated bots in the `general` channel should start flowing immediately.
4. Type a message and send; it appears in the feed with your user ID.
5. Refresh or stop the Go server to see auto-reconnect behavior.

![Kabaw Logo](logo.jpeg)

# Kabaw Chat WebSocket Server - Technical Evaluation

This repository is a **technical evaluation document** to test the competency of a candidate to implement frontend with real-time messaging and WebSockets.

## 🎯 Candidate Task

**You are required to create your own React single-page application that will connect to the WebSocket server running on `localhost:8080`.**

The provided Go server is a complete WebSocket implementation for real-time messaging. Your task is to build a React frontend that demonstrates your ability to:

- Connect to WebSocket servers
- Handle real-time message streaming  
- Manage WebSocket connection states
- Display live chat functionality
- Handle user interactions and message sending

## 📋 Evaluation Criteria

Your React application will be evaluated on:

1. **WebSocket Integration** - Proper connection and message handling
2. **Real-time Updates** - Live message display without page refresh
3. **User Experience** - Intuitive chat interface design
4. **Error Handling** - Connection failures and reconnection logic
5. **Code Quality** - Clean, maintainable React code structure

## Features

- Real-time messaging via WebSockets
- Multiple channel support
- User identification
- Connection statistics
- Health check endpoint
- Simple web interface for testing
- **Cross-Origin Resource Sharing (CORS) support** for WebSocket connections
- **Automatic message simulation** in the "general" channel for testing
- **Console logging** for all messages, connections, and disconnections (server-side)
- **Frontend console logging** in browser developer tools (client-side)
- **Automatic user ID generation** - each connection receives a unique user ID

## 🚀 Getting Started - Running the WebSocket Server

### Prerequisites

- Go 1.21 or higher
- Internet connection (for downloading dependencies)

### Installing Go

#### On Linux (Ubuntu/Debian)
```bash
# Method 1: Using package manager (may not have latest version)
sudo apt update
sudo apt install golang-go

# Method 2: Download latest version from official site
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### On macOS
```bash
# Method 1: Using Homebrew
brew install go

# Method 2: Download from official site
# Visit https://go.dev/dl/ and download the macOS installer
# Double-click the .pkg file and follow the installation prompts
```

#### Verify Installation
```bash
go version
# Should output: go version go1.21.0 linux/amd64 (or similar)
```

### Step-by-Step Server Setup

#### Step 1: Clone the Repository
```bash
git clone <repository-url>
cd kabawDiscord
```

#### Step 2: Install Go Dependencies
```bash
go mod tidy
```
This command will download the required WebSocket library (`gorilla/websocket`).

#### Step 3: Build the Server (Optional)
```bash
go build -o kabaw-discord main.go
```
This creates an executable binary. You can run it directly with `./kabaw-discord`.

#### Step 4: Run the WebSocket Server
```bash
# Option 1: Run directly with Go
go run main.go

# Option 2: Run the compiled binary
./kabaw-discord
```

#### Step 5: Verify Server is Running
The server will start and display:
```
2025/08/15 15:33:05 Server starting on port :8080
```

Test the server health:
```bash
curl http://localhost:8080/health
# Should return: {"service":"kabaw-discord-server","status":"ok"}
```

🎉 **Your WebSocket server is now running on `ws://localhost:8080/ws`**

## 💻 For Candidates: Building Your React Application

Now that the server is running, you need to create a React application that connects to it.

### Required React App Features

Your React application must implement:

1. **WebSocket Connection Management**
   - Connect to `ws://localhost:8080/ws?username=YourName&channel=general`
   - Handle connection states (connecting, connected, disconnected)
   - Implement reconnection logic

2. **Real-time Message Display**
   - Show incoming messages in real-time
   - Display message metadata (username, timestamp, user ID)
   - Handle different message types (`message`, `system`, `user_connected`)

3. **Message Sending**
   - Input field for typing messages
   - Send button functionality
   - Handle Enter key for sending

4. **User Interface Requirements**
   - Clean, modern chat interface
   - Message history display
   - Connection status indicator
   - User identification (show your assigned user ID)

5. **Error Handling**
   - Display connection errors
   - Handle WebSocket disconnections gracefully
   - Show loading states

### WebSocket Connection Details

**Endpoint**: `ws://localhost:8080/ws`

**Query Parameters**:
- `username` (optional): Your display name (default: "Anonymous")
- `channel` (optional): Channel to join (default: "general")

**Example Connection**:
```javascript
const ws = new WebSocket('ws://localhost:8080/ws?username=CandidateName&channel=general');
```

### ⚠️ Important Notes for Candidates

**YOU SHOULD**:
- Create a new React project (using `create-react-app` or similar)
- Build your React app in **your own separate repository**
- Connect your React app to the provided WebSocket server
- Focus on frontend implementation only

**REFERENCE ONLY**:
- `index.html` and `styles.css` are provided as reference implementations
- You can view them to understand the expected functionality
- Your React implementation should be original, not a copy of the reference client

### 🛠️ Candidate Troubleshooting

**Common Issues and Solutions**:

1. **WebSocket Connection Refused**
   ```bash
   # Make sure the Go server is running first
   go run main.go
   # Should see: "Server starting on port :8080"
   ```

2. **CORS Issues**
   - The server is configured to allow all origins
   - Your React app can run on any port (3000, 3001, etc.)
   - Cross-origin connections are fully supported

3. **Connection Testing**
   ```bash
   # Test server health
   curl http://localhost:8080/health

   # View server stats
   curl http://localhost:8080/stats
   ```

4. **Message Format Debugging**
   - Check browser console for WebSocket messages
   - All messages are logged in JSON format
   - Server logs show all connections and messages

## Testing Cross-Origin WebSocket Connections

To test cross-origin WebSocket functionality:

1. **Start the WebSocket server (port 8080):**
   ```bash
   go run main.go
   ```

2. **Start the HTTP client server (port 6969):**
   ```bash
   npm install
   npm start
   ```

3. **Access the test client:**
   - Open `http://localhost:6969/` in your browser (default route)
   - The HTML is served from port 6969, but connects to WebSocket on port 8080
   - This demonstrates cross-origin WebSocket support

### How Cross-Origin WebSocket Works

The cross-origin functionality is enabled through two key mechanisms:

#### Server-Side Configuration (Go)
```go
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        // Allow connections from any origin (for development)
        return true
    },
}
```
- **Disables origin checking**: Normally WebSocket connections are restricted to same-origin
- **Allows any domain/port** to connect to the WebSocket server
- **Returns `true`** for all origin requests, bypassing browser security restrictions

#### Client-Side Connection (JavaScript)
```javascript
// Cross-origin WebSocket connection to port 8080 from port 6969
const wsUrl = `ws://localhost:8080/ws?username=${username}&channel=${channel}`;
const ws = new WebSocket(wsUrl);
```
- HTML served from `http://localhost:6969` (npm server)
- WebSocket connects to `ws://localhost:8080` (Go server)
- **Different ports = different origins** → Cross-origin request
- Browser allows connection because server's `CheckOrigin` returns `true`

#### Additional CORS Support
- npm http-server provides CORS headers: `access-control-allow-origin: *`
- Enables cross-origin CSS loading: `http://localhost:8080/static/styles.css`
- Demonstrates real-world scenario: frontend from CDN connecting to API

**Security Note**: For production, restrict origins:
```go
CheckOrigin: func(r *http.Request) bool {
    origin := r.Header.Get("Origin")
    return origin == "https://yourdomain.com"
}
```

## Message Simulation

The server automatically generates simulated messages in the **"general"** channel to keep the chat active during testing:

- **Automatic messages**: Random messages are sent every 10 seconds
- **Only when active**: Messages are only sent when users are connected to the "general" channel
- **Variety**: 18+ different message templates with 3 different simulated users (each with unique UUIDs)
- **Consistent timing**: Messages sent every 10 seconds for predictable testing

### Simulated Test Users

When connected to the **"general"** channel, you'll see messages from these automated test users (each with unique UUIDs):

| Username | UUID | Description |
|----------|------|-------------|
| `ChatBot` | `550e8400-e29b-41d4-a716-446655440001` | Friendly bot that welcomes users and provides helpful tips |
| `Developer` | `550e8400-e29b-41d4-a716-446655440002` | Programming-focused user discussing code and development |
| `SystemHelper` | `550e8400-e29b-41d4-a716-446655440003` | Helper bot offering guidance and support |

### Sample Message Types

The simulated users will send various types of messages including:
- **Welcome messages**: "Welcome to the chat! 👋"
- **Programming discussions**: "Anyone else working on Go projects?"
- **General chat**: "How's everyone doing today?"
- **System updates**: "The server is running smoothly"
- **Friendly interactions**: "Hope you're having a great day! 😊"
- **Conversation starters**: "Coffee or tea? ☕"

This feature makes it easy to test the real-time messaging functionality without needing multiple users connected simultaneously.

**To see simulated messages:**
1. Start the WebSocket server: `go run main.go`
2. Connect to the "general" channel via the test client
3. Watch as simulated messages appear automatically every 10 seconds
4. The simulation stops when no users are connected to "general"

## Console Logging

The server provides detailed console logging for monitoring activity:

### Log Types

- **`[CONNECT]`**: User connections with channel and total client count
- **`[DISCONNECT]`**: User disconnections with updated client count
- **`[MESSAGE]`**: Real user messages with channel, username, and content
- **`[SIMULATED]`**: Automated test messages with user ID and content

### Example Console Output

```
2025/08/15 15:16:35 Server starting on port :8080
2025/08/15 15:16:40 [CONNECT] User: TestUser | Channel: general | Total clients: 1
2025/08/15 15:16:45 [SIMULATED] Channel: general | User: ChatBot (ID: 550e8400-e29b-41d4-a716-446655440001) | Content: Welcome to the chat! 👋
2025/08/15 15:16:50 [MESSAGE] Channel: general | User: TestUser | Content: Hello everyone!
2025/08/15 15:17:00 [SIMULATED] Channel: general | User: Developer (ID: 550e8400-e29b-41d4-a716-446655440002) | Content: Anyone else working on Go projects?
2025/08/15 15:17:10 [DISCONNECT] User: TestUser | Channel: general | Total clients: 0
```

This logging makes it easy to monitor server activity, debug issues, and understand message flow patterns.

## Frontend Console Logging

The test client also provides detailed logging in the browser's developer console:

### Frontend Log Types

- **`[FRONTEND-CONNECT]`**: Connection attempts and successful connections
- **`[FRONTEND-DISCONNECT]`**: User-initiated disconnections and connection closures
- **`[FRONTEND-MESSAGE]`**: All incoming messages in pretty-printed JSON format
- **`[FRONTEND-SEND]`**: Outgoing messages in JSON format before sending
- **`[FRONTEND-ERROR]`**: Connection errors and failures

### Example Browser Console Output

```javascript
[FRONTEND-CONNECT] Attempting to connect to: ws://localhost:8080/ws?username=TestUser&channel=general
[FRONTEND-CONNECT] Connected to WebSocket as TestUser in channel general
[FRONTEND-MESSAGE] {
  "type": "system",
  "username": "System",
  "content": "Welcome to the chat!",
  "timestamp": "2025-08-15T15:20:00Z",
  "channel": "general"
}
[FRONTEND-MESSAGE] {
  "type": "message",
  "username": "ChatBot",
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "content": "Welcome to the chat! 👋",
  "timestamp": "2025-08-15T15:20:10Z",
  "channel": "general"
}
[FRONTEND-SEND] {
  "type": "message",
  "content": "Hello everyone!"
}
[FRONTEND-MESSAGE] {
  "type": "message",
  "username": "TestUser",
  "content": "Hello everyone!",
  "timestamp": "2025-08-15T15:20:15Z",
  "channel": "general"
}
[FRONTEND-DISCONNECT] User initiated disconnect
[FRONTEND-DISCONNECT] Connection closed. Code: 1000, Reason:
```

### How to View Frontend Logs

1. Open the test client: `http://localhost:6969/`
2. Press `F12` or right-click → "Inspect" → "Console" tab
3. Connect to the chat and watch the console for detailed logging
4. All WebSocket activity will be logged with clear prefixes

This dual logging (server + frontend) provides complete visibility into the WebSocket communication flow.

## User ID Management

The server automatically generates and assigns a unique user ID to each connecting user:

### How It Works

1. **Connection**: When a user connects, the server generates a unique 32-character hexadecimal user ID
2. **Assignment**: The user ID is sent to the client in a `user_connected` message
3. **Storage**: The frontend stores the user ID in the `currentUserID` variable
4. **Usage**: All user messages include the user ID in the `user_id` field
5. **Cleanup**: The user ID is cleared when the connection closes

### User ID Features

- **Unique identification**: Each connection gets a cryptographically random user ID
- **Persistent during session**: User ID remains constant throughout the WebSocket connection
- **Automatic inclusion**: All user messages automatically include the user ID
- **Frontend tracking**: Client-side JavaScript maintains the current user ID
- **Console logging**: User ID assignment and clearing are logged in the browser console

### Example User ID Flow

```javascript
// 1. User connects
[FRONTEND-CONNECT] Connected to WebSocket as TestUser in channel general

// 2. Server assigns user ID
[FRONTEND-MESSAGE] {
  "type": "user_connected",
  "username": "System",
  "user_id": "a1b2c3d4e5f67890abcdef1234567890",
  "content": "Welcome to the chat!",
  "timestamp": "2025-08-15T15:30:00Z",
  "channel": "general"
}

// 3. Frontend logs user ID assignment
[FRONTEND-USER-ID] Assigned user ID: a1b2c3d4e5f67890abcdef1234567890

// 4. User messages include the ID
[FRONTEND-MESSAGE] {
  "type": "message",
  "username": "TestUser",
  "user_id": "a1b2c3d4e5f67890abcdef1234567890",
  "content": "Hello everyone!",
  "timestamp": "2025-08-15T15:30:15Z",
  "channel": "general"
}

// 5. On disconnect, ID is cleared
[FRONTEND-USER-ID] User ID cleared
```

## API Endpoints

### WebSocket Connection
- **Endpoint:** `ws://localhost:8080/ws`
- **Query Parameters:**
  - `username` (optional): Display name for the user (default: "Anonymous")
  - `channel` (optional): Channel name to join (default: "general")

**Example:**
```
ws://localhost:8080/ws?username=JohnDoe&channel=general
```

### HTTP Endpoints

#### Health Check
- **Endpoint:** `GET /health`
- **Response:** JSON status of the server

#### Statistics
- **Endpoint:** `GET /stats`
- **Response:** JSON with current connection count

#### Web Interface
- **Endpoint:** `GET /`
- **Response:** Simple HTML page with usage instructions

#### Test Client
- **Endpoint:** `GET /test`
- **Response:** Interactive HTML test client with WebSocket functionality

#### Static Files
- **Endpoint:** `GET /static/styles.css`
- **Response:** CSS stylesheet for the test client

## Message Format

### Incoming Messages (Client to Server)
```json
{
  "type": "message",
  "content": "Hello, world!"
}
```

### Outgoing Messages (Server to Client)
```json
{
  "type": "message",
  "username": "JohnDoe",
  "user_id": "a1b2c3d4e5f6789012345678",
  "content": "Hello, world!",
  "timestamp": "2025-08-15T15:30:00Z",
  "channel": "general"
}
```

### Message Types
- `message`: Regular chat message
- `system`: System notification (welcome messages, etc.)
- `user_connected`: Special message sent when user connects, includes their assigned user ID

## Usage Examples

### Using WebSocket Client (JavaScript)

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?username=TestUser&channel=general');

ws.onopen = function(event) {
    console.log('Connected to server');
};

ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log(`${message.username}: ${message.content}`);
};

// Send a message
ws.send(JSON.stringify({
    type: 'message',
    content: 'Hello, everyone!'
}));
```

### Using curl for HTTP endpoints

```bash
# Health check
curl http://localhost:8080/health

# Get statistics
curl http://localhost:8080/stats
```

## Testing the Server

1. **Start the server:**
   ```bash
   go run main.go
   ```

2. **Open your browser and go to:**
   ```
   http://localhost:8080
   ```

3. **Test WebSocket connection using browser developer tools:**
   ```javascript
   const ws = new WebSocket('ws://localhost:8080/ws?username=TestUser');
   ws.onmessage = (e) => console.log(JSON.parse(e.data));
   ws.send(JSON.stringify({type: 'message', content: 'Hello!'}));
   ```

## Architecture

- **Hub**: Central message broker that manages all client connections
- **Client**: Represents individual WebSocket connections
- **Message**: Data structure for chat messages
- **Channels**: Logical separation of conversations

## Development

### Project Structure
```
kabawDiscord/
├── logo.jpeg             # Project logo
├── .gitignore           # Git ignore file for build artifacts and dependencies
├── main.go              # WebSocket server implementation (MAIN FILE)
├── go.mod               # Go module definition
├── index.html           # Reference test client (for testing only)
├── styles.css           # CSS for reference client
└── README.md            # This evaluation document

# Files generated after setup:
# ├── kabaw-discord        # compiled Go binary (after go build)
# └── go.sum               # dependency checksums (after go mod tidy)
```

### Key Components

1. **Hub**: Manages client connections and message broadcasting
2. **Client**: Handles individual WebSocket connections
3. **Message Handling**: JSON-based message protocol
4. **Channel System**: Basic channel separation for messages

## Configuration

The server runs on port 8080 by default. To change the port, modify the `port` variable in `main.go`:

```go
port := ":8080"  // Change this to your desired port
```

## Limitations

- No message persistence (messages are not stored)
- No authentication system
- Basic channel system (no permissions)
- Simplified timestamp handling
- No rate limiting


## License

---

## 📝 Summary for Candidates

This technical evaluation tests your ability to:

1. **Set up and run a Go WebSocket server** (following the provided instructions)
2. **Create a React application** that connects to real-time WebSocket services
3. **Implement proper WebSocket handling** in a modern frontend framework
4. **Provide Feasible Timelines** Be able to accurately access the timetline for implementation
5. **Alternative Solutions** To be able to provide an alternative solution if timelines are not feasible

### Submission Requirements

**Create your own repository** and submit your React application with:

1. **Your React Application Code**
   - Complete React project in your own GitHub repository
   - Clean, well-organized code structure
   - Proper component organization

2. **README.md in your repository** containing:
   - Instructions on how to install and run your React app
   - Prerequisites (Node.js version, npm/yarn, etc.)
   - Step-by-step setup guide
   - How to connect to the WebSocket server

3. **Documentation**
   - Screenshots or video demonstration of the working chat
   - Brief explanation of your implementation approach
   - Any challenges faced and how you solved them
   - Technologies and libraries used

4. **Repository Structure Example**:
   ```
   your-chat-app-repo/
   ├── README.md              # Your setup instructions
   ├── package.json           # Dependencies
   ├── src/
   │   ├── components/        # React components
   │   ├── hooks/            # Custom hooks (WebSocket, etc.)
   │   ├── utils/            # Utility functions
   │   └── App.js            # Main app component
   └── public/               # Static assets
   ```

**Submission Method**:
- submit in provided github link

**⚠️ Important**: Do NOT fork this repository. Create your own independent repository for your React application.
**⚠️ Important**: You are allowed to use any tools you deem necessary including generative AI, the only caveat is that you are able to understand any code you commit.

**Good luck! 🚀**

import { useEffect, useMemo, useRef, useState } from 'react'
import type { FormEvent } from 'react'
import './App.css'

type Message = {
  type: 'message' | 'system' | 'user_connected' | string
  username?: string
  user_id?: string
  content: string
  timestamp?: string
  channel?: string
}

type ConnectionState = 'disconnected' | 'connecting' | 'connected' | 'reconnecting' | 'error'

const DEFAULT_WS_URL = import.meta.env.VITE_WS_URL ?? 'ws://localhost:8080/ws'

const formatTime = (value?: string) => {
  if (!value) return new Date().toLocaleTimeString()
  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime())) return value
  return parsed.toLocaleTimeString()
}

export default function App() {
  const [serverUrl, setServerUrl] = useState(DEFAULT_WS_URL)
  const [username, setUsername] = useState(`candidate-${Math.floor(Math.random() * 9000) + 1000}`)
  const [channel, setChannel] = useState('general')
  const [messages, setMessages] = useState<Message[]>([])
  const [messageBody, setMessageBody] = useState('')
  const [connectionState, setConnectionState] = useState<ConnectionState>('disconnected')
  const [userId, setUserId] = useState<string | null>(null)
  const [connectionError, setConnectionError] = useState<string | null>(null)
  const [reconnectCount, setReconnectCount] = useState(0)

  const socketRef = useRef<WebSocket | null>(null)
  const reconnectAttempts = useRef(0)
  const reconnectTimer = useRef<number | null>(null)
  const manualClose = useRef(false)
  const messageListRef = useRef<HTMLDivElement | null>(null)
  const hasOpened = useRef(false)
  const autoConnectStarted = useRef(false)

  const connectionUrl = useMemo(() => {
    const normalized = serverUrl.trim().replace(/\/$/, '')
    const safeUsername = username.trim() || 'Anonymous'
    const safeChannel = channel.trim() || 'general'
    return `${normalized}?username=${encodeURIComponent(safeUsername)}&channel=${encodeURIComponent(safeChannel)}`
  }, [serverUrl, username, channel])

  useEffect(() => {
    messageListRef.current?.scrollTo({
      top: messageListRef.current.scrollHeight,
      behavior: 'smooth',
    })
  }, [messages])

  useEffect(() => {
    // Avoid StrictMode double-invocation creating a second, short-lived socket
    if (!autoConnectStarted.current) {
      autoConnectStarted.current = true
      handleConnect()
    }
    return () => {
      manualClose.current = true
      clearReconnectTimer()
      socketRef.current?.close()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const clearReconnectTimer = () => {
    if (reconnectTimer.current !== null) {
      clearTimeout(reconnectTimer.current)
      reconnectTimer.current = null
    }
  }

  const appendMessage = (message: Message) => {
    setMessages((prev) => [
      ...prev,
      {
        ...message,
        timestamp: message.timestamp ?? new Date().toISOString(),
      },
    ])
  }

  const handleConnect = () => {
    clearReconnectTimer()
    manualClose.current = false
    setConnectionError(null)
    setConnectionState(reconnectAttempts.current ? 'reconnecting' : 'connecting')

    socketRef.current?.close()
    const socket = new WebSocket(connectionUrl)
    socketRef.current = socket

    socket.onopen = () => {
      hasOpened.current = true
      reconnectAttempts.current = 0
      setReconnectCount(0)
      setConnectionState('connected')
      appendMessage({
        type: 'system',
        username: 'System',
        content: `Connected as ${username} in #${channel}`,
        channel,
      })
    }

    socket.onmessage = (event) => {
      try {
        const parsed = JSON.parse(event.data) as Message
        if (parsed.type === 'user_connected' && parsed.user_id) {
          setUserId(parsed.user_id)
        }
        appendMessage(parsed)
      } catch (error) {
        setConnectionError('Received a message the client could not parse.')
      }
    }

    socket.onerror = () => {
      setConnectionState('error')
      setConnectionError('WebSocket reported an error.')
    }

    socket.onclose = (event) => {
      socketRef.current = null
      setUserId(null)
      const wasEstablished = hasOpened.current
      hasOpened.current = false

      // Only surface a close message if we had an established session
      if (wasEstablished) {
        appendMessage({
          type: 'system',
          username: 'System',
          content: `Connection closed (${event.code || 'no code'})`,
          channel,
        })
      }

      if (!manualClose.current) {
        scheduleReconnect()
      } else {
        setConnectionState('disconnected')
      }
    }
  }

  const scheduleReconnect = () => {
    reconnectAttempts.current += 1
    setReconnectCount(reconnectAttempts.current)
    const delay = Math.min(30000, 1000 * Math.pow(2, reconnectAttempts.current))
    setConnectionState('reconnecting')
    reconnectTimer.current = window.setTimeout(() => handleConnect(), delay)
  }

  const handleDisconnect = () => {
    manualClose.current = true
    reconnectAttempts.current = 0
    setReconnectCount(0)
    clearReconnectTimer()
    socketRef.current?.close()
    socketRef.current = null
    setConnectionState('disconnected')
    appendMessage({
      type: 'system',
      username: 'System',
      content: 'Disconnected from server',
      channel,
    })
  }

  const handleSend = (event?: FormEvent<HTMLFormElement>) => {
    event?.preventDefault()
    const trimmed = messageBody.trim()
    if (!trimmed) return

    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(
        JSON.stringify({
          type: 'message',
          content: trimmed,
        })
      )
      setMessageBody('')
    } else {
      setConnectionError('You are not connected yet.')
      setConnectionState('error')
    }
  }

  const stateLabel: Record<ConnectionState, string> = {
    disconnected: 'Disconnected',
    connecting: 'Connecting…',
    connected: 'Live',
    reconnecting: 'Reconnecting…',
    error: 'Error',
  }

  return (
    <div className="page">
      <header className="hero">
        <div>
          <p className="eyebrow">Kabaw sockets · Live chat demo</p>
          <h1>React client for the Kabaw WebSocket server</h1>
          <p className="lede">
            Connect to <code>ws://localhost:8080/ws</code>, join a channel, and watch messages stream in real time.
            Auto-reconnect keeps the feed alive if the server restarts.
          </p>
          <div className="chips">
            <span className={`chip state ${connectionState}`}>{stateLabel[connectionState]}</span>
            <span className="chip subtle">Channel: #{channel || 'general'}</span>
            <span className="chip subtle">User ID: {userId ?? 'waiting for assignment'}</span>
          </div>
        </div>
        <div className="badge-card">
          <p className="badge-title">Connection quick view</p>
          <ul>
            <li>
              <span>Status</span>
              <strong>{stateLabel[connectionState]}</strong>
            </li>
            <li>
              <span>Reconnect attempts</span>
              <strong>{reconnectCount}</strong>
            </li>
            <li>
              <span>Endpoint</span>
              <strong>{serverUrl}</strong>
            </li>
          </ul>
        </div>
      </header>

      <main className="grid">
        <section className="panel">
          <div className="panel-header">
            <div>
              <p className="eyebrow">Connection</p>
              <h2>Join a channel</h2>
            </div>
            <div className="actions">
              <button type="button" className="ghost" onClick={handleDisconnect} disabled={connectionState === 'disconnected'}>
                Disconnect
              </button>
              <button type="button" onClick={handleConnect}>
                {connectionState === 'connected' ? 'Reconnect' : 'Connect'}
              </button>
            </div>
          </div>

          <div className="form-grid">
            <label>
              <span>Server URL</span>
              <input value={serverUrl} onChange={(e) => setServerUrl(e.target.value)} placeholder="ws://localhost:8080/ws" />
            </label>
            <label>
              <span>Username</span>
              <input value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Your handle" />
            </label>
            <label>
              <span>Channel</span>
              <input value={channel} onChange={(e) => setChannel(e.target.value)} placeholder="general" />
            </label>
          </div>

          <div className="info-row">
            <div>
              <p className="label">Connection URL</p>
              <code className="code-pill">{connectionUrl}</code>
            </div>
            {connectionError && <div className="error">{connectionError}</div>}
          </div>
        </section>

        <section className="panel messages">
          <div className="panel-header">
            <div>
              <p className="eyebrow">Live feed</p>
              <h2>Messages</h2>
            </div>
            <span className={`chip state ${connectionState}`}>{stateLabel[connectionState]}</span>
          </div>

          <div className="message-list" ref={messageListRef}>
            {messages.length === 0 ? (
              <div className="empty">Connect to see automated messages or chat with others in the same channel.</div>
            ) : (
              messages.map((message, index) => {
                const isSystem = message.type === 'system' || message.type === 'user_connected'
                const isSelf = message.username === username

                return (
                  <article
                    key={`${message.timestamp}-${index}`}
                    className={`message-card ${isSystem ? 'system' : isSelf ? 'self' : 'other'}`}
                  >
                    <div className="message-meta">
                      <span className="author">
                        {message.username ?? 'Unknown'} {isSelf && <span className="tag">you</span>}
                      </span>
                      {message.user_id && <span className="pill">id {message.user_id.slice(0, 8)}</span>}
                      {message.channel && <span className="pill">#{message.channel}</span>}
                      <span className="time">{formatTime(message.timestamp)}</span>
                    </div>
                    <p className="message-body">{message.content}</p>
                  </article>
                )
              })
            )}
          </div>

          <form className="composer" onSubmit={handleSend}>
            <div>
              <label className="label">Message</label>
              <input
                value={messageBody}
                onChange={(e) => setMessageBody(e.target.value)}
                placeholder="Type and press Enter to send"
                onKeyDown={(event) => {
                  if (event.key === 'Enter' && !event.shiftKey) {
                    event.preventDefault()
                    handleSend()
                  }
                }}
              />
            </div>
            <button type="submit" disabled={connectionState === 'disconnected' || connectionState === 'error'}>
              Send
            </button>
          </form>
        </section>
      </main>
    </div>
  )
}

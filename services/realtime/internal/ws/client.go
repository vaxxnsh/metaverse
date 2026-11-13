package ws

import (
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultReadLimit  int64 = 4 << 10
	defaultWriteWait        = 10 * time.Second
	defaultPongWait         = 60 * time.Second
	defaultSendBuffer       = 32
)

// TransportConfig controls connection limits and health checks.
type TransportConfig struct {
	ReadLimit  int64
	WriteWait  time.Duration
	PongWait   time.Duration
	PingPeriod time.Duration
	SendBuffer int
}

func (c TransportConfig) normalized() TransportConfig {
	if c.ReadLimit <= 0 {
		c.ReadLimit = defaultReadLimit
	}
	if c.WriteWait <= 0 {
		c.WriteWait = defaultWriteWait
	}
	if c.PongWait <= 0 {
		c.PongWait = defaultPongWait
	}
	if c.PingPeriod <= 0 || c.PingPeriod >= c.PongWait {
		c.PingPeriod = c.PongWait * 9 / 10
	}
	if c.SendBuffer <= 0 {
		c.SendBuffer = defaultSendBuffer
	}
	return c
}

// Client is a connection with exactly one reader and one writer goroutine.
type Client struct {
	conn    *websocket.Conn
	config  TransportConfig
	send    chan Event
	closed  chan struct{}
	done    sync.Once
	onEvent func(*Client, Event)
	onClose func(*Client)
}

func newClient(conn *websocket.Conn, config TransportConfig, onEvent func(*Client, Event), onClose func(*Client)) *Client {
	config = config.normalized()
	return &Client{conn: conn, config: config, send: make(chan Event, config.SendBuffer), closed: make(chan struct{}), onEvent: onEvent, onClose: onClose}
}

func (c *Client) Start() { go c.readLoop(); go c.writeLoop() }

// Send never blocks a broadcast; false means closed or saturated.
func (c *Client) Send(event Event) bool {
	select {
	case <-c.closed:
		return false
	default:
	}
	select {
	case c.send <- event:
		return true
	case <-c.closed:
		return false
	default:
		return false
	}
}

func (c *Client) Close() {
	c.done.Do(func() {
		close(c.closed)
		_ = c.conn.Close()
		if c.onClose != nil {
			c.onClose(c)
		}
	})
}

func (c *Client) RemoteAddr() net.Addr { return c.conn.RemoteAddr() }

func (c *Client) readLoop() {
	defer c.Close()
	c.conn.SetReadLimit(c.config.ReadLimit)
	_ = c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait))
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(c.config.PongWait)) })
	for {
		messageType, raw, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		if messageType != websocket.TextMessage {
			c.protocolError("INVALID_MESSAGE_TYPE", "only text messages are supported")
			continue
		}
		event, err := DecodeEvent(raw)
		if err != nil {
			c.protocolError("INVALID_EVENT", "invalid event envelope")
			continue
		}
		if err := ValidateClientEvent(event); err != nil {
			c.protocolError("UNKNOWN_EVENT", "unsupported event type")
			continue
		}
		if c.onEvent != nil {
			c.onEvent(c, event)
		}
	}
}

func (c *Client) writeLoop() {
	ticker := time.NewTicker(c.config.PingPeriod)
	defer ticker.Stop()
	defer c.Close()
	for {
		select {
		case <-c.closed:
			return
		case event := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteWait))
			if err := c.conn.WriteJSON(event); err != nil {
				return
			}
		case <-ticker.C:
			deadline := time.Now().Add(c.config.WriteWait)
			if err := c.conn.WriteControl(websocket.PingMessage, nil, deadline); err != nil {
				return
			}
		}
	}
}

func (c *Client) protocolError(code, message string) {
	event, err := NewEvent(EventError, ErrorPayload{Code: code, Message: message})
	if err != nil || !c.Send(event) {
		c.Close()
	}
}

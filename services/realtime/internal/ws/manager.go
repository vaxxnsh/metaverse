package ws

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrInvalidRoom         = errors.New("invalid room")
	ErrClientNotRegistered = errors.New("client is not registered")
	ErrClientNotJoined     = errors.New("client has not joined a room")
	ErrOutsideBounds       = errors.New("position is outside room bounds")
	ErrStaticCollision     = errors.New("position collides with a static element")
	ErrUserCollision       = errors.New("position collides with another user")
	ErrMoveTooFar          = errors.New("movement exceeds maximum distance")
	ErrMoveTooFast         = errors.New("movement is rate limited")
)

type Manager struct {
	mu          sync.RWMutex
	clients     map[*Client]struct{}
	rooms       map[string]*Room
	clientRooms map[*Client]*Room
	movement    MovementConfig
}

// Departure is the room-presence state captured before a disconnected client
// is removed. Recipients are a snapshot and can be notified after unlocking.
type Departure struct {
	RoomID     string
	UserID     string
	Recipients []*Client
}

type MovementConfig struct {
	PlayerWidth, PlayerHeight int32
	MaxMoveDistance           int32
	MinMoveInterval           time.Duration
	PreventPlayerOverlap      bool
}

func NewManager() *Manager {
	return &Manager{
		clients:     make(map[*Client]struct{}),
		rooms:       make(map[string]*Room),
		clientRooms: make(map[*Client]*Room),
		movement:    MovementConfig{PlayerWidth: 1, PlayerHeight: 1, MaxMoveDistance: 5, MinMoveInterval: 50 * time.Millisecond, PreventPlayerOverlap: true},
	}
}

// SetMovementConfig is intended for server setup and tests, before clients move.
func (m *Manager) SetMovementConfig(config MovementConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if config.PlayerWidth <= 0 {
		config.PlayerWidth = 1
	}
	if config.PlayerHeight <= 0 {
		config.PlayerHeight = 1
	}
	if config.MaxMoveDistance <= 0 {
		config.MaxMoveDistance = 5
	}
	m.movement = config
}

// Join stores the authoritative room geometry, snapshots existing occupants,
// and registers an authenticated client in one atomic state transition.
func (m *Manager) Join(client *Client, userID string, space Space) ([]Occupant, error) {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(space.ID) == "" {
		return nil, ErrInvalidRoom
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[client]; !ok {
		return nil, ErrClientNotRegistered
	}
	if current, ok := m.clientRooms[client]; ok {
		m.removeFromRoomLocked(client, current.ID)
	}
	room, ok := m.rooms[space.ID]
	if !ok {
		room = &Room{ID: space.ID, Width: space.Width, Height: space.Height, Elements: space.Elements, users: make(map[*Client]*User)}
		m.rooms[space.ID] = room
	}
	occupants := make([]Occupant, 0, len(room.users))
	for _, user := range room.users {
		occupants = append(occupants, Occupant{UserID: user.ID, Position: user.Position})
	}
	room.users[client] = &User{ID: userID}
	m.clientRooms[client] = room
	return occupants, nil
}

// Move updates a joined client's position after validating only the room
// boundary. Static-element and player collisions are added in Phase 6.
func (m *Manager) Move(client *Client, target Position) (MovementPayload, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, ok := m.clientRooms[client]
	if !ok {
		return MovementPayload{}, ErrClientNotJoined
	}
	user, ok := room.users[client]
	if !ok {
		return MovementPayload{}, ErrClientNotJoined
	}
	current := user.Position
	if target.X < 0 || target.Y < 0 || target.X+m.movement.PlayerWidth > room.Width || target.Y+m.movement.PlayerHeight > room.Height {
		return MovementPayload{UserID: user.ID, Position: current}, ErrOutsideBounds
	}
	if abs(target.X-current.X) > m.movement.MaxMoveDistance || abs(target.Y-current.Y) > m.movement.MaxMoveDistance {
		return MovementPayload{UserID: user.ID, Position: current}, ErrMoveTooFar
	}
	if !user.LastMove.IsZero() && time.Since(user.LastMove) < m.movement.MinMoveInterval {
		return MovementPayload{UserID: user.ID, Position: current}, ErrMoveTooFast
	}
	for _, element := range room.Elements {
		if element.Static && overlaps(target, m.movement.PlayerWidth, m.movement.PlayerHeight, Position{X: element.X, Y: element.Y}, element.Width, element.Height) {
			return MovementPayload{UserID: user.ID, Position: current}, ErrStaticCollision
		}
	}
	if m.movement.PreventPlayerOverlap {
		for member, otherUser := range room.users {
			if member != client && overlaps(target, m.movement.PlayerWidth, m.movement.PlayerHeight, otherUser.Position, m.movement.PlayerWidth, m.movement.PlayerHeight) {
				return MovementPayload{UserID: user.ID, Position: current}, ErrUserCollision
			}
		}
	}

	user.Position = target
	user.LastMove = time.Now()
	return MovementPayload{UserID: user.ID, Position: target}, nil
}

func (m *Manager) Register(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[client] = struct{}{}
}

// Unregister removes a client globally and from its current room. It is safe
// to call more than once from concurrent connection cleanup paths.
func (m *Manager) Unregister(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.removeClientLocked(client)
}

// Disconnect atomically removes a client and captures the remaining room
// members for a user-left broadcast. Calling it more than once is safe.
func (m *Manager) Disconnect(client *Client) (Departure, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, registered := m.clients[client]; !registered {
		return Departure{}, false
	}

	if room, ok := m.clientRooms[client]; ok {
		departure := Departure{RoomID: room.ID}
		if user := room.users[client]; user != nil {
			departure.UserID = user.ID
		}
		departure.Recipients = make([]*Client, 0, len(room.users)-1)
		for member := range room.users {
			if member != client {
				departure.Recipients = append(departure.Recipients, member)
			}
		}
		m.removeClientLocked(client)
		return departure, true
	}
	m.removeClientLocked(client)
	return Departure{}, true
}

// Shutdown closes all active clients. Each client executes its normal close
// callback, making cleanup and user-left broadcasts follow one code path.
func (m *Manager) Shutdown() {
	m.mu.RLock()
	clients := make([]*Client, 0, len(m.clients))
	for client := range m.clients {
		clients = append(clients, client)
	}
	m.mu.RUnlock()
	for _, client := range clients {
		client.Close()
	}
}

// AssignRoom moves a registered client into a room, creating it when needed.
// Phase 4 calls it only after authentication and space loading succeed.
func (m *Manager) AssignRoom(client *Client, roomID string) error {
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return ErrInvalidRoom
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[client]; !ok {
		return ErrClientNotRegistered
	}

	if currentRoom, ok := m.clientRooms[client]; ok {
		if currentRoom.ID == roomID {
			return nil
		}
		m.removeFromRoomLocked(client, currentRoom.ID)
	}

	room, ok := m.rooms[roomID]
	if !ok {
		room = &Room{ID: roomID, users: make(map[*Client]*User)}
		m.rooms[roomID] = room
	}
	room.users[client] = &User{}
	m.clientRooms[client] = room
	return nil
}

func (m *Manager) LeaveRoom(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if room, ok := m.clientRooms[client]; ok {
		m.removeFromRoomLocked(client, room.ID)
	}
}

// Broadcast sends an event to members of one room without holding the state
// lock. The excluded client is useful when broadcasting its own movement.
func (m *Manager) Broadcast(roomID string, event Event, excluded *Client) int {
	recipients := m.RoomMembers(roomID, excluded)
	sent := 0
	for _, client := range recipients {
		if client.Send(event) {
			sent++
		}
	}
	return sent
}

func (m *Manager) RoomMembers(roomID string, excluded *Client) []*Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	room, ok := m.rooms[roomID]
	if !ok {
		return nil
	}

	clients := make([]*Client, 0, len(room.users))
	for client := range room.users {
		if client != excluded {
			clients = append(clients, client)
		}
	}
	return clients
}

func (m *Manager) ClientRoom(client *Client) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	room, ok := m.clientRooms[client]
	if !ok {
		return "", false
	}
	return room.ID, true
}

func (m *Manager) ClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

func (m *Manager) RoomClientCount(roomID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if room, ok := m.rooms[roomID]; ok {
		return len(room.users)
	}
	return 0
}

func (m *Manager) removeClientLocked(client *Client) {
	delete(m.clients, client)
	if room, ok := m.clientRooms[client]; ok {
		m.removeFromRoomLocked(client, room.ID)
	}
}

func overlaps(a Position, aw, ah int32, b Position, bw, bh int32) bool {
	return a.X < b.X+bw && a.X+aw > b.X && a.Y < b.Y+bh && a.Y+ah > b.Y
}

func abs(value int32) int32 {
	if value < 0 {
		return -value
	}
	return value
}

func (m *Manager) removeFromRoomLocked(client *Client, roomID string) {
	delete(m.clientRooms, client)
	room, ok := m.rooms[roomID]
	if !ok {
		return
	}
	delete(room.users, client)
	if len(room.users) == 0 {
		delete(m.rooms, roomID)
	}
}

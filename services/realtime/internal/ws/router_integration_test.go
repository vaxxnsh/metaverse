package ws

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestRealtimeFlow(t *testing.T) {
	manager := NewManager()
	manager.SetMovementConfig(MovementConfig{MaxMoveDistance: 10, PreventPlayerOverlap: true})
	space := Space{ID: "space-1", Width: 10, Height: 10}
	join := NewJoinHandler(manager, fakeVerifier{claims: Claims{Subject: "user"}}, fakeSpaces{space: space})
	router := NewEventRouter(join, NewMovementHandler(manager))

	handler := NewHTTPHandler(HTTPHandlerOptions{
		Transport: TransportConfig{PingPeriod: time.Hour},
		OnConnect: manager.Register,
		OnEvent:   router.Handle,
		OnClose: func(client *Client) {
			departure, ok := manager.Disconnect(client)
			if !ok || departure.RoomID == "" {
				return
			}
			event, _ := NewEvent(EventUserLeft, UserLeftPayload{UserID: departure.UserID})
			for _, recipient := range departure.Recipients {
				recipient.Send(event)
			}
		},
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	first := dialRealtimeServer(t, server.URL)
	defer first.Close()
	second := dialRealtimeServer(t, server.URL)
	defer second.Close()

	writeJoin(t, first, "first-token")
	assertEventType(t, first, EventSpaceJoined)
	writeJoin(t, second, "second-token")
	assertEventType(t, second, EventSpaceJoined)
	assertEventType(t, first, EventUserJoined)

	if err := first.WriteJSON(Event{Type: EventMove, Payload: []byte(`{"x":1,"y":0}`)}); err != nil {
		t.Fatal(err)
	}
	assertEventType(t, second, EventMovement)

	if err := first.WriteJSON(Event{Type: EventMove, Payload: []byte(`{"x":99,"y":0}`)}); err != nil {
		t.Fatal(err)
	}
	assertEventType(t, first, EventMovementRejected)

	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	assertEventType(t, second, EventUserLeft)
}

func writeJoin(t *testing.T, conn *websocket.Conn, token string) {
	t.Helper()
	if err := conn.WriteJSON(Event{Type: EventJoin, Payload: []byte(`{"spaceId":"space-1","token":"` + token + `"}`)}); err != nil {
		t.Fatal(err)
	}
}

func dialRealtimeServer(t *testing.T, serverURL string) *websocket.Conn {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial(testWebsocketURL(serverURL), nil)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	return conn
}

func assertEventType(t *testing.T, conn *websocket.Conn, expected string) {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(time.Second))
	var event Event
	if err := conn.ReadJSON(&event); err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}
	if event.Type != expected {
		t.Fatalf("event type = %q, want %q", event.Type, expected)
	}
}

package ws

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestHTTPHandlerReceivesAndSendsEvents(t *testing.T) {
	received := make(chan Event, 1)
	server := httptest.NewServer(NewHTTPHandler(HTTPHandlerOptions{
		Transport: TransportConfig{PingPeriod: time.Hour},
		OnConnect: func(client *Client) {
			event, _ := NewEvent(EventError, ErrorPayload{Code: "CONNECTED", Message: "ready"})
			client.Send(event)
		},
		OnEvent: func(_ *Client, event Event) { received <- event },
	}))
	defer server.Close()

	conn, _, err := websocket.DefaultDialer.Dial(testWebsocketURL(server.URL), nil)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	defer conn.Close()

	var outbound Event
	if err := conn.ReadJSON(&outbound); err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}
	if outbound.Type != EventError {
		t.Fatalf("outbound type = %q, want %q", outbound.Type, EventError)
	}
	if err := conn.WriteJSON(Event{Type: EventMove, Payload: []byte(`{"x":2,"y":3}`)}); err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}
	select {
	case event := <-received:
		if event.Type != EventMove {
			t.Fatalf("received type = %q, want %q", event.Type, EventMove)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for inbound event")
	}
}

func TestOriginAllowList(t *testing.T) {
	check := OriginAllowList([]string{"https://app.example.com"})
	if !check(httptest.NewRequest("GET", "/ws", nil)) {
		t.Fatal("request without an origin was rejected")
	}
	allowed := httptest.NewRequest("GET", "/ws", nil)
	allowed.Header.Set("Origin", "https://app.example.com")
	if !check(allowed) {
		t.Fatal("request with an allowed origin was rejected")
	}
	rejected := httptest.NewRequest("GET", "/ws", nil)
	rejected.Header.Set("Origin", "https://untrusted.example.com")
	if check(rejected) {
		t.Fatal("request with an untrusted origin was accepted")
	}
}

func testWebsocketURL(serverURL string) string {
	u, _ := url.Parse(serverURL)
	u.Scheme = "ws"
	return strings.TrimSuffix(u.String(), "/") + "/ws"
}

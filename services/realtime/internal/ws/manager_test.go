package ws

import "testing"

func TestManagerScopesBroadcastsToRoom(t *testing.T) {
	manager := NewManager()
	first := testClient()
	second := testClient()
	otherRoom := testClient()

	for _, client := range []*Client{first, second, otherRoom} {
		manager.Register(client)
	}
	if err := manager.AssignRoom(first, "space-a"); err != nil {
		t.Fatal(err)
	}
	if err := manager.AssignRoom(second, "space-a"); err != nil {
		t.Fatal(err)
	}
	if err := manager.AssignRoom(otherRoom, "space-b"); err != nil {
		t.Fatal(err)
	}

	event, err := NewEvent(EventMovement, MovementPayload{UserID: "user-1", Position: Position{X: 2, Y: 3}})
	if err != nil {
		t.Fatal(err)
	}
	if got := manager.Broadcast("space-a", event, first); got != 1 {
		t.Fatalf("Broadcast() = %d, want 1", got)
	}

	select {
	case received := <-second.send:
		if received.Type != EventMovement {
			t.Fatalf("received type = %q", received.Type)
		}
	default:
		t.Fatal("room member did not receive broadcast")
	}
	select {
	case <-otherRoom.send:
		t.Fatal("client in another room received broadcast")
	default:
	}
}

func TestManagerUnregisterRemovesEmptyRoom(t *testing.T) {
	manager := NewManager()
	client := testClient()
	manager.Register(client)
	if err := manager.AssignRoom(client, "space-a"); err != nil {
		t.Fatal(err)
	}
	manager.Unregister(client)
	manager.Unregister(client)
	if manager.ClientCount() != 0 {
		t.Fatalf("ClientCount() = %d, want 0", manager.ClientCount())
	}
	if manager.RoomClientCount("space-a") != 0 {
		t.Fatal("empty room was retained")
	}
	if _, ok := manager.ClientRoom(client); ok {
		t.Fatal("removed client retained a room")
	}
}

func TestManagerDisconnectReturnsRoomRecipientsOnce(t *testing.T) {
	manager := NewManager()
	first, second := testClient(), testClient()
	manager.Register(first)
	manager.Register(second)
	space := Space{ID: "space-a", Width: 10, Height: 10}
	_, _ = manager.Join(first, "user-one", space)
	_, _ = manager.Join(second, "user-two", space)

	departure, ok := manager.Disconnect(first)
	if !ok {
		t.Fatal("Disconnect() returned false")
	}
	if departure.UserID != "user-one" || departure.RoomID != "space-a" || len(departure.Recipients) != 1 || departure.Recipients[0] != second {
		t.Fatalf("unexpected departure: %#v", departure)
	}
	if _, ok := manager.Disconnect(first); ok {
		t.Fatal("second Disconnect() returned true")
	}
	if manager.RoomClientCount("space-a") != 1 {
		t.Fatal("remaining room member was removed")
	}
}

func testClient() *Client {
	return &Client{send: make(chan Event, 4), closed: make(chan struct{})}
}

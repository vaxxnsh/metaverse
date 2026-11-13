package ws

import "testing"

func TestMovementHandlerBroadcastsValidMovement(t *testing.T) {
	manager := NewManager()
	moving, observer := testClient(), testClient()
	manager.Register(moving)
	manager.Register(observer)
	space := Space{ID: "space-1", Width: 10, Height: 10}
	if _, err := manager.Join(moving, "user-one", space); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Join(observer, "user-two", space); err != nil {
		t.Fatal(err)
	}
	// Ignore the join-presence state; this test exercises only movement.
	handler := NewMovementHandler(manager)
	handler.Handle(moving, Event{Type: EventMove, Payload: []byte(`{"x":2,"y":3}`)})

	event := <-observer.send
	if event.Type != EventMovement {
		t.Fatalf("event type = %q", event.Type)
	}
	payload, err := DecodePayload[MovementPayload](event)
	if err != nil {
		t.Fatal(err)
	}
	if payload.UserID != "user-one" || payload.X != 2 || payload.Y != 3 {
		t.Fatalf("unexpected movement: %#v", payload)
	}
}

func TestMovementHandlerRejectsOutOfBoundsMove(t *testing.T) {
	manager := NewManager()
	client := testClient()
	manager.Register(client)
	if _, err := manager.Join(client, "user-one", Space{ID: "space-1", Width: 2, Height: 2}); err != nil {
		t.Fatal(err)
	}

	NewMovementHandler(manager).Handle(client, Event{Type: EventMove, Payload: []byte(`{"x":2,"y":0}`)})
	event := <-client.send
	if event.Type != EventMovementRejected {
		t.Fatalf("event type = %q", event.Type)
	}
	payload, err := DecodePayload[MovementRejectedPayload](event)
	if err != nil {
		t.Fatal(err)
	}
	if payload.Reason != "outside_bounds" || payload.X != 0 || payload.Y != 0 {
		t.Fatalf("unexpected rejection: %#v", payload)
	}
}

func TestMovementHandlerRejectsMoveBeforeJoin(t *testing.T) {
	manager := NewManager()
	client := testClient()
	manager.Register(client)
	NewMovementHandler(manager).Handle(client, Event{Type: EventMove, Payload: []byte(`{"x":1,"y":1}`)})
	event := <-client.send
	payload, err := DecodePayload[MovementRejectedPayload](event)
	if err != nil {
		t.Fatal(err)
	}
	if event.Type != EventMovementRejected || payload.Reason != "not_joined" {
		t.Fatalf("unexpected rejection: %#v", payload)
	}
}

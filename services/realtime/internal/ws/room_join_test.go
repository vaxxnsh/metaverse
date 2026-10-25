package ws

import (
	"context"
	"errors"
	"testing"
)

type fakeVerifier struct {
	claims Claims
	err    error
}

func (f fakeVerifier) Verify(string) (Claims, error) { return f.claims, f.err }

type fakeSpaces struct {
	space Space
	err   error
}

func (f fakeSpaces) GetRealtimeSpace(context.Context, string) (Space, error) { return f.space, f.err }

func TestJoinHandlerSendsSnapshotAndPresence(t *testing.T) {
	manager := NewManager()
	existing, joining := testClient(), testClient()
	manager.Register(existing)
	manager.Register(joining)
	if _, err := manager.Join(existing, "user-one", Space{ID: "space-1", Width: 10, Height: 10}); err != nil {
		t.Fatal(err)
	}
	handler := NewJoinHandler(manager, fakeVerifier{claims: Claims{Subject: "user-two"}}, fakeSpaces{space: Space{ID: "space-1", Width: 10, Height: 10}})
	handler.Handle(joining, Event{Type: EventJoin, Payload: []byte(`{"spaceId":"space-1","token":"jwt"}`)})
	joined := <-joining.send
	if joined.Type != EventSpaceJoined {
		t.Fatalf("got %s", joined.Type)
	}
	presence := <-existing.send
	if presence.Type != EventUserJoined {
		t.Fatalf("got %s", presence.Type)
	}
}

func TestJoinHandlerRejectsInvalidToken(t *testing.T) {
	manager := NewManager()
	client := testClient()
	manager.Register(client)
	handler := NewJoinHandler(manager, fakeVerifier{err: errors.New("bad token")}, fakeSpaces{})
	handler.Handle(client, Event{Type: EventJoin, Payload: []byte(`{"spaceId":"space-1","token":"bad"}`)})
	if event := <-client.send; event.Type != EventError {
		t.Fatalf("got %s", event.Type)
	}
}

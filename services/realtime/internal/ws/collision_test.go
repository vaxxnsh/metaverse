package ws

import (
	"testing"
	"time"
)

func TestMovementCollisionAndLimits(t *testing.T) {
	manager := NewManager()
	manager.SetMovementConfig(MovementConfig{PlayerWidth: 1, PlayerHeight: 1, MaxMoveDistance: 10, MinMoveInterval: time.Hour, PreventPlayerOverlap: true})
	first, second := testClient(), testClient()
	manager.Register(first)
	manager.Register(second)
	space := Space{ID: "space", Width: 10, Height: 10, Elements: []StaticElement{{X: 2, Y: 2, Width: 1, Height: 1, Static: true}}}
	if _, err := manager.Join(first, "one", space); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Join(second, "two", space); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Move(first, Position{X: 2, Y: 2}); err != ErrStaticCollision {
		t.Fatalf("static collision = %v", err)
	}
	if _, err := manager.Move(first, Position{X: 1, Y: 0}); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Move(second, Position{X: 1, Y: 0}); err != ErrUserCollision {
		t.Fatalf("user collision = %v", err)
	}
	if _, err := manager.Move(first, Position{X: 2, Y: 0}); err != ErrMoveTooFast {
		t.Fatalf("rate limit = %v", err)
	}
}

func TestMovementRejectsLongMove(t *testing.T) {
	manager := NewManager()
	manager.SetMovementConfig(MovementConfig{MaxMoveDistance: 2, PreventPlayerOverlap: false})
	client := testClient()
	manager.Register(client)
	_, _ = manager.Join(client, "user", Space{ID: "space", Width: 20, Height: 20})
	if _, err := manager.Move(client, Position{X: 3, Y: 0}); err != ErrMoveTooFar {
		t.Fatalf("long move = %v", err)
	}
}

package ws

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestDecodeEventAndJoinPayload(t *testing.T) {
	raw := []byte(`{"type":"join","payload":{"spaceId":"space-1","token":"jwt"}}`)

	event, err := DecodeEvent(raw)
	if err != nil {
		t.Fatalf("DecodeEvent() error = %v", err)
	}
	if err := ValidateClientEvent(event); err != nil {
		t.Fatalf("ValidateClientEvent() error = %v", err)
	}

	payload, err := DecodePayload[JoinPayload](event)
	if err != nil {
		t.Fatalf("DecodePayload() error = %v", err)
	}
	if payload.SpaceID != "space-1" || payload.Token != "jwt" {
		t.Fatalf("unexpected join payload: %#v", payload)
	}
}

func TestDecodeEventRejectsMalformedEnvelope(t *testing.T) {
	tests := [][]byte{
		[]byte(`not json`),
		[]byte(`{"payload":{}}`),
		[]byte(`{"type":"join"}`),
		[]byte(`{"type":"join","payload":null}`),
	}

	for _, raw := range tests {
		_, err := DecodeEvent(raw)
		if !errors.Is(err, ErrInvalidEvent) {
			t.Errorf("DecodeEvent(%s) error = %v, want ErrInvalidEvent", raw, err)
		}
	}
}

func TestValidateClientEventRejectsUnknownType(t *testing.T) {
	err := ValidateClientEvent(Event{Type: "teleport", Payload: json.RawMessage(`{}`)})
	if !errors.Is(err, ErrUnknownClientEvent) {
		t.Fatalf("ValidateClientEvent() error = %v, want ErrUnknownClientEvent", err)
	}
}

func TestNewEventRoundTrip(t *testing.T) {
	event, err := NewEvent(EventMovement, MovementPayload{
		UserID:   "user-1",
		Position: Position{X: 2, Y: 3},
	})
	if err != nil {
		t.Fatalf("NewEvent() error = %v", err)
	}

	encoded, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded struct {
		Type    string          `json:"type"`
		Payload MovementPayload `json:"payload"`
	}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if decoded.Type != EventMovement || decoded.Payload.UserID != "user-1" || decoded.Payload.X != 2 || decoded.Payload.Y != 3 {
		t.Fatalf("unexpected round trip: %#v", decoded)
	}
}

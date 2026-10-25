package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidEvent is returned when a message cannot be interpreted as the
	// common event envelope.
	ErrInvalidEvent = errors.New("invalid websocket event")
	// ErrUnknownClientEvent is returned when an inbound event has no handler.
	ErrUnknownClientEvent = errors.New("unknown client event")
)

const (
	EventJoin             = "join"
	EventMove             = "move"
	EventSpaceJoined      = "space-joined"
	EventUserJoined       = "user-joined"
	EventMovement         = "movement"
	EventMovementRejected = "movement-rejected"
	EventUserLeft         = "user-left"
	EventError            = "error"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type JoinPayload struct {
	SpaceID string `json:"spaceId"`
	Token   string `json:"token"`
}

type MovePayload struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type Position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type Occupant struct {
	UserID string `json:"userId"`
	Position
}

type SpaceJoinedPayload struct {
	SpaceID string     `json:"spaceId"`
	Spawn   Position   `json:"spawn"`
	Users   []Occupant `json:"users"`
}

type UserJoinedPayload struct {
	UserID string `json:"userId"`
	Position
}

type MovementPayload struct {
	UserID string `json:"userId"`
	Position
}

type MovementRejectedPayload struct {
	Position
	Reason string `json:"reason"`
}

type UserLeftPayload struct {
	UserID string `json:"userId"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func DecodeEvent(raw []byte) (Event, error) {
	var event Event
	if err := json.Unmarshal(raw, &event); err != nil {
		return Event{}, fmt.Errorf("%w: %v", ErrInvalidEvent, err)
	}

	event.Type = strings.TrimSpace(event.Type)
	if event.Type == "" || len(event.Payload) == 0 || string(event.Payload) == "null" {
		return Event{}, ErrInvalidEvent
	}

	return event, nil
}

func DecodePayload[T any](event Event) (T, error) {
	var payload T
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return payload, fmt.Errorf("%w: %v", ErrInvalidEvent, err)
	}
	return payload, nil
}

func ValidateClientEvent(event Event) error {
	switch event.Type {
	case EventJoin, EventMove:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrUnknownClientEvent, event.Type)
	}
}

func NewEvent(eventType string, payload any) (Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Event{}, err
	}
	return Event{Type: eventType, Payload: data}, nil
}

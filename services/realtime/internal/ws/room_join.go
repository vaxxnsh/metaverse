package ws

import (
	"context"
	"errors"
)

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrSpaceNotFound = errors.New("space not found")
)

type Claims struct{ Subject string }

type TokenVerifier interface{ Verify(string) (Claims, error) }

type SpaceProvider interface {
	GetRealtimeSpace(context.Context, string) (Space, error)
}

type StaticElement struct {
	ID            string
	X, Y          int32
	Width, Height int32
	Static        bool
}

type Space struct {
	ID            string
	Width, Height int32
	Elements      []StaticElement
}

// JoinHandler handles the authentication and room-entry parts of the protocol.
type JoinHandler struct {
	manager  *Manager
	verifier TokenVerifier
	spaces   SpaceProvider
}

func NewJoinHandler(manager *Manager, verifier TokenVerifier, spaces SpaceProvider) *JoinHandler {
	return &JoinHandler{manager: manager, verifier: verifier, spaces: spaces}
}

func (h *JoinHandler) Handle(client *Client, event Event) {
	if event.Type != EventJoin {
		return
	}
	payload, err := DecodePayload[JoinPayload](event)
	if err != nil {
		h.error(client, "INVALID_PAYLOAD", "invalid join payload")
		return
	}
	claims, err := h.verifier.Verify(payload.Token)
	if err != nil || claims.Subject == "" {
		h.error(client, "UNAUTHORIZED", "invalid or expired token")
		return
	}
	space, err := h.spaces.GetRealtimeSpace(context.Background(), payload.SpaceID)
	if err != nil {
		h.error(client, "SPACE_NOT_FOUND", "space does not exist")
		return
	}
	occupants, err := h.manager.Join(client, claims.Subject, space)
	if err != nil {
		h.error(client, "JOIN_FAILED", "unable to join space")
		return
	}
	joined, _ := NewEvent(EventSpaceJoined, SpaceJoinedPayload{SpaceID: space.ID, Spawn: Position{}, Users: occupants})
	if !client.Send(joined) {
		client.Close()
		return
	}
	presence, _ := NewEvent(EventUserJoined, UserJoinedPayload{UserID: claims.Subject, Position: Position{}})
	h.manager.Broadcast(space.ID, presence, client)
}

func (h *JoinHandler) error(client *Client, code, message string) {
	event, _ := NewEvent(EventError, ErrorPayload{Code: code, Message: message})
	if !client.Send(event) {
		client.Close()
	}
}

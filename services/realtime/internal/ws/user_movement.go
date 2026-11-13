package ws

import "errors"

// MovementHandler applies boundary-validated movement and publishes accepted
// state only to the other clients in the same room.
type MovementHandler struct {
	manager *Manager
}

func NewMovementHandler(manager *Manager) *MovementHandler {
	return &MovementHandler{manager: manager}
}

func (h *MovementHandler) Handle(client *Client, event Event) {
	if event.Type != EventMove {
		return
	}
	payload, err := DecodePayload[MovePayload](event)
	if err != nil {
		h.reject(client, Position{}, "invalid_payload")
		return
	}

	movement, err := h.manager.Move(client, Position{X: payload.X, Y: payload.Y})
	if err != nil {
		reason := "not_joined"
		if errors.Is(err, ErrOutsideBounds) {
			reason = "outside_bounds"
		}
		if errors.Is(err, ErrStaticCollision) {
			reason = "static_element_collision"
		}
		if errors.Is(err, ErrUserCollision) {
			reason = "user_collision"
		}
		if errors.Is(err, ErrMoveTooFar) {
			reason = "move_too_far"
		}
		if errors.Is(err, ErrMoveTooFast) {
			reason = "rate_limited"
		}
		h.reject(client, movement.Position, reason)
		return
	}

	eventToBroadcast, err := NewEvent(EventMovement, movement)
	if err != nil {
		return
	}
	roomID, ok := h.manager.ClientRoom(client)
	if !ok {
		return
	}
	h.manager.Broadcast(roomID, eventToBroadcast, client)
}

func (h *MovementHandler) reject(client *Client, position Position, reason string) {
	event, err := NewEvent(EventMovementRejected, MovementRejectedPayload{Position: position, Reason: reason})
	if err != nil || !client.Send(event) {
		client.Close()
	}
}

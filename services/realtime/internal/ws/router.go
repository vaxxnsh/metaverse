package ws

// EventRouter is the runtime dispatcher for supported client events.
type EventRouter struct {
	join     *JoinHandler
	movement *MovementHandler
}

func NewEventRouter(join *JoinHandler, movement *MovementHandler) *EventRouter {
	return &EventRouter{join: join, movement: movement}
}

func (r *EventRouter) Handle(client *Client, event Event) {
	switch event.Type {
	case EventJoin:
		r.join.Handle(client, event)
	case EventMove:
		r.movement.Handle(client, event)
	}
}

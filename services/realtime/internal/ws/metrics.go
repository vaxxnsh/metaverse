package ws

import "sync/atomic"

type Metrics struct {
	active                                 atomic.Int64
	connections, disconnects, joins, moves atomic.Uint64
}
type MetricsSnapshot struct {
	ActiveClients int64  `json:"activeClients"`
	Connections   uint64 `json:"connections"`
	Disconnects   uint64 `json:"disconnects"`
	JoinEvents    uint64 `json:"joinEvents"`
	MoveEvents    uint64 `json:"moveEvents"`
}

func NewMetrics() *Metrics       { return &Metrics{} }
func (m *Metrics) Connected()    { m.active.Add(1); m.connections.Add(1) }
func (m *Metrics) Disconnected() { m.active.Add(-1); m.disconnects.Add(1) }
func (m *Metrics) RecordEvent(eventType string) {
	if eventType == EventJoin {
		m.joins.Add(1)
	}
	if eventType == EventMove {
		m.moves.Add(1)
	}
}
func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{m.active.Load(), m.connections.Load(), m.disconnects.Load(), m.joins.Load(), m.moves.Load()}
}

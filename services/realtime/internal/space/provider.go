package space

import (
	"context"

	"github.com/vaxxnsh/metaverse/realtime/internal/ws"
	sharedrealtime "github.com/vaxxnsh/shared/realtime"
)

type Provider struct {
	postgres *sharedrealtime.PostgresProvider
}

func NewProvider(postgres *sharedrealtime.PostgresProvider) *Provider {
	return &Provider{postgres: postgres}
}

func (p *Provider) GetRealtimeSpace(ctx context.Context, spaceID string) (ws.Space, error) {
	space, err := p.postgres.GetRealtimeSpace(ctx, spaceID)
	if err != nil {
		return ws.Space{}, err
	}
	elements := make([]ws.StaticElement, 0, len(space.Elements))
	for _, element := range space.Elements {
		elements = append(elements, ws.StaticElement{ID: element.ID, X: element.X, Y: element.Y, Width: element.Width, Height: element.Height, Static: element.Static})
	}
	return ws.Space{ID: space.ID, Width: space.Width, Height: space.Height, Elements: elements}, nil
}

package domain

import (
	"fmt"
	"time"

	"github.com/vaxxnsh/metaverse/api/internal/db"
)

type Space struct {
	ID        string
	CreatorID string
	Name      string
	Width     int32
	Height    int32
	Thumbnail string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SpaceSummary struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Dimensions string `json:"dimensions"`
	Thumbnail  string `json:"thumbnail"`
}

type Element struct {
	ID        string
	Width     int32
	Height    int32
	ImageUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Map struct {
	ID        string
	Name      string
	Width     int32
	Height    int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SpaceElement struct {
	SpaceID   string
	ElementID string
	X         int32
	Y         int32
}

type MapElement struct {
	MapID     string
	ElementID string
	X         int32
	Y         int32
}

func DBSpaceToDomain(s db.Space) *Space {
	return &Space{
		ID:        s.ID.String(),
		CreatorID: s.CreatorID.String(),
		Name:      s.Name,
		Width:     s.Width,
		Height:    s.Height,
		Thumbnail: s.Thumbnail.String,
		CreatedAt: s.CreatedAt.Time,
		UpdatedAt: s.UpdatedAt.Time,
	}
}

func DBSpaceToSummary(s db.GetSpacesByCreatorRow) SpaceSummary {
	return SpaceSummary{
		ID:         s.ID.String(),
		Name:       s.Name,
		Dimensions: fmt.Sprintf("%dx%d", s.Width, s.Height),
		Thumbnail:  s.Thumbnail.String,
	}
}

func MapSpacesToSummary(spaces []db.GetSpacesByCreatorRow) []SpaceSummary {
	result := make([]SpaceSummary, 0, len(spaces))
	for _, s := range spaces {
		result = append(result, DBSpaceToSummary(s))
	}
	return result
}

func DBElementToDomain(e db.Element) *Element {
	return &Element{
		ID:        e.ID.String(),
		Width:     e.Width,
		Height:    e.Height,
		ImageUrl:  e.ImageUrl,
		CreatedAt: e.CreatedAt.Time,
		UpdatedAt: e.UpdatedAt.Time,
	}
}

func DBMapToDomain(m db.Map) *Map {
	return &Map{
		ID:        m.ID.String(),
		Name:      m.Name,
		Width:     m.Width,
		Height:    m.Height,
		CreatedAt: m.CreatedAt.Time,
		UpdatedAt: m.UpdatedAt.Time,
	}
}

func DBSpaceElementToDomain(se db.SpaceElement) *SpaceElement {
	return &SpaceElement{
		SpaceID:   se.SpaceID.String(),
		ElementID: se.ElementID.String(),
		X:         se.X,
		Y:         se.Y,
	}
}

func DBMapElementToDomain(me db.MapElement) *MapElement {
	return &MapElement{
		MapID:     me.MapID.String(),
		ElementID: me.ElementID.String(),
		X:         me.X,
		Y:         me.Y,
	}
}

func MapSpacesToDomain(spaces []db.Space) []Space {
	result := make([]Space, 0, len(spaces))
	for _, s := range spaces {
		result = append(result, *DBSpaceToDomain(s))
	}
	return result
}

func MapElementsToDomain(elements []db.Element) []Element {
	result := make([]Element, 0, len(elements))
	for _, e := range elements {
		result = append(result, *DBElementToDomain(e))
	}
	return result
}

func MapSpaceElementsToDomain(ses []db.SpaceElement) []SpaceElement {
	result := make([]SpaceElement, 0, len(ses))
	for _, se := range ses {
		result = append(result, *DBSpaceElementToDomain(se))
	}
	return result
}

func MapMapElementsToDomain(mes []db.MapElement) []MapElement {
	result := make([]MapElement, 0, len(mes))
	for _, me := range mes {
		result = append(result, *DBMapElementToDomain(me))
	}
	return result
}

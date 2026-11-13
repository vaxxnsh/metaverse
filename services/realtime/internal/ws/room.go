package ws

// Room contains the live users and immutable geometry for one active space.
// The manager owns its lifecycle; callers must access it through Manager.
type Room struct {
	ID       string
	Width    int32
	Height   int32
	Elements []StaticElement
	users    map[*Client]*User
}

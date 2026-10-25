package ws

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type HTTPHandlerOptions struct {
	Transport   TransportConfig
	CheckOrigin func(*http.Request) bool
	OnConnect   func(*Client)
	OnEvent     func(*Client, Event)
	OnClose     func(*Client)
}

func NewHTTPHandler(options HTTPHandlerOptions) http.Handler {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: options.CheckOrigin}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		client := newClient(conn, options.Transport, options.OnEvent, options.OnClose)
		if options.OnConnect != nil {
			options.OnConnect(client)
		}
		client.Start()
	})
}

// OriginAllowList rejects browser upgrades whose Origin is not configured.
func OriginAllowList(origins []string) func(*http.Request) bool {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowed[origin] = struct{}{}
		}
	}
	return func(r *http.Request) bool { _, ok := allowed[r.Header.Get("Origin")]; return ok }
}

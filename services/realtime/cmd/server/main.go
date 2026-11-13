package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vaxxnsh/metaverse/realtime/internal/auth"
	"github.com/vaxxnsh/metaverse/realtime/internal/space"
	"github.com/vaxxnsh/metaverse/realtime/internal/ws"
	sharedrealtime "github.com/vaxxnsh/shared/realtime"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	origins := os.Getenv("WS_ALLOWED_ORIGINS")
	if origins == "" {
		origins = "http://localhost:3000,http://localhost:8080"
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		slog.Error("DATABASE_URL is required")
		return
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecret"
	}
	pool, err := sharedrealtime.OpenPool(context.Background(), databaseURL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		return
	}
	defer pool.Close()
	manager := ws.NewManager()
	metrics := ws.NewMetrics()
	joinHandler := ws.NewJoinHandler(manager, auth.NewJWTVerifier(jwtSecret), space.NewProvider(sharedrealtime.NewPostgresProvider(pool)))
	router := ws.NewEventRouter(joinHandler, ws.NewMovementHandler(manager))

	handler := ws.NewHTTPHandler(ws.HTTPHandlerOptions{
		CheckOrigin: ws.OriginAllowList(strings.Split(origins, ",")),
		OnEvent:     func(client *ws.Client, event ws.Event) { metrics.RecordEvent(event.Type); router.Handle(client, event) },
		OnConnect: func(client *ws.Client) {
			manager.Register(client)
			metrics.Connected()
			slog.Info("websocket connected", "remoteAddr", client.RemoteAddr().String())
		},
		OnClose: func(client *ws.Client) {
			departure, ok := manager.Disconnect(client)
			if ok && departure.RoomID != "" {
				event, _ := ws.NewEvent(ws.EventUserLeft, ws.UserLeftPayload{UserID: departure.UserID})
				for _, recipient := range departure.Recipients {
					recipient.Send(event)
				}
			}
			metrics.Disconnected()
			slog.Info("websocket disconnected", "remoteAddr", client.RemoteAddr().String())
		},
	})
	mux := http.NewServeMux()
	mux.Handle("/ws", handler)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			http.Error(w, `{"status":"unavailable"}`, http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(metrics.Snapshot())
	})

	server := &http.Server{Addr: ":" + port, Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		slog.Info("realtime server listening", "port", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("realtime server error", "error", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	manager.Shutdown()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown error", "error", err)
	}
}

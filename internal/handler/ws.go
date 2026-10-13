package handler

import (
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/effect707/csinder/internal/auth"
	wsHub "github.com/effect707/csinder/internal/ws"
)

type WSHandler struct {
	hub *wsHub.Hub
}

func NewWSHandler(hub *wsHub.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) Handle(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("websocket accept error: %v", err)
		return
	}

	client := wsHub.NewClient(conn)
	h.hub.Register(user.ID, client)

	go client.WritePump(r.Context())

	client.ReadPump(r.Context())
	h.hub.Unregister(user.ID, client)
}

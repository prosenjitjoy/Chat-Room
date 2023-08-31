package ws

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"nhooyr.io/websocket"
)

type Handler struct {
	cs *ChatServer
}

func NewHandler(cs *ChatServer) *Handler {
	return &Handler{
		cs: cs,
	}
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var roomReq RoomRe
	err := json.NewDecoder(r.Body).Decode(&roomReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	h.cs.Rooms[roomReq.ID] = &Room{
		ID:          roomReq.ID,
		Name:        roomReq.Name,
		Subscribers: make(map[string]*Subscriber),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roomReq)
}

func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := make([]RoomRe, 0)
	for _, r := range h.cs.Rooms {
		rooms = append(rooms, RoomRe{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}

	roomID := chi.URLParam(r, "roomId")
	clientID := r.URL.Query().Get("userId")
	userName := r.URL.Query().Get("userName")

	subscriber := &Subscriber{
		ID:       clientID,
		RoomID:   roomID,
		UserName: userName,
		Messages: make(chan *Message, h.cs.MessageBuffer),
		Conn:     conn,
	}

	msg := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		UserName: userName,
	}

	h.cs.Register <- subscriber
	h.cs.Broadcast <- msg

	go subscriber.WriteMessage()
	subscriber.ReadMessage(h.cs)
}

func (h *Handler) GetSubscriber(w http.ResponseWriter, r *http.Request) {
	var subscribers []SubscriberRe
	roomID := chi.URLParam(r, "roomId")

	if _, ok := h.cs.Rooms[roomID]; !ok {
		subscribers = make([]SubscriberRe, 0)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(subscribers)
	}

	for _, s := range h.cs.Rooms[roomID].Subscribers {
		subscribers = append(subscribers, SubscriberRe{
			ID:       s.ID,
			UserName: s.UserName,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subscribers)
}

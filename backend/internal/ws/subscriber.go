package ws

import (
	"context"
	"log"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Subscriber struct {
	ID       string `json:""`
	RoomID   string `json:""`
	UserName string `json:""`
	Messages chan *Message
	Conn     *websocket.Conn
}

func (s *Subscriber) WriteMessage() {
	defer func() {
		s.Conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		message, ok := <-s.Messages
		if !ok {
			return
		}

		wsjson.Write(context.TODO(), s.Conn, message)
	}
}

func (s *Subscriber) ReadMessage(cs *ChatServer) {
	defer func() {
		cs.Unregister <- s
		s.Conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, m, err := s.Conn.Read(context.TODO())
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   s.RoomID,
			UserName: s.UserName,
		}

		cs.Broadcast <- msg
	}
}

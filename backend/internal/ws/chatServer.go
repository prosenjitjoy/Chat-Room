package ws

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ChatServer struct {
	MessageBuffer  int
	PublishLimiter *rate.Limiter

	Rooms   map[string]*Room
	RoomsMu sync.Mutex

	Register   chan *Subscriber
	Unregister chan *Subscriber
	Broadcast  chan *Message
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		MessageBuffer:  64,
		PublishLimiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
		Rooms:          make(map[string]*Room),
		Register:       make(chan *Subscriber),
		Unregister:     make(chan *Subscriber),
		Broadcast:      make(chan *Message, 5),
	}
}

func (cs *ChatServer) Run() {
	for {
		select {
		case s := <-cs.Register:
			if _, ok := cs.Rooms[s.RoomID]; ok {
				if _, ok := cs.Rooms[s.RoomID].Subscribers[s.ID]; !ok {
					cs.Rooms[s.RoomID].SubscribersMu.Lock()
					cs.Rooms[s.RoomID].Subscribers[s.ID] = s
					cs.Rooms[s.RoomID].SubscribersMu.Unlock()
				}
			}
		case s := <-cs.Unregister:
			if _, ok := cs.Rooms[s.RoomID]; ok {
				if _, ok := cs.Rooms[s.RoomID].Subscribers[s.ID]; ok {
					if len(cs.Rooms[s.RoomID].Subscribers) != 0 {
						cs.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomID:   s.RoomID,
							UserName: s.UserName,
						}
					}
					cs.Rooms[s.RoomID].SubscribersMu.Lock()
					delete(cs.Rooms[s.RoomID].Subscribers, s.ID)
					cs.Rooms[s.RoomID].SubscribersMu.Unlock()
				}
			}
		case m := <-cs.Broadcast:
			if _, ok := cs.Rooms[m.RoomID]; ok {
				cs.PublishLimiter.Wait(context.TODO())
				for _, s := range cs.Rooms[m.RoomID].Subscribers {
					s.Messages <- m
				}
			}
		}
	}
}

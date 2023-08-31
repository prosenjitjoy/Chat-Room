package ws

import "sync"

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	UserName string `json:"userName"`
}

type Room struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Subscribers   map[string]*Subscriber `json:"subscribers"`
	SubscribersMu sync.Mutex
}

type RoomRe struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SubscriberRe struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
}

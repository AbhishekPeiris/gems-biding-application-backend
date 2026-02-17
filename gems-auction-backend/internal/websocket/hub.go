package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Room manages clients for a single auction
type Room struct {
	auctionID  int64
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func newRoom(auctionID int64) *Room {
	return &Room{
		auctionID:  auctionID,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

func (r *Room) run() {
	for {
		select {
		case c := <-r.register:
			r.clients[c] = true

		case c := <-r.unregister:
			if _, ok := r.clients[c]; ok {
				delete(r.clients, c)
				close(c.send)
			}

		case msg := <-r.broadcast:
			for c := range r.clients {
				select {
				case c.send <- msg:
				default:
					// if client is slow, drop it
					delete(r.clients, c)
					close(c.send)
				}
			}
		}
	}
}

// Manager controls all rooms (auction rooms)
type Manager struct {
	mu    sync.RWMutex
	rooms map[int64]*Room
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[int64]*Room),
	}
}

func (m *Manager) getOrCreateRoom(auctionID int64) *Room {
	m.mu.RLock()
	room, ok := m.rooms[auctionID]
	m.mu.RUnlock()
	if ok {
		return room
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// re-check after lock
	room, ok = m.rooms[auctionID]
	if ok {
		return room
	}

	room = newRoom(auctionID)
	m.rooms[auctionID] = room
	go room.run()
	return room
}

// ---- WebSocket Upgrade ----

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production, restrict origins properly!
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeAuctionWS upgrades to websocket and joins auction room
func (m *Manager) ServeAuctionWS(c *gin.Context, auctionID int64) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// If auth middleware exists, it will set user_id.
	var userID int64 = 0
	if v, ok := c.Get("user_id"); ok {
		if id, ok2 := v.(int64); ok2 {
			userID = id
		}
	}

	room := m.getOrCreateRoom(auctionID)

	client := &Client{
		room:      room,
		conn:      conn,
		send:      make(chan []byte, 256),
		userID:    userID,
		auctionID: auctionID,
	}

	room.register <- client

	// send a welcome event
	_ = m.sendSystemEvent(auctionID, "WS_CONNECTED", gin.H{
		"user_id": userID,
	})

	go client.writePump()
	go client.readPump()
}

// BroadcastToAuction is used by your services to push events to frontend
func (m *Manager) BroadcastToAuction(auctionID int64, eventType string, payload any) {
	ev := Event{
		Type:      eventType,
		AuctionID: auctionID,
		Payload:   payload,
		Timestamp: time.Now(),
	}
	b, _ := json.Marshal(ev)
	room := m.getOrCreateRoom(auctionID)
	room.broadcast <- b
}

// helper: system events
func (m *Manager) sendSystemEvent(auctionID int64, eventType string, payload any) error {
	ev := Event{
		Type:      eventType,
		AuctionID: auctionID,
		Payload:   payload,
		Timestamp: time.Now(),
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	room := m.getOrCreateRoom(auctionID)
	room.broadcast <- b
	return nil
}

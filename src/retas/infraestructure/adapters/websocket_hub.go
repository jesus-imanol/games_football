package adapters

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	PongWait   = 120 * time.Second // Tiempo máximo sin recibir pong del cliente
	pingPeriod = 30 * time.Second  // Intervalo de envío de pings (debe ser < PongWait)
	writeWait  = 10 * time.Second  // Tiempo máximo para escribir un mensaje
)

// Client representa un cliente conectado al WebSocket
type Client struct {
	Conn   *websocket.Conn
	ZonaID string
	Send   chan []byte
}

// Hub mantiene el conjunto de clientes activos y difunde mensajes por zona
type Hub struct {
	// Clientes registrados agrupados por zona_id
	clients map[string]map[*Client]bool

	// Canal para registrar clientes
	register chan *Client

	// Canal para desregistrar clientes
	unregister chan *Client

	// Canal para broadcast de mensajes
	broadcast chan *BroadcastRequest

	// Mutex para sincronización
	mu sync.RWMutex
}

// BroadcastRequest contiene el mensaje y la zona a la que se enviará
type BroadcastRequest struct {
	ZonaID  string
	Message []byte
}

// NewHub crea una nueva instancia del Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastRequest),
	}
}

// Run ejecuta el hub en un goroutine
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.ZonaID]; !ok {
				h.clients[client.ZonaID] = make(map[*Client]bool)
			}
			h.clients[client.ZonaID][client] = true
			h.mu.Unlock()
			log.Printf("Cliente registrado en zona: %s. Total clientes en zona: %d", client.ZonaID, len(h.clients[client.ZonaID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.ZonaID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.clients, client.ZonaID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Cliente desregistrado de zona: %s", client.ZonaID)

		case broadcastReq := <-h.broadcast:
			h.mu.Lock()
			if clients, ok := h.clients[broadcastReq.ZonaID]; ok {
				for client := range clients {
					select {
					case client.Send <- broadcastReq.Message:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// RegisterClient registra un nuevo cliente en el hub
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient desregistra un cliente del hub
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastToZone envía un mensaje a todos los clientes de una zona específica
func (h *Hub) BroadcastToZone(zonaID string, message interface{}) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.broadcast <- &BroadcastRequest{
		ZonaID:  zonaID,
		Message: messageBytes,
	}

	return nil
}

// WritePump envía mensajes del hub al cliente websocket y mantiene la conexión viva con pings
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// El hub cerró el canal
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error escribiendo mensaje: %v", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error enviando ping: %v", err)
				return
			}
		}
	}
}

package controllers

import (
	"encoding/json"
	"games-football-api/src/retas/application"
	"games-football-api/src/retas/domain/entities"
	"games-football-api/src/retas/infraestructure/adapters"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir todas las conexiones en desarrollo
	},
}

type WebSocketController struct {
	hub                    *adapters.Hub
	unirseUseCase          *application.UnirseRetaUseCase
	crearRetaUseCase       *application.CrearRetaUseCase
	obtenerRetasUseCase    *application.ObtenerRetasPorZonaUseCase
	enviarMensajeUseCase   *application.EnviarMensajeUseCase
}

func NewWebSocketController(hub *adapters.Hub, unirseUseCase *application.UnirseRetaUseCase, crearRetaUseCase *application.CrearRetaUseCase, obtenerRetasUseCase *application.ObtenerRetasPorZonaUseCase, enviarMensajeUseCase *application.EnviarMensajeUseCase) *WebSocketController {
	return &WebSocketController{
		hub:                    hub,
		unirseUseCase:          unirseUseCase,
		crearRetaUseCase:       crearRetaUseCase,
		obtenerRetasUseCase:    obtenerRetasUseCase,
		enviarMensajeUseCase:   enviarMensajeUseCase,
	}
}

// HandleWebSocket maneja las conexiones WebSocket
func (wsc *WebSocketController) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error al actualizar a WebSocket: %v", err)
		return
	}

	client := &adapters.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	defer func() {
		if client.ZonaID != "" {
			wsc.hub.UnregisterClient(client)
		}
		conn.Close()
	}()

	// Iniciar escritura en goroutine
	go client.WritePump()

	// Configurar timeouts y pong handler
	conn.SetReadDeadline(time.Now().Add(adapters.PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(adapters.PongWait))
		return nil
	})

	// Leer mensajes del cliente
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Printf("Error inesperado: %v", err)
			}
			break
		}

		// Resetear el deadline con cada mensaje recibido (mantiene la conexión viva)
		conn.SetReadDeadline(time.Now().Add(adapters.PongWait))

		// Parsear el mensaje
		var wsMsg entities.WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Error al parsear mensaje: %v", err)
			wsc.sendError(client, "Formato de mensaje inválido")
			continue
		}

		// Si es la primera vez que el cliente envía un mensaje, registrarlo en su zona
		if client.ZonaID == "" && wsMsg.ZonaID != "" {
			client.ZonaID = wsMsg.ZonaID
			wsc.hub.RegisterClient(client)

			// Enviar las retas existentes de esta zona al cliente recién conectado
			retas, err := wsc.obtenerRetasUseCase.Execute(client.ZonaID)
			if err != nil {
				log.Printf("Error al obtener retas de zona %s: %v", client.ZonaID, err)
			} else {
				initMsg := entities.BroadcastMessage{
					Status: "retas_zona",
					Retas:  retas,
				}
				msgBytes, _ := json.Marshal(initMsg)
				select {
				case client.Send <- msgBytes:
				default:
				}
			}
		}

		// Enrutar según la acción
		switch wsMsg.Accion {
		case "unirse":
			wsc.handleUnirse(client, wsMsg)
		case "crear":
			wsc.handleCrear(client, wsMsg)
		case "enviar_mensaje":
			wsc.handleEnviarMensaje(client, wsMsg)
		default:
			wsc.sendError(client, "Acción no reconocida")
		}
	}
}

// handleUnirse maneja la acción de unirse a una reta
func (wsc *WebSocketController) handleUnirse(client *adapters.Client, msg entities.WebSocketMessage) {
	// Validar campos necesarios
	if msg.RetaID == "" || msg.UsuarioID == "" || msg.Nombre == "" {
		wsc.sendError(client, "Campos requeridos: reta_id, usuario_id, nombre")
		return
	}

	// Ejecutar el caso de uso
	jugadoresActuales, listaJugadores, err := wsc.unirseUseCase.Execute(msg.RetaID, msg.UsuarioID, msg.Nombre)
	if err != nil {
		wsc.sendError(client, err.Error())
		return
	}

	// Broadcast a todos los clientes de la zona
	broadcastMsg := entities.BroadcastMessage{
		Status:            "actualizacion",
		RetaID:            msg.RetaID,
		JugadoresActuales: jugadoresActuales,
		ListaJugadores:    listaJugadores,
	}

	if err := wsc.hub.BroadcastToZone(client.ZonaID, broadcastMsg); err != nil {
		log.Printf("Error al hacer broadcast: %v", err)
	}
}

// handleCrear maneja la acción de crear una nueva reta
func (wsc *WebSocketController) handleCrear(client *adapters.Client, msg entities.WebSocketMessage) {
	// Validar campos necesarios
	if msg.Titulo == "" || msg.FechaHora == "" || msg.MaxJugadores == 0 || msg.CreadorNombre == "" {
		wsc.sendError(client, "Campos requeridos: titulo, fecha_hora, max_jugadores, creador_nombre")
		return
	}

	// Generar creador_id automáticamente si no se envía
	creadorID := msg.CreadorID
	if creadorID == "" {
		creadorID = uuid.New().String()
	}

	// Ejecutar el caso de uso
	retaCreada, primerJugador, err := wsc.crearRetaUseCase.Execute(
		msg.ZonaID,
		msg.Titulo,
		msg.FechaHora,
		msg.MaxJugadores,
		creadorID,
		msg.CreadorNombre,
	)
	if err != nil {
		wsc.sendError(client, err.Error())
		return
	}

	// Preparar mensaje de broadcast
	listaJugadores := []entities.Jugador{*primerJugador}

	broadcastMsg := entities.BroadcastMessage{
		Status: "nueva_reta",
		Reta: &entities.RetaInfo{
			ID:                retaCreada.ID,
			Titulo:            retaCreada.Titulo,
			FechaHora:         retaCreada.FechaHora.Format("2006-01-02 15:04:05"),
			MaxJugadores:      retaCreada.MaxJugadores,
			JugadoresActuales: retaCreada.JugadoresActuales,
			ListaJugadores:    listaJugadores,
		},
	}

	// Broadcast a todos los clientes de la zona
	if err := wsc.hub.BroadcastToZone(client.ZonaID, broadcastMsg); err != nil {
		log.Printf("Error al hacer broadcast: %v", err)
	}
}

// sendError envía un mensaje de error solo al cliente específico
func (wsc *WebSocketController) sendError(client *adapters.Client, mensaje string) {
	errorMsg := entities.BroadcastMessage{
		Status:  "error",
		Mensaje: mensaje,
	}

	msgBytes, err := json.Marshal(errorMsg)
	if err != nil {
		log.Printf("Error al serializar mensaje de error: %v", err)
		return
	}

	select {
	case client.Send <- msgBytes:
	default:
		log.Printf("No se pudo enviar mensaje de error al cliente")
	}
}

// handleEnviarMensaje maneja la acción de enviar un mensaje al chat en vivo de una reta
func (wsc *WebSocketController) handleEnviarMensaje(client *adapters.Client, msg entities.WebSocketMessage) {
	// Validar campos necesarios
	if msg.RetaID == "" || msg.UsuarioID == "" || msg.Texto == "" {
		wsc.sendError(client, "Campos requeridos: reta_id, usuario_id, texto")
		return
	}

	// Ejecutar el caso de uso
	mensaje, err := wsc.enviarMensajeUseCase.Execute(msg.RetaID, msg.UsuarioID, msg.Texto)
	if err != nil {
		wsc.sendError(client, err.Error())
		return
	}

	// Broadcast a todos los clientes de la zona
	broadcastMsg := entities.BroadcastMessage{
		Status:      "nuevo_mensaje",
		RetaID:      msg.RetaID,
		MensajeChat: mensaje,
	}

	if err := wsc.hub.BroadcastToZone(client.ZonaID, broadcastMsg); err != nil {
		log.Printf("Error al hacer broadcast de mensaje: %v", err)
	}
}

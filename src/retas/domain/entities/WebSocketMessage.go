package entities

// WebSocketMessage representa el mensaje que se recibe del cliente
type WebSocketMessage struct {
	Accion    string `json:"accion"` // "unirse", "crear" o "enviar_mensaje"
	UsuarioID string `json:"usuario_id,omitempty"`
	Nombre    string `json:"nombre,omitempty"`
	RetaID    string `json:"reta_id,omitempty"`
	ZonaID    string `json:"zona_id"`

	// Campos específicos para "crear"
	Titulo        string `json:"titulo,omitempty"`
	FechaHora     string `json:"fecha_hora,omitempty"`
	MaxJugadores  int    `json:"max_jugadores,omitempty"`
	CreadorID     string `json:"creador_id,omitempty"`
	CreadorNombre string `json:"creador_nombre,omitempty"`

	// Campos específicos para "enviar_mensaje"
	Texto string `json:"texto,omitempty"`
}

// BroadcastMessage representa los mensajes de broadcast
type BroadcastMessage struct {
	Status            string     `json:"status"`
	RetaID            string     `json:"reta_id,omitempty"`
	JugadoresActuales int        `json:"jugadores_actuales,omitempty"`
	ListaJugadores    []Jugador  `json:"lista_jugadores,omitempty"`
	Mensaje           string     `json:"mensaje,omitempty"`
	Reta              *RetaInfo  `json:"reta,omitempty"`
	Retas             []RetaInfo `json:"retas,omitempty"`
	MensajeChat       *Mensaje   `json:"mensaje_chat,omitempty"`
}

// RetaInfo para el mensaje de nueva reta
type RetaInfo struct {
	ID                string    `json:"id"`
	Titulo            string    `json:"titulo"`
	FechaHora         string    `json:"fecha_hora"`
	MaxJugadores      int       `json:"max_jugadores"`
	JugadoresActuales int       `json:"jugadores_actuales"`
	ListaJugadores    []Jugador `json:"lista_jugadores"`
	HistorialChat     []Mensaje `json:"historial_chat"`
}

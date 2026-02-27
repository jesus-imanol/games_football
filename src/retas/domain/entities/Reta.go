package entities

import "time"

type Reta struct {
	ID                string    `json:"id"`
	ZonaID            string    `json:"zona_id"`
	Titulo            string    `json:"titulo"`
	FechaHora         time.Time `json:"fecha_hora"`
	MaxJugadores      int       `json:"max_jugadores"`
	JugadoresActuales int       `json:"jugadores_actuales"`
	CreadorID         string    `json:"creador_id"`
	CreadorNombre     string    `json:"creador_nombre"`
	CreatedAt         time.Time `json:"created_at"`
	HistorialChat     []Mensaje `json:"historial_chat,omitempty"`
}

func NewReta(zonaID, titulo, fechaHoraStr string, maxJugadores int, creadorID, creadorNombre string) (*Reta, error) {
	fechaHora, err := time.Parse("2006-01-02 15:04:05", fechaHoraStr)
	if err != nil {
		return nil, err
	}

	return &Reta{
		ZonaID:            zonaID,
		Titulo:            titulo,
		FechaHora:         fechaHora,
		MaxJugadores:      maxJugadores,
		JugadoresActuales: 0,
		CreadorID:         creadorID,
		CreadorNombre:     creadorNombre,
		CreatedAt:         time.Now(),
	}, nil
}

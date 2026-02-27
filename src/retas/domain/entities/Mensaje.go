package entities

import "time"

// Mensaje representa un mensaje del chat en vivo de una reta
type Mensaje struct {
	ID            string    `json:"id"`
	RetaID        string    `json:"reta_id"`
	UsuarioID     string    `json:"usuario_id"`
	NombreUsuario string    `json:"nombre"`
	Texto         string    `json:"texto"`
	Timestamp     time.Time `json:"timestamp"`
}

func NewMensaje(retaID, usuarioID, texto string) *Mensaje {
	return &Mensaje{
		RetaID:    retaID,
		UsuarioID: usuarioID,
		Texto:     texto,
		Timestamp: time.Now(),
	}
}

package entities

type Jugador struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	RetaID    string `json:"reta_id,omitempty"`
	UsuarioID string `json:"usuario_id,omitempty"`
}

func NewJugador(usuarioID, nombre string) *Jugador {
	return &Jugador{
		UsuarioID: usuarioID,
		Nombre:    nombre,
	}
}

package repositories

import (
	"games-football-api/src/retas/domain/entities"
)

// IRetaRepository define la interfaz para operaciones de retas
type IRetaRepository interface {
	// UnirseReta realiza la lógica de unirse a una reta con transacción y bloqueo
	UnirseReta(retaID, usuarioID, nombreJugador string) (jugadoresActuales int, listaJugadores []entities.Jugador, err error)

	// CrearReta crea una nueva reta e inserta al creador como primer jugador
	CrearReta(reta *entities.Reta) (retaCreada *entities.Reta, primerJugador *entities.Jugador, err error)

	// ObtenerJugadoresDeReta obtiene la lista de jugadores confirmados de una reta
	ObtenerJugadoresDeReta(retaID string) ([]entities.Jugador, error)
}

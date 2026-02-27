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

	// ObtenerRetasPorZona obtiene todas las retas activas de una zona con sus jugadores
	ObtenerRetasPorZona(zonaID string) ([]entities.RetaInfo, error)

	// GuardarMensaje persiste un mensaje de chat y retorna el mensaje enriquecido con nombre de usuario
	GuardarMensaje(mensaje entities.Mensaje) (*entities.Mensaje, error)

	// ObtenerMensajesDeReta obtiene el historial de mensajes de una reta
	ObtenerMensajesDeReta(retaID string) ([]entities.Mensaje, error)
}

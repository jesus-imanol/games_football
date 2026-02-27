package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"games-football-api/src/retas/domain/entities"
	"time"

	"github.com/google/uuid"
)

type MySQLRetaRepository struct {
	db *sql.DB
}

func NewMySQLRetaRepository(db *sql.DB) *MySQLRetaRepository {
	return &MySQLRetaRepository{
		db: db,
	}
}

// UnirseReta implementa la lógica de unirse a una reta con transacción y bloqueo
func (repo *MySQLRetaRepository) UnirseReta(retaID, usuarioID, nombreJugador string) (int, []entities.Jugador, error) {
	// Iniciar transacción
	tx, err := repo.db.Begin()
	if err != nil {
		return 0, nil, fmt.Errorf("error al iniciar transacción: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// SELECT FOR UPDATE para bloquear la fila
	var jugadoresActuales, maxJugadores int
	query := "SELECT jugadores_actuales, max_jugadores FROM retas WHERE id = ? FOR UPDATE"
	err = tx.QueryRow(query, retaID).Scan(&jugadoresActuales, &maxJugadores)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil, errors.New("reta no encontrada")
		}
		return 0, nil, fmt.Errorf("error al consultar reta: %w", err)
	}

	// Verificar si la reta está llena
	if jugadoresActuales >= maxJugadores {
		tx.Rollback()
		return 0, nil, errors.New("reta llena")
	}

	// Validar que el usuario_id exista en la tabla usuarios
	var existeUsuario int
	checkUsuarioQuery := "SELECT COUNT(*) FROM usuarios WHERE id = ?"
	err = tx.QueryRow(checkUsuarioQuery, usuarioID).Scan(&existeUsuario)
	if err != nil {
		return 0, nil, fmt.Errorf("error al verificar usuario: %w", err)
	}
	if existeUsuario == 0 {
		tx.Rollback()
		return 0, nil, errors.New("el usuario no existe")
	}

	// Verificar si el usuario ya está inscrito en esta reta
	var existeJugador int
	checkQuery := "SELECT COUNT(*) FROM reta_jugadores WHERE reta_id = ? AND usuario_id = ?"
	err = tx.QueryRow(checkQuery, retaID, usuarioID).Scan(&existeJugador)
	if err != nil {
		return 0, nil, fmt.Errorf("error al verificar jugador: %w", err)
	}

	if existeJugador > 0 {
		tx.Rollback()
		return 0, nil, errors.New("el usuario ya está inscrito en esta reta")
	}

	// Incrementar el contador de jugadores
	updateQuery := "UPDATE retas SET jugadores_actuales = jugadores_actuales + 1 WHERE id = ?"
	_, err = tx.Exec(updateQuery, retaID)
	if err != nil {
		return 0, nil, fmt.Errorf("error al actualizar contador: %w", err)
	}

	// Insertar al jugador
	jugadorID := uuid.New().String()
	insertQuery := "INSERT INTO reta_jugadores (id, reta_id, usuario_id, nombre_jugador) VALUES (?, ?, ?, ?)"
	_, err = tx.Exec(insertQuery, jugadorID, retaID, usuarioID, nombreJugador)
	if err != nil {
		return 0, nil, fmt.Errorf("error al insertar jugador: %w", err)
	}

	// Commit de la transacción
	err = tx.Commit()
	if err != nil {
		return 0, nil, fmt.Errorf("error al hacer commit: %w", err)
	}

	// Obtener la lista actualizada de jugadores
	listaJugadores, err := repo.ObtenerJugadoresDeReta(retaID)
	if err != nil {
		return 0, nil, fmt.Errorf("error al obtener lista de jugadores: %w", err)
	}

	return jugadoresActuales + 1, listaJugadores, nil
}

// CrearReta crea una nueva reta e inserta al creador como primer jugador
func (repo *MySQLRetaRepository) CrearReta(reta *entities.Reta) (*entities.Reta, *entities.Jugador, error) {
	// Iniciar transacción
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("error al iniciar transacción: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Generar UUID para la reta
	retaID := uuid.New().String()
	reta.ID = retaID

	// Insertar la reta
	insertRetaQuery := `
		INSERT INTO retas (id, zona_id, titulo, fecha_hora, max_jugadores, jugadores_actuales, creador_id, creador_nombre, created_at)
		VALUES (?, ?, ?, ?, ?, 1, ?, ?, NOW())
	`
	_, err = tx.Exec(insertRetaQuery, reta.ID, reta.ZonaID, reta.Titulo, reta.FechaHora, reta.MaxJugadores, reta.CreadorID, reta.CreadorNombre)
	if err != nil {
		return nil, nil, fmt.Errorf("error al insertar reta: %w", err)
	}

	// Insertar al creador como primer jugador
	jugadorID := uuid.New().String()
	insertJugadorQuery := "INSERT INTO reta_jugadores (id, reta_id, usuario_id, nombre_jugador) VALUES (?, ?, ?, ?)"
	_, err = tx.Exec(insertJugadorQuery, jugadorID, retaID, reta.CreadorID, reta.CreadorNombre)
	if err != nil {
		return nil, nil, fmt.Errorf("error al insertar primer jugador: %w", err)
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return nil, nil, fmt.Errorf("error al hacer commit: %w", err)
	}

	reta.JugadoresActuales = 1

	primerJugador := &entities.Jugador{
		ID:        jugadorID,
		UsuarioID: reta.CreadorID,
		Nombre:    reta.CreadorNombre,
		RetaID:    retaID,
	}

	return reta, primerJugador, nil
}

// ObtenerRetasPorZona obtiene todas las retas de una zona con sus jugadores
func (repo *MySQLRetaRepository) ObtenerRetasPorZona(zonaID string) ([]entities.RetaInfo, error) {
	query := `
		SELECT r.id, r.titulo, r.fecha_hora, r.max_jugadores, r.jugadores_actuales,
		       rj.id as jugador_id, rj.usuario_id, u.nombre
		FROM retas r
		LEFT JOIN reta_jugadores rj ON r.id = rj.reta_id
		LEFT JOIN usuarios u ON rj.usuario_id = u.id
		WHERE r.zona_id = ?
		ORDER BY r.created_at DESC, rj.created_at ASC
	`
	rows, err := repo.db.Query(query, zonaID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar retas: %w", err)
	}
	defer rows.Close()

	retasMap := make(map[string]*entities.RetaInfo)
	orden := []string{}

	for rows.Next() {
		var retaID, titulo string
		var fechaHora time.Time
		var maxJugadores, jugadoresActuales int
		var jugadorID, usuarioID, nombreJugador sql.NullString

		err := rows.Scan(&retaID, &titulo, &fechaHora, &maxJugadores, &jugadoresActuales,
			&jugadorID, &usuarioID, &nombreJugador)
		if err != nil {
			return nil, fmt.Errorf("error al escanear reta: %w", err)
		}

		if _, exists := retasMap[retaID]; !exists {
			retasMap[retaID] = &entities.RetaInfo{
				ID:                retaID,
				Titulo:            titulo,
				FechaHora:         fechaHora.Format("2006-01-02 15:04:05"),
				MaxJugadores:      maxJugadores,
				JugadoresActuales: jugadoresActuales,
				ListaJugadores:    []entities.Jugador{},
			}
			orden = append(orden, retaID)
		}

		if jugadorID.Valid {
			retasMap[retaID].ListaJugadores = append(retasMap[retaID].ListaJugadores, entities.Jugador{
				ID:        jugadorID.String,
				UsuarioID: usuarioID.String,
				Nombre:    nombreJugador.String,
				RetaID:    retaID,
			})
		}
	}

	result := make([]entities.RetaInfo, 0, len(orden))
	for _, id := range orden {
		// Obtener historial de chat para cada reta
		mensajes, err := repo.ObtenerMensajesDeReta(id)
		if err != nil {
			mensajes = []entities.Mensaje{}
		}
		retasMap[id].HistorialChat = mensajes
		result = append(result, *retasMap[id])
	}
	return result, nil
}

// ObtenerJugadoresDeReta obtiene la lista de jugadores confirmados con nombre real de usuarios
func (repo *MySQLRetaRepository) ObtenerJugadoresDeReta(retaID string) ([]entities.Jugador, error) {
	query := `
		SELECT rj.id, rj.usuario_id, u.nombre
		FROM reta_jugadores rj
		INNER JOIN usuarios u ON rj.usuario_id = u.id
		WHERE rj.reta_id = ?
		ORDER BY rj.created_at ASC
	`
	rows, err := repo.db.Query(query, retaID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar jugadores: %w", err)
	}
	defer rows.Close()

	jugadores := make([]entities.Jugador, 0)
	for rows.Next() {
		var jugador entities.Jugador
		err := rows.Scan(&jugador.ID, &jugador.UsuarioID, &jugador.Nombre)
		if err != nil {
			return nil, fmt.Errorf("error al escanear jugador: %w", err)
		}
		jugadores = append(jugadores, jugador)
	}

	return jugadores, nil
}

// GuardarMensaje inserta un mensaje de chat y retorna el mensaje con el nombre real del usuario (JOIN)
func (repo *MySQLRetaRepository) GuardarMensaje(mensaje entities.Mensaje) (*entities.Mensaje, error) {
	mensajeID := uuid.New().String()

	insertQuery := "INSERT INTO mensajes_reta (id, reta_id, usuario_id, texto) VALUES (?, ?, ?, ?)"
	_, err := repo.db.Exec(insertQuery, mensajeID, mensaje.RetaID, mensaje.UsuarioID, mensaje.Texto)
	if err != nil {
		return nil, fmt.Errorf("error al guardar mensaje: %w", err)
	}

	// Recuperar el mensaje con JOIN a usuarios para obtener el nombre real
	selectQuery := `
		SELECT m.id, m.reta_id, m.usuario_id, u.nombre, m.texto, m.creado_en
		FROM mensajes_reta m
		INNER JOIN usuarios u ON m.usuario_id = u.id
		WHERE m.id = ?
	`
	var resultado entities.Mensaje
	err = repo.db.QueryRow(selectQuery, mensajeID).Scan(
		&resultado.ID, &resultado.RetaID, &resultado.UsuarioID,
		&resultado.NombreUsuario, &resultado.Texto, &resultado.Timestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("error al recuperar mensaje enriquecido: %w", err)
	}

	return &resultado, nil
}

// ObtenerMensajesDeReta obtiene el historial completo de mensajes de una reta
func (repo *MySQLRetaRepository) ObtenerMensajesDeReta(retaID string) ([]entities.Mensaje, error) {
	query := `
		SELECT m.id, m.reta_id, m.usuario_id, u.nombre, m.texto, m.creado_en
		FROM mensajes_reta m
		INNER JOIN usuarios u ON m.usuario_id = u.id
		WHERE m.reta_id = ?
		ORDER BY m.creado_en ASC
	`
	rows, err := repo.db.Query(query, retaID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar mensajes: %w", err)
	}
	defer rows.Close()

	mensajes := make([]entities.Mensaje, 0)
	for rows.Next() {
		var msg entities.Mensaje
		err := rows.Scan(&msg.ID, &msg.RetaID, &msg.UsuarioID, &msg.NombreUsuario, &msg.Texto, &msg.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("error al escanear mensaje: %w", err)
		}
		mensajes = append(mensajes, msg)
	}

	return mensajes, nil
}

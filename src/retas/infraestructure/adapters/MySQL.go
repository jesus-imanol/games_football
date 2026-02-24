package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"games-football-api/src/retas/domain/entities"

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

// ObtenerJugadoresDeReta obtiene la lista de jugadores confirmados
func (repo *MySQLRetaRepository) ObtenerJugadoresDeReta(retaID string) ([]entities.Jugador, error) {
	query := "SELECT id, usuario_id, nombre_jugador FROM reta_jugadores WHERE reta_id = ? ORDER BY created_at ASC"
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

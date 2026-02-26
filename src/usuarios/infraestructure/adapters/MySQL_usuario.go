package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"games-football-api/src/usuarios/domain/entities"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MySQLUsuarioRepository struct {
	db *sql.DB
}

func NewMySQLUsuarioRepository(db *sql.DB) *MySQLUsuarioRepository {
	return &MySQLUsuarioRepository{
		db: db,
	}
}

// Login busca un usuario por username y compara el hash de la password
func (repo *MySQLUsuarioRepository) Login(username, password string) (*entities.Usuario, error) {
	query := "SELECT id, username, password, nombre FROM usuarios WHERE username = ?"
	row := repo.db.QueryRow(query, username)

	var usuario entities.Usuario
	var hashedPassword string
	err := row.Scan(&usuario.ID, &usuario.Username, &hashedPassword, &usuario.Nombre)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("credenciales inválidas")
		}
		return nil, fmt.Errorf("error al consultar usuario: %w", err)
	}

	// Comparar la password ingresada con el hash almacenado
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	return &usuario, nil
}

// Register crea un nuevo usuario con password hasheada en la base de datos
func (repo *MySQLUsuarioRepository) Register(username, password, nombre string) (*entities.Usuario, error) {
	id := uuid.New().String()

	// Hashear la password con bcrypt (cost por defecto = 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al hashear password: %w", err)
	}

	query := "INSERT INTO usuarios (id, username, password, nombre) VALUES (?, ?, ?, ?)"
	_, err = repo.db.Exec(query, id, username, string(hashedPassword), nombre)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return nil, errors.New("el username ya está registrado")
		}
		return nil, fmt.Errorf("error al registrar usuario: %w", err)
	}

	return &entities.Usuario{
		ID:       id,
		Username: username,
		Nombre:   nombre,
	}, nil
}

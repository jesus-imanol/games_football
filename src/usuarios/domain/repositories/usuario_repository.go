package repositories

import (
	"games-football-api/src/usuarios/domain/entities"
)

// IUsuarioRepository define la interfaz para operaciones de usuarios
type IUsuarioRepository interface {
	// Login busca un usuario por username y password, retorna el usuario si hace match
	Login(username, password string) (*entities.Usuario, error)

	// Register crea un nuevo usuario y retorna el usuario creado
	Register(username, password, nombre string) (*entities.Usuario, error)
}

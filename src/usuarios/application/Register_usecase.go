package application

import (
	"errors"
	"games-football-api/src/usuarios/domain/entities"
	"games-football-api/src/usuarios/domain/repositories"
)

type RegisterUseCase struct {
	usuarioRepo repositories.IUsuarioRepository
}

func NewRegisterUseCase(usuarioRepo repositories.IUsuarioRepository) *RegisterUseCase {
	return &RegisterUseCase{
		usuarioRepo: usuarioRepo,
	}
}

func (uc *RegisterUseCase) Execute(username, password, nombre string) (*entities.Usuario, error) {
	if username == "" || password == "" || nombre == "" {
		return nil, errors.New("username, password y nombre son requeridos")
	}

	usuario, err := uc.usuarioRepo.Register(username, password, nombre)
	if err != nil {
		return nil, err
	}

	return usuario, nil
}

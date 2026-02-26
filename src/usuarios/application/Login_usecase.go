package application

import (
	"errors"
	"games-football-api/src/usuarios/domain/entities"
	"games-football-api/src/usuarios/domain/repositories"
)

type LoginUseCase struct {
	usuarioRepo repositories.IUsuarioRepository
}

func NewLoginUseCase(usuarioRepo repositories.IUsuarioRepository) *LoginUseCase {
	return &LoginUseCase{
		usuarioRepo: usuarioRepo,
	}
}

func (uc *LoginUseCase) Execute(username, password string) (*entities.Usuario, error) {
	if username == "" || password == "" {
		return nil, errors.New("username y password son requeridos")
	}

	usuario, err := uc.usuarioRepo.Login(username, password)
	if err != nil {
		return nil, err
	}

	return usuario, nil
}

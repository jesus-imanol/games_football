package application

import (
	"games-football-api/src/retas/domain/entities"
	"games-football-api/src/retas/domain/repositories"
)

type UnirseRetaUseCase struct {
	retaRepo repositories.IRetaRepository
}

func NewUnirseRetaUseCase(retaRepo repositories.IRetaRepository) *UnirseRetaUseCase {
	return &UnirseRetaUseCase{
		retaRepo: retaRepo,
	}
}

func (uc *UnirseRetaUseCase) Execute(retaID, usuarioID, nombreJugador string) (int, []entities.Jugador, error) {
	// El repositorio maneja la transacción con SELECT FOR UPDATE y toda la lógica
	jugadoresActuales, listaJugadores, err := uc.retaRepo.UnirseReta(retaID, usuarioID, nombreJugador)
	if err != nil {
		return 0, nil, err
	}

	return jugadoresActuales, listaJugadores, nil
}

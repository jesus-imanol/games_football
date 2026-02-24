package application

import (
	"games-football-api/src/retas/domain/entities"
	"games-football-api/src/retas/domain/repositories"
)

type CrearRetaUseCase struct {
	retaRepo repositories.IRetaRepository
}

func NewCrearRetaUseCase(retaRepo repositories.IRetaRepository) *CrearRetaUseCase {
	return &CrearRetaUseCase{
		retaRepo: retaRepo,
	}
}

func (uc *CrearRetaUseCase) Execute(zonaID, titulo, fechaHora string, maxJugadores int, creadorID, creadorNombre string) (*entities.Reta, *entities.Jugador, error) {
	// Crear la entidad Reta
	reta, err := entities.NewReta(zonaID, titulo, fechaHora, maxJugadores, creadorID, creadorNombre)
	if err != nil {
		return nil, nil, err
	}

	// El repositorio crea la reta e inserta al creador como primer jugador
	retaCreada, primerJugador, err := uc.retaRepo.CrearReta(reta)
	if err != nil {
		return nil, nil, err
	}

	return retaCreada, primerJugador, nil
}

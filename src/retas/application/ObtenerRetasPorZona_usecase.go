package application

import (
	"games-football-api/src/retas/domain/entities"
	"games-football-api/src/retas/domain/repositories"
)

type ObtenerRetasPorZonaUseCase struct {
	retaRepo repositories.IRetaRepository
}

func NewObtenerRetasPorZonaUseCase(retaRepo repositories.IRetaRepository) *ObtenerRetasPorZonaUseCase {
	return &ObtenerRetasPorZonaUseCase{
		retaRepo: retaRepo,
	}
}

func (uc *ObtenerRetasPorZonaUseCase) Execute(zonaID string) ([]entities.RetaInfo, error) {
	return uc.retaRepo.ObtenerRetasPorZona(zonaID)
}

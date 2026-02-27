package application

import (
	"errors"
	"games-football-api/src/retas/domain/entities"
	"games-football-api/src/retas/domain/repositories"
)

type ObtenerHistorialChatUseCase struct {
	retaRepo repositories.IRetaRepository
}

func NewObtenerHistorialChatUseCase(retaRepo repositories.IRetaRepository) *ObtenerHistorialChatUseCase {
	return &ObtenerHistorialChatUseCase{
		retaRepo: retaRepo,
	}
}

// Execute obtiene el historial de mensajes de una reta
func (uc *ObtenerHistorialChatUseCase) Execute(retaID string) ([]entities.Mensaje, error) {
	if retaID == "" {
		return nil, errors.New("reta_id es requerido")
	}

	return uc.retaRepo.ObtenerMensajesDeReta(retaID)
}

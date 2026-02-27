package application

import (
	"errors"
	"games-football-api/src/retas/domain/entities"
	"games-football-api/src/retas/domain/repositories"
)

type EnviarMensajeUseCase struct {
	retaRepo repositories.IRetaRepository
}

func NewEnviarMensajeUseCase(retaRepo repositories.IRetaRepository) *EnviarMensajeUseCase {
	return &EnviarMensajeUseCase{
		retaRepo: retaRepo,
	}
}

// Execute guarda el mensaje en BD y retorna el mensaje enriquecido con el nombre real del usuario
func (uc *EnviarMensajeUseCase) Execute(retaID, usuarioID, texto string) (*entities.Mensaje, error) {
	if retaID == "" || usuarioID == "" || texto == "" {
		return nil, errors.New("reta_id, usuario_id y texto son requeridos")
	}

	mensaje := *entities.NewMensaje(retaID, usuarioID, texto)

	// El repositorio guarda el mensaje y hace JOIN con usuarios para obtener el nombre real
	mensajeEnriquecido, err := uc.retaRepo.GuardarMensaje(mensaje)
	if err != nil {
		return nil, err
	}

	return mensajeEnriquecido, nil
}

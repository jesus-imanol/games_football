package controllers

import (
	"games-football-api/src/usuarios/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterController struct {
	registerUseCase *application.RegisterUseCase
}

func NewRegisterController(registerUseCase *application.RegisterUseCase) *RegisterController {
	return &RegisterController{
		registerUseCase: registerUseCase,
	}
}

// RegisterRequest representa el cuerpo de la petición de registro
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nombre   string `json:"nombre" binding:"required"`
}

// HandleRegister maneja la petición POST de registro
func (rc *RegisterController) HandleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"mensaje": "Campos requeridos: username, password, nombre",
		})
		return
	}

	usuario, err := rc.registerUseCase.Execute(req.Username, req.Password, req.Nombre)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"mensaje": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"mensaje": "Usuario registrado exitosamente",
		"usuario": usuario,
	})
}

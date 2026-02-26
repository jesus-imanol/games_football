package entities

// Usuario representa a un usuario registrado en el sistema
type Usuario struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Nombre   string `json:"nombre"`
}

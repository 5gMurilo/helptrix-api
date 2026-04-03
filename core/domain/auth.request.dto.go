package domain

// AddressInputDTO represents the address fields for user registration.
//
//	@name	AddressInputDTO
type AddressInputDTO struct {
	Street       string `json:"street" binding:"required"`
	Number       string `json:"number" binding:"required"`
	Complement   string `json:"complement,omitempty"`
	Neighborhood string `json:"neighborhood" binding:"required"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required,len=2"`
}

// LoginRequestDTO representa o payload para autenticacao de usuario.
//
//	@name	LoginRequestDTO
type LoginRequestDTO struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequestDTO represents the payload for registering a new user.
//
//	@name	RegisterRequestDTO
type RegisterRequestDTO struct {
	Name       string          `json:"name" binding:"required"`
	Email      string          `json:"email" binding:"required,email"`
	Password   string          `json:"password" binding:"required,min=6"`
	UserType   string          `json:"user_type" binding:"required"`
	Categories []uint          `json:"categories" binding:"required,min=1"`
	Address    AddressInputDTO `json:"address" binding:"required"`
	Phone      string          `json:"phone" binding:"required"`
	Document   string          `json:"document" binding:"required"`
}

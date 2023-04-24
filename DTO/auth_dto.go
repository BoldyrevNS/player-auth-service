package DTO

type AuthRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequestDTO struct {
	Token string `json:"token"`
}

type TokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh-token"`
}

type RegistrationRequestDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UserDTO struct {
	Id        uint   `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Role      string `json:"role"`
}

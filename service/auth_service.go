package service

import (
	"auth-ms/DTO"
	"auth-ms/model"
	"auth-ms/provider"
	"auth-ms/shared/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Auth(data DTO.AuthRequestDTO) (DTO.TokenResponseDTO, error)
	Refresh(data DTO.RefreshRequestDTO) (DTO.TokenResponseDTO, error)
	Registration(data DTO.RegistrationRequestDTO) error
	DeleteUser(id uint) error
	GetAllUsers(currentUserId uint) ([]DTO.UserDTO, error)
}

type authServiceImpl struct {
	provider provider.UserProvider
}

func NewAuthService(userProvider provider.UserProvider) AuthService {
	return &authServiceImpl{provider: userProvider}
}

func (s *authServiceImpl) Auth(data DTO.AuthRequestDTO) (DTO.TokenResponseDTO, error) {
	user, err := s.provider.FindOneByEmail(data.Email)
	if err != nil {
		return DTO.TokenResponseDTO{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return DTO.TokenResponseDTO{}, err
	}

	tokens, err := token.CreateTokenPair(user)
	if err != nil {
		return DTO.TokenResponseDTO{}, err
	}
	return DTO.TokenResponseDTO{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (s *authServiceImpl) Registration(data DTO.RegistrationRequestDTO) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		return err
	}
	user := model.User{
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
		Password:  string(hashedPassword),
		Role:      "user",
	}
	return s.provider.Create(user)
}

func (s *authServiceImpl) Refresh(data DTO.RefreshRequestDTO) (DTO.TokenResponseDTO, error) {
	parsedToken, err := token.ParseRefreshToken(data.Token)
	if err != nil {
		return DTO.TokenResponseDTO{}, err
	}
	user, err := s.provider.FindOneById(parsedToken.Id)
	if err != nil {
		return DTO.TokenResponseDTO{}, err
	}
	newTokenPair, err := token.CreateTokenPair(user)
	if err != nil {
		return DTO.TokenResponseDTO{}, err
	}

	return DTO.TokenResponseDTO{AccessToken: newTokenPair.AccessToken, RefreshToken: newTokenPair.RefreshToken}, nil
}

func (s *authServiceImpl) DeleteUser(id uint) error {
	err := s.provider.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *authServiceImpl) GetAllUsers(currentUserId uint) ([]DTO.UserDTO, error) {
	var users []DTO.UserDTO
	rawUsers, err := s.provider.FindAllExcludeCurrent(currentUserId)
	if err != nil {
		return users, err
	}
	for _, user := range rawUsers {
		users = append(users, DTO.UserDTO{
			Id:        user.Id,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Email,
			Role:      user.Role,
		})
	}
	return users, nil
}

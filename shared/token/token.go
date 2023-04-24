package token

import (
	"auth-ms/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
	"time"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	Id        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Id uint `json:"id"`
}

type Pair struct {
	AccessToken  string
	RefreshToken string
}

func CreateTokenPair(data model.User) (Pair, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	accessTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, &AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
		Id:               data.Id,
		Firstname:        data.Firstname,
		Lastname:         data.Lastname,
		Email:            data.Email,
		Role:             data.Role,
	})
	accessToken, err := accessTokenRaw.SignedString(secret)
	if err != nil {
		return Pair{}, err
	}
	refreshTokenRaw := jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2))},
		Id:               data.Id,
	})
	refreshToken, err := refreshTokenRaw.SignedString(secret)
	if err != nil {
		return Pair{}, err
	}
	return Pair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func GetTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", fmt.Errorf("wrong token")
	}
	if headerParts[0] != "Bearer" {
		return "", fmt.Errorf("wrong token")
	}
	return headerParts[1], nil
}

func ParseAccessToken(token string) (AccessClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	parsedToken, err := jwt.ParseWithClaims(token, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return AccessClaims{}, err
	}
	if claims, ok := parsedToken.Claims.(*AccessClaims); ok && parsedToken.Valid {
		return AccessClaims{
			RegisteredClaims: claims.RegisteredClaims,
			Id:               claims.Id,
			Email:            claims.Email,
			Firstname:        claims.Firstname,
			Lastname:         claims.Lastname,
			Role:             claims.Role,
		}, nil
	}
	return AccessClaims{}, err
}

func ParseRefreshToken(token string) (RefreshClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	parsedToken, err := jwt.ParseWithClaims(token, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return RefreshClaims{}, err
	}
	if claims, ok := parsedToken.Claims.(*RefreshClaims); ok && parsedToken.Valid {
		return RefreshClaims{
			RegisteredClaims: claims.RegisteredClaims,
			Id:               claims.Id,
		}, nil
	}

	return RefreshClaims{}, err
}

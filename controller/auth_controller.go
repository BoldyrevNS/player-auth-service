package controller

import (
	"auth-ms/DTO"
	"auth-ms/service"
	"auth-ms/shared/response"
	"auth-ms/shared/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AuthController interface {
	Auth(ctx *gin.Context)
	Registration(ctx *gin.Context)
	Refresh(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
}

type authControllerImpl struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authControllerImpl{
		authService: authService,
	}
}

// Auth				godoc
// @Tags			Auth
// @Summary			Auth user
// @Description 	Check user credentials and auth in service
// @Param			user body DTO.AuthRequestDTO true "User authorization"
// @Success			200 {object} service.DataJSON{data=DTO.TokenResponseDTO{}}
// @Failure      	404  {object}  service.MessageJSON{}
// @Failure      	400
// @Router			/auth [post]
func (c *authControllerImpl) Auth(ctx *gin.Context) {
	authRequest := DTO.AuthRequestDTO{}
	err := ctx.ShouldBindJSON(&authRequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	tokens, err := c.authService.Auth(authRequest)
	if err != nil {
		if err.Error() == "record not found" {
			ctx.AbortWithStatusJSON(http.StatusNotFound, response.MessageJSON{Message: "User not found"})
			return
		}
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	response.SendJSON(ctx, http.StatusOK, response.DataJSON{
		Data: tokens,
	})
}

// Registration		godoc
// @Tags			Auth
// @Summary			User registration
// @Description 	Create new user in database, gave default permissions
// @Param			user body DTO.RegistrationRequestDTO true "User registration"
// @Success			201
// @Failure      	409 {object}  service.MessageJSON{}
// @Failure      	400
// @Router			/auth/registration [post]
func (c *authControllerImpl) Registration(ctx *gin.Context) {
	registrationRequest := DTO.RegistrationRequestDTO{}
	err := ctx.ShouldBindJSON(&registrationRequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.authService.Registration(registrationRequest)
	if err != nil {
		if err.Error() == "duplicated key not allowed" {
			ctx.AbortWithStatusJSON(http.StatusConflict, response.MessageJSON{
				Message: "Email already exists",
			})
		}
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.AbortWithStatus(http.StatusCreated)
}

// Refresh			godoc
// @Tags			Auth
// @Summary			Refresh tokens
// @Description 	Gave new token pair
// @Param			user body DTO.RefreshRequestDTO true "Refresh tokens"
// @Success			200 {object} service.DataJSON{data=DTO.TokenResponseDTO{}}
// @Failure      	400
// @Failure      	401
// @Router			/auth/refresh [post]
func (c *authControllerImpl) Refresh(ctx *gin.Context) {
	refreshRequest := DTO.RefreshRequestDTO{}
	err := ctx.ShouldBindJSON(&refreshRequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	tokens, err := c.authService.Refresh(refreshRequest)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	response.SendJSON(ctx, http.StatusOK, response.DataJSON{
		Data: tokens,
	})
}

// DeleteUser 		godoc
// @Tags			Auth
// @Summary			Delete user
// @Security 		BearerAuth
// @Description		Remove user data by id.
// @Param			userId   path   uint  true  "User ID"
// @Success			200
// @Failure      	400 {object} service.MessageJSON{}
// @Failure      	500 {object} service.MessageJSON{}
// @Router			/auth/{userId} [delete]
func (c *authControllerImpl) DeleteUser(ctx *gin.Context) {
	userIdParam, find := ctx.Params.Get("userId")
	if !find {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.MessageJSON{Message: "Provide id param"})
		return
	}
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.MessageJSON{Message: "Wrong format"})
	}
	err = c.authService.DeleteUser(uint(userId))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.MessageJSON{Message: "Delete error"})
	}
	ctx.AbortWithStatus(http.StatusOK)
}

// GetAllUsers 		godoc
// @Tags			Auth
// @Summary			Get all users
// @Security 		BearerAuth
// @Success			200 {object} service.DataJSON{data=[]DTO.UserDTO}
// @Failure      	401
// @Failure      	500
// @Router			/auth/allUsers [get]
func (c *authControllerImpl) GetAllUsers(ctx *gin.Context) {
	headerToken, err := token.GetTokenFromHeader(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	accessClaims, err := token.ParseAccessToken(headerToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	users, err := c.authService.GetAllUsers(accessClaims.Id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	response.SendJSON(ctx, http.StatusOK, response.DataJSON{
		Data: users,
	})
}

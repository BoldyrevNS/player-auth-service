package router

import (
	"auth-ms/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminPermissionMiddleware(c *gin.Context) {
	token, err := service.GetTokenFromHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, service.ResponseMessageJSON{
			Message: err.Error(),
		})
		return
	}
	accessTokenClaims, err := service.ParseAccessToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, service.ResponseMessageJSON{
			Message: err.Error(),
		})
		return
	}
	if accessTokenClaims.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, service.ResponseMessageJSON{
			Message: "permission denied",
		})
		return
	}
}

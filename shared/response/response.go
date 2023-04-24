package response

import "github.com/gin-gonic/gin"

type MessageJSON struct {
	Message string `json:"message" example:"msg"`
}

type DataJSON struct {
	Data interface{} `json:"data"`
}

func SendJSON(ctx *gin.Context, status int, res interface{}) {
	ctx.Header("Content-type", "application/json")
	ctx.JSON(status, res)
}

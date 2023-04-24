package router

import (
	"auth-ms/controller"
	"auth-ms/db"
	"auth-ms/provider"
	"auth-ms/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func authRoutes(router *gin.RouterGroup) {
	public := router.Group("/auth")
	adminPermission := router.Group("/auth")
	adminPermission.Use(AdminPermissionMiddleware)

	userProvider := provider.NewUserProvider(db.Instance)
	authService := service.NewAuthService(userProvider)
	authController := controller.NewAuthController(authService)

	public.POST("/", authController.Auth)
	public.POST("/registration", authController.Registration)
	public.POST("/refresh", authController.Refresh)

	adminPermission.DELETE(":userId", authController.DeleteUser)
	adminPermission.GET("/allUsers", authController.GetAllUsers)
}

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	baseRouter := router.Group("/api/v1")
	authRoutes(baseRouter)

	return router
}

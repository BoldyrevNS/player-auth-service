package main

import (
	"auth-ms/db"
	docs "auth-ms/docs"
	"auth-ms/model"
	"auth-ms/router"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func migrate() error {
	err := db.Instance.AutoMigrate(&model.User{})
	return err
}

// @title		Player auth-service
// @version		1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host 		localhost:8080
// @BasePath	/api/v1/
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Enviroment error")
	}
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", os.Getenv("DISPLAY_HOST"), os.Getenv("DISPLAY_PORT"))
	db.NewDatabase().Init()
	err = migrate()
	if err != nil {
		log.Fatalf("Migration error")
	}
	routes := router.NewRouter()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("DISPLAY_PORT")),
		Handler: routes,
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Serve error")
	}
}

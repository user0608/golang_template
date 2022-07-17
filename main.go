package main

import (
	"goqrs/database"
	"goqrs/envs"
	"goqrs/handlers"
	"goqrs/security"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Eror loading .env file:", err)
		os.Exit(0)
	}
}
func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if err := security.LoadRSAKeys(); err != nil {
		log.Println("Err RSA Kes:", err)
		os.Exit(0)
	}
	log.Println("RSA KEYS OK!")
	if err := database.PrepareConnection(); err != nil {
		log.Println("Err database:", err)
		os.Exit(0)
	}
	log.Println("DATABASE OK!")
	e := echo.New()
	e.HideBanner = true
	e.Use(database.GormMiddleware)
	handlers.StartRoutes(e)
	go func() {
		for sig := range c {
			if err := e.Close(); err != nil {
				log.Println("error trying to stop echo service:", err)
				os.Exit(0)
			}
			log.Println("stopping echo service:", sig)
			os.Exit(1)
		}
	}()
	err := e.Start(envs.FindEnv("GOQRS_ADDRESS", "localhost:8080"))
	if err != nil {
		log.Println("Error staring echo service:", err)
		os.Exit(0)
	}
}

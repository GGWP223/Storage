package main

import (
	proto "File_Storage/internal/api/jwt"
	"File_Storage/internal/app"
	"File_Storage/internal/domain"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtConnection, err := GetConnection("localhost:8083")

	if err != nil {
		log.Fatal(err)
	}

	clients := domain.Clients{
		Jwt: proto.NewJwtServiceClient(jwtConnection),
	}

	appCors := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8086"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	baseApp := domain.NewBaseApp(&clients, appCors)

	apps := []domain.App{
		app.NewAuthApp(baseApp, "localhost:8081", "localhost:8082"),
		app.NewJwtApp(baseApp, "localhost:8083"),
		app.NewStorageApp(baseApp, "localhost:8084", "localhost:8085"),
		app.NewWebApp("localhost:8086"),
	}

	errCh := make(chan error, len(apps))

	for _, a := range apps {
		go func(a domain.App) {
			if err := a.Start(); err != nil {
				errCh <- err
			}
		}(a)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case sig := <-quit:
		log.Printf("Получен сигнал %v, начинаем остановку...", sig)
	case err := <-errCh:
		log.Printf("Ошибка при старте сервиса: %v", err)
	}

	GracefulShutdown(apps, []*grpc.ClientConn{
		jwtConnection,
	})
}

func GetConnection(address string) (*grpc.ClientConn, error) {
	jwtClient, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Failed to dial: " + err.Error())
	}

	return jwtClient, nil
}

func GracefulShutdown(apps []domain.App, connections []*grpc.ClientConn) {
	for _, a := range apps {
		a.GracefulShutdown()
	}

	for _, conn := range connections {
		if err := conn.Close(); err != nil {
			log.Fatal("Failed to close connection: " + err.Error())
		}
	}

	log.Printf("Успешно!")
}

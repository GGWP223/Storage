package app

import (
	proto "File_Storage/internal/api/jwt"
	"File_Storage/internal/domain"
	"File_Storage/internal/handlers"
	services "File_Storage/internal/services/jwt"
	"errors"
	"net"
	"os"

	"google.golang.org/grpc"
)

type JwtApp struct {
	app     *domain.BaseApp
	srv     *grpc.Server
	address string
}

func NewJwtApp(app *domain.BaseApp, address string) *JwtApp {
	return &JwtApp{app: app, address: address}
}

func (app *JwtApp) Start() error {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return errors.New("JWT_SECRET env variable not set")
	}

	service := services.NewService(secret)
	handler := handlers.NewJwtServer(service)

	app.srv = grpc.NewServer()

	proto.RegisterJwtServiceServer(app.srv, handler)

	lis, err := net.Listen("tcp", app.address)

	if err != nil {
		return errors.New("failed to listen" + err.Error())
	}

	if err := app.srv.Serve(lis); err != nil {
		return errors.New("Failed to serve gRPC server: " + err.Error())
	}

	return nil
}

func (app *JwtApp) GracefulShutdown() {
	if app.srv != nil {
		app.srv.GracefulStop()
	}
}

package app

import (
	proto "File_Storage/internal/api/auth"
	"File_Storage/internal/db"
	"File_Storage/internal/domain"
	handlers "File_Storage/internal/handlers"
	services "File_Storage/internal/services/auth"
	"context"
	"errors"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type AuthApp struct {
	app         *domain.BaseApp
	srv         *grpc.Server
	grpcAddress string
	restAddress string
}

func NewAuthApp(app *domain.BaseApp, grpcAddress, restAddress string) *AuthApp {
	return &AuthApp{
		app:         app,
		grpcAddress: grpcAddress,
		restAddress: restAddress,
	}
}

func (app *AuthApp) Start() error {
	dsn := os.Getenv("AUTH_DSN")

	database, err := db.InitDB(dsn, &services.User{})
	if err != nil {
		return errors.New("failed to connect to database")
	}

	repo := services.NewRepository(database)
	service := services.NewService(repo, app.app.Clients.Jwt)
	handler := handlers.NewAuthHandler(service)

	app.srv = grpc.NewServer()
	proto.RegisterAuthServiceServer(app.srv, handler)

	lis, err := net.Listen("tcp", app.grpcAddress)

	if err != nil {
		return errors.New("failed to listen: " + err.Error())
	}

	go func() {
		if err := app.srv.Serve(lis); err != nil {
			panic("Failed to serve gRPC server: " + err.Error())
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := proto.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, app.grpcAddress, opts); err != nil {
		return errors.New("Failed to register auth rest handler: " + err.Error())
	}

	srv := http.Server{
		Addr:    app.restAddress,
		Handler: app.app.Cors.Handler(mux),
	}

	if err := srv.ListenAndServe(); err != nil {
		return errors.New("Failed to serve gRPC server: " + err.Error())
	}

	return nil
}

func (app *AuthApp) GracefulShutdown() {
	if app.srv != nil {
		app.srv.GracefulStop()
	}
}

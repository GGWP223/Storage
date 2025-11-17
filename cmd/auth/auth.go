package main

import (
	"File_Storage/internal/db"
	handlers "File_Storage/internal/handlers/auth"
	services "File_Storage/internal/services/auth"
	"log"
	"net"
	"os"

	proto "File_Storage/internal/api/auth"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	JWTSecret := os.Getenv("JWT_SECRET")
	dsn := os.Getenv("AUTH_DSN")

	database, err := db.InitDB(dsn, &services.User{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	repo := services.NewRepository(database)
	service := services.NewService(repo, JWTSecret)

	grpcHandler := handlers.NewGRPCAuthHandler(service)
	restHandler := handlers.NewRestHandler(service)

	grpcServer := grpc.NewServer()
	echoServer := echo.New()

	echoServer.Use(middleware.Logger())

	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:8084"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	echoServer.POST("/login", restHandler.Login)
	echoServer.POST("/register", restHandler.Register)

	proto.RegisterAuthServiceServer(grpcServer, grpcHandler)

	lis, err := net.Listen("tcp", "localhost:8083")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		if err := echoServer.Start("localhost:8082"); err != nil {
			log.Fatalf("Failed to start REST server: %v", err)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

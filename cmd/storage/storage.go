package main

import (
	authProto "File_Storage/internal/api/auth"
	"File_Storage/internal/db"
	handlers "File_Storage/internal/handlers/storage"
	"File_Storage/internal/services/storage"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("STORAGE_DSN")
	database, err := db.InitDB(dsn, &storage.FileMetadata{})

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	conn, err := grpc.Dial("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect to auth server: %v", err)
	}

	defer conn.Close()

	client := authProto.NewAuthServiceClient(conn)

	repo := storage.NewRepository(database)
	service := storage.NewService(repo, client)

	restHandler := handlers.NewRestStorageServer(service)

	echoServer := echo.New()
	echoServer.Use(middleware.Logger())

	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"http://localhost:8084"},
		AllowMethods:  []string{echo.GET, echo.POST, echo.OPTIONS, echo.DELETE},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Disposition"},
	}))

	echoServer.GET("/GetFiles", restHandler.GetFiles)
	echoServer.POST("/DownloadFile", restHandler.DownloadFile)
	echoServer.POST("/UploadFile", restHandler.UploadFile)
	echoServer.DELETE("/DeleteFile", restHandler.DeleteFile)

	if err != nil {
		log.Fatalf("Failed to listen on port 8080: %v", err)
	}

	if err := echoServer.Start("localhost:8081"); err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}
}

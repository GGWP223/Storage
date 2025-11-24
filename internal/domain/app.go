package domain

import (
	"File_Storage/internal/api/jwt"

	"github.com/rs/cors"
)

type App interface {
	Start() error
	GracefulShutdown()
}

type Clients struct {
	Jwt jwt.JwtServiceClient
}

type BaseApp struct {
	Clients *Clients
	Cors    *cors.Cors
}

func NewBaseApp(clients *Clients, cors *cors.Cors) *BaseApp {
	return &BaseApp{Clients: clients, Cors: cors}
}

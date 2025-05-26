package server

import (
	"context"
	"log"
	"shofy/app/api/config"
	db "shofy/db/sqlc"
	"shofy/modules/azure/service"
	azureService "shofy/modules/azure/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DBPool  *pgxpool.Pool
	Queries *db.Queries
	Ctx     context.Context
	OpenAI  *service.AzureOpenAI
}

func NewServer(ctx context.Context) *Server {
	dbPool, err := config.LoadDbConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Init OpenAI
	openAI := azureService.NewOpenAI(ctx)

	// Set JWT key
	// secretBaseKey := os.Getenv("SECRET_BASE_KEY")
	// if secretBaseKey == "" {
	// 	log.Fatalf("SECRET_BASE_KEY is not set")
	// }
	// validation.SetJWTKey(secretBaseKey)

	return &Server{
		DBPool:  dbPool,
		Queries: db.New(dbPool),
		Ctx:     ctx,
		OpenAI:  openAI,
	}
}

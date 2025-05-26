package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shofy/app/api/router"
	"shofy/app/api/server"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

var (
	//      *server.Server
	apiRouter *gin.Engine
	ctx       context.Context
	srv       *server.Server
)

func initServer(ctx context.Context) {
	gotenv.Load()

	srv = server.NewServer(ctx)
	apiRouter = router.InitRouter(ctx, srv)
}

func main() {
	ctx = context.Background()

	initServer(ctx)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: apiRouter,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // Wait for interrupt signal

	// Create a context with a timeout for the shutdown
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		panic(err)
	}
}

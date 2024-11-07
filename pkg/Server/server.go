package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
	"ws-chat/configs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	Server struct {
		App *echo.Echo
		Db  *mongo.Client
		Cfg *configs.Config
	}
)

var (
	serverInstance Server
	once           sync.Once
)

func (s *Server) httpListening(ctx context.Context) {
	go func() {
		if err := s.App.Start(s.Cfg.Server.Url); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: %v", err)
		}
	}()

	// Wait for context cancellation (e.g., from OS signal)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.App.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

}

func Start(ctx context.Context, cfg *configs.Config, db *mongo.Client) {
	s := Server{
		App: echo.New(),
		Db:  db,
		Cfg: cfg,
	}

	s.App.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	s.roomService()

	s.httpListening(ctx)

}

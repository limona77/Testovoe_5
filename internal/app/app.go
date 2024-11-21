package app

import (
	"Testovoe_5/internal/config"
	"Testovoe_5/internal/controller"
	"Testovoe_5/internal/pkg/postgres"
	"Testovoe_5/internal/pkg/slogger"
	"Testovoe_5/internal/repository"
	"Testovoe_5/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ReadTimeout  = 3 * time.Second
	WriteTimeout = 3 * time.Second
)

func Run() {
	slogger.SetLogger()
	cfg := config.NewConfig()
	slog.Info("config ok", cfg)
	router := gin.Default()

	slog.Info("connecting to postgres")
	db := postgres.New(cfg.PG.URL)
	defer db.Close()
	slog.Info("connect to postgres ok")

	slog.Info("init repositories")
	repositories := repository.NewRepositories(db)

	slog.Info("init services")
	deps := service.ServicesDeps{
		Repository: repositories,
	}

	services := service.NewServices(deps)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router.Handler(),
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}
	controller.NewRouter(router, services)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		slog.Info("timeout of 5 seconds.")
	}
	slog.Info("Server exiting")
}

package main

import (
	"authMicro/internal/api/rest"
	"authMicro/internal/api/rest/handler"
	"authMicro/internal/config"
	"authMicro/utlis/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	l := logger.NewLogger()
	cfg, errors := config.NewConfig()
	if len(errors) > 0 {
		for _, err := range errors {
			l.Error(err.Error())
		}
		return
	}

	router := rest.NewRouter(
		handler.NewRegisterHandler(l),
		handler.NewLoginHandler(l),
		handler.NewRefreshTokenHandler(l),
	)

	go func() {
		l.Infof("Starting REST server on %s", cfg.HTTPConfig.ServerAddr)
		if err := router.Start(cfg.HTTPConfig.ServerAddr); err != nil {
			l.Fatalf("Failed to start REST server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	l.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		l.Errorf("Error during server shutdown: %v", err)
	}

	l.Info("Server exited properly")
}

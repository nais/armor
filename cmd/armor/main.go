package main

import (
	"context"
	"fmt"
	"github.com/nais/armor/pkg/google"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nais/armor/config"
	"github.com/nais/armor/pkg/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}

	cfg, err := SetupConfig()
	if err != nil {
		log.WithError(err).Fatal("new config")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	var opts []option.ClientOption
	googleClient := google.NewClient(cfg, ctx, log.WithField("component", "armor-client"), opts...)
	app := handler.NewApp(ctx, googleClient, log.WithField("system", "armor"))

	h := handler.NewHandler(app)
	router := app.SetupHttpRouter(h)

	server := http.Server{
		Addr:              cfg.Port,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       10 * time.Minute,
	}

	log.WithField("addr", fmt.Sprintf("%s", cfg.Port)).Info("starting server")
	ctx, cancel := context.WithCancel(ctx)
	go LogError(log, cancel, func() error { return http.ListenAndServe(cfg.Port, router) })

	<-ctx.Done()

	stop()
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Error(err)
	}
}

func LogError(log *logrus.Logger, cancel context.CancelFunc, fn func() error) {
	if err := fn(); err != nil {
		cancel()
		log.WithError(err).Error("error")
	}
}

func SetupConfig() (*config.Config, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	if err = cfg.Validate([]string{
		config.DevelopmentMode,
		config.Port,
	}); err != nil {
		return nil, err
	}
	return cfg, nil
}

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

	cfg, err := config.SetupConfig()
	if err != nil {
		log.WithError(err).Fatal("setting up new config")
	}

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("setting log level")
	}

	log.Level = logLevel

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	var opts []option.ClientOption
	gSecurityClient := google.NewSecurityClient(cfg, ctx, log.WithField("component", "armor-security-client"), opts...)
	gServiceClient := google.NewServiceClient(ctx, log.WithField("component", "armor-serice-client"))

	h := handler.NewHandler(ctx, cfg, gSecurityClient, gServiceClient, log.WithField("system", "armor"))
	router := handler.SetupHttpRouter(h)

	server := http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       10 * time.Minute,
	}

	log.WithField("addr", fmt.Sprintf("%s", ":8080")).Info("starting server")
	ctx, cancel := context.WithCancel(ctx)
	go LogError(log, cancel, func() error { return http.ListenAndServe(":8080", router) })

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

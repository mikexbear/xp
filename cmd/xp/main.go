package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"xp/internal/config"
	"xp/internal/logging"
	"xp/internal/probe"
	"xp/internal/subscription"
)

func main() {
	logging.InitLogging()
	slog.Info("Logging configured")

	appConfigPath := parseArgs()
	appConfig, err := config.Load(appConfigPath)
	if err != nil {
		slog.Error("Failed to load application configuration", "err", err)
		panic(err)
	}
	slog.Info(fmt.Sprintf("Application configuration created (path: %s)", appConfigPath))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	subs, err := subscription.NewSubscriptions(appConfig.SubscriptionUrls)
	if err != nil {
		slog.Error("Failed to create subscriptions", "err", err)
		panic(err)
	}
	slog.Info("Subscriptions created")
	go subscription.InitJob(subs, appConfig.SubscriptionFetchingIntervalSec, ctx)

	prober := probe.NewProber(subs, appConfig.ProbeUrls, appConfig.ProbeTimeoutSec)
	slog.Info("Prober created")
	go probe.InitJob(prober, appConfig.ProbeIntervalSec, ctx)

	// TODO: start HTTP server
}

func parseArgs() string {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "Path to YAML configuration")

	flag.Parse()

	return configPath
}

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notifier"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	logger := log.New(os.Stdout, "portwatch: ", log.LstdFlags)

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		logger.Fatalf("failed to load config: %v", err)
	}

	var notifiers []monitor.Notifier

	if cfg.WebhookURL != "" {
		wh := notifier.NewWebhookNotifier(cfg.WebhookURL)
		notifiers = append(notifiers, wh)
		logger.Printf("webhook notifier enabled: %s", cfg.WebhookURL)
	}

	if cfg.DesktopNotify {
		dn := notifier.NewDesktopNotifier(cfg.AppName)
		notifiers = append(notifiers, dn)
		logger.Println("desktop notifier enabled")
	}

	if len(notifiers) == 0 {
		logger.Println("warning: no notifiers configured, changes will only be logged")
	}

	mon := monitor.New(cfg.Interval, logger, notifiers...)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Printf("starting monitor (interval: %s)", cfg.Interval)
	if err := mon.Run(ctx); err != nil && err != context.Canceled {
		logger.Fatalf("monitor error: %v", err)
	}
	logger.Println("shutdown complete")
}

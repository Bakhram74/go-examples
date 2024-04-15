package main

import (
	"log"
	"single-window/config"
	"single-window/internal/app"
	"single-window/pkg/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	l := logger.New(cfg.Level)

	err = app.RunMigration(cfg, l)
	if err != nil {
		log.Fatalf("Migration error: %s", err)
	}

	app.Run(cfg, l)
}

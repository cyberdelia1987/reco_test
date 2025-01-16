package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/cyber/test-project/app"
	"github.com/cyber/test-project/shutdown"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	application, err := app.InitApplication(*configPath)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	go func() {
		err = application.Start()
		if err != nil {
			log.Fatalf("Failed to start application: %v", err)
		}
	}()

	shutdown.ListenForSignals([]os.Signal{os.Interrupt, syscall.SIGTERM}, application)
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	appName    = "Go File Notifier"
	configFile = "config.json"
	stateFile  = "state.json"
)

func main() {
	// Load configuration
	config, err := LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Load state
	stateManager := NewStateManager(stateFile)
	if err := stateManager.Load(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	// Create notifier and monitor
	notifier := NewNotifier(appName)
	monitor := NewFileMonitor(config, stateManager, notifier)

	// Start
	if err := monitor.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run in goroutine
	done := make(chan error, 1)
	go func() {
		done <- monitor.Run()
	}()

	// Wait for signal or error
	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	case <-sigChan:
		fmt.Println("\nShutting down...")
	}

	monitor.Stop()
}

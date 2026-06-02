package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"MONITORING-TOOL/internal/collector"
	"MONITORING-TOOL/internal/models"
	"MONITORING-TOOL/internal/tracer"
	"MONITORING-TOOL/internal/client"
)

func main() {

	defer recoverPanic()

	err := tracer.InitializeBufferedWriter()
	if err != nil {
		panic(err)
	}

	

	stopFlusher := make(chan struct{})

	go tracer.StartCSVFlusher(stopFlusher)

	log.Println("====================================")
	log.Println("Monitoring Agent Started")
	log.Println("====================================")

	currentOS := runtime.GOOS

	log.Println("Running OS:", currentOS)

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("hostname error:", err)
		hostname = "unknown-host"
	}

	log.Println("Hostname:", hostname)

	config, err := models.LoadConfig("config.json")
	if err != nil {
		log.Println("config.json not found or error loading, using defaults:", err)
		config = &models.Config{
			ServerURL:                 "http://localhost:8081",
			CollectionIntervalSeconds: 60,
			UserID:                    0,
			EnableCPU:                 true,
			EnableMemory:              true,
			EnableDisk:                true,
			EnableNetwork:             true,
		}
	}

	// Phase 4: Device Registration
	if config.UserID == 0 {
		log.Println("UserID is 0. Attempting to register device with backend...")
		userID, err := client.RegisterDevice(config.ServerURL)
		if err != nil {
			log.Fatalf("Failed to register device: %v", err)
		}
		log.Printf("Device registered successfully! UserID: %d\n", userID)
		
		config.UserID = userID
		err = models.SaveConfig("config.json", config)
		if err != nil {
			log.Printf("Failed to save config.json: %v\n", err)
		}
	} else {
		log.Printf("Device already registered. UserID: %d\n", config.UserID)
	}

	interval := time.Duration(config.CollectionIntervalSeconds) * time.Second

	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	signals := make(chan os.Signal, 1)

	signal.Notify(
		signals,
		os.Interrupt,
		syscall.SIGTERM,
	)

	for {

		select {

		case <-ticker.C:

			runCollectionCycle(config)

		case sig := <-signals:

			log.Println("shutdown signal received:", sig)
			log.Println("Stopping Monitoring Agent")

			close(stopFlusher)

			time.Sleep(2 * time.Second)

			return
		}
	}
}

func runCollectionCycle(config *models.Config) {

	defer recoverPanic()

	log.Println("Collecting metrics...")

	metrics, err := collector.CollectMetrics()
	if err != nil {
		log.Println("collector error:", err)
		return
	}

	if metrics == nil {
		log.Println("metrics are nil")
		return
	}

	err = tracer.BufferedWriteMetrics(metrics)
	if err != nil {
		log.Println("csv write error:", err)
		// continue even if CSV fails
	} else {
		log.Println("metrics buffered successfully (CSV)")
	}

	// Send to Backend API
	err = client.SendMetrics(metrics, config)
	if err != nil {
		log.Println("metrics upload failed:", err)
	}

	// Send Heartbeat
	err = client.SendHeartbeat(config.UserID, config.ServerURL)
	if err != nil {
		log.Println("heartbeat failed:", err)
	}
}

func recoverPanic() {

	if r := recover(); r != nil {

		log.Println("panic recovered:", r)
	}
}
package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"MONITORING-TOOL/internal/collector"
	"MONITORING-TOOL/internal/tracer"
	"MONITORING-TOOL/internal/models"
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

	config, err := models.LoadConfig("config/config.json")
if err != nil {
	panic(err)
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

			runCollectionCycle()

		case sig := <-signals:

			log.Println("shutdown signal received:", sig)
			log.Println("Stopping Monitoring Agent")

			close(stopFlusher)

			time.Sleep(2 * time.Second)

			return
		}
	}
}

func runCollectionCycle() {

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
		return
	}

	log.Println("metrics buffered successfully")
}

func recoverPanic() {

	if r := recover(); r != nil {

		log.Println("panic recovered:", r)
	}
}
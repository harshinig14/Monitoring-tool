package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"MONITORING-TOOL/internal/tracer"
)

func main() {

	traceDirectory := "traces"

	files, err := os.ReadDir(traceDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		fileName := file.Name()

		if !strings.HasSuffix(fileName, ".csv") {
			continue
		}

		sourcePath := filepath.Join(traceDirectory, fileName)

		targetPath := sourcePath + ".gz"

		log.Println("Compressing:", sourcePath)

		err := tracer.CompressFile(sourcePath, targetPath)
		if err != nil {
			log.Println("compression failed:", err)
			continue
		}

		log.Println("compressed successfully:", targetPath)
	}

	log.Println("compression completed")
}
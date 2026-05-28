package tracer

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sync"
	"time"

	"MONITORING-TOOL/internal/models"
)

var (
	csvFile   *os.File
	csvWriter *csv.Writer
	mutex     sync.Mutex
)

func InitializeBufferedWriter() error {

	err := ensureTraceDirectory()
	if err != nil {
		return err
	}

	filePath := getTraceFilePath()

	isNewFile := !fileExists(filePath)

	file, err := os.OpenFile(
		filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return err
	}

	csvFile = file

	csvWriter = csv.NewWriter(file)

	if isNewFile {

		err = csvWriter.Write(models.CSVHeaders)
		if err != nil {
			return err
		}

		csvWriter.Flush()
	}

	return nil
}

func BufferedWriteMetrics(metrics *models.Metrics) error {

	mutex.Lock()
	defer mutex.Unlock()

	err := csvWriter.Write(metrics.ToCSVRow())
	if err != nil {
		return err
	}

	return nil
}

func StartCSVFlusher(stopChan chan struct{}) {

	ticker := time.NewTicker(10 * time.Second)

	defer ticker.Stop()

	for {

		select {

		case <-ticker.C:

			FlushCSV()

		case <-stopChan:

			FlushCSV()

			CloseCSV()

			return
		}
	}
}

func FlushCSV() {

	mutex.Lock()
	defer mutex.Unlock()

	if csvWriter != nil {
		csvWriter.Flush()
	}
}

func CloseCSV() {

	mutex.Lock()
	defer mutex.Unlock()

	if csvFile != nil {
		csvFile.Close()
	}
}

func getTraceFilePath() string {

	fileName := "trace_" + time.Now().Format("2006_01_02") + ".csv"

	return filepath.Join(TraceDirectory, fileName)
}
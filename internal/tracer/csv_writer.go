package tracer

import (
	"os"
)

const TraceDirectory = "traces"

func ensureTraceDirectory() error {

	return os.MkdirAll(TraceDirectory, os.ModePerm)
}

func fileExists(path string) bool {

	_, err := os.Stat(path)

	return err == nil
}
package tracer

import (
	"os"
)

const TraceDirectory = "traces"

func createTraceFile() {

	err := os.MkdirAll("traces", 0755)
	if err != nil {
		panic(err)
	}

}
	func ensureTraceDirectory() error {

		return os.MkdirAll(TraceDirectory, os.ModePerm)
	}

func fileExists(path string) bool {

	_, err := os.Stat(path)

	return err == nil
}
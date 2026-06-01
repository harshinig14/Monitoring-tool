package main

import (
	"os"
	"os/exec"
	"fmt"
)

func main() {
	fmt.Println("Starting cross-platform builds...")
	build("windows", "amd64", "MONITORING-TOOL.exe")
	build("linux", "amd64", "MONITORING-TOOL")
	build("darwin", "amd64", "MONITORING-TOOL-MAC")
	build("darwin", "arm64", "MONITORING-TOOL-MAC-ARM")
	fmt.Println("All builds completed successfully!")
}

func build(goos string, goarch string, output string) {
	fmt.Printf("Building %s (%s) -> %s\n", goos, goarch, output)
	cmd := exec.Command(
		"go",
		"build",
		"-ldflags=-s -w",
		"-o",
		output,
		"./cmd/agent",
	)

	cmd.Env = append(
		os.Environ(),
		"GOOS="+goos,
		"GOARCH="+goarch,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to build %s: %v\n", output, err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

func installWindows() {
	fmt.Println("Installing Windows Agent...")
	targetDir := `C:\Program Files\MonitoringTool`

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Println("Failed to create target directory:", err)
		return
	}

	agentData, err := fs.ReadFile(embeddedFiles, "binaries/windows/agent.exe")
	if err != nil {
		fmt.Println("Failed to read embedded Windows agent:", err)
		return
	}

	err = os.WriteFile(targetDir+"\\agent.exe", agentData, 0755)
	if err != nil {
		fmt.Println("Failed to write agent.exe:", err)
		return
	}

	os.MkdirAll(targetDir+"\\logs", 0755)
	
	defaultConfig := `{"CollectionIntervalSeconds": 60}`
	os.WriteFile(targetDir+"\\config.json", []byte(defaultConfig), 0644)

	exec.Command(
		"sc",
		"create",
		"MonitoringTool",
		"binPath=",
		targetDir+"\\agent.exe",
	).Run()

	exec.Command(
		"sc",
		"start",
		"MonitoringTool",
	).Run()
	
	fmt.Println("Windows Agent installed and started successfully.")
}

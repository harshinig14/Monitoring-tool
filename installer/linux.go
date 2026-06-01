package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

func installLinux() {
	fmt.Println("Installing Linux Agent...")
	targetDir := "/opt/monitoring-tool"

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Println("Failed to create target directory:", err)
		return
	}

	agentData, err := fs.ReadFile(embeddedFiles, "binaries/linux/agent")
	if err != nil {
		fmt.Println("Failed to read embedded Linux agent:", err)
		return
	}

	err = os.WriteFile(targetDir+"/agent", agentData, 0755)
	if err != nil {
		fmt.Println("Failed to write agent:", err)
		return
	}

	os.Chmod(targetDir+"/agent", 0755)

	service := `[Unit]
Description=Monitoring Tool

[Service]
ExecStart=/opt/monitoring-tool/agent
Restart=always

[Install]
WantedBy=multi-user.target
`

	os.WriteFile(
		"/etc/systemd/system/monitoring-tool.service",
		[]byte(service),
		0644,
	)

	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "enable", "monitoring-tool").Run()
	exec.Command("systemctl", "start", "monitoring-tool").Run()
	
	fmt.Println("Linux Agent installed and started successfully.")
}

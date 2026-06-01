package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
)

func installMac() {
	fmt.Println("Installing macOS Agent...")
	targetDir := "/Library/Application Support/MonitoringTool"

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Println("Failed to create target directory:", err)
		return
	}

	var source string
	if runtime.GOARCH == "arm64" {
		source = "binaries/mac/agent-arm"
	} else {
		source = "binaries/mac/agent-intel"
	}

	agentData, err := fs.ReadFile(embeddedFiles, source)
	if err != nil {
		fmt.Println("Failed to read embedded macOS agent:", err)
		return
	}

	err = os.WriteFile(targetDir+"/agent", agentData, 0755)
	if err != nil {
		fmt.Println("Failed to write agent:", err)
		return
	}

	os.Chmod(targetDir+"/agent", 0755)

	plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.monitoring.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Library/Application Support/MonitoringTool/agent</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
`
	plistPath := "/Library/LaunchDaemons/com.monitoring.agent.plist"
	os.WriteFile(plistPath, []byte(plist), 0644)

	exec.Command("launchctl", "load", plistPath).Run()
	exec.Command("launchctl", "start", "com.monitoring.agent").Run()
	
	fmt.Println("macOS Agent installed and started successfully.")
}

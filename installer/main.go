package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Monitoring Agent Installer")
	fmt.Println("--------------------------")

	switch runtime.GOOS {
	case "windows":
		installWindows()
	case "linux":
		installLinux()
	case "darwin":
		installMac()
	default:
		fmt.Println("Unsupported OS")
	}
}

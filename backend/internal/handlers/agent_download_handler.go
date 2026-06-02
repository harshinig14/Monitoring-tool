package handlers

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DownloadAgentWindows(c *gin.Context) {
	serveAgentFile(c, "setup.exe", "agent-windows.exe")
}

func DownloadAgentLinux(c *gin.Context) {
	serveAgentFile(c, "setup-linux", "agent-linux")
}

func DownloadAgentMac(c *gin.Context) {
	serveAgentFile(c, "setup-mac", "agent-macos.pkg")
}

func serveAgentFile(c *gin.Context, fileName string, downloadName string) {
	// Look in both current directory release/ and parent directory release/ to support different CWDs
	paths := []string{
		filepath.Join("release", fileName),
		filepath.Join("..", "release", fileName),
		filepath.Join("backend", "release", fileName),
	}

	var foundPath string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			foundPath = p
			break
		}
	}

	if foundPath == "" {
		c.JSON(404, gin.H{"error": "Agent installer file not found on server"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(foundPath)
}

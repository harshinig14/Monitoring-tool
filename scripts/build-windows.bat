@echo off

echo Building Windows Binary...

SET GOOS=windows
SET GOARCH=amd64

go build -o dist/windows/MONITORING-TOOL.exe ./cmd/agent

echo Windows Build Completed
pause
package main

import "embed"

//go:embed binaries/**
var embeddedFiles embed.FS

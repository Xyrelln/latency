package main

import (
	"embed"
	"op-latency-mobile/internal/app"
	"os"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	os.Exit(app.Run(assets))
}

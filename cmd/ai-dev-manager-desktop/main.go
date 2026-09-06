package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/desktop"
	"ai-dev-manager-v2/internal/management"
	"ai-dev-manager-v2/internal/store"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var embeddedFrontend embed.FS

func main() {
	if err := runDesktop(); err != nil {
		fmt.Fprintln(os.Stderr, "desktop error:", err)
		os.Exit(1)
	}
}

func runDesktop() error {
	statePath, err := store.DefaultPath()
	if err != nil {
		return err
	}
	assets, err := frontendAssets()
	if err != nil {
		return err
	}
	adapter := newDesktopAdapter(statePath)
	return wails.Run(&options.App{
		Title:     "AI Dev Manager V2",
		Width:     1120,
		Height:    760,
		MinWidth:  820,
		MinHeight: 560,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 246, G: 247, B: 249, A: 1},
		Bind: []interface{}{
			adapter,
		},
	})
}

func newDesktopAdapter(statePath string) *desktop.Adapter {
	application := app.New(statePath)
	return desktop.NewAdapter(management.New(application))
}

func frontendAssets() (fs.FS, error) {
	return fs.Sub(embeddedFrontend, "frontend")
}

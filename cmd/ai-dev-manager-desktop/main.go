package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"syscall"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/desktop"
	"ai-dev-manager-v2/internal/gateway"
	"ai-dev-manager-v2/internal/store"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var embeddedFrontend embed.FS

func main() {
	var err error
	if len(os.Args) > 1 && os.Args[1] == "--gateway-child" {
		err = runGatewayChild(os.Args[2:])
	} else {
		err = runDesktop()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "desktop error:", err)
		os.Exit(1)
	}
}

func runGatewayChild(args []string) error {
	fs := flag.NewFlagSet("gateway-child", flag.ContinueOnError)
	listen := fs.String("listen", gateway.DefaultHTTPListen, "Gateway listen address")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("gateway child accepts only --listen")
	}
	statePath, err := store.DefaultPath()
	if err != nil {
		return err
	}
	service := app.New(statePath)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return gateway.RunHTTP(ctx, service, *listen)
}

func runDesktop() error {
	assets, err := frontendAssets()
	if err != nil {
		return err
	}
	adapter := desktop.NewClientAdapter()
	return wails.Run(&options.App{
		Title:     "AI Dev Manager V2 — 1.0 RC3",
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

func frontendAssets() (fs.FS, error) {
	return fs.Sub(embeddedFrontend, "frontend")
}

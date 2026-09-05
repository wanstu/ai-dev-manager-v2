package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/gateway"
	"ai-dev-manager-v2/internal/store"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	statePath, err := store.DefaultPath()
	if err != nil {
		return err
	}
	service := app.New(statePath)
	if len(args) == 0 {
		return usageError()
	}

	switch args[0] {
	case "workspace":
		return runWorkspace(service, args[1:])
	case "env", "environment":
		return runEnvironment(service, args[1:])
	case "exec":
		return runExec(service, args[1:])
	case "gateway":
		return runGateway(service, args[1:])
	case "state":
		if len(args) == 2 && args[1] == "path" {
			fmt.Println(statePath)
			return nil
		}
		return fmt.Errorf("usage: ai-dev-manager state path")
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return usageError()
	}
}

func runWorkspace(service *app.Service, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ai-dev-manager workspace <add|list>")
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("workspace add", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		name := fs.String("name", "", "workspace display name")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if fs.NArg() != 1 {
			return fmt.Errorf("usage: ai-dev-manager workspace add [--name NAME] PATH")
		}
		ws, err := service.Workspaces.Add(fs.Arg(0), *name)
		if err != nil {
			return err
		}
		return writeJSON(ws)
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("usage: ai-dev-manager workspace list")
		}
		items, err := service.Workspaces.List()
		if err != nil {
			return err
		}
		return writeJSON(items)
	default:
		return fmt.Errorf("unknown workspace command %q", args[0])
	}
}

func runEnvironment(service *app.Service, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ai-dev-manager env <create|list|inspect|writer>")
	}
	switch args[0] {
	case "create":
		fs := flag.NewFlagSet("env create", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		workspaceID := fs.String("workspace", "", "workspace id")
		name := fs.String("name", "", "environment name")
		root := fs.String("root", "", "optional root inside workspace")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if fs.NArg() != 0 || strings.TrimSpace(*workspaceID) == "" || strings.TrimSpace(*name) == "" {
			return fmt.Errorf("usage: ai-dev-manager env create --workspace WS_ID --name NAME [--root PATH]")
		}
		env, err := service.Environments.Create(*workspaceID, *name, *root)
		if err != nil {
			return err
		}
		return writeJSON(env)
	case "list":
		items, err := service.Environments.List()
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "inspect":
		if len(args) != 2 {
			return fmt.Errorf("usage: ai-dev-manager env inspect ENV_ID")
		}
		env, err := service.Environments.Get(args[1])
		if err != nil {
			return err
		}
		caps, err := service.Capabilities(context.Background(), env.ID)
		if err != nil {
			return err
		}
		return writeJSON(map[string]any{"environment": env, "capabilities": caps})
	case "writer":
		return runWriter(service, args[1:])
	default:
		return fmt.Errorf("unknown env command %q", args[0])
	}
}

func runWriter(service *app.Service, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ai-dev-manager env writer <acquire|release>")
	}
	switch args[0] {
	case "acquire":
		fs := flag.NewFlagSet("env writer acquire", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		owner := fs.String("owner", "", "stable agent/session owner")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if fs.NArg() != 1 || strings.TrimSpace(*owner) == "" {
			return fmt.Errorf("usage: ai-dev-manager env writer acquire --owner OWNER ENV_ID")
		}
		env, err := service.Environments.AcquireWriter(fs.Arg(0), *owner)
		if err != nil {
			return err
		}
		return writeJSON(env)
	case "release":
		fs := flag.NewFlagSet("env writer release", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		owner := fs.String("owner", "", "current writer owner")
		force := fs.Bool("force", false, "release regardless of owner")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if fs.NArg() != 1 || (!*force && strings.TrimSpace(*owner) == "") {
			return fmt.Errorf("usage: ai-dev-manager env writer release [--owner OWNER|--force] ENV_ID")
		}
		env, err := service.Environments.ReleaseWriter(fs.Arg(0), *owner, *force)
		if err != nil {
			return err
		}
		return writeJSON(env)
	default:
		return fmt.Errorf("unknown writer command %q", args[0])
	}
}

func runExec(service *app.Service, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ai-dev-manager exec <allow|list>")
	}
	switch args[0] {
	case "allow":
		if len(args) != 2 {
			return fmt.Errorf("usage: ai-dev-manager exec allow EXECUTABLE")
		}
		if err := service.AllowExecutable(args[1]); err != nil {
			return err
		}
		items, err := service.AllowedExecutables()
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "list":
		items, err := service.AllowedExecutables()
		if err != nil {
			return err
		}
		return writeJSON(items)
	default:
		return fmt.Errorf("unknown exec command %q", args[0])
	}
}

func runGateway(service *app.Service, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ai-dev-manager gateway <stdio|http>")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	switch args[0] {
	case "stdio":
		if len(args) != 1 {
			return fmt.Errorf("usage: ai-dev-manager gateway stdio")
		}
		return gateway.RunStdio(ctx, service)
	case "http":
		fs := flag.NewFlagSet("gateway http", flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
		listen := fs.String("listen", "", "loopback listen address, for example 127.0.0.1:41137")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if fs.NArg() != 0 || strings.TrimSpace(*listen) == "" {
			return fmt.Errorf("usage: ai-dev-manager gateway http --listen 127.0.0.1:PORT")
		}
		return gateway.RunHTTP(ctx, service, *listen)
	default:
		return fmt.Errorf("unknown gateway command %q", args[0])
	}
}

func writeJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func usageError() error {
	printUsage()
	return fmt.Errorf("invalid command")
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `ai-dev-manager-v2

Commands:
  workspace add [--name NAME] PATH
  workspace list
  env create --workspace WS_ID --name NAME [--root PATH]
  env list
  env inspect ENV_ID
  env writer acquire --owner OWNER ENV_ID
  env writer release [--owner OWNER|--force] ENV_ID
  exec allow EXECUTABLE
  exec list
  gateway stdio
  gateway http --listen 127.0.0.1:PORT
  state path`)
}

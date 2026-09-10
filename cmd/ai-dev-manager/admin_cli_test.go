package main

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/gateway"
)

func startCLIAdminTestServer(t *testing.T, home string) *app.Service {
	t.Helper()
	service := app.New(filepath.Join(home, "state.json"))
	server := httptest.NewServer(gateway.NewHTTPHandler(service))
	t.Cleanup(server.Close)
	t.Setenv(admBaseURLEnv, server.URL)
	return service
}

func TestTopLevelCLIManagementUsesAdminMCPWithoutLocalStateFallback(t *testing.T) {
	localHome := t.TempDir()
	t.Setenv("ADM_V2_HOME", localHome)

	remoteState := filepath.Join(t.TempDir(), "remote-state.json")
	remote := app.New(remoteState)
	server := httptest.NewServer(gateway.NewHTTPHandler(remote))
	defer server.Close()
	t.Setenv(admBaseURLEnv, server.URL)

	root := t.TempDir()
	captureStdout(t, func() {
		if err := run([]string{"workspace", "add", "--path", root, "--name", "remote-only"}); err != nil {
			t.Fatal(err)
		}
	})
	items, err := remote.Workspaces.List()
	if err != nil || len(items) != 1 || items[0].Path != root {
		t.Fatalf("remote Admin MCP workspace state=%+v err=%v", items, err)
	}
	if _, err := os.Stat(filepath.Join(localHome, "state.json")); !os.IsNotExist(err) {
		t.Fatalf("normal CLI management unexpectedly touched local state.json: %v", err)
	}
}

func TestTopLevelCLIManagementDoesNotFallbackWhenAdminMCPStops(t *testing.T) {
	localHome := t.TempDir()
	t.Setenv("ADM_V2_HOME", localHome)
	local := app.New(filepath.Join(localHome, "state.json"))
	if _, err := local.Workspaces.Add(t.TempDir(), "local-only"); err != nil {
		t.Fatal(err)
	}

	remote := app.New(filepath.Join(t.TempDir(), "remote-state.json"))
	server := httptest.NewServer(gateway.NewHTTPHandler(remote))
	baseURL := server.URL
	server.Close()
	t.Setenv(admBaseURLEnv, baseURL)

	err := run([]string{"workspace", "list"})
	if err == nil || !strings.Contains(err.Error(), "Admin MCP") {
		t.Fatalf("stopped Admin MCP should fail without local fallback, got %v", err)
	}
}

func TestCLIADMURLFlagOverridesEnvironmentAndSupportsBasePath(t *testing.T) {
	remote := app.New(filepath.Join(t.TempDir(), "remote-state.json"))
	server := httptest.NewServer(gateway.NewHTTPHandler(remote))
	defer server.Close()
	t.Setenv(admBaseURLEnv, "http://127.0.0.1:1")

	output := captureStdout(t, func() {
		if err := run([]string{"--adm-url", server.URL, "workspace", "list"}); err != nil {
			t.Fatal(err)
		}
	})
	if strings.TrimSpace(output) != "[]" {
		t.Fatalf("workspace list output=%q", output)
	}

	if got, err := normalizeCLIADMBaseURL("https://adm.example.test:8443/control/"); err != nil || got != "https://adm.example.test:8443/control" {
		t.Fatalf("normalized base URL=%q err=%v", got, err)
	}
}

func TestCLIManagementHelpDoesNotRequireRunningADM(t *testing.T) {
	t.Setenv(admBaseURLEnv, "http://127.0.0.1:1")
	output := captureStdout(t, func() {
		if err := run([]string{"workspace", "-h"}); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(output, "workspace add --path") {
		t.Fatalf("workspace help missing expected usage: %s", output)
	}
}

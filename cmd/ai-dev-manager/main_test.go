package main

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/app"
)

func TestTopLevelHelpExplainsQuickStartAndGatewayLifecycle(t *testing.T) {
	output := captureStdout(t, func() {
		if err := run([]string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{
		"快速开始（HTTP Gateway）",
		"ai-dev-manager-v2 gateway start",
		"gateway status",
		"gateway stop",
		"gateway restart",
		"doctor",
		"仅供 MCP 客户端使用",
	} {
		if !strings.Contains(output, required) {
			t.Fatalf("top-level help missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayHelpExplainsForegroundHTTPAndClientOnlyStdio(t *testing.T) {
	output := captureStdout(t, func() {
		if err := runGateway(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{
		"当前终端前台启动",
		"--detach",
		"Ctrl+C 停止",
		"运行状态",
		"人不要手动运行",
	} {
		if !strings.Contains(output, required) {
			t.Fatalf("gateway help missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayStartDetachRoutesToDetachedLauncher(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	original := startGatewayDetached
	t.Cleanup(func() { startGatewayDetached = original })

	for _, flagName := range []string{"--detach", "-d"} {
		var gotListen string
		startGatewayDetached = func(listen string) error {
			gotListen = listen
			return nil
		}
		if err := runGateway(service, []string{"start", flagName, "--listen", "127.0.0.1:45555"}); err != nil {
			t.Fatal(err)
		}
		if gotListen != "127.0.0.1:45555" {
			t.Fatalf("%s listen = %q", flagName, gotListen)
		}
	}
}

func TestCLIRejectsOldPositionalWorkspaceAddGrammar(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	err := runWorkspace(service, []string{"add", t.TempDir()})
	if err == nil || !strings.Contains(err.Error(), "--path") {
		t.Fatalf("old positional grammar should explain --path requirement, got %v", err)
	}
}

func TestDoctorRunsWithoutArguments(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	service := app.New(statePath)
	output := captureStdout(t, func() {
		if err := runDoctor(service, statePath, nil); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"ADM V2 诊断", statePath, "Gateway", "Workspace（0）", "Environment（0）", "执行白名单（0）"} {
		if !strings.Contains(output, required) {
			t.Fatalf("doctor output missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayStatusReportsRunningGatewayDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"name":"ai-dev-manager-v2","version":"test-version","status":"ok","pid":43210,"transport":"http"}`)
	}))
	defer server.Close()
	listen := strings.TrimPrefix(server.URL, "http://")

	output := captureStdout(t, func() {
		if err := printGatewayStatus(listen); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"运行中", "43210", "test-version", "/mcp"} {
		if !strings.Contains(output, required) {
			t.Fatalf("gateway status missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayStatusReportsIncompatibleListenerClearly(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	listen := strings.TrimPrefix(server.URL, "http://")

	output := captureStdout(t, func() {
		if err := printGatewayStatus(listen); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"版本不兼容", "gateway restart"} {
		if !strings.Contains(output, required) {
			t.Fatalf("incompatible status missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayStatusReportsStoppedWhenNothingListens(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	listen := listener.Addr().String()
	_ = listener.Close()

	output := captureStdout(t, func() {
		if err := printGatewayStatus(listen); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(output, "已停止") || !strings.Contains(output, "gateway start") {
		t.Fatalf("stopped status is not actionable:\n%s", output)
	}
}

func TestSameADMExecutableAllowsSiblingUpgradeBinary(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows executable-path semantics")
	}
	base := filepath.Join(`D:\projects\ai-dev-manager-v2`, "ai-dev-manager-v2.exe")
	next := filepath.Join(`D:\projects\ai-dev-manager-v2`, "ai-dev-manager-v2.next.exe")
	if !sameADMExecutable(base, next) {
		t.Fatal("next build in the same directory should be allowed to replace the official ADM V2 executable")
	}
	if sameADMExecutable(`D:\other\ai-dev-manager-v2.exe`, next) {
		t.Fatal("different directories must not be treated as the same ADM installation")
	}
	if sameADMExecutable(`D:\projects\ai-dev-manager-v2\other.exe`, next) {
		t.Fatal("another executable in the same directory must not be treated as ADM V2")
	}
}

func TestStopHTTPGatewayCanStopLegacySameExecutableOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("legacy Gateway process discovery is Windows-specific")
	}
	originalMatcher := matchesADMExecutable
	matchesADMExecutable = func(targetPath, currentPath string) bool {
		return strings.EqualFold(filepath.Clean(targetPath), filepath.Clean(currentPath))
	}
	t.Cleanup(func() { matchesADMExecutable = originalMatcher })

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	listen := listener.Addr().String()
	_ = listener.Close()

	cmd := exec.Command(os.Args[0], "-test.run=^TestLegacyGatewayHelperProcess$")
	cmd.Env = append(os.Environ(), "ADM_TEST_LEGACY_GATEWAY=1", "ADM_TEST_LEGACY_LISTEN="+listen)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	deadline := time.Now().Add(5 * time.Second)
	ready := false
	for time.Now().Before(deadline) {
		response, requestErr := http.Get("http://" + listen + "/healthz")
		if requestErr == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusNotFound {
				ready = true
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	if !ready {
		t.Fatal("legacy Gateway helper did not start")
	}

	output := captureStdout(t, func() {
		if err := stopHTTPGateway(listen); err != nil {
			t.Fatal(err)
		}
	})
	_ = cmd.Wait()
	if !strings.Contains(output, "检测到旧版 ADM V2 Gateway") || !strings.Contains(output, "已停止") {
		t.Fatalf("legacy Gateway stop output is not clear:\n%s", output)
	}
	if connection, err := net.DialTimeout("tcp", listen, 200*time.Millisecond); err == nil {
		_ = connection.Close()
		t.Fatalf("legacy Gateway still listens on %s", listen)
	}
}

func TestLegacyGatewayHelperProcess(t *testing.T) {
	if os.Getenv("ADM_TEST_LEGACY_GATEWAY") != "1" {
		return
	}
	listen := os.Getenv("ADM_TEST_LEGACY_LISTEN")
	server := &http.Server{
		Addr: listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" {
				http.NotFound(w, r)
				return
			}
			w.WriteHeader(http.StatusOK)
		}),
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		os.Exit(2)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	defer func() { os.Stdout = original }()

	fn()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	_ = reader.Close()
	return string(data)
}

package gateway

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestInspectHTTPReportsRunningCompatibleGateway(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"name":%q,"version":"test","status":"ok","pid":%d,"transport":"http"}`, serverName, os.Getpid())
	}))
	defer server.Close()

	listen := strings.TrimPrefix(server.URL, "http://")
	status, err := InspectHTTP(listen)
	if err != nil {
		t.Fatal(err)
	}
	if status.State != HTTPStateRunning || status.PID != os.Getpid() || status.Version != "test" {
		t.Fatalf("status = %+v", status)
	}
	if status.MCPURL != server.URL+"/mcp" {
		t.Fatalf("mcp_url=%q want %q", status.MCPURL, server.URL+"/mcp")
	}
}

func TestInspectHTTPDistinguishesStoppedAndIncompatibleEndpoints(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	listen := listener.Addr().String()
	_ = listener.Close()
	stopped, err := InspectHTTP(listen)
	if err != nil {
		t.Fatal(err)
	}
	if stopped.State != HTTPStateStopped {
		t.Fatalf("stopped status = %+v", stopped)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()
	incompatibleListen := strings.TrimPrefix(server.URL, "http://")
	incompatible, err := InspectHTTP(incompatibleListen)
	if err != nil {
		t.Fatal(err)
	}
	if incompatible.State != HTTPStateIncompatible || incompatible.Detail == "" {
		t.Fatalf("incompatible status = %+v", incompatible)
	}
	if _, err := StopHTTP(incompatibleListen); err == nil || !strings.Contains(err.Error(), "refusing") {
		t.Fatalf("StopHTTP must refuse incompatible endpoint, got %v", err)
	}
}

func TestHTTPBaseURLRejectsInvalidListen(t *testing.T) {
	if _, err := HTTPBaseURL("not-an-endpoint"); err == nil {
		t.Fatal("invalid listen must fail")
	}
	got, err := HTTPBaseURL("127.0.0.1:41137")
	if err != nil || got != "http://127.0.0.1:41137" {
		t.Fatalf("baseURL=%q err=%v", got, err)
	}
}

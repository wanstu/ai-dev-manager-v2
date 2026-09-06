package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const DefaultHTTPListen = "127.0.0.1:41137"

const (
	HTTPStateStopped      = "stopped"
	HTTPStateRunning      = "running"
	HTTPStateIncompatible = "incompatible"
)

type HTTPHealth struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	PID       int    `json:"pid"`
	Transport string `json:"transport"`
}

type HTTPStatus struct {
	State   string `json:"state"`
	Listen  string `json:"listen"`
	BaseURL string `json:"base_url"`
	MCPURL  string `json:"mcp_url"`
	PID     int    `json:"pid,omitempty"`
	Version string `json:"version,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func HTTPBaseURL(listen string) (string, error) {
	listen = strings.TrimSpace(listen)
	host, port, err := net.SplitHostPort(listen)
	if err != nil {
		return "", fmt.Errorf("invalid Gateway listen address %q; expected host:port", listen)
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return "http://" + net.JoinHostPort(host, port), nil
}

func InspectHTTP(listen string) (HTTPStatus, error) {
	listen = strings.TrimSpace(listen)
	baseURL, err := HTTPBaseURL(listen)
	if err != nil {
		return HTTPStatus{}, err
	}
	status := HTTPStatus{State: HTTPStateStopped, Listen: listen, BaseURL: baseURL, MCPURL: baseURL + "/mcp"}
	client := &http.Client{Timeout: 1200 * time.Millisecond}
	response, err := client.Get(baseURL + "/healthz")
	if err != nil {
		if isHTTPConnectionFailure(err) {
			return status, nil
		}
		return HTTPStatus{}, fmt.Errorf("check Gateway %s: %w", baseURL, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		status.State = HTTPStateIncompatible
		status.Detail = fmt.Sprintf("port %s responded with %s instead of an ADM V2 health response", listen, response.Status)
		return status, nil
	}
	var health HTTPHealth
	if err := json.NewDecoder(response.Body).Decode(&health); err != nil {
		status.State = HTTPStateIncompatible
		status.Detail = fmt.Sprintf("port %s returned invalid ADM V2 health JSON: %v", listen, err)
		return status, nil
	}
	if health.Name != serverName || health.Status != "ok" || health.Transport != "http" {
		status.State = HTTPStateIncompatible
		status.Detail = fmt.Sprintf("port %s is not the expected ADM V2 HTTP Gateway", listen)
		return status, nil
	}
	status.State = HTTPStateRunning
	status.PID = health.PID
	status.Version = health.Version
	return status, nil
}

func WaitHTTPReady(listen string, timeout time.Duration) (HTTPStatus, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		status, err := InspectHTTP(listen)
		if err != nil {
			return HTTPStatus{}, err
		}
		switch status.State {
		case HTTPStateRunning:
			return status, nil
		case HTTPStateIncompatible:
			return status, fmt.Errorf("Gateway endpoint is incompatible: %s", status.Detail)
		}
		time.Sleep(100 * time.Millisecond)
	}
	baseURL, _ := HTTPBaseURL(listen)
	return HTTPStatus{}, fmt.Errorf("Gateway did not become ready within %s: %s/healthz", timeout, baseURL)
}

func StopHTTP(listen string) (HTTPStatus, error) {
	status, err := InspectHTTP(listen)
	if err != nil {
		return HTTPStatus{}, err
	}
	switch status.State {
	case HTTPStateStopped:
		return status, nil
	case HTTPStateIncompatible:
		return status, fmt.Errorf("refusing to stop incompatible process on %s: %s", listen, status.Detail)
	}
	if status.PID <= 0 {
		return status, fmt.Errorf("Gateway %s did not provide a usable PID", status.BaseURL)
	}
	if err := TerminateHTTPProcess(status.PID, listen); err != nil {
		return status, err
	}
	return InspectHTTP(listen)
}

func TerminateHTTPProcess(pid int, listen string) error {
	baseURL, err := HTTPBaseURL(listen)
	if err != nil {
		return err
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find Gateway process %d: %w", pid, err)
	}
	if err := process.Kill(); err != nil {
		return fmt.Errorf("stop Gateway process %d: %w", pid, err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		connection, dialErr := net.DialTimeout("tcp", listen, 200*time.Millisecond)
		if dialErr != nil {
			return nil
		}
		_ = connection.Close()
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("Gateway process %d was terminated but endpoint %s still responds", pid, baseURL)
}

func isHTTPConnectionFailure(err error) bool {
	var netErr net.Error
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "connection refused") || strings.Contains(message, "actively refused") || (errors.As(err, &netErr) && netErr.Timeout())
}

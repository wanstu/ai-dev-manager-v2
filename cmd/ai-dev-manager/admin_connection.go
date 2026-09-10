package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"

	"ai-dev-manager-v2/internal/adminmcp"
	"ai-dev-manager-v2/internal/gateway"
)

const admBaseURLEnv = "ADM_V2_URL"

func defaultADMBaseURL() string {
	baseURL, err := gateway.HTTPBaseURL(gateway.DefaultHTTPListen)
	if err != nil {
		return "http://127.0.0.1:43137"
	}
	return baseURL
}

func parseADMTarget(args []string) (string, []string, error) {
	baseURL := strings.TrimSpace(os.Getenv(admBaseURLEnv))
	if baseURL == "" {
		baseURL = defaultADMBaseURL()
	}
	if len(args) == 0 {
		return baseURL, args, nil
	}

	filtered := make([]string, 0, len(args))
	seen := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--adm-url":
			if seen {
				return "", nil, fmt.Errorf("--adm-url may only be provided once")
			}
			if i+1 >= len(args) || strings.TrimSpace(args[i+1]) == "" {
				return "", nil, fmt.Errorf("--adm-url requires an ADM base URL")
			}
			baseURL = strings.TrimSpace(args[i+1])
			seen = true
			i++
		case strings.HasPrefix(arg, "--adm-url="):
			if seen {
				return "", nil, fmt.Errorf("--adm-url may only be provided once")
			}
			baseURL = strings.TrimSpace(strings.TrimPrefix(arg, "--adm-url="))
			if baseURL == "" {
				return "", nil, fmt.Errorf("--adm-url requires an ADM base URL")
			}
			seen = true
		default:
			filtered = append(filtered, arg)
		}
	}
	return baseURL, filtered, nil
}

func localGatewayListenFromBaseURL(raw string) (string, error) {
	baseURL, err := normalizeCLIADMBaseURL(raw)
	if err != nil {
		return "", err
	}
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "http" {
		return "", fmt.Errorf("local Gateway lifecycle requires an http ADM base URL; use gateway status for remote HTTPS health checks")
	}
	if strings.Trim(parsed.Path, "/") != "" {
		return "", fmt.Errorf("local Gateway lifecycle requires an ADM base URL without a base path")
	}
	host := parsed.Hostname()
	port := parsed.Port()
	if host == "" || port == "" {
		return "", fmt.Errorf("local Gateway lifecycle requires an explicit loopback host and port")
	}
	loopback := strings.EqualFold(host, "localhost")
	if !loopback {
		ip := net.ParseIP(host)
		loopback = ip != nil && ip.IsLoopback()
	}
	if !loopback {
		return "", fmt.Errorf("local Gateway lifecycle is only available for loopback ADM URLs")
	}
	return net.JoinHostPort(host, port), nil
}

func newCLIAdminClient(baseURL string) (*adminmcp.Client, error) {
	normalized, err := normalizeCLIADMBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	return adminmcp.New(normalized + "/admin/mcp"), nil
}

func normalizeCLIADMBaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("ADM base URL is required")
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid ADM base URL %q: %w", raw, err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("unsupported ADM URL scheme %q; expected http or https", parsed.Scheme)
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("ADM base URL %q requires a host", raw)
	}
	if parsed.User != nil {
		return "", fmt.Errorf("ADM base URL must not contain userinfo")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("ADM base URL must not contain query or fragment")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	parsed.RawPath = ""
	return parsed.String(), nil
}

func withCLIAdmin(baseURL string, fn func(cliManagementBackend) error) error {
	client, err := newCLIAdminClient(baseURL)
	if err != nil {
		return err
	}
	return fn(client)
}

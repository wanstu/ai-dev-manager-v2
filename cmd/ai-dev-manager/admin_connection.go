package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"ai-dev-manager-v2/internal/adminmcp"
)

const (
	defaultADMBaseURL = "http://127.0.0.1:41137"
	admBaseURLEnv     = "ADM_V2_URL"
)

func parseADMTarget(args []string) (string, []string, error) {
	baseURL := strings.TrimSpace(os.Getenv(admBaseURLEnv))
	if baseURL == "" {
		baseURL = defaultADMBaseURL
	}
	if len(args) == 0 {
		return baseURL, args, nil
	}
	if args[0] == "--adm-url" {
		if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
			return "", nil, fmt.Errorf("--adm-url requires an ADM base URL")
		}
		baseURL = strings.TrimSpace(args[1])
		args = args[2:]
	} else if strings.HasPrefix(args[0], "--adm-url=") {
		baseURL = strings.TrimSpace(strings.TrimPrefix(args[0], "--adm-url="))
		if baseURL == "" {
			return "", nil, fmt.Errorf("--adm-url requires an ADM base URL")
		}
		args = args[1:]
	}
	return baseURL, args, nil
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

//go:build windows

package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/desktop"
)

func TestActiveConnectionProfile(t *testing.T) {
	profiles := desktop.ConnectionProfiles{
		Profiles: []desktop.ConnectionProfile{
			{ID: "a", Name: "A", BaseURL: "http://127.0.0.1:43137"},
			{ID: "b", Name: "B", BaseURL: "http://127.0.0.1:43138"},
		},
		ActiveID: "b",
	}
	profile, ok := activeConnectionProfile(profiles)
	if !ok || profile.ID != "b" {
		t.Fatalf("active profile = %#v, %v; want b", profile, ok)
	}
	profiles.ActiveID = "missing"
	if _, ok := activeConnectionProfile(profiles); ok {
		t.Fatal("missing active profile should not resolve")
	}
}

func TestBuildTrayQuitPromptOffersStopOnlyForRunningLocalADM(t *testing.T) {
	profile := desktop.ConnectionProfile{ID: "local", Name: "本地 ADM", BaseURL: "http://127.0.0.1:43137"}
	prompt := buildTrayQuitPrompt(profile, desktop.ADMConnectionStatus{
		State: "running", BaseURL: profile.BaseURL, LocalBootstrapEligible: true,
	}, nil)
	wantButtons := []string{trayQuitStopBackground, trayQuitKeepBackground, trayQuitCancel}
	if !prompt.canStopLocal || !reflect.DeepEqual(prompt.buttons, wantButtons) {
		t.Fatalf("local running prompt = %#v; want stop/keep/cancel", prompt)
	}
	if prompt.defaultButton != trayQuitKeepBackground {
		t.Fatalf("default button = %q; want keep background", prompt.defaultButton)
	}
	if !strings.Contains(prompt.message, "/mcp") || !strings.Contains(prompt.message, "/admin/mcp") {
		t.Fatalf("prompt should explain MCP endpoints: %q", prompt.message)
	}
}

func TestBuildTrayQuitPromptNeverOffersStopForRemoteOrUnknown(t *testing.T) {
	remote := desktop.ConnectionProfile{ID: "remote", Name: "Remote", BaseURL: "https://adm.example.test"}
	prompt := buildTrayQuitPrompt(remote, desktop.ADMConnectionStatus{
		State: "running", BaseURL: remote.BaseURL, LocalBootstrapEligible: false,
	}, nil)
	if prompt.canStopLocal || containsStringValue(prompt.buttons, trayQuitStopBackground) {
		t.Fatalf("remote prompt must not offer stop: %#v", prompt)
	}

	prompt = buildTrayQuitPrompt(remote, desktop.ADMConnectionStatus{}, errors.New("health check failed"))
	if prompt.canStopLocal || containsStringValue(prompt.buttons, trayQuitStopBackground) {
		t.Fatalf("unknown prompt must not offer stop: %#v", prompt)
	}
	if !strings.Contains(prompt.message, "不会尝试停止") {
		t.Fatalf("unknown prompt must explain safe fallback: %q", prompt.message)
	}
}

func containsStringValue(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

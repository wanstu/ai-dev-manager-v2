//go:build windows

package main

import (
	"os"
	"strings"
	"testing"
)

func TestTrayExitOnlyQuitsDesktop(t *testing.T) {
	source, err := os.ReadFile("tray_manager_windows.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	for _, required := range []string{
		`trayQuitLabel = "退出"`,
		`menu.Add(trayQuitLabel, t.quitDesktop)`,
		`func (t *trayManager) quitDesktop()`,
		`wailsruntime.Quit(ctx)`,
	} {
		if !strings.Contains(text, required) {
			t.Fatalf("tray source missing %q", required)
		}
	}
	for _, forbidden := range []string{
		`StopLocalADM`,
		`stopLocalBackgroundAndQuit`,
		`activeConnectionStatusForQuit`,
		`canSafelyStopLocalBackground`,
		`停止本地后台服务并退出`,
		`停止后台`,
		`QuestionDialog`,
		`保留后台服务`,
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("tray exit must not include %q; stopping CLI/MCP remains a separate explicit main-UI action", forbidden)
		}
	}
}

//go:build windows

package main

import (
	"os"
	"strings"
	"testing"
)

func TestTrayExitOnlyQuitsDesktopAndPreservesBackgroundService(t *testing.T) {
	source, err := os.ReadFile("tray_manager_windows.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	for _, required := range []string{
		`trayQuitKeepBackgroundLabel = "退出 Desktop（保留后台服务）"`,
		`menu.Add(trayQuitKeepBackgroundLabel, t.quitKeepBackground)`,
		`func (t *trayManager) quitKeepBackground()`,
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
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("tray exit must not include %q; CLI service stop is manual from the main UI", forbidden)
		}
	}
}

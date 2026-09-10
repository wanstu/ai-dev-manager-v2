//go:build windows

package main

import (
	"context"
	"sync"

	"ai-dev-manager-v2/internal/desktop"

	"fyne.io/systray"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const traySupported = true

type trayManager struct {
	adapter *desktop.Adapter
	icon    []byte

	mu        sync.RWMutex
	ctx       context.Context
	start     func()
	end       func()
	startOnce sync.Once
	endOnce   sync.Once
}

func newTrayManager(icon []byte, adapter *desktop.Adapter) *trayManager {
	return &trayManager{adapter: adapter, icon: icon}
}

func (t *trayManager) Startup(ctx context.Context) {
	t.mu.Lock()
	t.ctx = ctx
	t.mu.Unlock()
	t.startOnce.Do(func() {
		// External-loop setup can invoke onReady immediately. Wails must have
		// supplied its runtime context and accepted the single-instance lock first.
		t.start, t.end = systray.RunWithExternalLoop(t.onReady, func() {})
		if t.start != nil {
			t.start()
		}
	})
}

func (t *trayManager) Shutdown(context.Context) {
	t.endOnce.Do(func() {
		if t.end != nil {
			t.end()
		}
	})
}

func (t *trayManager) ShowWindow() {
	ctx := t.runtimeContext()
	if ctx == nil {
		return
	}
	wailsruntime.WindowShow(ctx)
	wailsruntime.WindowUnminimise(ctx)
}

func (t *trayManager) hideWindow() {
	ctx := t.runtimeContext()
	if ctx != nil {
		wailsruntime.WindowHide(ctx)
	}
}

func (t *trayManager) quit() {
	ctx := t.runtimeContext()
	if ctx != nil {
		wailsruntime.Quit(ctx)
	}
}

func (t *trayManager) runtimeContext() context.Context {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.ctx
}

func (t *trayManager) onReady() {
	if len(t.icon) != 0 {
		systray.SetIcon(t.icon)
	}
	systray.SetTooltip("AI Dev Manager V2")
	systray.SetOnTapped(t.ShowWindow)

	showItem := systray.AddMenuItem("显示主窗口", "显示 AI Dev Manager V2")
	hideItem := systray.AddMenuItem("隐藏主窗口", "隐藏到系统托盘")
	systray.AddSeparator()
	launchItem := systray.AddMenuItemCheckbox("开机启动", "登录 Windows 后在托盘中启动", false)
	if t.adapter == nil {
		launchItem.Disable()
	} else if preferences, err := t.adapter.GetDesktopPreferences(); err != nil || !preferences.LaunchAtLoginSupported {
		launchItem.Disable()
	} else if preferences.LaunchAtLogin {
		launchItem.Check()
	}
	wailsruntime.EventsOn(t.runtimeContext(), "desktop:preferences-changed", func(...interface{}) {
		if t.adapter == nil {
			return
		}
		preferences, err := t.adapter.GetDesktopPreferences()
		if err != nil {
			return
		}
		if preferences.LaunchAtLogin {
			launchItem.Check()
		} else {
			launchItem.Uncheck()
		}
	})
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("退出", "退出 AI Dev Manager V2 Desktop")

	go func() {
		for range showItem.ClickedCh {
			t.ShowWindow()
		}
	}()
	go func() {
		for range hideItem.ClickedCh {
			t.hideWindow()
		}
	}()
	go func() {
		for range launchItem.ClickedCh {
			if t.adapter == nil {
				continue
			}
			preferences, err := t.adapter.GetDesktopPreferences()
			if err != nil || !preferences.LaunchAtLoginSupported {
				continue
			}
			preferences, err = t.adapter.SetLaunchAtLogin(!preferences.LaunchAtLogin)
			if err != nil {
				t.ShowWindow()
				_, _ = wailsruntime.MessageDialog(t.runtimeContext(), wailsruntime.MessageDialogOptions{
					Type: wailsruntime.ErrorDialog, Title: "开机启动设置失败", Message: err.Error(),
				})
				continue
			}
			wailsruntime.EventsEmit(t.runtimeContext(), "desktop:preferences-changed")
			if preferences.LaunchAtLogin {
				launchItem.Check()
			} else {
				launchItem.Uncheck()
			}
		}
	}()
	go func() {
		for range quitItem.ClickedCh {
			t.quit()
		}
	}()
}

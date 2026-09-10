//go:build windows

package main

import (
	"context"
	"runtime"
	"sync"

	"ai-dev-manager-v2/internal/desktop"

	"github.com/gogpu/systray"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const traySupported = true

type trayManager struct {
	adapter *desktop.Adapter
	icon    []byte

	mu         sync.RWMutex
	ctx        context.Context
	tray       *systray.SystemTray
	launchItem *systray.MenuItem
	stopEvents func()
	started    bool
	stopping   bool
}

func newTrayManager(icon []byte, adapter *desktop.Adapter) *trayManager {
	return &trayManager{adapter: adapter, icon: icon}
}

// Startup only captures the Wails runtime context. The tray message loop is
// intentionally started from DomReady, matching the proven ime-lock-v2
// lifecycle and avoiding tray callbacks before the WebView/runtime is ready.
func (t *trayManager) Startup(ctx context.Context) {
	t.mu.Lock()
	t.ctx = ctx
	t.mu.Unlock()
}

func (t *trayManager) DomReady(ctx context.Context) {
	t.mu.Lock()
	t.ctx = ctx
	if t.started {
		t.mu.Unlock()
		return
	}
	t.started = true
	t.stopping = false
	t.mu.Unlock()

	go t.run()
}

func (t *trayManager) Shutdown(context.Context) {
	if t == nil {
		return
	}
	t.mu.Lock()
	t.stopping = true
	tray := t.tray
	stopEvents := t.stopEvents
	t.stopEvents = nil
	t.mu.Unlock()
	if stopEvents != nil {
		stopEvents()
	}
	if tray != nil {
		tray.Remove()
	}
}

func (t *trayManager) run() {
	// gogpu/systray creates a hidden Win32 window and then pumps that
	// window's message queue. Both operations must stay on the same OS
	// thread; otherwise the tray icon can be visible while left/right click
	// messages are never dispatched.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tray := systray.New()
	menu := systray.NewMenu()
	menu.Add("显示主窗口", t.ShowWindow)
	menu.Add("隐藏主窗口", t.hideWindow)
	menu.AddSeparator()

	launchEnabled := false
	launchSupported := false
	if t.adapter != nil {
		if preferences, err := t.adapter.GetDesktopPreferences(); err == nil {
			launchSupported = preferences.LaunchAtLoginSupported
			launchEnabled = preferences.LaunchAtLogin
		}
	}
	var launchItem *systray.MenuItem
	launchItem = menu.AddCheckbox("开机启动", launchEnabled, func() {
		t.toggleLaunchAtLogin(launchItem)
	})
	if !launchSupported {
		launchItem.SetDisabled(true)
	}
	menu.AddSeparator()
	menu.Add("退出", t.quit)

	tray.SetIcon(t.icon).
		SetTooltip("adm-desktop").
		SetMenu(menu)
	tray.OnClick(t.ShowWindow)
	tray.OnDoubleClick(t.ShowWindow)
	tray.OnRightClick(t.syncLaunchAtLogin)
	tray.Show()

	t.mu.Lock()
	if t.stopping {
		t.started = false
		t.mu.Unlock()
		tray.Remove()
		return
	}
	t.tray = tray
	t.launchItem = launchItem
	ctx := t.ctx
	t.mu.Unlock()

	if ctx != nil {
		stopEvents := wailsruntime.EventsOn(ctx, "desktop:preferences-changed", func(...interface{}) {
			t.syncLaunchAtLogin()
		})
		t.mu.Lock()
		if t.stopping {
			t.mu.Unlock()
			stopEvents()
			tray.Remove()
			return
		}
		t.stopEvents = stopEvents
		t.mu.Unlock()
	}

	_ = tray.Run()

	t.mu.Lock()
	stopEvents := t.stopEvents
	t.stopEvents = nil
	t.tray = nil
	t.launchItem = nil
	t.started = false
	t.mu.Unlock()
	if stopEvents != nil {
		stopEvents()
	}
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
	if t == nil {
		return nil
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.ctx
}

func (t *trayManager) toggleLaunchAtLogin(item *systray.MenuItem) {
	if t.adapter == nil || item == nil {
		return
	}
	preferences, err := t.adapter.GetDesktopPreferences()
	if err != nil || !preferences.LaunchAtLoginSupported {
		item.SetDisabled(true)
		return
	}
	preferences, err = t.adapter.SetLaunchAtLogin(!preferences.LaunchAtLogin)
	if err != nil {
		t.ShowWindow()
		if ctx := t.runtimeContext(); ctx != nil {
			_, _ = wailsruntime.MessageDialog(ctx, wailsruntime.MessageDialogOptions{
				Type: wailsruntime.ErrorDialog, Title: "开机启动设置失败", Message: err.Error(),
			})
		}
		return
	}
	item.SetChecked(preferences.LaunchAtLogin)
	if ctx := t.runtimeContext(); ctx != nil {
		wailsruntime.EventsEmit(ctx, "desktop:preferences-changed")
	}
}

func (t *trayManager) syncLaunchAtLogin() {
	if t.adapter == nil {
		return
	}
	preferences, err := t.adapter.GetDesktopPreferences()
	if err != nil {
		return
	}
	t.mu.RLock()
	item := t.launchItem
	t.mu.RUnlock()
	if item == nil {
		return
	}
	item.SetDisabled(!preferences.LaunchAtLoginSupported)
	item.SetChecked(preferences.LaunchAtLogin)
}

package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/desktop"
)

func TestDesktopClientAdapterStartsDisconnected(t *testing.T) {
	adapter := desktop.NewClientAdapter()
	if _, err := adapter.GetSnapshot(); err == nil || !strings.Contains(err.Error(), "not connected") {
		t.Fatalf("disconnected Desktop snapshot error = %v", err)
	}
}

func TestProductionDesktopUsesDisconnectedAdminMCPClient(t *testing.T) {
	source, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	if !strings.Contains(text, "desktop.NewClientAdapter()") {
		t.Fatal("production Desktop must construct the Admin MCP client adapter")
	}
	if strings.Contains(text, "newDesktopAdapter(statePath)") || strings.Contains(text, "management.New(application)") {
		t.Fatal("production Desktop must not silently construct a direct state-file management adapter")
	}
}

func TestProductionDesktopUsesTrayLifecycleAndSingleInstance(t *testing.T) {
	source, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	for _, required := range []string{"newTrayManager", "StartHidden", "HideWindowOnClose", "SingleInstanceLock", "OnSecondInstanceLaunch", "--autostart"} {
		if !strings.Contains(text, required) {
			t.Fatalf("production Desktop missing lifecycle marker %q", required)
		}
	}
	traySource, err := os.ReadFile("tray_manager_windows.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"RunWithExternalLoop", "显示主窗口", "隐藏主窗口", "开机启动", "退出", "SetOnTapped"} {
		if !strings.Contains(string(traySource), required) {
			t.Fatalf("Windows tray manager missing %q", required)
		}
	}
}

func TestGatewayChildRejectsPositionalArgumentsBeforeStartingServer(t *testing.T) {
	if err := runGatewayChild([]string{"unexpected"}); err == nil || !strings.Contains(err.Error(), "only --listen") {
		t.Fatalf("unexpected gateway child error: %v", err)
	}
}

func TestDesktopLaunchOptionsRecogniseAutostart(t *testing.T) {
	hidden, err := desktopLaunchOptions([]string{"--autostart"})
	if err != nil {
		t.Fatal(err)
	}
	if !hidden {
		t.Fatal("--autostart must start the Desktop hidden")
	}
	if _, err := desktopLaunchOptions([]string{"unexpected"}); err == nil || !strings.Contains(err.Error(), "only --autostart") {
		t.Fatalf("unexpected Desktop launch option error: %v", err)
	}
}

func TestWailsProjectConfigLivesWithDesktopCommand(t *testing.T) {
	config, err := os.ReadFile("wails.json")
	if err != nil {
		t.Fatalf("read desktop wails.json: %v", err)
	}
	text := string(config)
	for _, required := range []string{`"name": "AI Dev Manager V2"`, `"frontend:dir": "frontend"`, `"wailsjsdir": "frontend/wailsjs"`} {
		if !strings.Contains(text, required) {
			t.Fatalf("desktop wails.json missing %s", required)
		}
	}
	if _, err := os.Stat(filepath.Join("..", "..", "wails.json")); !os.IsNotExist(err) {
		t.Fatalf("root wails.json must not exist; Wails must run from the desktop command directory: %v", err)
	}
}

func TestDesktopIconAssetsAreWired(t *testing.T) {
	for _, source := range []string{
		filepath.Join("..", "..", "assets", "icons", "ai-dev-manager-app.png"),
		filepath.Join("..", "..", "assets", "icons", "ai-dev-manager-window.png"),
		filepath.Join("..", "..", "assets", "icons", "ai-dev-manager-tray.png"),
	} {
		if info, err := os.Stat(source); err != nil || info.Size() == 0 {
			t.Fatalf("missing Desktop icon source %s: %v", source, err)
		}
	}

	tray, err := os.ReadFile(filepath.Join("assets", "tray.ico"))
	if err != nil {
		t.Fatal(err)
	}
	if len(tray) < 4 || !bytes.Equal(tray[:4], []byte{0, 0, 1, 0}) {
		t.Fatal("embedded tray icon is not a Windows ICO")
	}

	assets, err := frontendAssets()
	if err != nil {
		t.Fatal(err)
	}
	brand, err := fs.ReadFile(assets, "assets/ai-dev-manager-window.png")
	if err != nil {
		t.Fatal(err)
	}
	sourceBrand, err := os.ReadFile(filepath.Join("..", "..", "assets", "icons", "ai-dev-manager-window.png"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(brand, sourceBrand) {
		t.Fatal("embedded Desktop brand mark must match assets/icons/ai-dev-manager-window.png")
	}

	buildScript, err := os.ReadFile(filepath.Join("..", "..", "scripts", "build-desktop.ps1"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"ai-dev-manager-app.png", "appicon.png", "windows\\icon.ico", "Remove-Item"} {
		if !strings.Contains(string(buildScript), required) {
			t.Fatalf("Desktop build script missing icon preparation marker %q", required)
		}
	}
}

func TestEmbeddedFrontendUsesDesktopManagementAndGatewayBindings(t *testing.T) {
	assets, err := frontendAssets()
	if err != nil {
		t.Fatal(err)
	}
	index, err := fs.ReadFile(assets, "index.html")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"workspaceCount", "environmentCount", "execCount", "mcpCount", "skillCount", "memoryCount", "refreshButton", "launchAtLogin", "desktopShellHint", "app-brand-mark", "ai-dev-manager-window.png",
		"gatewayState", "gatewayBaseURL", "gatewayHealthURL", "gatewayURL", "gatewayAdminURL", "gatewayRefreshButton", "gatewayStartButton", "gatewayStopButton",
		"workspaceForm", "environmentForm", "environmentDetailPanel", "aria-modal",
		"execForm", "managementEnvironment", "mcpEditorFlow", "mcpEditorSummary", "mcpEditorHint", "mcpForm", "mcpTransport", "mcpEndpoint", "mcpExecutable", "mcpReconnectInterval", "mcpEditCancelButton", "mcpSubmitButton", "mcpImportForm", "mcpImportApplyButton", "generic-mcpservers", "mcpFilter", "mcpStateFilter", "mcpVisibleCount", "status-legend", "skillSourceForm", "skillSourceRoot", "skillSupportRoots", "skillSourceList", "skillList", "skillFilter", "skillStateFilter", "skillVisibleCount", "loadGlobalMemory", "globalMemoryForm",
		"environmentMCPSelections", "environmentSkillSelections", "loadEnvironmentMemory", "environmentMemoryForm",
		"runtimeRefreshButton", "runtimeHint", "verifierList", "processList", "runList", "runtimeOutput",
	} {
		if !strings.Contains(string(index), required) {
			t.Fatalf("desktop index missing %q", required)
		}
	}
	javascript, err := fs.ReadFile(assets, "app.js")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"window.go?.desktop?.Adapter", "GetSnapshot", "refreshSnapshot", "global_memory_count",
		"GetDesktopPreferences", "SetLaunchAtLogin", "loadDesktopPreferences", "updateLaunchAtLogin",
		"ConnectADM", "StartLocalADM", "StopLocalADM", "refreshConnectedADM", "initializeConnectionProfiles",
		"AddWorkspace", "RenameWorkspace", "RemoveWorkspace",
		"CreateEnvironment", "RenameEnvironment", "RemoveEnvironment", "InspectEnvironment",
		"AllowExecutable", "RemoveExecutable",
		"AddMCP", "UpdateMCP", "PreviewMCPImport", "ApplyMCPImport", "ProbeMCPHealth", "SetMCPDefault", "RemoveMCP", "beginMCPEdit", "resetMCPEditor", "mcpReconnectInterval", "mcpReferenceVariableNames", "配置引用：", "配置 ·", "运行 ·", "尚未探测", "可用性 ·", "mcpStateFilter", "mcpVisibleCount", "reference-text", "reference_name", "没有符合当前筛选条件的 MCP",
		"AddSkillSource", "ListSkillSources", "RefreshSkillSource", "RemoveSkillSource", "ListEnvironmentSkills", "SetSkillDefault", "RemoveSkill", "skillStateFilter", "skillVisibleCount", "没有符合当前筛选条件的 Skill", "endpoint", "artifact_path", "source_root", "unconfigured",
		"SetEnvironmentMCP", "SetEnvironmentSkill",
		"ListGlobalMemory", "WriteGlobalMemory", "DeleteGlobalMemory",
		"ListEnvironmentMemory", "WriteEnvironmentMemory", "DeleteEnvironmentMemory",
		"ListVerifiers", "RunVerifier", "ListProcesses", "GetProcessLogs", "StopProcess", "ListRuns", "CancelRun", "refreshRuntimeContext",
		"正在显式读取 Global Memory", "正在显式读取 Environment-private Memory",
		"environmentDetailBackdrop", "statusPanel.hidden = false", "statusPanel.hidden = true",
		"只移除 ADM Workspace 记录，不删除目录", "只移除 ADM Environment 记录，不删除 root 或项目文件",
		"这是全局删除，不是只从当前 Environment 禁用", "删除全局 Skill source", "预览不会修改 catalog 或 Environment",
	} {
		if !strings.Contains(string(javascript), required) {
			t.Fatalf("desktop app.js missing %q", required)
		}
	}
	if strings.Contains(string(javascript), "stateBadge(config") {
		t.Fatal("desktop app.js still contains the undefined MCP config badge variable")
	}
	connections, err := fs.ReadFile(assets, "connections.js")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"GetConnectionProfiles", "SaveConnectionProfile", "SelectConnectionProfile", "DeleteConnectionProfile", "DisconnectADM", "withConnectionTransition", "desktopPendingRequests", "desktopRequestQueue"} {
		if !strings.Contains(string(connections), required) {
			t.Fatalf("desktop connections.js missing %q", required)
		}
	}
	dialogs, err := fs.ReadFile(assets, "dialogs.js")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"showModal", "preventScroll", "data-dialog-open", "data-dialog-close", "activeEditorDialog"} {
		if !strings.Contains(string(dialogs), required) {
			t.Fatalf("desktop dialogs.js missing %q", required)
		}
	}
	styles, err := fs.ReadFile(assets, "styles.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{".topbar-brand", ".app-brand-mark", ".topbar-actions", ".desktop-toggle", ".list-toolbar", ".status-legend", ".editor-hint", ".filtered-resource-list", ".resource-row[data-editing", ".resource-actions .check-field", ".editor-dialog", ".dialog-message"} {
		if !strings.Contains(string(styles), required) {
			t.Fatalf("desktop styles.css missing %q", required)
		}
	}
	if strings.Contains(string(styles), "鈥?") {
		t.Fatal("desktop styles.css contains the previously broken manager-flow marker encoding")
	}
}

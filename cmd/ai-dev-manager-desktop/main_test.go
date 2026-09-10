package main

import (
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

func TestGatewayChildRejectsPositionalArgumentsBeforeStartingServer(t *testing.T) {
	if err := runGatewayChild([]string{"unexpected"}); err == nil || !strings.Contains(err.Error(), "only --listen") {
		t.Fatalf("unexpected gateway child error: %v", err)
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
		"workspaceCount", "environmentCount", "execCount", "mcpCount", "skillCount", "memoryCount", "refreshButton",
		"gatewayState", "gatewayBaseURL", "gatewayHealthURL", "gatewayURL", "gatewayAdminURL", "gatewayRefreshButton", "gatewayStartButton", "gatewayStopButton",
		"workspaceForm", "environmentForm", "environmentDetailPanel", "aria-modal",
		"execForm", "managementEnvironment", "mcpForm", "mcpTransport", "mcpEndpoint", "mcpExecutable", "mcpImportForm", "mcpImportApplyButton", "generic-mcpservers", "mcpFilter", "mcpStateFilter", "mcpVisibleCount", "skillSourceForm", "skillSourceRoot", "skillSupportRoots", "skillSourceList", "skillList", "skillFilter", "skillStateFilter", "skillVisibleCount", "loadGlobalMemory", "globalMemoryForm",
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
		"ConnectADM", "StartLocalADM", "StopLocalADM", "refreshConnectedADM", "adm-v2.desktop.base-url",
		"AddWorkspace", "RenameWorkspace", "RemoveWorkspace",
		"CreateEnvironment", "RenameEnvironment", "RemoveEnvironment", "InspectEnvironment",
		"AllowExecutable", "RemoveExecutable",
		"AddMCP", "PreviewMCPImport", "ApplyMCPImport", "ProbeMCPHealth", "SetMCPDefault", "RemoveMCP", "mcpStateFilter", "mcpVisibleCount", "reference-text", "reference_name", "没有符合当前筛选条件的 MCP",
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
	styles, err := fs.ReadFile(assets, "styles.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{".manager-flow > summary::before", "content: '\\203A'", ".list-toolbar", ".filtered-resource-list", ".resource-actions .check-field"} {
		if !strings.Contains(string(styles), required) {
			t.Fatalf("desktop styles.css missing %q", required)
		}
	}
	if strings.Contains(string(styles), "鈥?") {
		t.Fatal("desktop styles.css contains the previously broken manager-flow marker encoding")
	}
}

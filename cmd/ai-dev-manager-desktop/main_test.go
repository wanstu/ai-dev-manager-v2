package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/desktop"
)

func TestNewDesktopAdapterUsesProvidedADMState(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	adapter := newDesktopAdapter(statePath)
	root := t.TempDir()

	ws, err := adapter.AddWorkspace(desktop.WorkspaceInput{Path: root, Name: "desktop-test"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := adapter.CreateEnvironment(desktop.EnvironmentInput{WorkspaceID: ws.ID, Name: "main"}); err != nil {
		t.Fatal(err)
	}
	snapshot, err := adapter.GetSnapshot()
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Workspaces) != 1 || len(snapshot.Environments) != 1 {
		t.Fatalf("desktop snapshot = %+v", snapshot)
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
		"gatewayState", "gatewayRefreshButton", "gatewayStartButton", "gatewayStopButton",
		"workspaceForm", "environmentForm", "environmentDetailPanel", "aria-modal",
		"execForm", "managementEnvironment", "mcpForm", "mcpTransport", "mcpEndpoint", "mcpExecutable", "mcpImportForm", "mcpImportApplyButton", "skillSourceForm", "skillSourceRoot", "skillSupportRoots", "skillSourceList", "skillList", "loadGlobalMemory", "globalMemoryForm",
		"environmentMCPSelections", "environmentSkillSelections", "loadEnvironmentMemory", "environmentMemoryForm",
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
		"GetGatewayStatus", "StartGateway", "StopGateway", "refreshGatewayStatus",
		"AddWorkspace", "RenameWorkspace", "RemoveWorkspace",
		"CreateEnvironment", "RenameEnvironment", "RemoveEnvironment", "InspectEnvironment",
		"AllowExecutable", "RemoveExecutable",
		"AddMCP", "PreviewMCPImport", "ApplyMCPImport", "ProbeMCPHealth", "SetMCPDefault", "RemoveMCP",
		"AddSkillSource", "ListSkillSources", "RefreshSkillSource", "RemoveSkillSource", "ListEnvironmentSkills", "SetSkillDefault", "RemoveSkill", "endpoint", "artifact_path", "source_root", "unconfigured",
		"SetEnvironmentMCP", "SetEnvironmentSkill",
		"ListGlobalMemory", "WriteGlobalMemory", "DeleteGlobalMemory",
		"ListEnvironmentMemory", "WriteEnvironmentMemory", "DeleteEnvironmentMemory",
		"正在显式读取 Global Memory", "正在显式读取 Environment-private Memory",
		"environmentDetailBackdrop", "statusPanel.hidden = false", "statusPanel.hidden = true",
		"只移除 ADM Workspace 记录，不删除目录", "只移除 ADM Environment 记录，不删除 root 或项目文件",
		"这是全局删除，不是只从当前 Environment 禁用", "删除全局 Skill source", "预览不会修改 catalog 或 Environment",
	} {
		if !strings.Contains(string(javascript), required) {
			t.Fatalf("desktop app.js missing %q", required)
		}
	}
}

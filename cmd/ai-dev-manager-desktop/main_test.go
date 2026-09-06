package main

import (
	"io/fs"
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

func TestEmbeddedFrontendUsesWorkspaceAndEnvironmentManagementBindings(t *testing.T) {
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
		"workspaceForm", "environmentForm", "environmentDetailPanel",
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
		"AddWorkspace", "RenameWorkspace", "RemoveWorkspace",
		"CreateEnvironment", "RenameEnvironment", "RemoveEnvironment", "InspectEnvironment",
		"只移除 ADM Workspace 记录，不删除目录", "只移除 ADM Environment 记录，不删除 root 或项目文件",
	} {
		if !strings.Contains(string(javascript), required) {
			t.Fatalf("desktop app.js missing %q", required)
		}
	}
}

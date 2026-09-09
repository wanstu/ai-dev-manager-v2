package management_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/management"
)

func TestSnapshotAggregatesPersistedStateWithoutMemoryValuesOrGit(t *testing.T) {
	root := t.TempDir()
	if _, err := os.Stat(filepath.Join(root, ".git")); !os.IsNotExist(err) {
		t.Fatalf("fixture must not be a Git repository")
	}
	application := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := application.Workspaces.Add(root, "projects")
	if err != nil {
		t.Fatal(err)
	}
	mcpEntry, err := application.MCPs.Add("filesystem", true)
	if err != nil {
		t.Fatal(err)
	}
	skillEntry, err := application.Skills.Add("go-project", false)
	if err != nil {
		t.Fatal(err)
	}
	env, err := application.Environments.Create(ws.ID, "main", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := application.SetEnvironmentSkill(env.ID, skillEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := application.Memory.GlobalWrite("global-secret", "do-not-expose-global"); err != nil {
		t.Fatal(err)
	}
	if err := application.Memory.EnvironmentWrite(env.ID, "private-secret", "do-not-expose-private"); err != nil {
		t.Fatal(err)
	}
	if err := application.AllowExecutable("go"); err != nil {
		t.Fatal(err)
	}

	service := management.New(application)
	snapshot, err := service.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Workspaces) != 1 || snapshot.Workspaces[0].ID != ws.ID {
		t.Fatalf("snapshot Workspaces = %+v", snapshot.Workspaces)
	}
	if len(snapshot.Environments) != 1 || snapshot.Environments[0].ID != env.ID {
		t.Fatalf("snapshot Environments = %+v", snapshot.Environments)
	}
	if snapshot.Environments[0].PrivateMemoryCount != 1 || snapshot.Environments[0].PrivateMemory != nil {
		t.Fatalf("snapshot Environment leaked private Memory state: %+v", snapshot.Environments[0])
	}
	if len(snapshot.AllowedExecutables) != 1 || snapshot.AllowedExecutables[0] != "go" {
		t.Fatalf("snapshot allowlist = %+v", snapshot.AllowedExecutables)
	}
	if len(snapshot.MCPs) != 1 || snapshot.MCPs[0].ID != mcpEntry.ID {
		t.Fatalf("snapshot MCPs = %+v", snapshot.MCPs)
	}
	if len(snapshot.Skills) != 1 || snapshot.Skills[0].ID != skillEntry.ID {
		t.Fatalf("snapshot Skills = %+v", snapshot.Skills)
	}
	if snapshot.GlobalMemoryCount != 1 {
		t.Fatalf("global_memory_count = %d, want 1", snapshot.GlobalMemoryCount)
	}
	encoded, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	for _, secret := range []string{"do-not-expose-global", "do-not-expose-private"} {
		if strings.Contains(text, secret) {
			t.Fatalf("management snapshot leaked Memory value %q: %s", secret, text)
		}
	}
}

func TestManagementMutationsDelegateToExistingServicesAndRemainSafe(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "keep.txt")
	if err := os.WriteFile(marker, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	application := app.New(filepath.Join(t.TempDir(), "state.json"))
	service := management.New(application)

	ws, err := service.WorkspaceAdd(root, "before")
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.EnvironmentCreate(ws.ID, "before", "")
	if err != nil {
		t.Fatal(err)
	}
	if env.Root != root || env.PrivateMemory != nil {
		t.Fatalf("management Environment create result = %+v", env)
	}
	if _, err := service.WorkspaceRemove(ws.ID); err == nil || !strings.Contains(err.Error(), env.ID) {
		t.Fatalf("Workspace removal must keep existing Environment guard, got %v", err)
	}

	mcpEntry, err := service.MCPAdd("filesystem", "http://127.0.0.1:9999/mcp", false)
	if err != nil {
		t.Fatal(err)
	}
	skillRoot := filepath.Join(t.TempDir(), "skills")
	if err := os.MkdirAll(filepath.Join(skillRoot, "go-project"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillRoot, "go-project", "SKILL.md"), []byte("# Go project\nUse Go tooling and run tests.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	skillEntries, err := service.SkillAdd(skillRoot, "", false)
	if err != nil || len(skillEntries) != 1 {
		t.Fatalf("SkillAdd entries=%+v err=%v", skillEntries, err)
	}
	skillEntry := skillEntries[0]
	if _, err := service.EnvironmentMCPSet(env.ID, mcpEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	if _, err := service.EnvironmentSkillSet(env.ID, skillEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := service.GlobalMemoryWrite("global-key", "global-secret-value"); err != nil {
		t.Fatal(err)
	}
	if err := service.EnvironmentMemoryWrite(env.ID, "private-key", "private-secret-value"); err != nil {
		t.Fatal(err)
	}
	allowed, err := service.ExecAllow("go")
	if err != nil || len(allowed) != 1 || allowed[0] != "go" {
		t.Fatalf("ExecAllow result = %v err=%v", allowed, err)
	}

	renamedWorkspace, err := service.WorkspaceRename(ws.ID, "after-workspace")
	if err != nil || renamedWorkspace.Name != "after-workspace" || renamedWorkspace.Path != root {
		t.Fatalf("WorkspaceRename result = %+v err=%v", renamedWorkspace, err)
	}
	renamedEnvironment, err := service.EnvironmentRename(env.ID, "after-environment")
	if err != nil {
		t.Fatal(err)
	}
	if renamedEnvironment.Name != "after-environment" || renamedEnvironment.Root != root || renamedEnvironment.PrivateMemory != nil || renamedEnvironment.PrivateMemoryCount != 1 {
		t.Fatalf("EnvironmentRename result = %+v", renamedEnvironment)
	}
	encoded, err := json.Marshal(renamedEnvironment)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "private-secret-value") {
		t.Fatalf("Environment mutation result leaked private Memory: %s", encoded)
	}

	if _, err := service.MCPSetDefault(mcpEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SkillSetDefault(skillEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	snapshot, err := service.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.GlobalMemoryCount != 1 || len(snapshot.Environments) != 1 || snapshot.Environments[0].PrivateMemoryCount != 1 {
		t.Fatalf("snapshot after management mutations = %+v", snapshot)
	}
	if !snapshot.MCPs[0].DefaultIncludeInEnv || !snapshot.Skills[0].DefaultIncludeInEnv {
		t.Fatalf("catalog defaults did not flow through management boundary: MCP=%+v Skill=%+v", snapshot.MCPs, snapshot.Skills)
	}

	if _, err := service.EnvironmentMCPSet(env.ID, "mcp_missing", true); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("management MCP selection must preserve existing validation, got %v", err)
	}
	if _, err := service.EnvironmentRename(env.ID, "   "); err == nil || !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("management Environment rename must preserve existing validation, got %v", err)
	}
	if err := service.GlobalMemoryWrite("   ", "x"); err == nil || !strings.Contains(err.Error(), "memory key is required") {
		t.Fatalf("management Global Memory must preserve existing validation, got %v", err)
	}

	if err := service.EnvironmentMemoryDelete(env.ID, "private-key"); err != nil {
		t.Fatal(err)
	}
	if err := service.GlobalMemoryDelete("global-key"); err != nil {
		t.Fatal(err)
	}
	if err := service.MCPRemove(mcpEntry.ID); err != nil {
		t.Fatal(err)
	}
	if err := service.SkillRemove(skillEntry.ID); err != nil {
		t.Fatal(err)
	}
	allowed, err = service.ExecRemove("go")
	if err != nil || len(allowed) != 0 {
		t.Fatalf("ExecRemove result = %v err=%v", allowed, err)
	}
	if _, err := service.ExecRemove("go"); err == nil || !strings.Contains(err.Error(), "not allowlisted") {
		t.Fatalf("management ExecRemove must preserve existing missing-entry error, got %v", err)
	}

	if err := service.EnvironmentRemove(env.ID); err != nil {
		t.Fatal(err)
	}
	removedWorkspace, err := service.WorkspaceRemove(ws.ID)
	if err != nil {
		t.Fatal(err)
	}
	if removedWorkspace.ID != ws.ID {
		t.Fatalf("removed Workspace = %+v", removedWorkspace)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "keep\n" {
		t.Fatalf("management removal must not delete project files: data=%q err=%v", data, err)
	}
}

func TestManagementMCPAddConfigUsesTypedModelAndSecretBoundary(t *testing.T) {
	application := app.New(filepath.Join(t.TempDir(), "state.json"))
	service := management.New(application)

	headerEntry, err := service.MCPAddConfig("http-auth", catalog.MCPConfig{
		Transport:  catalog.MCPTransportStreamableHTTP,
		AuthMode:   catalog.MCPAuthHeaders,
		Endpoint:   "http://127.0.0.1:9000/mcp",
		HeaderRefs: map[string]string{"Authorization": "Bearer ${MCP_API_TOKEN}"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if headerEntry.Transport != catalog.MCPTransportStreamableHTTP || headerEntry.AuthMode != catalog.MCPAuthHeaders || headerEntry.HeaderRefs["Authorization"] != "Bearer ${MCP_API_TOKEN}" {
		t.Fatalf("typed HTTP MCP from management = %+v", headerEntry)
	}

	stdioEntry, err := service.MCPAddConfig("stdio", catalog.MCPConfig{
		Transport:  catalog.MCPTransportStdio,
		Executable: "test-mcp",
		Args:       []string{"serve"},
		EnvRefs:    map[string]string{"API_TOKEN": "${MCP_API_TOKEN}"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if stdioEntry.Transport != catalog.MCPTransportStdio || stdioEntry.Executable != "test-mcp" || stdioEntry.EnvRefs["API_TOKEN"] != "${MCP_API_TOKEN}" {
		t.Fatalf("typed stdio MCP from management = %+v", stdioEntry)
	}

	for _, tc := range []struct {
		name string
		cfg  catalog.MCPConfig
	}{
		{
			name: "literal-header",
			cfg: catalog.MCPConfig{
				Transport:  catalog.MCPTransportStreamableHTTP,
				AuthMode:   catalog.MCPAuthHeaders,
				Endpoint:   "http://127.0.0.1:9001/mcp",
				HeaderRefs: map[string]string{"Authorization": "Bearer plain-secret-token"},
			},
		},
		{
			name: "literal-env",
			cfg: catalog.MCPConfig{
				Transport:  catalog.MCPTransportStdio,
				Executable: "test-mcp",
				EnvRefs:    map[string]string{"API_TOKEN": "plain-secret-token"},
			},
		},
	} {
		_, err := service.MCPAddConfig(tc.name, tc.cfg)
		if err == nil {
			t.Fatalf("management MCPAddConfig(%s) must reject literal credential refs", tc.name)
		}
		if strings.Contains(err.Error(), "plain-secret-token") {
			t.Fatalf("management MCPAddConfig leaked literal secret: %v", err)
		}
	}

	mcps, err := application.MCPs.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(mcps) != 2 {
		t.Fatalf("rejected literal credential refs must not persist MCPs: %+v", mcps)
	}
}

func TestManagementMCPImportPreviewIsSanitizedAndReadOnly(t *testing.T) {
	application := app.New(filepath.Join(t.TempDir(), "state.json"))
	service := management.New(application)

	preview, err := service.MCPImportPreview(catalog.MCPImportFormatCodexPlugin, `{"mcpServers":{"remote":{"url":"http://127.0.0.1:9000/mcp","headers":{"Authorization":"Bearer plain-secret-token"}}}}`, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.Candidates) != 1 {
		t.Fatalf("preview candidates=%+v", preview.Candidates)
	}
	candidate := preview.Candidates[0]
	if candidate.Definition.HeaderRefs["Authorization"] != "Bearer ${REMOTE_AUTHORIZATION}" {
		t.Fatalf("literal credential was not converted to a bearer reference template: %+v", candidate)
	}
	if len(candidate.ReferenceRequirements) != 1 || candidate.ReferenceRequirements[0].Name != "REMOTE_AUTHORIZATION" {
		t.Fatalf("reference requirement missing: %+v", candidate.ReferenceRequirements)
	}
	encoded, err := json.Marshal(preview)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "plain-secret-token") {
		t.Fatalf("management preview leaked literal secret: %s", encoded)
	}
	mcps, err := application.MCPs.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(mcps) != 0 {
		t.Fatalf("preview must not persist MCP definitions: %+v", mcps)
	}
}

func TestSnapshotReflectsSubsequentPersistedChangesAndUsesArrays(t *testing.T) {
	application := app.New(filepath.Join(t.TempDir(), "state.json"))
	service := management.New(application)

	empty, err := service.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if empty.Workspaces == nil || empty.Environments == nil || empty.AllowedExecutables == nil || empty.MCPs == nil || empty.Skills == nil {
		t.Fatalf("empty snapshot must use arrays, not nil slices: %+v", empty)
	}

	root := t.TempDir()
	ws, err := application.Workspaces.Add(root, "before")
	if err != nil {
		t.Fatal(err)
	}
	env, err := application.Environments.Create(ws.ID, "before", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := application.Workspaces.Rename(ws.ID, "after-workspace"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.Environments.Rename(env.ID, "after-environment"); err != nil {
		t.Fatal(err)
	}
	if err := application.Memory.GlobalWrite("one", "value"); err != nil {
		t.Fatal(err)
	}

	updated, err := service.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.Workspaces) != 1 || updated.Workspaces[0].Name != "after-workspace" {
		t.Fatalf("snapshot did not reflect Workspace rename: %+v", updated.Workspaces)
	}
	if len(updated.Environments) != 1 || updated.Environments[0].Name != "after-environment" {
		t.Fatalf("snapshot did not reflect Environment rename: %+v", updated.Environments)
	}
	if updated.GlobalMemoryCount != 1 {
		t.Fatalf("snapshot did not reflect Global Memory change: %+v", updated)
	}
}

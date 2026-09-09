package gateway

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/verifier"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGatewayDevelopsPlainDirectoryWithoutGit(t *testing.T) {
	root := t.TempDir()
	if _, err := os.Stat(filepath.Join(root, ".git")); !os.IsNotExist(err) {
		t.Fatalf("fixture must not be a git repository")
	}
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	names := toolNames(tools.Tools)
	for _, required := range []string{"workspace_list", "workspace_add", "workspace_inspect", "workspace_rename", "workspace_remove", "exec_allow", "exec_allow_remove", "environment_create", "environment_rename", "environment_remove", "environment_writer_acquire", "environment_writer_heartbeat", "environment_verifier_list", "environment_verifier_run", "mcp_list", "mcp_add", "environment_mcp_set", "environment_mcp_tools", "environment_mcp_call", "skill_list", "skill_add", "environment_skill_set", "environment_skill_list", "environment_skill_read", "memory_global_write", "memory_environment_write", "tree", "read", "search", "write", "edit", "delete", "exec", "git_status"} {
		if !contains(names, required) {
			t.Fatalf("missing gateway tool %q in %v", required, names)
		}
	}
	if contains(names, "run_workflow_start") {
		t.Fatalf("mis-scoped workflow orchestration tool is still exposed: %v", names)
	}
	retiredWorkflow, retiredErr := session.CallTool(ctx, &mcp.CallToolParams{Name: "run_workflow_start", Arguments: map[string]any{}})
	if retiredErr == nil && (retiredWorkflow == nil || !retiredWorkflow.IsError) {
		t.Fatalf("retired workflow orchestration tool unexpectedly callable: result=%+v", retiredWorkflow)
	}

	created, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_create",
		Arguments: map[string]any{"workspace_id": ws.ID, "name": "mcp-task"},
	})
	if err != nil || created.IsError {
		t.Fatalf("environment_create failed: err=%v result=%+v", err, created)
	}
	envs, err := service.Environments.List()
	if err != nil || len(envs) != 1 {
		t.Fatalf("environment persistence after MCP create: envs=%+v err=%v", envs, err)
	}
	envID := envs[0].ID

	acquired, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_writer_acquire",
		Arguments: map[string]any{"environment_id": envID, "owner": "mcp-session"},
	})
	if err != nil || acquired.IsError {
		t.Fatalf("writer acquire failed: err=%v result=%+v", err, acquired)
	}
	if !strings.Contains(toolText(t, acquired), "expires_at") {
		t.Fatalf("writer acquire result must expose lease expiry: %s", toolText(t, acquired))
	}

	heartbeat, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_writer_heartbeat",
		Arguments: map[string]any{"environment_id": envID, "owner": "mcp-session"},
	})
	if err != nil || heartbeat.IsError {
		t.Fatalf("writer heartbeat failed: err=%v result=%+v", err, heartbeat)
	}

	written, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "write",
		Arguments: map[string]any{
			"environment_id": envID,
			"writer_owner":   "mcp-session",
			"path":           "plain.txt",
			"content":        "works without git\n",
		},
	})
	if err != nil || written.IsError {
		t.Fatalf("write failed: err=%v result=%+v", err, written)
	}

	read, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "read",
		Arguments: map[string]any{"environment_id": envID, "path": "plain.txt"},
	})
	if err != nil || read.IsError {
		t.Fatalf("read failed: err=%v result=%+v", err, read)
	}
	if !strings.Contains(toolText(t, read), "works without git") {
		t.Fatalf("unexpected read result: %s", toolText(t, read))
	}

	gitResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "git_status",
		Arguments: map[string]any{"environment_id": envID},
	})
	if err != nil {
		t.Fatalf("git_status transport error = %v", err)
	}
	if !gitResult.IsError {
		t.Fatalf("git_status should be a local tool error for non-git root: %+v", gitResult)
	}

	data, err := os.ReadFile(filepath.Join(root, "plain.txt"))
	if err != nil || string(data) != "works without git\n" {
		t.Fatalf("MCP write did not reach plain directory: data=%q err=%v", data, err)
	}

	released, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_writer_release",
		Arguments: map[string]any{"environment_id": envID, "owner": "mcp-session"},
	})
	if err != nil || released.IsError {
		t.Fatalf("writer release failed: err=%v result=%+v", err, released)
	}
	removed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_remove",
		Arguments: map[string]any{"environment_id": envID},
	})
	if err != nil || removed.IsError {
		t.Fatalf("environment_remove failed: err=%v result=%+v", err, removed)
	}
	if envs, err := service.Environments.List(); err != nil || len(envs) != 0 {
		t.Fatalf("environment_remove did not remove ADM record: envs=%+v err=%v", envs, err)
	}
	if data, err := os.ReadFile(filepath.Join(root, "plain.txt")); err != nil || string(data) != "works without git\n" {
		t.Fatalf("environment_remove must not delete project files: data=%q err=%v", data, err)
	}
}

func TestGatewayVerifierListAndRunAreEnvironmentScopedAndWriterGated(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "marker.txt"), []byte("gateway-still-works\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "verifier-gateway")
	if err != nil {
		t.Fatal(err)
	}
	envA, err := service.Environments.Create(ws.ID, "a", "")
	if err != nil {
		t.Fatal(err)
	}
	envB, err := service.Environments.Create(ws.ID, "b", "")
	if err != nil {
		t.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err := service.AllowExecutable(exe); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADM_V2_GATEWAY_VERIFIER_HELPER", "1")
	helperArgs := []string{"-test.run=^TestGatewayVerifierHelperProcess$"}
	passA, err := service.AddVerifier(envA.ID, model.VerifierDefinition{
		Name:           "pass-a",
		Kind:           verifier.KindTest,
		Enabled:        true,
		Executable:     exe,
		Args:           helperArgs,
		TimeoutSeconds: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	disabledA, err := service.AddVerifier(envA.ID, model.VerifierDefinition{
		Name:       "disabled-a",
		Kind:       verifier.KindCustom,
		Enabled:    false,
		Executable: exe,
		Args:       helperArgs,
	})
	if err != nil {
		t.Fatal(err)
	}
	brokenA, err := service.AddVerifier(envA.ID, model.VerifierDefinition{
		Name:       "broken-a",
		Kind:       verifier.KindCustom,
		Enabled:    true,
		Executable: "adm-v2-gateway-verifier-not-allowlisted",
	})
	if err != nil {
		t.Fatal(err)
	}
	passB, err := service.AddVerifier(envB.ID, model.VerifierDefinition{
		Name:       "pass-b",
		Kind:       verifier.KindTest,
		Enabled:    true,
		Executable: exe,
		Args:       helperArgs,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	listedA, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_verifier_list",
		Arguments: map[string]any{"environment_id": envA.ID},
	})
	if err != nil || listedA.IsError {
		t.Fatalf("environment_verifier_list A failed: err=%v result=%+v", err, listedA)
	}
	textA := toolText(t, listedA)
	if !strings.Contains(textA, passA.ID) || strings.Contains(textA, passB.ID) {
		t.Fatalf("Environment A verifier list is not scoped: %s", textA)
	}
	listedB, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_verifier_list",
		Arguments: map[string]any{"environment_id": envB.ID},
	})
	if err != nil || listedB.IsError {
		t.Fatalf("environment_verifier_list B failed: err=%v result=%+v", err, listedB)
	}
	textB := toolText(t, listedB)
	if !strings.Contains(textB, passB.ID) || strings.Contains(textB, passA.ID) {
		t.Fatalf("Environment B verifier list is not scoped: %s", textB)
	}

	withoutWriter, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "environment_verifier_run",
		Arguments: map[string]any{
			"environment_id": envA.ID,
			"writer_owner":   "missing-owner",
			"verifier_id":    passA.ID,
		},
	})
	if err != nil {
		t.Fatalf("environment_verifier_run without writer transport error: %v", err)
	}
	if !withoutWriter.IsError {
		t.Fatalf("environment_verifier_run without matching writer must fail locally: %+v", withoutWriter)
	}
	if _, err := service.Environments.AcquireWriter(envA.ID, "gateway-verifier-owner"); err != nil {
		t.Fatal(err)
	}

	passed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "environment_verifier_run",
		Arguments: map[string]any{
			"environment_id": envA.ID,
			"writer_owner":   "gateway-verifier-owner",
			"verifier_id":    passA.ID,
		},
	})
	if err != nil || passed.IsError {
		t.Fatalf("environment_verifier_run pass failed: err=%v result=%+v", err, passed)
	}
	passedText := toolText(t, passed)
	for _, required := range []string{passA.ID, "\"status\":\"passed\"", "GATEWAY_VERIFIER_OK"} {
		if !strings.Contains(strings.ReplaceAll(passedText, " ", ""), strings.ReplaceAll(required, " ", "")) {
			t.Fatalf("passing verifier result missing %q: %s", required, passedText)
		}
	}

	for name, verifierID := range map[string]string{
		"missing":  "vf_missing",
		"disabled": disabledA.ID,
		"broken":   brokenA.ID,
	} {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: "environment_verifier_run",
			Arguments: map[string]any{
				"environment_id": envA.ID,
				"writer_owner":   "gateway-verifier-owner",
				"verifier_id":    verifierID,
			},
		})
		if err != nil {
			t.Fatalf("%s verifier transport error: %v", name, err)
		}
		if !result.IsError {
			t.Fatalf("%s verifier must fail locally: %+v", name, result)
		}
	}

	gatewayInfo, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "gateway_info", Arguments: map[string]any{}})
	if err != nil || gatewayInfo.IsError {
		t.Fatalf("broken verifier must not break gateway_info: err=%v result=%+v", err, gatewayInfo)
	}
	readResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "read",
		Arguments: map[string]any{"environment_id": envA.ID, "path": "marker.txt"},
	})
	if err != nil || readResult.IsError || !strings.Contains(toolText(t, readResult), "gateway-still-works") {
		t.Fatalf("broken verifier must not break read: err=%v result=%+v", err, readResult)
	}
}

func TestGatewayVerifierHelperProcess(t *testing.T) {
	if os.Getenv("ADM_V2_GATEWAY_VERIFIER_HELPER") != "1" {
		return
	}
	_, _ = os.Stdout.WriteString("GATEWAY_VERIFIER_OK\n")
}

func TestGatewayEnvironmentRenamePreservesContextAndProjectData(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "keep.txt")
	if err := os.WriteFile(marker, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "projects")
	if err != nil {
		t.Fatal(err)
	}
	mcpEntry, err := service.MCPs.Add("filesystem", false)
	if err != nil {
		t.Fatal(err)
	}
	skillEntry, err := service.Skills.Add("go", false)
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "before", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentMCP(env.ID, mcpEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentSkill(env.ID, skillEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := service.Memory.EnvironmentWrite(env.ID, "note", "keep"); err != nil {
		t.Fatal(err)
	}
	before, err := service.Environments.AcquireWriter(env.ID, "rename-owner")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()
	renamed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_rename",
		Arguments: map[string]any{"environment_id": env.ID, "name": "after"},
	})
	if err != nil || renamed.IsError {
		t.Fatalf("environment_rename failed: err=%v result=%+v", err, renamed)
	}
	if !strings.Contains(toolText(t, renamed), "after") {
		t.Fatalf("environment_rename result missing new name: %s", toolText(t, renamed))
	}
	after, err := service.Environments.Get(env.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.ID != before.ID || after.WorkspaceID != before.WorkspaceID || after.Root != before.Root || !after.CreatedAt.Equal(before.CreatedAt) {
		t.Fatalf("environment_rename changed stable identity: before=%+v after=%+v", before, after)
	}
	if after.Writer == nil || after.Writer.Owner != "rename-owner" {
		t.Fatalf("environment_rename changed writer: %+v", after.Writer)
	}
	if len(after.EnabledMCPIDs) != 1 || after.EnabledMCPIDs[0] != mcpEntry.ID || len(after.EnabledSkillIDs) != 1 || after.EnabledSkillIDs[0] != skillEntry.ID {
		t.Fatalf("environment_rename changed selections: %+v", after)
	}
	memoryEntry, err := service.Memory.EnvironmentRead(env.ID, "note")
	if err != nil || memoryEntry.Value != "keep" {
		t.Fatalf("environment_rename changed private memory: entry=%+v err=%v", memoryEntry, err)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "keep\n" {
		t.Fatalf("environment_rename touched project data: data=%q err=%v", data, err)
	}

	blank, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_rename",
		Arguments: map[string]any{"environment_id": env.ID, "name": "   "},
	})
	if err != nil {
		t.Fatalf("blank environment_rename transport error: %v", err)
	}
	if !blank.IsError {
		t.Fatalf("blank environment name must be a tool error: %+v", blank)
	}
}

func TestGatewayEnvironmentManagementViewsDoNotLeakPrivateMemory(t *testing.T) {
	root := t.TempDir()
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "projects")
	if err != nil {
		t.Fatal(err)
	}
	mcpEntry, err := service.MCPs.Add("filesystem", false)
	if err != nil {
		t.Fatal(err)
	}
	removedMCP, err := service.MCPs.Add("removed-mcp", false)
	if err != nil {
		t.Fatal(err)
	}
	skillEntry, err := service.Skills.Add("go-project", false)
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "inspect", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{mcpEntry.ID, removedMCP.ID} {
		if _, err := service.SetEnvironmentMCP(env.ID, id, true); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := service.SetEnvironmentSkill(env.ID, skillEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := service.Memory.EnvironmentWrite(env.ID, "secret", "do-not-print"); err != nil {
		t.Fatal(err)
	}
	if err := service.MCPs.Remove(removedMCP.ID); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()
	listed, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "environment_list", Arguments: map[string]any{}})
	if err != nil || listed.IsError {
		t.Fatalf("environment_list failed: err=%v result=%+v", err, listed)
	}
	listText := toolText(t, listed)
	for _, required := range []string{env.ID, "private_memory_count"} {
		if !strings.Contains(listText, required) {
			t.Fatalf("environment_list missing %q: %s", required, listText)
		}
	}
	if strings.Contains(listText, "do-not-print") || strings.Contains(listText, "\"private_memory\"") {
		t.Fatalf("environment_list leaked private Memory values: %s", listText)
	}

	inspected, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_inspect",
		Arguments: map[string]any{"environment_id": env.ID},
	})
	if err != nil || inspected.IsError {
		t.Fatalf("environment_inspect failed: err=%v result=%+v", err, inspected)
	}
	inspectText := toolText(t, inspected)
	for _, required := range []string{ws.ID, "projects", "filesystem", "go-project", "unresolved_mcp_ids", removedMCP.ID, "private_memory_count"} {
		if !strings.Contains(inspectText, required) {
			t.Fatalf("environment_inspect missing %q: %s", required, inspectText)
		}
	}
	if strings.Contains(inspectText, "do-not-print") || strings.Contains(inspectText, "\"private_memory\"") {
		t.Fatalf("environment_inspect leaked private Memory values: %s", inspectText)
	}
}

func TestGatewayWorkspaceManagementGuardsProjectData(t *testing.T) {
	root := t.TempDir()
	sentinel := filepath.Join(root, "keep.txt")
	if err := os.WriteFile(sentinel, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "before")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	inspected, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_inspect",
		Arguments: map[string]any{"workspace_id": ws.ID},
	})
	if err != nil || inspected.IsError || !strings.Contains(toolText(t, inspected), ws.ID) {
		t.Fatalf("workspace_inspect failed: err=%v result=%+v", err, inspected)
	}
	renamed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_rename",
		Arguments: map[string]any{"workspace_id": ws.ID, "name": "after"},
	})
	if err != nil || renamed.IsError || !strings.Contains(toolText(t, renamed), "after") {
		t.Fatalf("workspace_rename failed: err=%v result=%+v", err, renamed)
	}

	env, err := service.Environments.Create(ws.ID, "guard", "")
	if err != nil {
		t.Fatal(err)
	}
	blocked, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_remove",
		Arguments: map[string]any{"workspace_id": ws.ID},
	})
	if err != nil {
		t.Fatalf("workspace_remove transport error: %v", err)
	}
	if !blocked.IsError {
		t.Fatalf("workspace_remove must be rejected while Environment references Workspace: %+v", blocked)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("blocked workspace_remove touched project data: %v", err)
	}

	if _, err := service.Environments.Remove(env.ID); err != nil {
		t.Fatal(err)
	}
	removed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_remove",
		Arguments: map[string]any{"workspace_id": ws.ID},
	})
	if err != nil || removed.IsError {
		t.Fatalf("workspace_remove failed after Environment removal: err=%v result=%+v", err, removed)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("workspace_remove must not delete project data: %v", err)
	}
}

func TestGatewayExecAllowlistRemoveRevokesEntry(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	if err := service.AllowExecutable("go"); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	removed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "exec_allow_remove",
		Arguments: map[string]any{"executable": "GO"},
	})
	if err != nil || removed.IsError {
		t.Fatalf("exec_allow_remove failed: err=%v result=%+v", err, removed)
	}
	allowed, err := service.AllowedExecutables()
	if err != nil {
		t.Fatal(err)
	}
	if len(allowed) != 0 {
		t.Fatalf("allowlist after gateway removal = %v; want empty", allowed)
	}

	missing, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "exec_allow_remove",
		Arguments: map[string]any{"executable": "go"},
	})
	if err != nil {
		t.Fatalf("second exec_allow_remove transport error: %v", err)
	}
	if !missing.IsError {
		t.Fatalf("removing a missing allowlist entry must be a tool error: %+v", missing)
	}
}

func TestGatewayMCPAddRejectsLiteralSecretRefs(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	cases := []struct {
		name   string
		args   map[string]any
		secret string
	}{
		{
			name: "literal header ref",
			args: map[string]any{
				"name":        "literal-header",
				"transport":   "streamable-http",
				"auth_mode":   "headers",
				"endpoint":    "http://127.0.0.1:9000/mcp",
				"header_refs": map[string]string{"Authorization": "Bearer plain-secret-token"},
			},
			secret: "plain-secret-token",
		},
		{
			name: "literal stdio env ref",
			args: map[string]any{
				"name":       "literal-env",
				"transport":  "stdio",
				"executable": "test-mcp",
				"env_refs":   map[string]string{"API_TOKEN": "plain-secret-token"},
			},
			secret: "plain-secret-token",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "mcp_add", Arguments: tc.args})
			if err != nil {
				t.Fatalf("mcp_add transport error: %v", err)
			}
			if !result.IsError {
				t.Fatalf("mcp_add must reject literal credential refs: %+v", result)
			}
			text := toolText(t, result)
			if strings.Contains(text, tc.secret) {
				t.Fatalf("mcp_add error leaked literal secret: %s", text)
			}
		})
	}

	mcps, err := service.MCPs.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(mcps) != 0 {
		t.Fatalf("rejected literal credential refs must not persist MCPs: %+v", mcps)
	}

	accepted, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "mcp_add",
		Arguments: map[string]any{
			"name":        "referenced-header",
			"transport":   "streamable-http",
			"auth_mode":   "headers",
			"endpoint":    "http://127.0.0.1:9000/mcp",
			"header_refs": map[string]string{"Authorization": "Bearer ${MCP_API_TOKEN}"},
		},
	})
	if err != nil || accepted.IsError {
		t.Fatalf("mcp_add must accept environment-backed header refs: err=%v result=%+v", err, accepted)
	}
}

func TestGatewayUsesOnlyEnabledConfiguredExternalMCPAndSkillRuntime(t *testing.T) {
	newExternal := func() *mcp.Server {
		external := mcp.NewServer(&mcp.Implementation{Name: "external-test", Version: "dev"}, nil)
		mcp.AddTool(external, &mcp.Tool{Name: "external_ping", Description: "Return a test pong."},
			func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
				return nil, map[string]any{"pong": "external-pong"}, nil
			})
		return external
	}
	externalHTTP := httptest.NewServer(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return newExternal()
	}, &mcp.StreamableHTTPOptions{Stateless: true, DisableLocalhostProtection: true}))
	defer externalHTTP.Close()

	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	root := t.TempDir()
	ws, err := service.Workspaces.Add(root, "external-context")
	if err != nil {
		t.Fatal(err)
	}
	mcpEntry, err := service.MCPs.AddMCP("external", externalHTTP.URL, false)
	if err != nil {
		t.Fatal(err)
	}

	skillRoot := filepath.Join(t.TempDir(), "skills")
	supportRoot := filepath.Join(t.TempDir(), "gsd-core")
	if err := os.MkdirAll(filepath.Join(skillRoot, "gsd-next"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(supportRoot, "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	supportFile := filepath.Join(supportRoot, "workflows", "smart-entry.md")
	if err := os.WriteFile(filepath.Join(skillRoot, "gsd-next", "SKILL.md"), []byte("# GSD smart entry\n@"+filepath.ToSlash(supportFile)+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(supportFile, []byte("# Smart entry workflow\nRead .planning/STATE.md first.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	discovered, err := service.Skills.AddSkillRoot(skillRoot, []string{supportRoot}, false)
	if err != nil || len(discovered) != 1 {
		t.Fatalf("discover Skill root: entries=%+v err=%v", discovered, err)
	}
	skillEntry := discovered[0]

	enabledEnv, err := service.Environments.Create(ws.ID, "enabled", "")
	if err != nil {
		t.Fatal(err)
	}
	disabledEnv, err := service.Environments.Create(ws.ID, "disabled", "")
	if err != nil {
		t.Fatal(err)
	}
	secondEnabledEnv, err := service.Environments.Create(ws.ID, "second-enabled", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, environmentID := range []string{enabledEnv.ID, secondEnabledEnv.ID} {
		if _, err := service.SetEnvironmentSkill(environmentID, skillEntry.ID, true); err != nil {
			t.Fatal(err)
		}
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	blocked, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_mcp_tools",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "mcp_id": mcpEntry.ID},
	})
	if err != nil {
		t.Fatalf("disabled external MCP transport error: %v", err)
	}
	if !blocked.IsError {
		t.Fatalf("external MCP must be blocked until enabled for the Environment: %+v", blocked)
	}

	if _, err := service.SetEnvironmentMCP(enabledEnv.ID, mcpEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	listed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_mcp_tools",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "mcp_id": mcpEntry.ID},
	})
	if err != nil || listed.IsError || !strings.Contains(toolText(t, listed), "external_ping") {
		t.Fatalf("enabled external MCP tools failed: err=%v result=%+v", err, listed)
	}
	called, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_mcp_call",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "mcp_id": mcpEntry.ID, "tool": "external_ping"},
	})
	if err != nil || called.IsError || !strings.Contains(toolText(t, called), "external-pong") {
		t.Fatalf("enabled external MCP call failed: err=%v result=%+v", err, called)
	}

	skillList, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_list",
		Arguments: map[string]any{"environment_id": enabledEnv.ID},
	})
	if err != nil || skillList.IsError || !strings.Contains(toolText(t, skillList), "gsd-next") {
		t.Fatalf("enabled Environment Skill list failed: err=%v result=%+v", err, skillList)
	}
	skillArtifact, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "skill_id": skillEntry.ID},
	})
	if err != nil || skillArtifact.IsError || !strings.Contains(toolText(t, skillArtifact), "GSD smart entry") {
		t.Fatalf("enabled Environment Skill artifact read failed: err=%v result=%+v", err, skillArtifact)
	}
	skillSupport, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "skill_id": skillEntry.ID, "path": supportFile},
	})
	if err != nil || skillSupport.IsError || !strings.Contains(toolText(t, skillSupport), "Smart entry workflow") {
		t.Fatalf("enabled Environment Skill support read failed: err=%v result=%+v", err, skillSupport)
	}
	disabledSkill, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": disabledEnv.ID, "skill_id": skillEntry.ID},
	})
	if err != nil {
		t.Fatalf("disabled Skill transport error: %v", err)
	}
	if !disabledSkill.IsError {
		t.Fatalf("disabled Environment must not read global Skill: %+v", disabledSkill)
	}
	secondRead, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": secondEnabledEnv.ID, "skill_id": skillEntry.ID},
	})
	if err != nil || secondRead.IsError || !strings.Contains(toolText(t, secondRead), "GSD smart entry") {
		t.Fatalf("second enabled Environment did not consume shared Skill: err=%v result=%+v", err, secondRead)
	}
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("not skill content"), 0o644); err != nil {
		t.Fatal(err)
	}
	escape, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "skill_id": skillEntry.ID, "path": outside},
	})
	if err != nil {
		t.Fatalf("Skill escape transport error: %v", err)
	}
	if !escape.IsError {
		t.Fatalf("Skill read outside configured roots must be rejected: %+v", escape)
	}

	if err := os.Remove(filepath.Join(skillRoot, "gsd-next", "SKILL.md")); err != nil {
		t.Fatal(err)
	}
	broken, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "skill_id": skillEntry.ID},
	})
	if err != nil {
		t.Fatalf("broken Skill transport error: %v", err)
	}
	if !broken.IsError {
		t.Fatalf("missing Skill artifact must be a local tool error: %+v", broken)
	}
	gatewayInfo, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "gateway_info", Arguments: map[string]any{}})
	if err != nil || gatewayInfo.IsError {
		t.Fatalf("broken Skill must not make unrelated Gateway tools unavailable: err=%v result=%+v", err, gatewayInfo)
	}

	if _, err := service.SetEnvironmentMCP(enabledEnv.ID, mcpEntry.ID, false); err != nil {
		t.Fatal(err)
	}
	blockedAgain, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_mcp_call",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "mcp_id": mcpEntry.ID, "tool": "external_ping"},
	})
	if err != nil {
		t.Fatalf("disabled external MCP call transport error: %v", err)
	}
	if !blockedAgain.IsError {
		t.Fatalf("external MCP call must be rejected immediately after Environment disable: %+v", blockedAgain)
	}
}

func TestRealHostGSDSkillThroughHTTPGateway(t *testing.T) {
	skillRoot := strings.TrimSpace(os.Getenv("ADM_REAL_GSD_SKILLS_ROOT"))
	supportRoot := strings.TrimSpace(os.Getenv("ADM_REAL_GSD_SUPPORT_ROOT"))
	if skillRoot == "" || supportRoot == "" {
		t.Skip("set ADM_REAL_GSD_SKILLS_ROOT and ADM_REAL_GSD_SUPPORT_ROOT for host acceptance")
	}

	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	workspaceRoot := t.TempDir()
	ws, err := service.Workspaces.Add(workspaceRoot, "real-gsd-acceptance")
	if err != nil {
		t.Fatal(err)
	}
	discovered, err := service.Skills.AddSkillRoot(skillRoot, []string{supportRoot}, false)
	if err != nil {
		t.Fatal(err)
	}
	var gsdNextID string
	for _, entry := range discovered {
		if entry.Name == "gsd-next" {
			gsdNextID = entry.ID
			break
		}
	}
	if gsdNextID == "" {
		t.Fatalf("real GSD root did not discover gsd-next; discovered %d Skills", len(discovered))
	}
	enabledEnv, err := service.Environments.Create(ws.ID, "enabled", "")
	if err != nil {
		t.Fatal(err)
	}
	disabledEnv, err := service.Environments.Create(ws.ID, "disabled", "")
	if err != nil {
		t.Fatal(err)
	}
	secondEnv, err := service.Environments.Create(ws.ID, "second", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, environmentID := range []string{enabledEnv.ID, secondEnv.ID} {
		if _, err := service.SetEnvironmentSkill(environmentID, gsdNextID, true); err != nil {
			t.Fatal(err)
		}
	}

	httpServer := httptest.NewServer(NewHTTPHandler(service))
	defer httpServer.Close()
	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "adm-v2-real-gsd-acceptance", Version: "dev"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: httpServer.URL + "/mcp"}, nil)
	if err != nil {
		t.Fatalf("connect real GSD HTTP Gateway: %v", err)
	}
	defer session.Close()

	listed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_list",
		Arguments: map[string]any{"environment_id": enabledEnv.ID},
	})
	if err != nil || listed.IsError || !strings.Contains(toolText(t, listed), "gsd-next") {
		t.Fatalf("real gsd-next list failed: err=%v result=%+v", err, listed)
	}
	artifact, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "skill_id": gsdNextID},
	})
	if err != nil || artifact.IsError {
		t.Fatalf("real gsd-next artifact read failed: err=%v result=%+v", err, artifact)
	}
	artifactText := toolText(t, artifact)
	for _, required := range []string{"name: gsd-next", "gsd-core/workflows/smart-entry.md"} {
		if !strings.Contains(strings.ReplaceAll(artifactText, "\\", "/"), required) {
			t.Fatalf("real gsd-next artifact missing %q: %s", required, artifactText)
		}
	}
	supportPath := filepath.Join(supportRoot, "workflows", "smart-entry.md")
	support, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": enabledEnv.ID, "skill_id": gsdNextID, "path": supportPath},
	})
	if err != nil || support.IsError || !strings.Contains(strings.ToLower(toolText(t, support)), "smart") {
		t.Fatalf("real gsd-next support read failed: err=%v result=%+v", err, support)
	}
	disabled, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": disabledEnv.ID, "skill_id": gsdNextID},
	})
	if err != nil {
		t.Fatalf("disabled real GSD Skill transport error: %v", err)
	}
	if !disabled.IsError {
		t.Fatalf("disabled Environment must not read real gsd-next: %+v", disabled)
	}
	second, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_skill_read",
		Arguments: map[string]any{"environment_id": secondEnv.ID, "skill_id": gsdNextID},
	})
	if err != nil || second.IsError || !strings.Contains(toolText(t, second), "name: gsd-next") {
		t.Fatalf("second Environment did not consume shared real gsd-next: err=%v result=%+v", err, second)
	}
}

func TestHTTPGatewayPersistsContextAcrossProcessRestart(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	projectRoot := t.TempDir()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	listen := listener.Addr().String()
	_ = listener.Close()
	endpoint := "http://" + listen + "/mcp"

	var running *exec.Cmd
	t.Cleanup(func() {
		if running != nil && running.Process != nil {
			_ = running.Process.Kill()
			_ = running.Wait()
		}
	})
	startGateway := func() *exec.Cmd {
		cmd := exec.Command(os.Args[0], "-test.run=^TestHTTPGatewayRestartHelperProcess$")
		cmd.Env = append(os.Environ(),
			"ADM_TEST_GATEWAY_RESTART_HELPER=1",
			"ADM_TEST_GATEWAY_RESTART_STATE="+statePath,
			"ADM_TEST_GATEWAY_RESTART_LISTEN="+listen,
		)
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		running = cmd
		return cmd
	}
	stopGateway := func(cmd *exec.Cmd) {
		if err := cmd.Process.Kill(); err != nil {
			t.Fatalf("kill temporary Gateway: %v", err)
		}
		if err := cmd.Wait(); err != nil {
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) {
				t.Fatalf("wait for temporary Gateway: %v", err)
			}
		}
		running = nil
	}

	first := startGateway()
	ctx := context.Background()
	firstSession := connectHTTPWithRetry(t, ctx, endpoint)
	added, err := firstSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_add",
		Arguments: map[string]any{"path": projectRoot, "name": "restart-persist"},
	})
	if err != nil || added.IsError {
		t.Fatalf("workspace_add before restart failed: err=%v result=%+v", err, added)
	}
	stateService := app.New(statePath)
	workspaces, err := stateService.Workspaces.List()
	if err != nil || len(workspaces) != 1 {
		t.Fatalf("Workspace state before restart = %+v err=%v", workspaces, err)
	}
	created, err := firstSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_create",
		Arguments: map[string]any{"workspace_id": workspaces[0].ID, "name": "restart-context"},
	})
	if err != nil || created.IsError {
		t.Fatalf("environment_create before restart failed: err=%v result=%+v", err, created)
	}
	environments, err := stateService.Environments.List()
	if err != nil || len(environments) != 1 {
		t.Fatalf("Environment state before restart = %+v err=%v", environments, err)
	}
	environmentID := environments[0].ID
	written, err := firstSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "memory_environment_write",
		Arguments: map[string]any{
			"environment_id": environmentID,
			"key":            "restart-check",
			"value":          "persisted",
		},
	})
	if err != nil || written.IsError {
		t.Fatalf("memory_environment_write before restart failed: err=%v result=%+v", err, written)
	}
	_ = firstSession.Close()
	stopGateway(first)

	second := startGateway()
	secondSession := connectHTTPWithRetry(t, ctx, endpoint)
	defer secondSession.Close()
	inspected, err := secondSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_inspect",
		Arguments: map[string]any{"environment_id": environmentID},
	})
	if err != nil || inspected.IsError || !strings.Contains(toolText(t, inspected), environmentID) {
		t.Fatalf("environment_inspect after restart failed: err=%v result=%+v", err, inspected)
	}
	memoryResult, err := secondSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "memory_environment_read",
		Arguments: map[string]any{
			"environment_id": environmentID,
			"key":            "restart-check",
		},
	})
	if err != nil || memoryResult.IsError || !strings.Contains(toolText(t, memoryResult), "persisted") {
		t.Fatalf("Environment-private Memory after restart failed: err=%v result=%+v", err, memoryResult)
	}
	_ = secondSession.Close()
	stopGateway(second)
}

func TestHTTPGatewayRestartHelperProcess(t *testing.T) {
	if os.Getenv("ADM_TEST_GATEWAY_RESTART_HELPER") != "1" {
		return
	}
	statePath := os.Getenv("ADM_TEST_GATEWAY_RESTART_STATE")
	listen := os.Getenv("ADM_TEST_GATEWAY_RESTART_LISTEN")
	if statePath == "" || listen == "" {
		t.Fatal("restart helper requires state path and listen address")
	}
	if err := RunHTTP(context.Background(), app.New(statePath), listen); err != nil {
		t.Fatal(err)
	}
}

func connectHTTPWithRetry(t *testing.T, ctx context.Context, endpoint string) *mcp.ClientSession {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		client := mcp.NewClient(&mcp.Implementation{Name: "adm-v2-restart-test", Version: "dev"}, nil)
		session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
		if err == nil {
			return session
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("connect temporary HTTP Gateway %s: %v", endpoint, lastErr)
	return nil
}

func TestHTTPGatewayUsesStreamableMCPAtMCPPath(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	httpServer := httptest.NewServer(NewHTTPHandler(service))
	defer httpServer.Close()

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "adm-v2-http-test", Version: "dev"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: httpServer.URL + "/mcp"}, nil)
	if err != nil {
		t.Fatalf("connect streamable HTTP gateway: %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(toolNames(tools.Tools), "gateway_info") || !contains(toolNames(tools.Tools), "workspace_add") {
		t.Fatalf("unexpected HTTP gateway tools: %v", toolNames(tools.Tools))
	}
	for _, tool := range tools.Tools {
		if (tool.Name == "gateway_info" || tool.Name == "workspace_add" || tool.Name == "environment_rename" || tool.Name == "environment_remove" || tool.Name == "environment_writer_acquire" || tool.Name == "environment_writer_heartbeat") && tool.OutputSchema != nil {
			t.Fatalf("generic tool %q must omit outputSchema for broad MCP client compatibility; got %#v", tool.Name, tool.OutputSchema)
		}
		if tool.Name == "environment_inspect" && tool.OutputSchema == nil {
			t.Fatal("environment_inspect must retain its structured outputSchema")
		}
	}
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "gateway_info", Arguments: map[string]any{}})
	if err != nil || result.IsError || !strings.Contains(toolText(t, result), "ai-dev-manager-v2") {
		t.Fatalf("gateway_info over HTTP failed: err=%v result=%+v", err, result)
	}
}

func connectInMemory(t *testing.T, ctx context.Context, server *mcp.Server) *mcp.ClientSession {
	t.Helper()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server.Connect: %v", err)
	}
	t.Cleanup(func() { _ = serverSession.Close() })
	client := mcp.NewClient(&mcp.Implementation{Name: "adm-v2-test", Version: "dev"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect: %v", err)
	}
	return clientSession
}

func toolNames(tools []*mcp.Tool) []string {
	out := make([]string, 0, len(tools))
	for _, tool := range tools {
		out = append(out, tool.Name)
	}
	sort.Strings(out)
	return out
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func toolText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) == 0 {
		t.Fatal("tool result has no content")
	}
	text, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("tool result content type = %T", result.Content[0])
	}
	return text.Text
}

package main

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/app"
)

func TestTopLevelHelpExplainsQuickStartAndGatewayLifecycle(t *testing.T) {
	output := captureStdout(t, func() {
		if err := run([]string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{
		"快速开始（HTTP Gateway）",
		"ai-dev-manager-v2 gateway start",
		"gateway status",
		"gateway stop",
		"gateway restart",
		"ai-dev-manager-v2 mcp -h",
		"ai-dev-manager-v2 skill -h",
		"ai-dev-manager-v2 memory -h",
		"doctor",
		"仅供 MCP 客户端使用",
	} {
		if !strings.Contains(output, required) {
			t.Fatalf("top-level help missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayHelpExplainsForegroundHTTPAndClientOnlyStdio(t *testing.T) {
	output := captureStdout(t, func() {
		if err := runGateway(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{
		"当前终端前台启动",
		"--detach",
		"Ctrl+C 停止",
		"运行状态",
		"人不要手动运行",
	} {
		if !strings.Contains(output, required) {
			t.Fatalf("gateway help missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayStartDetachRoutesToDetachedLauncher(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	original := startGatewayDetached
	t.Cleanup(func() { startGatewayDetached = original })

	for _, flagName := range []string{"--detach", "-d"} {
		var gotListen string
		startGatewayDetached = func(listen string) error {
			gotListen = listen
			return nil
		}
		if err := runGateway(service, []string{"start", flagName, "--listen", "127.0.0.1:45555"}); err != nil {
			t.Fatal(err)
		}
		if gotListen != "127.0.0.1:45555" {
			t.Fatalf("%s listen = %q", flagName, gotListen)
		}
	}
}

func TestWorkspaceHelpExplainsLifecycleAndDeletionSafety(t *testing.T) {
	output := captureStdout(t, func() {
		if err := runWorkspace(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{
		"workspace inspect --workspace-id",
		"workspace rename --workspace-id",
		"workspace remove --workspace-id",
		"不删除项目目录或文件",
		"Environment 引用",
	} {
		if !strings.Contains(output, required) {
			t.Fatalf("workspace help missing %q:\n%s", required, output)
		}
	}
}

func TestCLIRejectsOldPositionalWorkspaceAddGrammar(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	err := runWorkspace(service, []string{"add", t.TempDir()})
	if err == nil || !strings.Contains(err.Error(), "--path") {
		t.Fatalf("old positional grammar should explain --path requirement, got %v", err)
	}
}

func TestExecHelpIncludesAllowlistRemoval(t *testing.T) {
	output := captureStdout(t, func() {
		if err := runExec(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"exec allow --executable", "exec remove --executable", "exec list", "后续 exec 立即"} {
		if !strings.Contains(output, required) {
			t.Fatalf("exec help missing %q:\n%s", required, output)
		}
	}
}

func TestCatalogCLIManagesGlobalMCPAndSkill(t *testing.T) {
	home := t.TempDir()
	t.Setenv("ADM_V2_HOME", home)

	captureStdout(t, func() {
		if err := run([]string{"mcp", "add", "--name", "filesystem", "--endpoint", "http://127.0.0.1:9999/mcp", "--default"}); err != nil {
			t.Fatal(err)
		}
	})
	service := app.New(filepath.Join(home, "state.json"))
	mcps, err := service.MCPs.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(mcps) != 1 || mcps[0].Name != "filesystem" || mcps[0].Endpoint != "http://127.0.0.1:9999/mcp" || !mcps[0].DefaultIncludeInEnv {
		t.Fatalf("MCP catalog after CLI add = %+v", mcps)
	}

	captureStdout(t, func() {
		if err := run([]string{"mcp", "set-default", "--id", mcps[0].ID, "--enabled", "false"}); err != nil {
			t.Fatal(err)
		}
	})
	mcps, err = service.MCPs.List()
	if err != nil {
		t.Fatal(err)
	}
	if mcps[0].DefaultIncludeInEnv {
		t.Fatalf("MCP default was not updated: %+v", mcps[0])
	}

	skillRoot := filepath.Join(t.TempDir(), "skills")
	if err := os.MkdirAll(filepath.Join(skillRoot, "go-project"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillRoot, "go-project", "SKILL.md"), []byte("# Go project\nUse Go tooling and run tests.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	captureStdout(t, func() {
		if err := run([]string{"skill", "add", "--root", skillRoot}); err != nil {
			t.Fatal(err)
		}
	})
	skills, err := service.Skills.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "go-project" || skills[0].ArtifactPath != filepath.Join(skillRoot, "go-project", "SKILL.md") || skills[0].DefaultIncludeInEnv {
		t.Fatalf("Skill catalog after CLI add = %+v", skills)
	}

	captureStdout(t, func() {
		if err := run([]string{"mcp", "remove", "--id", mcps[0].ID}); err != nil {
			t.Fatal(err)
		}
		if err := run([]string{"skill", "remove", "--id", skills[0].ID}); err != nil {
			t.Fatal(err)
		}
	})
	if mcps, err := service.MCPs.List(); err != nil || len(mcps) != 0 {
		t.Fatalf("MCP catalog after CLI remove = %+v err=%v", mcps, err)
	}
	if skills, err := service.Skills.List(); err != nil || len(skills) != 0 {
		t.Fatalf("Skill catalog after CLI remove = %+v err=%v", skills, err)
	}
}

func TestCatalogCLIRejectsInvalidDefaultValue(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	err := runCatalog("mcp", service.MCPs, []string{"set-default", "--id", "mcp_x", "--enabled", "maybe"})
	if err == nil || !strings.Contains(err.Error(), "true 或 false") {
		t.Fatalf("invalid catalog default should fail clearly, got %v", err)
	}
}

func TestCatalogHelpIsDiscoverable(t *testing.T) {
	mcpOutput := captureStdout(t, func() {
		if err := runCatalog("mcp", nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"mcp add --name NAME --endpoint URL", "mcp list", "mcp remove --id ID", "mcp set-default --id ID --enabled true|false"} {
		if !strings.Contains(mcpOutput, required) {
			t.Fatalf("mcp help missing %q:\n%s", required, mcpOutput)
		}
	}

	skillOutput := captureStdout(t, func() {
		if err := runCatalog("skill", nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"skill add --root PATH", "--support-root PATH", "skill list", "skill remove --id ID", "skill set-default --id ID --enabled true|false"} {
		if !strings.Contains(skillOutput, required) {
			t.Fatalf("skill help missing %q:\n%s", required, skillOutput)
		}
	}
}

func TestGlobalMemoryCLIUsesExplicitGlobalScope(t *testing.T) {
	home := t.TempDir()
	t.Setenv("ADM_V2_HOME", home)

	captureStdout(t, func() {
		if err := run([]string{"memory", "global", "write", "--key", "machine", "--value", "windows"}); err != nil {
			t.Fatal(err)
		}
	})
	service := app.New(filepath.Join(home, "state.json"))
	entry, err := service.Memory.GlobalRead("machine")
	if err != nil || entry.Value != "windows" {
		t.Fatalf("Global Memory after CLI write = %+v err=%v", entry, err)
	}

	root := t.TempDir()
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "scope-check", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Memory.EnvironmentRead(env.ID, "machine"); err == nil {
		t.Fatal("Global Memory CLI write must not create Environment-private Memory")
	}

	listOutput := captureStdout(t, func() {
		if err := run([]string{"memory", "global", "list"}); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(listOutput, "machine") || !strings.Contains(listOutput, "windows") {
		t.Fatalf("memory global list missing persisted entry:\n%s", listOutput)
	}
	readOutput := captureStdout(t, func() {
		if err := run([]string{"memory", "global", "read", "--key", "machine"}); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(readOutput, "machine") || !strings.Contains(readOutput, "windows") {
		t.Fatalf("memory global read missing persisted entry:\n%s", readOutput)
	}

	captureStdout(t, func() {
		if err := run([]string{"memory", "global", "delete", "--key", "machine"}); err != nil {
			t.Fatal(err)
		}
	})
	if _, err := service.Memory.GlobalRead("machine"); err == nil {
		t.Fatal("memory global delete did not remove key")
	}
}

func TestGlobalMemoryCLIRequiresValueFlagButAllowsEmptyValue(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	if err := runGlobalMemory(service, []string{"write", "--key", "empty"}); err == nil || !strings.Contains(err.Error(), "--key 和 --value") {
		t.Fatalf("missing --value must fail clearly, got %v", err)
	}
	captureStdout(t, func() {
		if err := runGlobalMemory(service, []string{"write", "--key", "empty", "--value", ""}); err != nil {
			t.Fatal(err)
		}
	})
	entry, err := service.Memory.GlobalRead("empty")
	if err != nil || entry.Value != "" {
		t.Fatalf("explicit empty Global Memory value = %+v err=%v", entry, err)
	}
}

func TestGlobalMemoryHelpIsDiscoverable(t *testing.T) {
	output := captureStdout(t, func() {
		if err := runMemory(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(output, "memory global -h") || !strings.Contains(output, "显式选择作用域") {
		t.Fatalf("memory help is not scope-explicit:\n%s", output)
	}
	globalOutput := captureStdout(t, func() {
		if err := runGlobalMemory(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"memory global list", "memory global read --key", "memory global write --key", "memory global delete --key"} {
		if !strings.Contains(globalOutput, required) {
			t.Fatalf("global memory help missing %q:\n%s", required, globalOutput)
		}
	}
}

func TestEnvironmentRenameCLIChangesOnlyName(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	root := t.TempDir()
	marker := filepath.Join(root, "keep.txt")
	if err := os.WriteFile(marker, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "before", "")
	if err != nil {
		t.Fatal(err)
	}
	before := env
	output := captureStdout(t, func() {
		if err := runEnvironment(service, []string{"rename", "--environment-id", env.ID, "--name", "after"}); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(output, "after") {
		t.Fatalf("environment rename output missing new name:\n%s", output)
	}
	after, err := service.Environments.Get(env.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Name != "after" || after.ID != before.ID || after.WorkspaceID != before.WorkspaceID || after.Root != before.Root || !after.CreatedAt.Equal(before.CreatedAt) {
		t.Fatalf("environment rename changed stable identity: before=%+v after=%+v", before, after)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "keep" {
		t.Fatalf("environment rename touched project data: data=%q err=%v", data, err)
	}
	if err := runEnvironment(service, []string{"rename", "--environment-id", env.ID, "--name", "   "}); err == nil || !strings.Contains(err.Error(), "--name") {
		t.Fatalf("blank Environment rename must fail clearly, got %v", err)
	}
}

func TestEnvironmentHelpExplainsRenameSafety(t *testing.T) {
	output := captureStdout(t, func() {
		if err := runEnvironment(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"environment rename --environment-id", "不移动根目录", "不触碰项目文件"} {
		if !strings.Contains(output, required) {
			t.Fatalf("environment help missing %q:\n%s", required, output)
		}
	}
}

func TestEnvironmentSelectionCLIIsScopedToOneEnvironment(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	root := t.TempDir()
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}
	env1, err := service.Environments.Create(ws.ID, "one", "")
	if err != nil {
		t.Fatal(err)
	}
	env2, err := service.Environments.Create(ws.ID, "two", "")
	if err != nil {
		t.Fatal(err)
	}
	mcpEntry, err := service.MCPs.Add("filesystem", false)
	if err != nil {
		t.Fatal(err)
	}
	skillEntry, err := service.Skills.Add("go-project", false)
	if err != nil {
		t.Fatal(err)
	}

	captureStdout(t, func() {
		if err := runEnvironment(service, []string{"mcp", "enable", "--environment-id", env1.ID, "--mcp-id", mcpEntry.ID}); err != nil {
			t.Fatal(err)
		}
		if err := runEnvironment(service, []string{"skill", "enable", "--environment-id", env1.ID, "--skill-id", skillEntry.ID}); err != nil {
			t.Fatal(err)
		}
	})

	env1, err = service.Environments.Get(env1.ID)
	if err != nil {
		t.Fatal(err)
	}
	env2, err = service.Environments.Get(env2.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(env1.EnabledMCPIDs) != 1 || env1.EnabledMCPIDs[0] != mcpEntry.ID {
		t.Fatalf("env1 MCP selections = %+v", env1.EnabledMCPIDs)
	}
	if len(env1.EnabledSkillIDs) != 1 || env1.EnabledSkillIDs[0] != skillEntry.ID {
		t.Fatalf("env1 Skill selections = %+v", env1.EnabledSkillIDs)
	}
	if len(env2.EnabledMCPIDs) != 0 || len(env2.EnabledSkillIDs) != 0 {
		t.Fatalf("env2 selections changed unexpectedly: MCP=%+v Skill=%+v", env2.EnabledMCPIDs, env2.EnabledSkillIDs)
	}
	mcps, err := service.MCPs.List()
	if err != nil {
		t.Fatal(err)
	}
	skills, err := service.Skills.List()
	if err != nil {
		t.Fatal(err)
	}
	if mcps[0].DefaultIncludeInEnv || skills[0].DefaultIncludeInEnv {
		t.Fatalf("Environment selection must not change catalog defaults: MCP=%+v Skill=%+v", mcps[0], skills[0])
	}

	captureStdout(t, func() {
		if err := runEnvironment(service, []string{"mcp", "disable", "--environment-id", env1.ID, "--mcp-id", mcpEntry.ID}); err != nil {
			t.Fatal(err)
		}
		if err := runEnvironment(service, []string{"skill", "disable", "--environment-id", env1.ID, "--skill-id", skillEntry.ID}); err != nil {
			t.Fatal(err)
		}
	})
	env1, err = service.Environments.Get(env1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(env1.EnabledMCPIDs) != 0 || len(env1.EnabledSkillIDs) != 0 {
		t.Fatalf("disable did not clear selections: MCP=%+v Skill=%+v", env1.EnabledMCPIDs, env1.EnabledSkillIDs)
	}
}

func TestEnvironmentSelectionCLIRejectsUnknownCatalogID(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	root := t.TempDir()
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "one", "")
	if err != nil {
		t.Fatal(err)
	}
	err = runEnvironmentSelection("mcp", service, []string{"enable", "--environment-id", env.ID, "--mcp-id", "mcp_missing"})
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("unknown MCP ID must fail clearly, got %v", err)
	}
}

func TestEnvironmentSelectionHelpIsDiscoverable(t *testing.T) {
	output := captureStdout(t, func() {
		if err := runEnvironment(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"environment mcp -h", "environment skill -h"} {
		if !strings.Contains(output, required) {
			t.Fatalf("environment help missing %q:\n%s", required, output)
		}
	}
	for _, kind := range []string{"mcp", "skill"} {
		selectionOutput := captureStdout(t, func() {
			if err := runEnvironmentSelection(kind, nil, []string{"-h"}); err != nil {
				t.Fatal(err)
			}
		})
		for _, action := range []string{"enable", "disable"} {
			required := "environment " + kind + " " + action
			if !strings.Contains(selectionOutput, required) {
				t.Fatalf("%s selection help missing %q:\n%s", kind, required, selectionOutput)
			}
		}
		if !strings.Contains(selectionOutput, "不修改 catalog 默认值或其他 Environment") {
			t.Fatalf("%s selection help does not explain scope:\n%s", kind, selectionOutput)
		}
	}
}

func TestEnvironmentMemoryCLIIsScopedByExplicitEnvironmentID(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	root := t.TempDir()
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}
	env1, err := service.Environments.Create(ws.ID, "one", "")
	if err != nil {
		t.Fatal(err)
	}
	env2, err := service.Environments.Create(ws.ID, "two", "")
	if err != nil {
		t.Fatal(err)
	}

	captureStdout(t, func() {
		if err := runMemory(service, []string{"environment", "write", "--environment-id", env1.ID, "--key", "task", "--value", "one-only"}); err != nil {
			t.Fatal(err)
		}
	})
	entry, err := service.Memory.EnvironmentRead(env1.ID, "task")
	if err != nil || entry.Value != "one-only" {
		t.Fatalf("env1 private Memory = %+v err=%v", entry, err)
	}
	if _, err := service.Memory.EnvironmentRead(env2.ID, "task"); err == nil {
		t.Fatal("Environment-private Memory must not be readable through another Environment ID")
	}
	if _, err := service.Memory.GlobalRead("task"); err == nil {
		t.Fatal("Environment-private Memory write must not create Global Memory")
	}

	listOutput := captureStdout(t, func() {
		if err := runMemory(service, []string{"environment", "list", "--environment-id", env1.ID}); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(listOutput, "task") || !strings.Contains(listOutput, "one-only") {
		t.Fatalf("memory environment list missing private entry:\n%s", listOutput)
	}
	readOutput := captureStdout(t, func() {
		if err := runMemory(service, []string{"environment", "read", "--environment-id", env1.ID, "--key", "task"}); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(readOutput, "task") || !strings.Contains(readOutput, "one-only") {
		t.Fatalf("memory environment read missing private entry:\n%s", readOutput)
	}

	captureStdout(t, func() {
		if err := runMemory(service, []string{"environment", "delete", "--environment-id", env1.ID, "--key", "task"}); err != nil {
			t.Fatal(err)
		}
	})
	if _, err := service.Memory.EnvironmentRead(env1.ID, "task"); err == nil {
		t.Fatal("memory environment delete did not remove private key")
	}
}

func TestEnvironmentMemoryCLIRequiresExplicitScopeArgumentsAndAllowsEmptyValue(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	root := t.TempDir()
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "one", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := runEnvironmentMemory(service, []string{"write", "--key", "empty", "--value", ""}); err == nil || !strings.Contains(err.Error(), "--environment-id") {
		t.Fatalf("missing Environment ID must fail clearly, got %v", err)
	}
	if err := runEnvironmentMemory(service, []string{"write", "--environment-id", env.ID, "--key", "empty"}); err == nil || !strings.Contains(err.Error(), "--value") {
		t.Fatalf("missing --value must fail clearly, got %v", err)
	}
	captureStdout(t, func() {
		if err := runEnvironmentMemory(service, []string{"write", "--environment-id", env.ID, "--key", "empty", "--value", ""}); err != nil {
			t.Fatal(err)
		}
	})
	entry, err := service.Memory.EnvironmentRead(env.ID, "empty")
	if err != nil || entry.Value != "" {
		t.Fatalf("explicit empty Environment-private Memory value = %+v err=%v", entry, err)
	}
}

func TestEnvironmentMemoryHelpIsDiscoverable(t *testing.T) {
	output := captureStdout(t, func() {
		if err := runMemory(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"memory global -h", "memory environment -h", "显式 Environment ID"} {
		if !strings.Contains(output, required) {
			t.Fatalf("memory help missing %q:\n%s", required, output)
		}
	}
	environmentOutput := captureStdout(t, func() {
		if err := runEnvironmentMemory(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"memory environment list --environment-id", "memory environment read --environment-id", "memory environment write --environment-id", "memory environment delete --environment-id"} {
		if !strings.Contains(environmentOutput, required) {
			t.Fatalf("Environment Memory help missing %q:\n%s", required, environmentOutput)
		}
	}
}

func TestEnvironmentListAndInspectShowManagementContextWithoutMemoryValues(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	root := t.TempDir()
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

	listOutput := captureStdout(t, func() {
		if err := runEnvironment(service, []string{"list"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{env.ID, "private_memory_count", "1"} {
		if !strings.Contains(listOutput, required) {
			t.Fatalf("environment list missing %q:\n%s", required, listOutput)
		}
	}
	if strings.Contains(listOutput, "do-not-print") || strings.Contains(listOutput, "\"private_memory\"") {
		t.Fatalf("environment list leaked private Memory values:\n%s", listOutput)
	}

	inspectOutput := captureStdout(t, func() {
		if err := runEnvironment(service, []string{"inspect", "--environment-id", env.ID}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"\"workspace\"", ws.ID, "projects", "filesystem", "go-project", "unresolved_mcp_ids", removedMCP.ID, "private_memory_count"} {
		if !strings.Contains(inspectOutput, required) {
			t.Fatalf("environment inspect missing %q:\n%s", required, inspectOutput)
		}
	}
	if strings.Contains(inspectOutput, "do-not-print") || strings.Contains(inspectOutput, "\"private_memory\"") {
		t.Fatalf("environment inspect leaked private Memory values:\n%s", inspectOutput)
	}
}

func TestEnvironmentInspectHelpExplainsResolvedContextAndMemoryPrivacy(t *testing.T) {
	output := captureStdout(t, func() {
		if err := runEnvironment(nil, []string{"-h"}); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"Workspace 关系", "已解析/未解析 MCP/Skill", "不展开 Memory 值"} {
		if !strings.Contains(output, required) {
			t.Fatalf("environment help missing %q:\n%s", required, output)
		}
	}
}

func TestDoctorRunsWithoutArguments(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	service := app.New(statePath)
	output := captureStdout(t, func() {
		if err := runDoctor(service, statePath, nil); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"ADM V2 诊断", statePath, "Gateway", "Workspace（0）", "Environment（0）", "执行白名单（0）"} {
		if !strings.Contains(output, required) {
			t.Fatalf("doctor output missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayStatusReportsRunningGatewayDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"name":"ai-dev-manager-v2","version":"test-version","status":"ok","pid":43210,"transport":"http"}`)
	}))
	defer server.Close()
	listen := strings.TrimPrefix(server.URL, "http://")

	output := captureStdout(t, func() {
		if err := printGatewayStatus(listen); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"运行中", "43210", "test-version", "/mcp"} {
		if !strings.Contains(output, required) {
			t.Fatalf("gateway status missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayStatusReportsIncompatibleListenerClearly(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	listen := strings.TrimPrefix(server.URL, "http://")

	output := captureStdout(t, func() {
		if err := printGatewayStatus(listen); err != nil {
			t.Fatal(err)
		}
	})
	for _, required := range []string{"版本不兼容", "gateway restart"} {
		if !strings.Contains(output, required) {
			t.Fatalf("incompatible status missing %q:\n%s", required, output)
		}
	}
}

func TestGatewayStatusReportsStoppedWhenNothingListens(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	listen := listener.Addr().String()
	_ = listener.Close()

	output := captureStdout(t, func() {
		if err := printGatewayStatus(listen); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(output, "已停止") || !strings.Contains(output, "gateway start") {
		t.Fatalf("stopped status is not actionable:\n%s", output)
	}
}

func TestSameADMExecutableAllowsSiblingUpgradeBinary(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows executable-path semantics")
	}
	base := filepath.Join(`D:\projects\ai-dev-manager-v2`, "ai-dev-manager-v2.exe")
	next := filepath.Join(`D:\projects\ai-dev-manager-v2`, "ai-dev-manager-v2.next.exe")
	if !sameADMExecutable(base, next) {
		t.Fatal("next build in the same directory should be allowed to replace the official ADM V2 executable")
	}
	if sameADMExecutable(`D:\other\ai-dev-manager-v2.exe`, next) {
		t.Fatal("different directories must not be treated as the same ADM installation")
	}
	if sameADMExecutable(`D:\projects\ai-dev-manager-v2\other.exe`, next) {
		t.Fatal("another executable in the same directory must not be treated as ADM V2")
	}
}

func TestStopHTTPGatewayCanStopLegacySameExecutableOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("legacy Gateway process discovery is Windows-specific")
	}
	originalMatcher := matchesADMExecutable
	matchesADMExecutable = func(targetPath, currentPath string) bool {
		return strings.EqualFold(filepath.Clean(targetPath), filepath.Clean(currentPath))
	}
	t.Cleanup(func() { matchesADMExecutable = originalMatcher })

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	listen := listener.Addr().String()
	_ = listener.Close()

	cmd := exec.Command(os.Args[0], "-test.run=^TestLegacyGatewayHelperProcess$")
	cmd.Env = append(os.Environ(), "ADM_TEST_LEGACY_GATEWAY=1", "ADM_TEST_LEGACY_LISTEN="+listen)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill() })

	deadline := time.Now().Add(5 * time.Second)
	ready := false
	for time.Now().Before(deadline) {
		response, requestErr := http.Get("http://" + listen + "/healthz")
		if requestErr == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusNotFound {
				ready = true
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	if !ready {
		t.Fatal("legacy Gateway helper did not start")
	}

	output := captureStdout(t, func() {
		if err := stopHTTPGateway(listen); err != nil {
			t.Fatal(err)
		}
	})
	_ = cmd.Wait()
	if !strings.Contains(output, "检测到旧版 ADM V2 Gateway") || !strings.Contains(output, "已停止") {
		t.Fatalf("legacy Gateway stop output is not clear:\n%s", output)
	}
	if connection, err := net.DialTimeout("tcp", listen, 200*time.Millisecond); err == nil {
		_ = connection.Close()
		t.Fatalf("legacy Gateway still listens on %s", listen)
	}
}

func TestLegacyGatewayHelperProcess(t *testing.T) {
	if os.Getenv("ADM_TEST_LEGACY_GATEWAY") != "1" {
		return
	}
	listen := os.Getenv("ADM_TEST_LEGACY_LISTEN")
	server := &http.Server{
		Addr: listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" {
				http.NotFound(w, r)
				return
			}
			w.WriteHeader(http.StatusOK)
		}),
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		os.Exit(2)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	defer func() { os.Stdout = original }()

	fn()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	_ = reader.Close()
	return string(data)
}

package app_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/runtime"
)

func TestEmptyDirectoryCanBeRegisteredAndOpened(t *testing.T) {
	root := t.TempDir()
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "empty")
	if err != nil {
		t.Fatalf("register empty directory: %v", err)
	}
	env, err := service.Environments.Create(ws.ID, "empty-task", "")
	if err != nil {
		t.Fatalf("create environment for empty directory: %v", err)
	}
	items, err := service.Tree(env.ID, ".", 2, 20)
	if err != nil {
		t.Fatalf("tree empty directory: %v", err)
	}
	if got := items.([]runtime.FileInfo); len(got) != 0 {
		t.Fatalf("empty directory tree = %#v", got)
	}
}

func TestPlainDirectoryDevelopmentDoesNotRequireGit(t *testing.T) {
	root := t.TempDir()
	if _, err := os.Stat(filepath.Join(root, ".git")); !os.IsNotExist(err) {
		t.Fatalf("fixture must not be a git repository")
	}
	if err := os.WriteFile(filepath.Join(root, "hello.txt"), []byte("hello world\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	statePath := filepath.Join(t.TempDir(), "state.json")
	service := app.New(statePath)
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatalf("register plain directory: %v", err)
	}
	env, err := service.Environments.Create(ws.ID, "task-a", "")
	if err != nil {
		t.Fatalf("create environment without git: %v", err)
	}
	if env.Root != root {
		t.Fatalf("root=%q want %q", env.Root, root)
	}

	caps, err := service.Capabilities(context.Background(), env.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"files.tree", "files.read", "search.text", "files.write", "files.edit", "files.delete"} {
		if !contains(caps, required) {
			t.Fatalf("missing core capability %q in %v", required, caps)
		}
	}
	for _, gitCap := range []string{"git.status", "git.diff", "git.branch"} {
		if contains(caps, gitCap) {
			t.Fatalf("plain directory unexpectedly exposes %q", gitCap)
		}
	}

	if _, err := service.Environments.AcquireWriter(env.ID, "session-a"); err != nil {
		t.Fatalf("acquire writer: %v", err)
	}
	if _, err := service.Write(env.ID, "session-a", "nested/new.txt", "alpha beta\n", true); err != nil {
		t.Fatalf("write: %v", err)
	}
	readValue, err := service.Read(env.ID, "nested/new.txt", 0)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if readValue.(string) != "alpha beta\n" {
		t.Fatalf("read=%q", readValue)
	}
	if _, err := service.Edit(env.ID, "session-a", "nested/new.txt", "beta", "gamma", 1); err != nil {
		t.Fatalf("edit: %v", err)
	}
	searchValue, err := service.Search(env.ID, ".", "gamma", 0, 0, 0)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if !strings.Contains(strings.TrimSpace(toText(searchValue)), "gamma") {
		t.Fatalf("search result does not contain edited text: %#v", searchValue)
	}
	if _, err := service.Delete(env.ID, "session-a", "nested/new.txt"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "nested", "new.txt")); !os.IsNotExist(err) {
		t.Fatalf("deleted file still exists")
	}

	if _, err := service.GitStatus(context.Background(), env.ID); err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("git status should fail locally as unsupported, got %v", err)
	}

	reloaded := app.New(statePath)
	persisted, err := reloaded.Environments.Get(env.ID)
	if err != nil {
		t.Fatalf("environment did not survive restart: %v", err)
	}
	if persisted.Writer == nil || persisted.Writer.Owner != "session-a" {
		t.Fatalf("writer lease did not survive restart: %#v", persisted.Writer)
	}
}

func TestDevelopmentContextCatalogSelectionsAndMemoryPersistWithoutLeaking(t *testing.T) {
	root := t.TempDir()
	statePath := filepath.Join(t.TempDir(), "state.json")
	service := app.New(statePath)
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}

	mcpEntry, err := service.MCPs.Add("filesystem-extra", true)
	if err != nil {
		t.Fatal(err)
	}
	skillEntry, err := service.Skills.Add("go-project", false)
	if err != nil {
		t.Fatal(err)
	}

	envA, err := service.Environments.Create(ws.ID, "a", "")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(envA.EnabledMCPIDs, mcpEntry.ID) || contains(envA.EnabledSkillIDs, skillEntry.ID) {
		t.Fatalf("environment A defaults = mcp:%v skill:%v", envA.EnabledMCPIDs, envA.EnabledSkillIDs)
	}

	if _, err := service.MCPs.SetDefault(mcpEntry.ID, false); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Skills.SetDefault(skillEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	envB, err := service.Environments.Create(ws.ID, "b", "")
	if err != nil {
		t.Fatal(err)
	}
	if contains(envB.EnabledMCPIDs, mcpEntry.ID) || !contains(envB.EnabledSkillIDs, skillEntry.ID) {
		t.Fatalf("environment B defaults = mcp:%v skill:%v", envB.EnabledMCPIDs, envB.EnabledSkillIDs)
	}

	envA, err = service.Environments.Get(envA.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(envA.EnabledMCPIDs, mcpEntry.ID) || contains(envA.EnabledSkillIDs, skillEntry.ID) {
		t.Fatalf("global default change rewrote existing environment A: %+v", envA)
	}

	if _, err := service.SetEnvironmentMCP(envB.ID, mcpEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentSkill(envA.ID, skillEntry.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := service.Memory.GlobalWrite("machine", "windows"); err != nil {
		t.Fatal(err)
	}
	if err := service.Memory.EnvironmentWrite(envA.ID, "task", "login"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Memory.EnvironmentRead(envB.ID, "task"); err == nil {
		t.Fatal("environment B must not read environment A private memory")
	}
	if got, err := service.Memory.GlobalRead("machine"); err != nil || got.Value != "windows" {
		t.Fatalf("global memory read = %+v err=%v", got, err)
	}

	if err := service.MCPs.Remove(mcpEntry.ID); err != nil {
		t.Fatal(err)
	}
	envB, err = service.Environments.Get(envB.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(envB.EnabledMCPIDs, mcpEntry.ID) {
		t.Fatalf("removed global MCP should leave explicit unresolved selection, got %v", envB.EnabledMCPIDs)
	}

	reloaded := app.New(statePath)
	persistedA, err := reloaded.Environments.Get(envA.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(persistedA.EnabledSkillIDs, skillEntry.ID) {
		t.Fatalf("environment skill selection did not persist: %v", persistedA.EnabledSkillIDs)
	}
	if got, err := reloaded.Memory.EnvironmentRead(envA.ID, "task"); err != nil || got.Value != "login" {
		t.Fatalf("environment memory did not persist: %+v err=%v", got, err)
	}
	if got, err := reloaded.Memory.GlobalRead("machine"); err != nil || got.Value != "windows" {
		t.Fatalf("global memory did not persist: %+v err=%v", got, err)
	}
}

func TestNoCatalogOrMemoryDataIsNotADevelopmentPrerequisite(t *testing.T) {
	root := t.TempDir()
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "no-context-data", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(env.EnabledMCPIDs) != 0 || len(env.EnabledSkillIDs) != 0 || len(env.PrivateMemory) != 0 {
		t.Fatalf("expected empty optional context data: %+v", env)
	}
	if _, err := service.Environments.AcquireWriter(env.ID, "session"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Write(env.ID, "session", "works.txt", "still develops\n", false); err != nil {
		t.Fatalf("development must work without MCP/Skill/Memory data: %v", err)
	}
}

func TestSingleWriterIsScopedByPhysicalRoot(t *testing.T) {
	root := t.TempDir()
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "plain")
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
	if _, err := service.Environments.AcquireWriter(envA.ID, "session-a"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Environments.AcquireWriter(envB.ID, "session-b"); err == nil {
		t.Fatal("expected writer conflict for two environments sharing one physical root")
	}
	if _, err := service.Write(envB.ID, "session-b", "should-not-exist.txt", "x", true); err == nil {
		t.Fatal("mutation without writer should fail")
	}
}

func TestExecIsOptionalAndExplicitlyAllowlisted(t *testing.T) {
	root := t.TempDir()
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "exec", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Environments.AcquireWriter(env.ID, "session-exec"); err != nil {
		t.Fatal(err)
	}

	caps, err := service.Capabilities(context.Background(), env.ID)
	if err != nil {
		t.Fatal(err)
	}
	if contains(caps, "shell.exec") {
		t.Fatalf("shell.exec should not exist before allowlisting: %v", caps)
	}

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err := service.AllowExecutable(exe); err != nil {
		t.Fatal(err)
	}
	caps, err = service.Capabilities(context.Background(), env.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(caps, "shell.exec") {
		t.Fatalf("shell.exec missing after allowlisting: %v", caps)
	}
	value, err := service.Exec(context.Background(), env.ID, "session-exec", exe, []string{"-test.run=TestExecHelperProcess"}, "", 10000, 20000)
	if err != nil {
		t.Fatalf("exec allowlisted test binary: %v", err)
	}
	if !strings.Contains(toText(value), "ADM_V2_EXEC_OK") {
		t.Fatalf("unexpected exec result: %#v", value)
	}
}

func TestExecHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "" {
		return
	}
	if strings.Contains(strings.Join(os.Args, " "), "TestExecHelperProcess") {
		fmtPrintln("ADM_V2_EXEC_OK")
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func toText(value any) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(fmtSprint(value)), "{", ""), "}", ""))
}

func fmtSprint(value any) string { return fmt.Sprint(value) }
func fmtPrintln(value any)       { fmt.Println(value) }

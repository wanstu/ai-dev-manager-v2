package management_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/app"
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

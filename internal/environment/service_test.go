package environment

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/store"
	"ai-dev-manager-v2/internal/workspace"
)

func TestWriterLeaseExpiresAndAnotherEnvironmentCanAcquireRoot(t *testing.T) {
	service, envA, envB, now := newLeaseTestService(t)

	acquired, err := service.AcquireWriter(envA.ID, "session-a")
	if err != nil {
		t.Fatal(err)
	}
	if acquired.Writer == nil {
		t.Fatal("writer lease was not created")
	}
	if got, want := acquired.Writer.ExpiresAt, now.Add(service.WriterLeaseTTL()); !got.Equal(want) {
		t.Fatalf("expires_at=%s want %s", got, want)
	}

	*now = now.Add(service.WriterLeaseTTL() + time.Second)
	if visible, err := service.Get(envA.ID); err != nil {
		t.Fatal(err)
	} else if visible.Writer != nil {
		t.Fatalf("expired writer must be hidden from Get: %#v", visible.Writer)
	}
	listed, err := service.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, env := range listed {
		if env.ID == envA.ID && env.Writer != nil {
			t.Fatalf("expired writer must be hidden from List: %#v", env.Writer)
		}
	}

	acquiredB, err := service.AcquireWriter(envB.ID, "session-b")
	if err != nil {
		t.Fatalf("expired writer should not block another environment on the same root: %v", err)
	}
	if acquiredB.Writer == nil || acquiredB.Writer.Owner != "session-b" {
		t.Fatalf("unexpected replacement writer: %#v", acquiredB.Writer)
	}
}

func TestWriterLeaseRenewsOnAcquireRequireAndTouch(t *testing.T) {
	service, envA, _, now := newLeaseTestService(t)

	first, err := service.AcquireWriter(envA.ID, "session-a")
	if err != nil {
		t.Fatal(err)
	}
	acquiredAt := first.Writer.AcquiredAt
	firstExpiry := first.Writer.ExpiresAt

	*now = now.Add(service.WriterLeaseTTL() / 2)
	renewed, err := service.AcquireWriter(envA.ID, "session-a")
	if err != nil {
		t.Fatal(err)
	}
	if !renewed.Writer.AcquiredAt.Equal(acquiredAt) {
		t.Fatalf("renewal changed acquired_at: got %s want %s", renewed.Writer.AcquiredAt, acquiredAt)
	}
	if !renewed.Writer.ExpiresAt.After(firstExpiry) {
		t.Fatalf("acquire renewal did not extend lease: first=%s renewed=%s", firstExpiry, renewed.Writer.ExpiresAt)
	}

	secondExpiry := renewed.Writer.ExpiresAt
	*now = now.Add(service.WriterLeaseTTL() / 2)
	required, err := service.RequireWriter(envA.ID, "session-a")
	if err != nil {
		t.Fatal(err)
	}
	if !required.Writer.ExpiresAt.After(secondExpiry) {
		t.Fatalf("RequireWriter did not renew lease: previous=%s renewed=%s", secondExpiry, required.Writer.ExpiresAt)
	}

	activityBefore := required.LastActivityAt
	*now = now.Add(time.Second)
	if err := service.Touch(envA.ID, "session-a"); err != nil {
		t.Fatal(err)
	}
	visible, err := service.Get(envA.ID)
	if err != nil {
		t.Fatal(err)
	}
	if visible.Writer == nil || !visible.Writer.LastSeenAt.Equal(*now) {
		t.Fatalf("Touch did not renew last_seen_at: %#v", visible.Writer)
	}
	if !visible.LastActivityAt.After(activityBefore) {
		t.Fatalf("Touch did not update last_activity_at: before=%s after=%s", activityBefore, visible.LastActivityAt)
	}
}

func TestExpiredWriterNoLongerAuthorizesMutation(t *testing.T) {
	service, envA, _, now := newLeaseTestService(t)
	if _, err := service.AcquireWriter(envA.ID, "session-a"); err != nil {
		t.Fatal(err)
	}

	*now = now.Add(service.WriterLeaseTTL())
	if _, err := service.RequireWriter(envA.ID, "session-a"); err == nil {
		t.Fatal("expired writer must not pass RequireWriter")
	}
	if released, err := service.ReleaseWriter(envA.ID, "some-other-owner", false); err != nil {
		t.Fatalf("expired writer should be releasable without force: %v", err)
	} else if released.Writer != nil {
		t.Fatalf("expired writer was not cleared on release: %#v", released.Writer)
	}
}

func TestWriterWithoutExpiresAtIsImmediatelyInvalid(t *testing.T) {
	service, envA, envB, _ := newLeaseTestService(t)
	if _, err := service.AcquireWriter(envA.ID, "session-a"); err != nil {
		t.Fatal(err)
	}
	if err := service.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, envA.ID)
		state.Environments[idx].Writer.ExpiresAt = time.Time{}
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if visible, err := service.Get(envA.ID); err != nil {
		t.Fatal(err)
	} else if visible.Writer != nil {
		t.Fatalf("writer without expires_at must be invalid: %#v", visible.Writer)
	}
	if _, err := service.AcquireWriter(envB.ID, "session-b"); err != nil {
		t.Fatalf("invalid writer should not block the same physical root: %v", err)
	}
}

func TestCreateReusesSameWorkspaceNameAndRoot(t *testing.T) {
	root := t.TempDir()
	stateStore := store.New(filepath.Join(t.TempDir(), "state.json"))
	workspaces := workspace.New(stateStore)
	ws, err := workspaces.Add(root, "projects")
	if err != nil {
		t.Fatal(err)
	}
	service := New(stateStore, workspaces)

	first, err := service.Create(ws.ID, "main", "")
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Create(ws.ID, "MAIN", root)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID {
		t.Fatalf("duplicate create must reuse existing environment: first=%s second=%s", first.ID, second.ID)
	}
	items, err := service.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("duplicate create produced %d environments, want 1", len(items))
	}

	other, err := service.Create(ws.ID, "isolated", root)
	if err != nil {
		t.Fatal(err)
	}
	if other.ID == first.ID {
		t.Fatal("different environment names on the same root must remain distinct contexts")
	}
}

func TestRemoveEnvironmentProtectsActiveWriterAndNeverDeletesRoot(t *testing.T) {
	service, envA, _, now := newLeaseTestService(t)
	marker := filepath.Join(envA.Root, "keep.txt")
	if err := os.WriteFile(marker, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := service.AcquireWriter(envA.ID, "session-a"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Remove(envA.ID); err == nil {
		t.Fatal("active writer must block environment removal")
	}

	*now = now.Add(service.WriterLeaseTTL() + time.Second)
	removed, err := service.Remove(envA.ID)
	if err != nil {
		t.Fatalf("expired writer should not block removal: %v", err)
	}
	if removed.ID != envA.ID {
		t.Fatalf("removed environment id=%s want %s", removed.ID, envA.ID)
	}
	if _, err := service.Get(envA.ID); err == nil {
		t.Fatal("removed environment must no longer be retrievable")
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "keep" {
		t.Fatalf("environment removal must not delete root files: data=%q err=%v", data, err)
	}
}

func TestRenameEnvironmentChangesOnlyMetadata(t *testing.T) {
	service, envA, _, _ := newLeaseTestService(t)
	marker := filepath.Join(envA.Root, "keep.txt")
	if err := os.WriteFile(marker, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetMCP(envA.ID, "mcp_keep", true); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetSkill(envA.ID, "skill_keep", true); err != nil {
		t.Fatal(err)
	}
	if err := service.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, envA.ID)
		if state.Environments[idx].PrivateMemory == nil {
			state.Environments[idx].PrivateMemory = map[string]string{}
		}
		state.Environments[idx].PrivateMemory["note"] = "keep"
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	before, err := service.AcquireWriter(envA.ID, "session-a")
	if err != nil {
		t.Fatal(err)
	}

	renamed, err := service.Rename(envA.ID, "  renamed  ")
	if err != nil {
		t.Fatal(err)
	}
	if renamed.Name != "renamed" || renamed.ID != before.ID || renamed.WorkspaceID != before.WorkspaceID || renamed.Root != before.Root || !renamed.CreatedAt.Equal(before.CreatedAt) {
		t.Fatalf("rename changed stable environment identity: before=%+v after=%+v", before, renamed)
	}
	if len(renamed.EnabledMCPIDs) != 1 || renamed.EnabledMCPIDs[0] != "mcp_keep" || len(renamed.EnabledSkillIDs) != 1 || renamed.EnabledSkillIDs[0] != "skill_keep" {
		t.Fatalf("rename changed selections: %+v", renamed)
	}
	if renamed.PrivateMemory["note"] != "keep" {
		t.Fatalf("rename changed private memory: %+v", renamed.PrivateMemory)
	}
	if renamed.Writer == nil || renamed.Writer.Owner != "session-a" {
		t.Fatalf("rename changed active writer: %+v", renamed.Writer)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "keep" {
		t.Fatalf("rename touched project files: data=%q err=%v", data, err)
	}
	if _, err := service.Rename(envA.ID, "   "); err == nil {
		t.Fatal("blank environment name must be rejected")
	}
}

func newLeaseTestService(t *testing.T) (*Service, model.Environment, model.Environment, *time.Time) {
	t.Helper()
	root := t.TempDir()
	stateStore := store.New(filepath.Join(t.TempDir(), "state.json"))
	workspaces := workspace.New(stateStore)
	ws, err := workspaces.Add(root, "lease-test")
	if err != nil {
		t.Fatal(err)
	}
	service := New(stateStore, workspaces)
	current := time.Now().UTC()
	service.now = func() time.Time { return current }
	service.writerLeaseTTL = time.Minute

	envA, err := service.Create(ws.ID, "a", "")
	if err != nil {
		t.Fatal(err)
	}
	envB, err := service.Create(ws.ID, "b", "")
	if err != nil {
		t.Fatal(err)
	}
	current = time.Now().UTC().Add(time.Second)
	return service, envA, envB, &current
}

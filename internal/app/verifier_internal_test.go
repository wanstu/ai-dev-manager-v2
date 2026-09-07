package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/verifier"
)

func TestRunVerifierHeartbeatsWriterDuringLongRun(t *testing.T) {
	root := t.TempDir()
	service := New(filepath.Join(t.TempDir(), "state.json"))
	service.writerHeartbeatInterval = func(time.Duration) time.Duration { return 20 * time.Millisecond }
	ws, err := service.Workspaces.Add(root, "heartbeat")
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "heartbeat", "")
	if err != nil {
		t.Fatal(err)
	}
	acquired, err := service.Environments.AcquireWriter(env.ID, "heartbeat-owner")
	if err != nil {
		t.Fatal(err)
	}
	if acquired.Writer == nil {
		t.Fatal("writer lease missing after acquire")
	}
	initialLastSeen := acquired.Writer.LastSeenAt

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err := service.AllowExecutable(exe); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADM_V2_VERIFIER_HEARTBEAT_HELPER", "1")
	definition, err := service.AddVerifier(env.ID, model.VerifierDefinition{
		Kind:           verifier.KindTest,
		Enabled:        true,
		Executable:     exe,
		Args:           []string{"-test.run=^TestVerifierHeartbeatHelperProcess$"},
		TimeoutSeconds: 30,
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.RunVerifier(context.Background(), env.ID, "heartbeat-owner", definition.ID, 4096)
	if err != nil {
		t.Fatalf("run verifier with heartbeat: %v", err)
	}
	if result.Status != verifier.StatusPassed {
		t.Fatalf("heartbeat verifier result = %+v", result)
	}
	after, err := service.Environments.Get(env.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Writer == nil || !after.Writer.LastSeenAt.After(initialLastSeen) {
		t.Fatalf("writer heartbeat did not advance during verifier run: before=%s after=%+v", initialLastSeen, after.Writer)
	}
}

func TestVerifierHeartbeatHelperProcess(t *testing.T) {
	if os.Getenv("ADM_V2_VERIFIER_HEARTBEAT_HELPER") != "1" {
		return
	}
	time.Sleep(150 * time.Millisecond)
	fmt.Fprintln(os.Stdout, "HEARTBEAT_VERIFIER_OK")
}

package gateway

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/verifier"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestWorkflowLifecycleThroughRealStreamableHTTP(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "marker.txt"), []byte("phase9-http-marker\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	service := app.New(statePath)
	workspace, err := service.Workspaces.Add(root, "phase9-http")
	if err != nil {
		t.Fatal(err)
	}
	environment, err := service.Environments.Create(workspace.ID, "phase9-http", "")
	if err != nil {
		t.Fatal(err)
	}
	if err := service.AllowExecutable(os.Args[0]); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Environments.AcquireWriter(environment.ID, agentRunTestWriter); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADM_TEST_WORKFLOW_HELPER", "1")
	passing, err := service.AddVerifier(environment.ID, model.VerifierDefinition{
		Kind: verifier.KindCustom, Enabled: true, Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowVerifierPassHelper$"}, Name: "phase9-http-pass",
	})
	if err != nil {
		t.Fatal(err)
	}
	rejecting, err := service.AddVerifier(environment.ID, model.VerifierDefinition{
		Kind: verifier.KindCustom, Enabled: true, Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowVerifierFailHelper$"}, Name: "phase9-http-reject",
	})
	if err != nil {
		t.Fatal(err)
	}

	owner := newRuntimeOwner(service)
	defer owner.Close()
	httpServer := httptest.NewServer(newHTTPHandler(service, owner))
	defer httpServer.Close()
	endpoint := httpServer.URL + "/mcp"
	ctx := context.Background()

	first := connectHTTPWithRetry(t, ctx, endpoint)
	acceptedStart := callGatewayTool(t, ctx, first, "run_workflow_start", map[string]any{
		"environment_id": environment.ID,
		"writer_owner":   agentRunTestWriter,
		"goal":           "real HTTP accepted workflow",
		"steps": []map[string]any{{
			"name":       "execute",
			"executable": os.Args[0],
			"args":       []string{"-test.run=^TestWorkflowStepSuccessHelper$"},
		}},
		"verifier_ids":     []string{passing.ID},
		"max_output_bytes": 4096,
	})
	if acceptedStart.IsError {
		t.Fatalf("accepted workflow start failed: %s", toolText(t, acceptedStart))
	}
	acceptedID := runIDFromToolText(t, acceptedStart)
	_ = first.Close()

	second := connectHTTPWithRetry(t, ctx, endpoint)
	accepted := waitHTTPWorkflowTerminal(t, ctx, second, environment.ID, acceptedID)
	if accepted.State != agentRunSucceeded || accepted.Kind != agentRunKindWorkflow || accepted.Workflow == nil || accepted.Workflow.Outcome != workflowOutcomeAccepted || accepted.Workflow.Review.State != workflowReviewAccepted {
		_ = second.Close()
		t.Fatalf("accepted HTTP workflow=%+v", accepted)
	}
	if accepted.Workflow.Plan.Goal != "real HTTP accepted workflow" || len(accepted.Workflow.Steps) != 1 || accepted.Workflow.Steps[0].State != workflowStepSucceeded {
		_ = second.Close()
		t.Fatalf("accepted HTTP audit=%+v", accepted.Workflow)
	}

	rejectedStart := callGatewayTool(t, ctx, second, "run_workflow_start", map[string]any{
		"environment_id": environment.ID,
		"writer_owner":   agentRunTestWriter,
		"goal":           "real HTTP rejected workflow",
		"steps": []map[string]any{{
			"name":       "execute",
			"executable": os.Args[0],
			"args":       []string{"-test.run=^TestWorkflowStepSuccessHelper$"},
		}},
		"verifier_ids": []string{rejecting.ID},
	})
	if rejectedStart.IsError {
		_ = second.Close()
		t.Fatalf("rejected workflow start failed: %s", toolText(t, rejectedStart))
	}
	rejectedID := runIDFromToolText(t, rejectedStart)
	_ = second.Close()

	third := connectHTTPWithRetry(t, ctx, endpoint)
	defer third.Close()
	rejected := waitHTTPWorkflowTerminal(t, ctx, third, environment.ID, rejectedID)
	if rejected.State != agentRunSucceeded || rejected.ErrorKind != "" || rejected.Workflow == nil || rejected.Workflow.Outcome != workflowOutcomeRejected || rejected.Workflow.Review.State != workflowReviewRejected {
		t.Fatalf("review rejection was not distinct over HTTP: %+v", rejected)
	}
	listed := callGatewayTool(t, ctx, third, "run_list", map[string]any{"environment_id": environment.ID})
	listedText := toolText(t, listed)
	if listed.IsError || !strings.Contains(listedText, acceptedID) || !strings.Contains(listedText, rejectedID) || !strings.Contains(listedText, `"kind":"workflow"`) {
		t.Fatalf("workflow runs missing from later-client list: %s", listedText)
	}
	read := callGatewayTool(t, ctx, third, "read", map[string]any{"environment_id": environment.ID, "path": "marker.txt"})
	if read.IsError || !strings.Contains(toolText(t, read), "phase9-http-marker") {
		t.Fatalf("workflow capability broke ordinary read: %s", toolText(t, read))
	}

	stateBytes, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	persisted := string(stateBytes)
	for _, forbidden := range []string{acceptedID, rejectedID, "workflow", "outcome", "reviewer_error"} {
		if strings.Contains(persisted, forbidden) {
			t.Fatalf("workflow observation leaked into persisted state via %q", forbidden)
		}
	}
}

func waitHTTPWorkflowTerminal(t *testing.T, ctx context.Context, session *mcp.ClientSession, environmentID, runID string) agentRunStatus {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		result := callGatewayTool(t, ctx, session, "run_status", map[string]any{"environment_id": environmentID, "run_id": runID})
		if !result.IsError {
			var output struct {
				Result agentRunStatus `json:"result"`
			}
			if err := json.Unmarshal([]byte(toolText(t, result)), &output); err != nil {
				t.Fatalf("decode run_status: %v", err)
			}
			if output.Result.State != agentRunRunning {
				return output.Result
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("workflow run %s did not reach terminal state", runID)
	return agentRunStatus{}
}

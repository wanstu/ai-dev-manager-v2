package gateway

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/verifier"
)

func TestWorkflowRunAcceptedPlanExecutorReviewerAudit(t *testing.T) {
	service, environmentID, _ := agentRunTestService(t)
	t.Setenv("ADM_TEST_WORKFLOW_HELPER", "1")
	definition, err := service.AddVerifier(environmentID, model.VerifierDefinition{
		Kind:       verifier.KindCustom,
		Enabled:    true,
		Executable: os.Args[0],
		Args:       []string{"-test.run=^TestWorkflowVerifierPassHelper$"},
		Name:       "phase9-pass-review",
	})
	if err != nil {
		t.Fatal(err)
	}
	owner := newRuntimeOwner(service)
	defer owner.Close()

	started, err := owner.StartWorkflowRun(environmentID, agentRunTestWriter, "compile and review", []workflowStepRequest{
		{Name: "compile", Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowStepSuccessHelper$"}},
		{Name: "exercise", Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowStepSuccessHelper$"}},
	}, []string{definition.ID}, 4096)
	if err != nil {
		t.Fatal(err)
	}
	status := waitAgentRunTerminal(t, owner, environmentID, started.ID)
	if status.Kind != agentRunKindWorkflow || status.State != agentRunSucceeded || status.ErrorKind != "" {
		t.Fatalf("workflow run status=%+v", status)
	}
	if status.Workflow == nil {
		t.Fatal("workflow audit missing")
	}
	workflow := status.Workflow
	if workflow.Plan.Goal != "compile and review" || len(workflow.Plan.Steps) != 2 || workflow.Plan.Steps[0].ID != "step_01" || workflow.Plan.Steps[1].ID != "step_02" {
		t.Fatalf("plan=%+v", workflow.Plan)
	}
	if len(workflow.Plan.VerifierIDs) != 1 || workflow.Plan.VerifierIDs[0] != definition.ID {
		t.Fatalf("plan verifier refs=%v", workflow.Plan.VerifierIDs)
	}
	for _, step := range workflow.Steps {
		if step.State != workflowStepSucceeded || step.ExitCode == nil || *step.ExitCode != 0 || !strings.Contains(step.Stdout, "phase9-step-success") {
			t.Fatalf("step=%+v", step)
		}
	}
	if workflow.Review.State != workflowReviewAccepted || workflow.Outcome != workflowOutcomeAccepted || len(workflow.Review.Verifiers) != 1 || workflow.Review.Verifiers[0].Status != verifier.StatusPassed {
		t.Fatalf("review=%+v outcome=%q", workflow.Review, workflow.Outcome)
	}
}

func TestWorkflowRunReviewRejectionIsNotOrchestrationFailure(t *testing.T) {
	service, environmentID, _ := agentRunTestService(t)
	t.Setenv("ADM_TEST_WORKFLOW_HELPER", "1")
	definition, err := service.AddVerifier(environmentID, model.VerifierDefinition{
		Kind:       verifier.KindCustom,
		Enabled:    true,
		Executable: os.Args[0],
		Args:       []string{"-test.run=^TestWorkflowVerifierFailHelper$"},
		Name:       "phase9-reject-review",
	})
	if err != nil {
		t.Fatal(err)
	}
	owner := newRuntimeOwner(service)
	defer owner.Close()

	started, err := owner.StartWorkflowRun(environmentID, agentRunTestWriter, "review rejection", []workflowStepRequest{
		{Name: "execute", Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowStepSuccessHelper$"}},
	}, []string{definition.ID}, 4096)
	if err != nil {
		t.Fatal(err)
	}
	status := waitAgentRunTerminal(t, owner, environmentID, started.ID)
	if status.State != agentRunSucceeded || status.ErrorKind != "" || status.Workflow == nil {
		t.Fatalf("review rejection became orchestration failure: %+v", status)
	}
	if status.Workflow.Review.State != workflowReviewRejected || status.Workflow.Outcome != workflowOutcomeRejected {
		t.Fatalf("review=%+v outcome=%q", status.Workflow.Review, status.Workflow.Outcome)
	}
	if len(status.Workflow.Review.Verifiers) != 1 || status.Workflow.Review.Verifiers[0].Status != verifier.StatusFailed || status.Workflow.Review.Verifiers[0].ExitCode != 7 {
		t.Fatalf("rejected verifier evidence=%+v", status.Workflow.Review.Verifiers)
	}
}

func TestWorkflowRunExecutorAndReviewerFailuresAreDistinct(t *testing.T) {
	service, environmentID, _ := agentRunTestService(t)
	t.Setenv("ADM_TEST_WORKFLOW_HELPER", "1")
	pass, err := service.AddVerifier(environmentID, model.VerifierDefinition{
		Kind: verifier.KindCustom, Enabled: true, Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowVerifierPassHelper$"},
	})
	if err != nil {
		t.Fatal(err)
	}
	owner := newRuntimeOwner(service)
	defer owner.Close()

	executorFailed, err := owner.StartWorkflowRun(environmentID, agentRunTestWriter, "executor failure", []workflowStepRequest{
		{Name: "fail", Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowStepFailHelper$"}},
	}, []string{pass.ID}, 4096)
	if err != nil {
		t.Fatal(err)
	}
	executorStatus := waitAgentRunTerminal(t, owner, environmentID, executorFailed.ID)
	if executorStatus.State != agentRunFailed || executorStatus.ErrorKind != "executor_step_failed" || executorStatus.Workflow == nil || executorStatus.Workflow.Review.State != workflowReviewNotRun {
		t.Fatalf("executor failure status=%+v", executorStatus)
	}
	if executorStatus.Workflow.Steps[0].State != workflowStepFailed || executorStatus.Workflow.Steps[0].ExitCode == nil || *executorStatus.Workflow.Steps[0].ExitCode != 9 {
		t.Fatalf("executor step evidence=%+v", executorStatus.Workflow.Steps[0])
	}

	reviewerFailed, err := owner.StartWorkflowRun(environmentID, agentRunTestWriter, "reviewer infrastructure failure", []workflowStepRequest{
		{Name: "execute", Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowStepSuccessHelper$"}},
	}, []string{"vf_missing"}, 4096)
	if err != nil {
		t.Fatal(err)
	}
	reviewerStatus := waitAgentRunTerminal(t, owner, environmentID, reviewerFailed.ID)
	if reviewerStatus.State != agentRunFailed || reviewerStatus.ErrorKind != "reviewer_error" || reviewerStatus.Workflow == nil || reviewerStatus.Workflow.Review.State != workflowReviewError || reviewerStatus.Workflow.Review.ErrorKind != "reviewer_error" {
		t.Fatalf("reviewer failure status=%+v", reviewerStatus)
	}
}

func TestWorkflowPlannerAuthorityFailureDoesNotInstallRun(t *testing.T) {
	service, environmentID, _ := agentRunTestService(t)
	owner := newRuntimeOwner(service)
	defer owner.Close()
	before, err := owner.ListAgentRuns(environmentID)
	if err != nil {
		t.Fatal(err)
	}
	for name, steps := range map[string][]workflowStepRequest{
		"forbidden executable": {{Name: "bad", Executable: "definitely-not-allowed"}},
		"escaped cwd":          {{Name: "bad", Executable: os.Args[0], Args: []string{"-test.run=^$"}, Cwd: "../escape"}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := owner.StartWorkflowRun(environmentID, agentRunTestWriter, "unsafe", steps, []string{"vf_later"}, 1024); err == nil {
				t.Fatal("unsafe workflow unexpectedly started")
			}
		})
	}
	after, err := owner.ListAgentRuns(environmentID)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != len(before) {
		t.Fatalf("planner authority failure installed run: before=%d after=%d", len(before), len(after))
	}
}

func TestWorkflowRunCancellationReusesAgentRunBoundary(t *testing.T) {
	service, environmentID, root := agentRunTestService(t)
	portFile := filepath.Join(root, "workflow-cancel.port")
	t.Setenv("ADM_TEST_DEV_PROCESS_HELPER", "1")
	t.Setenv("ADM_TEST_DEV_PROCESS_PORT_FILE", portFile)
	pass, err := service.AddVerifier(environmentID, model.VerifierDefinition{
		Kind: verifier.KindCustom, Enabled: true, Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowVerifierPassHelper$"},
	})
	if err != nil {
		t.Fatal(err)
	}
	owner := newRuntimeOwner(service)
	defer owner.Close()
	started, err := owner.StartWorkflowRun(environmentID, agentRunTestWriter, "cancel workflow", []workflowStepRequest{
		{Name: "serve", Executable: os.Args[0], Args: []string{"-test.run=^TestDevProcessHelper$"}, TimeoutMS: 30000},
	}, []string{pass.ID}, 4096)
	if err != nil {
		t.Fatal(err)
	}
	port := waitDevProcessPortFile(t, portFile)
	canceled, err := owner.CancelAgentRun(environmentID, agentRunTestWriter, started.ID)
	if err != nil {
		t.Fatal(err)
	}
	if canceled.State != agentRunCanceled || canceled.Workflow == nil || canceled.Workflow.Steps[0].State != workflowStepCanceled || canceled.Workflow.Review.State != workflowReviewCanceled {
		t.Fatalf("canceled workflow=%+v", canceled)
	}
	waitPortReleased(t, port)
}

func TestWorkflowRunBoundsExecutorOutput(t *testing.T) {
	service, environmentID, _ := agentRunTestService(t)
	t.Setenv("ADM_TEST_AGENT_RUN_HELPER", "1")
	t.Setenv("ADM_TEST_AGENT_RUN_MODE", "large")
	t.Setenv("ADM_TEST_WORKFLOW_HELPER", "1")
	pass, err := service.AddVerifier(environmentID, model.VerifierDefinition{
		Kind: verifier.KindCustom, Enabled: true, Executable: os.Args[0], Args: []string{"-test.run=^TestWorkflowVerifierPassHelper$"},
	})
	if err != nil {
		t.Fatal(err)
	}
	owner := newRuntimeOwner(service)
	defer owner.Close()
	started, err := owner.StartWorkflowRun(environmentID, agentRunTestWriter, "bounded output", []workflowStepRequest{
		{Name: "large", Executable: os.Args[0], Args: []string{"-test.run=^TestAgentRunCommandHelper$"}},
	}, []string{pass.ID}, 64)
	if err != nil {
		t.Fatal(err)
	}
	status := waitAgentRunTerminal(t, owner, environmentID, started.ID)
	if status.State != agentRunSucceeded || status.Workflow == nil || len(status.Workflow.Steps[0].Stdout) != 64 {
		t.Fatalf("bounded workflow status=%+v", status)
	}
}

func TestWorkflowStepSuccessHelper(t *testing.T) {
	if os.Getenv("ADM_TEST_WORKFLOW_HELPER") != "1" {
		return
	}
	fmt.Fprintln(os.Stdout, "phase9-step-success")
}

func TestWorkflowStepFailHelper(t *testing.T) {
	if os.Getenv("ADM_TEST_WORKFLOW_HELPER") != "1" {
		return
	}
	fmt.Fprintln(os.Stderr, "phase9-step-fail")
	os.Exit(9)
}

func TestWorkflowVerifierPassHelper(t *testing.T) {
	if os.Getenv("ADM_TEST_WORKFLOW_HELPER") != "1" {
		return
	}
	fmt.Fprintln(os.Stdout, "phase9-review-pass")
}

func TestWorkflowVerifierFailHelper(t *testing.T) {
	if os.Getenv("ADM_TEST_WORKFLOW_HELPER") != "1" {
		return
	}
	fmt.Fprintln(os.Stderr, "phase9-review-reject")
	os.Exit(7)
}

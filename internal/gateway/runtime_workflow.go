package gateway

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/identity"
	runtimepkg "ai-dev-manager-v2/internal/runtime"
	"ai-dev-manager-v2/internal/verifier"
)

const (
	maxWorkflowSteps      = 32
	maxWorkflowVerifiers  = 32
	defaultWorkflowOutput = 120000
)

type workflowStepRequest struct {
	Name       string
	Executable string
	Args       []string
	Cwd        string
	TimeoutMS  int64
}

type workflowStepState string

const (
	workflowStepPending   workflowStepState = "pending"
	workflowStepRunning   workflowStepState = "running"
	workflowStepSucceeded workflowStepState = "succeeded"
	workflowStepFailed    workflowStepState = "failed"
	workflowStepCanceled  workflowStepState = "canceled"
)

type workflowReviewState string

const (
	workflowReviewPending  workflowReviewState = "pending"
	workflowReviewRunning  workflowReviewState = "running"
	workflowReviewAccepted workflowReviewState = "accepted"
	workflowReviewRejected workflowReviewState = "rejected"
	workflowReviewError    workflowReviewState = "error"
	workflowReviewNotRun   workflowReviewState = "not_run"
	workflowReviewCanceled workflowReviewState = "canceled"
)

type workflowOutcome string

const (
	workflowOutcomeAccepted workflowOutcome = "accepted"
	workflowOutcomeRejected workflowOutcome = "rejected"
)

type workflowPlanStep struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Executable string   `json:"executable"`
	Args       []string `json:"args,omitempty"`
	Cwd        string   `json:"cwd,omitempty"`
	TimeoutMS  int64    `json:"timeout_ms,omitempty"`
}

type workflowPlan struct {
	Goal        string             `json:"goal"`
	Steps       []workflowPlanStep `json:"steps"`
	VerifierIDs []string           `json:"verifier_ids"`
}

type workflowStepStatus struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	State       workflowStepState `json:"state"`
	StartedAt   *time.Time        `json:"started_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	ExitCode    *int              `json:"exit_code,omitempty"`
	Stdout      string            `json:"stdout,omitempty"`
	Stderr      string            `json:"stderr,omitempty"`
	ErrorKind   string            `json:"error_kind,omitempty"`
	Message     string            `json:"message,omitempty"`
}

type workflowReviewStatus struct {
	State     workflowReviewState `json:"state"`
	Verifiers []verifier.Result   `json:"verifiers,omitempty"`
	ErrorKind string              `json:"error_kind,omitempty"`
	Message   string              `json:"message,omitempty"`
}

type workflowRunStatus struct {
	Plan    workflowPlan         `json:"plan"`
	Steps   []workflowStepStatus `json:"steps"`
	Review  workflowReviewStatus `json:"review"`
	Outcome workflowOutcome      `json:"outcome,omitempty"`
}

type ownedWorkflowRun struct {
	plan           workflowPlan
	steps          []workflowStepStatus
	review         workflowReviewStatus
	outcome        workflowOutcome
	maxOutputBytes int
}

func (o *runtimeOwner) StartWorkflowRun(environmentID, writerOwner, goal string, steps []workflowStepRequest, verifierIDs []string, maxOutputBytes int) (agentRunStatus, error) {
	if o == nil || o.service == nil {
		return agentRunStatus{}, fmt.Errorf("runtime owner is not initialized")
	}
	if _, err := o.service.Environments.RequireWriter(environmentID, writerOwner); err != nil {
		return agentRunStatus{}, err
	}
	rt, _, err := o.service.Runtime(environmentID)
	if err != nil {
		return agentRunStatus{}, err
	}
	runCtx, cancel := context.WithCancel(o.processContext())
	plan, err := materializeWorkflowPlan(runCtx, rt, goal, steps, verifierIDs)
	if err != nil {
		cancel()
		return agentRunStatus{}, err
	}
	if maxOutputBytes <= 0 {
		maxOutputBytes = defaultWorkflowOutput
	}
	runID, err := identity.New("run")
	if err != nil {
		cancel()
		return agentRunStatus{}, err
	}
	stepStatuses := make([]workflowStepStatus, len(plan.Steps))
	for i, step := range plan.Steps {
		stepStatuses[i] = workflowStepStatus{ID: step.ID, Name: step.Name, State: workflowStepPending}
	}
	run := &ownedAgentRun{
		id:            runID,
		environmentID: environmentID,
		writerOwner:   writerOwner,
		kind:          agentRunKindWorkflow,
		workflow: &ownedWorkflowRun{
			plan:           plan,
			steps:          stepStatuses,
			review:         workflowReviewStatus{State: workflowReviewPending},
			maxOutputBytes: maxOutputBytes,
		},
		ctx:       runCtx,
		cancel:    cancel,
		done:      make(chan struct{}),
		state:     agentRunRunning,
		startedAt: time.Now().UTC(),
	}
	if err := o.installAgentRun(run); err != nil {
		cancel()
		return agentRunStatus{}, err
	}
	go o.executeWorkflowRun(run)
	go o.heartbeatAgentRun(run)
	return o.agentRunStatus(run), nil
}

func materializeWorkflowPlan(ctx context.Context, rt *runtimepkg.Runtime, goal string, steps []workflowStepRequest, verifierIDs []string) (workflowPlan, error) {
	goal = strings.TrimSpace(goal)
	if goal == "" {
		return workflowPlan{}, fmt.Errorf("workflow goal is required")
	}
	if len(steps) == 0 {
		return workflowPlan{}, fmt.Errorf("workflow requires at least one executor step")
	}
	if len(steps) > maxWorkflowSteps {
		return workflowPlan{}, fmt.Errorf("workflow exceeds %d executor steps", maxWorkflowSteps)
	}
	if len(verifierIDs) == 0 {
		return workflowPlan{}, fmt.Errorf("workflow requires at least one verifier")
	}
	if len(verifierIDs) > maxWorkflowVerifiers {
		return workflowPlan{}, fmt.Errorf("workflow exceeds %d verifiers", maxWorkflowVerifiers)
	}

	plan := workflowPlan{Goal: goal, Steps: make([]workflowPlanStep, 0, len(steps)), VerifierIDs: make([]string, 0, len(verifierIDs))}
	for i, request := range steps {
		name := strings.TrimSpace(request.Name)
		if name == "" {
			return workflowPlan{}, fmt.Errorf("workflow step %d name is required", i+1)
		}
		if request.TimeoutMS < 0 {
			return workflowPlan{}, fmt.Errorf("workflow step %d timeout_ms must be nonnegative", i+1)
		}
		executable := strings.TrimSpace(request.Executable)
		cwd := strings.TrimSpace(request.Cwd)
		args := append([]string(nil), request.Args...)
		if _, err := rt.PrepareCommand(ctx, executable, args, cwd); err != nil {
			return workflowPlan{}, fmt.Errorf("workflow step %d is not executable: %w", i+1, err)
		}
		plan.Steps = append(plan.Steps, workflowPlanStep{
			ID:         fmt.Sprintf("step_%02d", i+1),
			Name:       name,
			Executable: executable,
			Args:       args,
			Cwd:        cwd,
			TimeoutMS:  request.TimeoutMS,
		})
	}

	seenVerifiers := make(map[string]struct{}, len(verifierIDs))
	for i, verifierID := range verifierIDs {
		verifierID = strings.TrimSpace(verifierID)
		if verifierID == "" {
			return workflowPlan{}, fmt.Errorf("workflow verifier %d id is required", i+1)
		}
		if _, exists := seenVerifiers[verifierID]; exists {
			return workflowPlan{}, fmt.Errorf("workflow verifier %q is duplicated", verifierID)
		}
		seenVerifiers[verifierID] = struct{}{}
		plan.VerifierIDs = append(plan.VerifierIDs, verifierID)
	}
	return plan, nil
}

func (o *runtimeOwner) executeWorkflowRun(run *ownedAgentRun) {
	workflow := run.workflow
	if workflow == nil {
		o.failWorkflowRun(run, "workflow_missing", "workflow run has no workflow state")
		return
	}

	for i, step := range workflow.plan.Steps {
		if run.ctx.Err() != nil {
			o.cancelWorkflowRun(run, -1)
			return
		}
		startedAt := time.Now().UTC()
		run.mu.Lock()
		workflow.steps[i].State = workflowStepRunning
		workflow.steps[i].StartedAt = &startedAt
		run.mu.Unlock()

		if _, writerErr := o.service.Environments.RequireWriter(run.environmentID, run.writerOwner); writerErr != nil {
			completedAt := time.Now().UTC()
			run.mu.Lock()
			status := &workflow.steps[i]
			status.CompletedAt = &completedAt
			status.State = workflowStepFailed
			status.ErrorKind = "executor_error"
			status.Message = writerErr.Error()
			workflow.review.State = workflowReviewNotRun
			run.state = agentRunFailed
			run.errorKind = "executor_error"
			run.message = fmt.Sprintf("workflow step %q writer validation failed: %s", step.ID, writerErr.Error())
			run.completedAt = &completedAt
			run.mu.Unlock()
			finishWorkflowRun(run)
			return
		}
		stepRuntime, _, runtimeErr := o.service.Runtime(run.environmentID)
		if runtimeErr != nil {
			completedAt := time.Now().UTC()
			run.mu.Lock()
			status := &workflow.steps[i]
			status.CompletedAt = &completedAt
			status.State = workflowStepFailed
			status.ErrorKind = "executor_error"
			status.Message = runtimeErr.Error()
			workflow.review.State = workflowReviewNotRun
			run.state = agentRunFailed
			run.errorKind = "executor_error"
			run.message = fmt.Sprintf("workflow step %q runtime validation failed: %s", step.ID, runtimeErr.Error())
			run.completedAt = &completedAt
			run.mu.Unlock()
			finishWorkflowRun(run)
			return
		}
		result, err := stepRuntime.Exec(run.ctx, step.Executable, step.Args, step.Cwd, step.TimeoutMS, workflow.maxOutputBytes)
		completedAt := time.Now().UTC()
		run.mu.Lock()
		status := &workflow.steps[i]
		status.CompletedAt = &completedAt
		status.Stdout = result.Stdout
		status.Stderr = result.Stderr
		canceled := run.cancelRequested || run.ctx.Err() == context.Canceled
		if canceled {
			status.State = workflowStepCanceled
			workflow.review.State = workflowReviewCanceled
			run.state = agentRunCanceled
			run.completedAt = &completedAt
			run.mu.Unlock()
			finishWorkflowRun(run)
			return
		}
		if err != nil {
			status.State = workflowStepFailed
			status.ErrorKind = "executor_error"
			if strings.Contains(err.Error(), "timed out after") {
				status.ErrorKind = "executor_timeout"
			}
			status.Message = err.Error()
			workflow.review.State = workflowReviewNotRun
			run.state = agentRunFailed
			run.errorKind = status.ErrorKind
			run.message = fmt.Sprintf("workflow step %q failed: %s", step.ID, err.Error())
			run.completedAt = &completedAt
			run.mu.Unlock()
			finishWorkflowRun(run)
			return
		}
		exitCode := result.ExitCode
		status.ExitCode = &exitCode
		if result.ExitCode != 0 {
			status.State = workflowStepFailed
			status.ErrorKind = "executor_step_failed"
			status.Message = fmt.Sprintf("step exited with code %d", result.ExitCode)
			workflow.review.State = workflowReviewNotRun
			run.state = agentRunFailed
			run.errorKind = "executor_step_failed"
			run.message = fmt.Sprintf("workflow step %q exited with code %d", step.ID, result.ExitCode)
			run.completedAt = &completedAt
			run.mu.Unlock()
			finishWorkflowRun(run)
			return
		}
		status.State = workflowStepSucceeded
		run.mu.Unlock()
	}

	run.mu.Lock()
	workflow.review.State = workflowReviewRunning
	run.mu.Unlock()
	rejected := false
	for _, verifierID := range workflow.plan.VerifierIDs {
		result, err := o.service.RunVerifier(run.ctx, run.environmentID, run.writerOwner, verifierID, workflow.maxOutputBytes)
		now := time.Now().UTC()
		run.mu.Lock()
		canceled := run.cancelRequested || run.ctx.Err() == context.Canceled
		if canceled {
			workflow.review.State = workflowReviewCanceled
			run.state = agentRunCanceled
			run.completedAt = &now
			run.mu.Unlock()
			finishWorkflowRun(run)
			return
		}
		if err != nil {
			workflow.review.State = workflowReviewError
			workflow.review.ErrorKind = "reviewer_error"
			workflow.review.Message = err.Error()
			run.state = agentRunFailed
			run.errorKind = "reviewer_error"
			run.message = fmt.Sprintf("workflow reviewer %q failed: %s", verifierID, err.Error())
			run.completedAt = &now
			run.mu.Unlock()
			finishWorkflowRun(run)
			return
		}
		workflow.review.Verifiers = append(workflow.review.Verifiers, result)
		if result.Status != verifier.StatusPassed {
			rejected = true
		}
		run.mu.Unlock()
	}

	now := time.Now().UTC()
	run.mu.Lock()
	if rejected {
		workflow.review.State = workflowReviewRejected
		workflow.outcome = workflowOutcomeRejected
	} else {
		workflow.review.State = workflowReviewAccepted
		workflow.outcome = workflowOutcomeAccepted
	}
	run.state = agentRunSucceeded
	run.completedAt = &now
	run.mu.Unlock()
	finishWorkflowRun(run)
}

func (o *runtimeOwner) failWorkflowRun(run *ownedAgentRun, errorKind, message string) {
	now := time.Now().UTC()
	run.mu.Lock()
	run.state = agentRunFailed
	run.errorKind = errorKind
	run.message = message
	run.completedAt = &now
	if run.workflow != nil {
		run.workflow.review.State = workflowReviewNotRun
	}
	run.mu.Unlock()
	finishWorkflowRun(run)
}

func (o *runtimeOwner) cancelWorkflowRun(run *ownedAgentRun, stepIndex int) {
	now := time.Now().UTC()
	run.mu.Lock()
	if run.workflow != nil {
		if stepIndex >= 0 && stepIndex < len(run.workflow.steps) {
			run.workflow.steps[stepIndex].State = workflowStepCanceled
			run.workflow.steps[stepIndex].CompletedAt = &now
		}
		run.workflow.review.State = workflowReviewCanceled
	}
	run.state = agentRunCanceled
	run.completedAt = &now
	run.mu.Unlock()
	finishWorkflowRun(run)
}

func finishWorkflowRun(run *ownedAgentRun) {
	if run.cancel != nil {
		run.cancel()
	}
	close(run.done)
}

func workflowRunStatusLocked(workflow *ownedWorkflowRun) workflowRunStatus {
	status := workflowRunStatus{
		Plan:    cloneWorkflowPlan(workflow.plan),
		Steps:   make([]workflowStepStatus, len(workflow.steps)),
		Review:  cloneWorkflowReview(workflow.review),
		Outcome: workflow.outcome,
	}
	for i := range workflow.steps {
		status.Steps[i] = cloneWorkflowStepStatus(workflow.steps[i])
	}
	return status
}

func cloneWorkflowPlan(plan workflowPlan) workflowPlan {
	clone := workflowPlan{Goal: plan.Goal, VerifierIDs: append([]string(nil), plan.VerifierIDs...), Steps: make([]workflowPlanStep, len(plan.Steps))}
	for i, step := range plan.Steps {
		step.Args = append([]string(nil), step.Args...)
		clone.Steps[i] = step
	}
	return clone
}

func cloneWorkflowStepStatus(status workflowStepStatus) workflowStepStatus {
	status.StartedAt = cloneTime(status.StartedAt)
	status.CompletedAt = cloneTime(status.CompletedAt)
	if status.ExitCode != nil {
		exitCode := *status.ExitCode
		status.ExitCode = &exitCode
	}
	return status
}

func cloneWorkflowReview(review workflowReviewStatus) workflowReviewStatus {
	review.Verifiers = append([]verifier.Result(nil), review.Verifiers...)
	return review
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	clone := *value
	return &clone
}

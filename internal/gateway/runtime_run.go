package gateway

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"ai-dev-manager-v2/internal/identity"
	runtimepkg "ai-dev-manager-v2/internal/runtime"
)

const agentRunStopTimeout = 7 * time.Second

type agentRunKind string

const (
	agentRunKindCommand  agentRunKind = "command"
	agentRunKindWorkflow agentRunKind = "workflow"
)

type agentRunState string

const (
	agentRunRunning   agentRunState = "running"
	agentRunSucceeded agentRunState = "succeeded"
	agentRunFailed    agentRunState = "failed"
	agentRunCanceled  agentRunState = "canceled"
)

type agentRunStatus struct {
	ID            string             `json:"id"`
	EnvironmentID string             `json:"environment_id"`
	Kind          agentRunKind       `json:"kind"`
	State         agentRunState      `json:"state"`
	Executable    string             `json:"executable,omitempty"`
	Args          []string           `json:"args,omitempty"`
	Cwd           string             `json:"cwd,omitempty"`
	StartedAt     time.Time          `json:"started_at"`
	CompletedAt   *time.Time         `json:"completed_at,omitempty"`
	ExitCode      *int               `json:"exit_code,omitempty"`
	Stdout        string             `json:"stdout,omitempty"`
	Stderr        string             `json:"stderr,omitempty"`
	ErrorKind     string             `json:"error_kind,omitempty"`
	Message       string             `json:"message,omitempty"`
	Workflow      *workflowRunStatus `json:"workflow,omitempty"`
}

type ownedAgentRun struct {
	id            string
	environmentID string
	writerOwner   string
	kind          agentRunKind
	executable    string
	args          []string
	cwd           string
	timeoutMS     int64
	maxOutput     int
	workflow      *ownedWorkflowRun
	ctx           context.Context
	cancel        context.CancelFunc
	done          chan struct{}

	mu              sync.Mutex
	state           agentRunState
	startedAt       time.Time
	completedAt     *time.Time
	result          runtimepkg.CommandResult
	hasResult       bool
	hasExitCode     bool
	cancelRequested bool
	errorKind       string
	message         string
}

func (o *runtimeOwner) StartAgentRun(environmentID, writerOwner, executable string, args []string, cwd string, timeoutMS int64, maxOutputBytes int) (agentRunStatus, error) {
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
	if _, err := rt.PrepareCommand(runCtx, executable, args, cwd); err != nil {
		cancel()
		return agentRunStatus{}, err
	}
	runID, err := identity.New("run")
	if err != nil {
		cancel()
		return agentRunStatus{}, err
	}
	run := &ownedAgentRun{
		id:            runID,
		environmentID: environmentID,
		writerOwner:   writerOwner,
		kind:          agentRunKindCommand,
		executable:    executable,
		args:          append([]string(nil), args...),
		cwd:           cwd,
		timeoutMS:     timeoutMS,
		maxOutput:     maxOutputBytes,
		ctx:           runCtx,
		cancel:        cancel,
		done:          make(chan struct{}),
		state:         agentRunRunning,
		startedAt:     time.Now().UTC(),
	}
	if err := o.installAgentRun(run); err != nil {
		cancel()
		return agentRunStatus{}, err
	}
	go o.executeAgentRun(run, rt)
	go o.heartbeatAgentRun(run)
	return o.agentRunStatus(run), nil
}

func (o *runtimeOwner) ListAgentRuns(environmentID string) ([]agentRunStatus, error) {
	if _, err := o.service.Environments.Get(environmentID); err != nil {
		return nil, err
	}
	o.mu.Lock()
	items := make([]*ownedAgentRun, 0)
	for _, run := range o.runs {
		if run.environmentID == environmentID {
			items = append(items, run)
		}
	}
	o.mu.Unlock()
	statuses := make([]agentRunStatus, 0, len(items))
	for _, run := range items {
		statuses = append(statuses, o.agentRunStatus(run))
	}
	sort.Slice(statuses, func(i, j int) bool { return statuses[i].StartedAt.Before(statuses[j].StartedAt) })
	return statuses, nil
}

func (o *runtimeOwner) AgentRunStatus(environmentID, runID string) (agentRunStatus, error) {
	run, err := o.agentRun(environmentID, runID)
	if err != nil {
		return agentRunStatus{}, err
	}
	return o.agentRunStatus(run), nil
}

func (o *runtimeOwner) CancelAgentRun(environmentID, writerOwner, runID string) (agentRunStatus, error) {
	if _, err := o.service.Environments.RequireWriter(environmentID, writerOwner); err != nil {
		return agentRunStatus{}, err
	}
	run, err := o.agentRun(environmentID, runID)
	if err != nil {
		return agentRunStatus{}, err
	}
	if run.writerOwner != writerOwner {
		return agentRunStatus{}, fmt.Errorf("run %q belongs to writer %q", runID, run.writerOwner)
	}
	o.requestAgentRunCancel(run, "", "")
	if err := waitOwnedAgentRun(run, agentRunStopTimeout); err != nil {
		return agentRunStatus{}, err
	}
	return o.agentRunStatus(run), nil
}

func (o *runtimeOwner) executeAgentRun(run *ownedAgentRun, rt *runtimepkg.Runtime) {
	result, err := rt.Exec(run.ctx, run.executable, run.args, run.cwd, run.timeoutMS, run.maxOutput)
	now := time.Now().UTC()

	run.mu.Lock()
	run.result = result
	run.hasResult = true
	canceled := run.cancelRequested || run.ctx.Err() == context.Canceled
	switch {
	case canceled:
		run.state = agentRunCanceled
	case err != nil:
		run.state = agentRunFailed
		if run.errorKind == "" {
			if strings.Contains(err.Error(), "timed out after") {
				run.errorKind = "timeout"
			} else {
				run.errorKind = "run_exec_failed"
			}
		}
		if run.message == "" {
			run.message = err.Error()
		}
	case result.ExitCode != 0:
		run.state = agentRunFailed
		run.hasExitCode = true
		if run.errorKind == "" {
			run.errorKind = "command_failed"
		}
	default:
		run.state = agentRunSucceeded
		run.hasExitCode = true
	}
	run.completedAt = &now
	run.mu.Unlock()

	run.cancel()
	close(run.done)
}

func (o *runtimeOwner) heartbeatAgentRun(run *ownedAgentRun) {
	interval := o.service.Environments.WriterLeaseTTL() / 3
	if interval <= 0 {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-run.done:
			return
		case <-ticker.C:
			if _, err := o.service.Environments.HeartbeatWriter(run.environmentID, run.writerOwner); err != nil {
				o.requestAgentRunCancel(run, "writer_heartbeat_failed", err.Error())
				return
			}
		}
	}
}

func (o *runtimeOwner) installAgentRun(run *ownedAgentRun) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return fmt.Errorf("runtime owner is closed")
	}
	if o.runs == nil {
		o.runs = map[string]*ownedAgentRun{}
	}
	o.runs[run.id] = run
	return nil
}

func (o *runtimeOwner) agentRun(environmentID, runID string) (*ownedAgentRun, error) {
	o.mu.Lock()
	run := o.runs[runID]
	o.mu.Unlock()
	if run == nil || run.environmentID != environmentID {
		return nil, fmt.Errorf("run %q not found for environment %q", runID, environmentID)
	}
	return run, nil
}

func (o *runtimeOwner) agentRunStatus(run *ownedAgentRun) agentRunStatus {
	run.mu.Lock()
	defer run.mu.Unlock()
	status := agentRunStatus{
		ID:            run.id,
		EnvironmentID: run.environmentID,
		Kind:          run.kind,
		State:         run.state,
		Executable:    run.executable,
		Args:          append([]string(nil), run.args...),
		Cwd:           run.cwd,
		StartedAt:     run.startedAt,
		CompletedAt:   run.completedAt,
		ErrorKind:     run.errorKind,
		Message:       run.message,
	}
	if run.hasResult {
		status.Stdout = run.result.Stdout
		status.Stderr = run.result.Stderr
	}
	if run.hasExitCode {
		exitCode := run.result.ExitCode
		status.ExitCode = &exitCode
	}
	if run.workflow != nil {
		workflow := workflowRunStatusLocked(run.workflow)
		status.Workflow = &workflow
	}
	return status
}

func (o *runtimeOwner) requestAgentRunCancel(run *ownedAgentRun, errorKind, message string) {
	run.mu.Lock()
	if run.state != agentRunRunning {
		run.mu.Unlock()
		return
	}
	run.cancelRequested = true
	if errorKind != "" && run.errorKind == "" {
		run.errorKind = errorKind
	}
	if message != "" && run.message == "" {
		run.message = message
	}
	cancel := run.cancel
	run.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (o *runtimeOwner) dropAgentRunsForEnvironment(environmentID string) {
	o.mu.Lock()
	var runs []*ownedAgentRun
	for runID, run := range o.runs {
		if run.environmentID == environmentID {
			runs = append(runs, run)
			delete(o.runs, runID)
		}
	}
	o.mu.Unlock()
	for _, run := range runs {
		o.requestAgentRunCancel(run, "", "")
		_ = waitOwnedAgentRun(run, agentRunStopTimeout)
	}
}

func (o *runtimeOwner) closeAgentRuns(runs []*ownedAgentRun) error {
	var errs []string
	for _, run := range runs {
		o.requestAgentRunCancel(run, "", "")
	}
	for _, run := range runs {
		if err := waitOwnedAgentRun(run, agentRunStopTimeout); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("agent run cleanup failed: %s", strings.Join(errs, "; "))
	}
	return nil
}

func waitOwnedAgentRun(run *ownedAgentRun, timeout time.Duration) error {
	select {
	case <-run.done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("run %q did not stop within %s", run.id, timeout)
	}
}

func (o *runtimeOwner) runningAgentRunCountLocked() int {
	count := 0
	for _, run := range o.runs {
		run.mu.Lock()
		running := run.state == agentRunRunning
		run.mu.Unlock()
		if running {
			count++
		}
	}
	return count
}

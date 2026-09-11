package gateway

import (
	"context"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/model"
)

func (o *runtimeOwner) ResourceRetentionReport(ctx context.Context) (model.ResourceRetentionReport, error) {
	report, err := o.service.ResourceRetentionReport(ctx)
	if err != nil {
		return model.ResourceRetentionReport{}, err
	}
	now := time.Now().UTC()
	for i := range report.Resources {
		item := &report.Resources[i]
		if !hasRetentionBlocker(item.Blockers, app.RetentionRuntimeObservationBlocker) {
			continue
		}
		item.Blockers = removeRetentionBlocker(item.Blockers, app.RetentionRuntimeObservationBlocker)
		item.Blockers = append(item.Blockers, o.runtimeRetentionBlockers(item.Kind, item.ID)...)
		item.Blockers = uniqueGatewayStrings(item.Blockers)
		refreshRetentionCleanupState(item, now)
	}
	report.GeneratedAt = now
	return report, nil
}

func (o *runtimeOwner) ResourceRetentionCleanup(ctx context.Context, request model.ResourceRetentionCleanupRequest) (model.ResourceRetentionCleanupResult, error) {
	report, err := o.ResourceRetentionReport(ctx)
	if err != nil {
		return model.ResourceRetentionCleanupResult{}, err
	}
	if !request.Execute {
		return o.service.ResourceRetentionCleanupFromReport(report, request)
	}
	result, err := o.service.ResourceRetentionCleanupFromRuntimeReport(ctx, report, request)
	if err != nil {
		return model.ResourceRetentionCleanupResult{}, err
	}
	for _, mutation := range result.Removed {
		if mutation.Kind == model.RetentionResourceMCP {
			o.DropMCP(mutation.ID)
		}
	}
	return result, nil
}

func (o *runtimeOwner) runtimeRetentionBlockers(kind, id string) []string {
	o.mu.Lock()
	defer o.mu.Unlock()
	blockers := []string{}
	switch kind {
	case model.RetentionResourceEnvironment:
		for key, session := range o.sessions {
			if key.environmentID == id && session != nil {
				blockers = append(blockers, "active_mcp_session")
				break
			}
		}
		for key, inFlight := range o.inFlight {
			if key.environmentID == id && inFlight {
				blockers = append(blockers, "mcp_operation_in_flight")
				break
			}
		}
		for _, process := range o.processes {
			if process.environmentID != id {
				continue
			}
			process.mu.Lock()
			running := process.state == devProcessRunning
			process.mu.Unlock()
			if running {
				blockers = append(blockers, "active_process")
				break
			}
		}
		for _, run := range o.runs {
			if run.environmentID != id {
				continue
			}
			run.mu.Lock()
			running := run.state == agentRunRunning
			run.mu.Unlock()
			if running {
				blockers = append(blockers, "active_run")
				break
			}
		}
	case model.RetentionResourceMCP:
		for key, session := range o.sessions {
			if key.mcpID == id && session != nil {
				blockers = append(blockers, "active_mcp_session")
				break
			}
		}
		for key, inFlight := range o.inFlight {
			if key.mcpID == id && inFlight {
				blockers = append(blockers, "mcp_operation_in_flight")
				break
			}
		}
	}
	return blockers
}

func refreshRetentionCleanupState(item *model.ResourceRetentionItem, now time.Time) {
	if item == nil {
		return
	}
	item.CleanupEligible = false
	switch item.Retention.Persistence {
	case model.PersistenceDurable:
		item.CleanupState = model.RetentionCleanupDurable
		return
	case model.PersistenceTemporary:
	default:
		item.CleanupState = model.RetentionCleanupBlocked
		return
	}
	if item.Retention.ExpiresAt == nil {
		item.CleanupState = model.RetentionCleanupBlocked
		return
	}
	if now.Before(item.Retention.ExpiresAt.UTC()) {
		item.CleanupState = model.RetentionCleanupNotDue
		return
	}
	if len(item.Blockers) != 0 {
		item.CleanupState = model.RetentionCleanupBlocked
		return
	}
	item.CleanupState = model.RetentionCleanupEligible
	item.CleanupEligible = true
}

func hasRetentionBlocker(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), target) {
			return true
		}
	}
	return false
}

func removeRetentionBlocker(values []string, target string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), target) {
			continue
		}
		result = append(result, value)
	}
	return result
}

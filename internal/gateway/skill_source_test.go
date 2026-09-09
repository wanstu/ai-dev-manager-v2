package gateway

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/app"
)

func TestGatewaySkillSourceLifecyclePreservesUnresolvedSelections(t *testing.T) {
	root := filepath.Join(t.TempDir(), "skills")
	writeGatewaySkill(t, root, "agent-skill", "# agent skill\n")
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ctx := context.Background()
	owner := newRuntimeOwner(service)
	defer owner.Close()
	session := connectInMemory(t, ctx, newServer(service, owner))
	defer session.Close()

	added := callGatewayTool(t, ctx, session, "skill_source_add", map[string]any{
		"root":                           root,
		"default_include_in_environment": true,
	})
	if added.IsError || !strings.Contains(toolText(t, added), "skill_source_id") {
		t.Fatalf("skill_source_add failed: %s", toolText(t, added))
	}
	sources, err := service.Skills.ListSkillSources()
	if err != nil || len(sources) != 1 {
		t.Fatalf("sources=%+v err=%v", sources, err)
	}
	refreshed := callGatewayTool(t, ctx, session, "skill_source_refresh", map[string]any{"id": sources[0].ID})
	if refreshed.IsError || !strings.Contains(toolText(t, refreshed), `"added":1`) || !strings.Contains(toolText(t, refreshed), "agent-skill") {
		t.Fatalf("skill_source_refresh failed: %s", toolText(t, refreshed))
	}
	listedSources := callGatewayTool(t, ctx, session, "skill_source_list", map[string]any{})
	if listedSources.IsError || !strings.Contains(toolText(t, listedSources), sources[0].ID) {
		t.Fatalf("skill_source_list failed: %s", toolText(t, listedSources))
	}
	entries, err := service.Skills.List()
	if err != nil || len(entries) != 1 {
		t.Fatalf("entries=%+v err=%v", entries, err)
	}

	workspace, err := service.Workspaces.Add(t.TempDir(), "gateway-skill-source")
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(workspace.ID, "gateway-skill-source", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(env.EnabledSkillIDs) != 1 || env.EnabledSkillIDs[0] != entries[0].ID {
		t.Fatalf("source default selection missing: %+v", env.EnabledSkillIDs)
	}
	removed := callGatewayTool(t, ctx, session, "skill_source_remove", map[string]any{"id": sources[0].ID})
	if removed.IsError || !strings.Contains(toolText(t, removed), `"removed":1`) {
		t.Fatalf("skill_source_remove failed: %s", toolText(t, removed))
	}
	shown := callGatewayTool(t, ctx, session, "environment_skill_list", map[string]any{"environment_id": env.ID})
	if shown.IsError || !strings.Contains(toolText(t, shown), entries[0].ID) || !strings.Contains(toolText(t, shown), `"state":"unresolved"`) {
		t.Fatalf("removed Skill selection was not preserved as unresolved: %s", toolText(t, shown))
	}
}

func writeGatewaySkill(t *testing.T, root, name, content string) {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

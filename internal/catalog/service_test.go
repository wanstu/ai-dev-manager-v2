package catalog_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/skill"
	"ai-dev-manager-v2/internal/store"
)

func TestConfiguredCatalogEntriesRequireRuntimeSources(t *testing.T) {
	state := store.New(filepath.Join(t.TempDir(), "state.json"))
	mcps := catalog.New(state, catalog.KindMCP)
	skills := catalog.New(state, catalog.KindSkill)

	for _, endpoint := range []string{"", "filesystem", "ftp://example.test/mcp"} {
		if _, err := mcps.AddMCP("invalid-"+strings.ReplaceAll(endpoint, ":", "-"), endpoint, false); err == nil || !strings.Contains(err.Error(), "valid http(s) URL") {
			t.Fatalf("AddMCP(%q) error = %v", endpoint, err)
		}
	}
	mcpEntry, err := mcps.AddMCP("external", "http://127.0.0.1:9000/mcp", true)
	if err != nil {
		t.Fatal(err)
	}
	if mcpEntry.Endpoint != "http://127.0.0.1:9000/mcp" || !mcpEntry.DefaultIncludeInEnv {
		t.Fatalf("configured MCP = %+v", mcpEntry)
	}

	if entries, err := skills.List(); err != nil || len(entries) != 0 {
		t.Fatalf("Skill catalog must not scan host paths before explicit root configuration: entries=%+v err=%v", entries, err)
	}
	skillRoot := filepath.Join(t.TempDir(), "skills")
	supportRoot := filepath.Join(t.TempDir(), "gsd-core")
	if err := os.MkdirAll(filepath.Join(skillRoot, "gsd-next"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(supportRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	artifact := filepath.Join(skillRoot, "gsd-next", "SKILL.md")
	if err := os.WriteFile(artifact, []byte("# gsd-next\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	discovered, err := skills.AddSkillRoot(skillRoot, []string{supportRoot}, true)
	if err != nil || len(discovered) != 1 {
		t.Fatalf("AddSkillRoot entries=%+v err=%v", discovered, err)
	}
	entry := discovered[0]
	if entry.ID != skill.StableID("gsd-next") || entry.Name != "gsd-next" || entry.ArtifactPath != artifact || entry.SourceRoot != skillRoot || !entry.DefaultIncludeInEnv {
		t.Fatalf("configured Skill = %+v", entry)
	}
	if len(entry.SupportRoots) != 1 || entry.SupportRoots[0] != supportRoot {
		t.Fatalf("configured Skill support roots = %+v", entry.SupportRoots)
	}
}

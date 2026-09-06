package catalog_test

import (
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/store"
)

func TestConfiguredCatalogEntriesRequireRuntimeContent(t *testing.T) {
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

	if _, err := skills.AddSkill("blank", "   ", false); err == nil || !strings.Contains(err.Error(), "instructions are required") {
		t.Fatalf("blank Skill instructions error = %v", err)
	}
	skillEntry, err := skills.AddSkill("review", "Review the change and run focused tests.", false)
	if err != nil {
		t.Fatal(err)
	}
	if skillEntry.Instructions != "Review the change and run focused tests." {
		t.Fatalf("configured Skill = %+v", skillEntry)
	}
}

package catalog

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/store"
)

func TestPreviewMCPImportOpenCodeJSONC(t *testing.T) {
	preview, err := PreviewMCPImport(MCPImportPreviewRequest{
		Format: MCPImportFormatOpenCode,
		Content: `{
			// OpenCode-style config
			"mcp": {"servers": {
				"local fs": {
					"type": "local",
					"command": ["node", "server.js", "--root", "."],
					"environment": {"DEBUG": "1", "API_TOKEN": "plain-secret-token"},
					"disabled": true
				},
				"remote": {
					"type": "remote",
					"url": "http://127.0.0.1:9000/mcp",
					"headers": {"Authorization": "Bearer {env:REMOTE_TOKEN}"}
				}
			}}
		}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if preview.Format != MCPImportFormatOpenCode || len(preview.Candidates) != 2 {
		t.Fatalf("preview=%+v", preview)
	}
	local := preview.Candidates[0]
	if local.Name != "local fs" || local.Definition.Transport != MCPTransportStdio || local.Definition.Executable != "node" || len(local.Definition.Args) != 3 {
		t.Fatalf("local candidate=%+v", local)
	}
	if local.Definition.EnvRefs["DEBUG"] != "1" {
		t.Fatalf("non-credential env literal should be preserved: %+v", local.Definition.EnvRefs)
	}
	if local.Definition.EnvRefs["API_TOKEN"] != "${LOCAL_FS_API_TOKEN}" || len(local.ReferenceRequirements) != 1 {
		t.Fatalf("credential literal was not converted: refs=%+v requirements=%+v", local.Definition.EnvRefs, local.ReferenceRequirements)
	}
	if len(local.Warnings) != 1 || local.Warnings[0].Code != "source_disabled_ignored" {
		t.Fatalf("source disabled warning missing: %+v", local.Warnings)
	}
	remote := preview.Candidates[1]
	if remote.Name != "remote" || remote.Definition.Transport != MCPTransportStreamableHTTP || remote.Definition.AuthMode != MCPAuthHeaders {
		t.Fatalf("remote candidate=%+v", remote)
	}
	if remote.Definition.HeaderRefs["Authorization"] != "Bearer ${REMOTE_TOKEN}" {
		t.Fatalf("source env reference not normalized: %+v", remote.Definition.HeaderRefs)
	}
	if text := previewText(preview); strings.Contains(text, "plain-secret-token") {
		t.Fatalf("preview leaked literal secret: %s", text)
	}
}

func TestPreviewMCPImportCodexWrapperAndDirectMap(t *testing.T) {
	wrapper, err := PreviewMCPImport(MCPImportPreviewRequest{
		Format:         MCPImportFormatCodexPlugin,
		Content:        `{"mcpServers":{"stdio":{"command":"python","args":["-m","server"],"env":{"TOKEN":"${PY_TOKEN}"}},"http":{"url":"https://example.com/mcp"}}}`,
		DefaultInclude: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wrapper.Candidates) != 2 {
		t.Fatalf("wrapper preview=%+v", wrapper)
	}
	if !wrapper.Candidates[0].Definition.DefaultIncludeInEnv || !wrapper.Candidates[1].Definition.DefaultIncludeInEnv {
		t.Fatalf("default include option not applied: %+v", wrapper.Candidates)
	}

	direct, err := PreviewMCPImport(MCPImportPreviewRequest{
		Format:  MCPImportFormatCodexPlugin,
		Content: `{"direct":{"command":"go","args":["run","./cmd/mcp"]}}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(direct.Candidates) != 1 || direct.Candidates[0].Definition.Executable != "go" {
		t.Fatalf("direct preview=%+v", direct)
	}
}

func TestPreviewMCPImportWorkBuddyCodeBuddyClaudeAndMCPHub(t *testing.T) {
	workbuddy, err := PreviewMCPImport(MCPImportPreviewRequest{
		Format:  MCPImportFormatWorkBuddy,
		Content: `{"mcpServers":{"wb-http":{"type":"streamableHttp","url":"https://workbuddy.example/mcp","headers":{"Authorization":"Bearer {env:WB_TOKEN}"},"x-workbuddy":{"package":"ignored"}},"wb-stdio":{"command":"node","args":["server.js"],"env":{"DEBUG":"1"}}}}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(workbuddy.Candidates) != 2 || workbuddy.Candidates[0].Format != MCPImportFormatWorkBuddy {
		t.Fatalf("workbuddy preview=%+v", workbuddy)
	}
	if workbuddy.Candidates[0].Definition.HeaderRefs["Authorization"] != "Bearer ${WB_TOKEN}" {
		t.Fatalf("workbuddy header ref=%+v", workbuddy.Candidates[0].Definition.HeaderRefs)
	}
	if len(workbuddy.Candidates[0].Warnings) == 0 || workbuddy.Candidates[0].Warnings[0].Code != "unsupported_extension" {
		t.Fatalf("workbuddy extension warning missing: %+v", workbuddy.Candidates[0].Warnings)
	}

	codebuddy, err := PreviewMCPImport(MCPImportPreviewRequest{
		Format:  MCPImportFormatCodeBuddy,
		Content: `{"mcpServers":{"cb":{"command":["python","-m","cb"],"env":{"API_KEY":"literal-codebuddy-key"}}}}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(codebuddy.Candidates) != 1 || codebuddy.Candidates[0].Definition.Executable != "python" || codebuddy.Candidates[0].Definition.EnvRefs["API_KEY"] != "${CB_API_KEY}" {
		t.Fatalf("codebuddy preview=%+v", codebuddy)
	}
	if strings.Contains(previewText(codebuddy), "literal-codebuddy-key") {
		t.Fatalf("codebuddy preview leaked secret: %s", previewText(codebuddy))
	}

	claude, err := PreviewMCPImport(MCPImportPreviewRequest{
		Format:  MCPImportFormatClaudeCode,
		Content: `{"projects":{"/repo":{"mcpServers":{"claude-fs":{"command":"npx","args":["-y","@modelcontextprotocol/server-filesystem"],"env":{"TOKEN":"${CLAUDE_TOKEN:-fallback}"}}}}}}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(claude.Candidates) != 1 || !strings.Contains(claude.Candidates[0].SourcePath, "projects./repo.mcpServers") || claude.Candidates[0].Definition.EnvRefs["TOKEN"] != "${CLAUDE_TOKEN:-fallback}" {
		t.Fatalf("claude preview=%+v", claude)
	}
	_, err = PreviewMCPImport(MCPImportPreviewRequest{
		Format:  MCPImportFormatClaudeCode,
		Content: `{"projects":{"/a":{"mcpServers":{"one":{"command":"node"}}},"/b":{"mcpServers":{"two":{"command":"node"}}}}}`,
	})
	if err == nil || !strings.Contains(err.Error(), "scope_selector_required") {
		t.Fatalf("multi-project claude import should require scope selector, got %v", err)
	}

	mcphub, err := PreviewMCPImport(MCPImportPreviewRequest{
		Format:  MCPImportFormatMCPHub,
		Content: `{"servers":{"hub":{"url":"https://hub.example/mcp","enabled":true,"routing":{"tag":"ignored"}}}}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(mcphub.Candidates) != 1 || mcphub.Candidates[0].Definition.Endpoint != "https://hub.example/mcp" {
		t.Fatalf("mcphub preview=%+v", mcphub)
	}
	if len(mcphub.Candidates[0].Warnings) < 2 {
		t.Fatalf("mcphub warnings missing enabled/extension facts: %+v", mcphub.Candidates[0].Warnings)
	}
}

func TestPreviewMCPImportAutoDetectsSourceHints(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string
		want    string
	}{
		{name: "workbuddy", content: `{"x-workbuddy":true,"mcpServers":{"wb":{"command":"node"}}}`, want: MCPImportFormatWorkBuddy},
		{name: "claude", content: `{"projects":{"/repo":{"mcpServers":{"fs":{"command":"node"}}}}}`, want: MCPImportFormatClaudeCode},
		{name: "mcphub", content: `{"servers":{"hub":{"command":"node"}}}`, want: MCPImportFormatMCPHub},
	} {
		t.Run(tc.name, func(t *testing.T) {
			preview, err := PreviewMCPImport(MCPImportPreviewRequest{Format: MCPImportFormatAuto, Content: tc.content})
			if err != nil {
				t.Fatal(err)
			}
			if preview.Format != tc.want {
				t.Fatalf("format=%q want %q", preview.Format, tc.want)
			}
		})
	}
}

func TestPreviewMCPImportAutoDetectionAndAmbiguous(t *testing.T) {
	auto, err := PreviewMCPImport(MCPImportPreviewRequest{
		Format:  MCPImportFormatAuto,
		Content: `{"mcpServers":{"stdio":{"command":"node"}}}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if auto.Format != MCPImportFormatCodexPlugin {
		t.Fatalf("auto format=%q", auto.Format)
	}

	_, err = PreviewMCPImport(MCPImportPreviewRequest{
		Format:  MCPImportFormatAuto,
		Content: `{"mcp":{"servers":{"opencode":{"type":"local","command":["node"]}}},"mcpServers":{"codex":{"command":"python"}}}`,
	})
	if err == nil || !strings.Contains(err.Error(), "ambiguous_format") {
		t.Fatalf("expected ambiguous_format, got %v", err)
	}
}

func TestApplyMCPImportConflictPoliciesAndAtomicity(t *testing.T) {
	service := NewMCP(store.New(filepath.Join(t.TempDir(), "state.json")))
	content := `{
		"mcpServers": {
			"alpha": {"command": "alpha-mcp", "args": ["serve"]},
			"beta": {"url": "http://127.0.0.1:9010/mcp"}
		}
	}`
	if _, err := service.ApplyMCPImport(MCPImportApplyRequest{Content: content}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.ApplyMCPImport(MCPImportApplyRequest{Content: content}); err == nil || !strings.Contains(err.Error(), "conflict") {
		t.Fatalf("default conflict policy must reject existing names, got %v", err)
	}
	items, err := service.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("conflict error must not persist partial changes: %+v", items)
	}

	skipContent := `{"mcpServers":{"alpha":{"command":"ignored"},"gamma":{"command":"gamma-mcp"}}}`
	skipped, err := service.ApplyMCPImport(MCPImportApplyRequest{Content: skipContent, ConflictPolicy: MCPImportConflictSkip})
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped.Imported) != 1 || skipped.Imported[0].Name != "gamma" || len(skipped.Skipped) != 1 || skipped.Skipped[0].Name != "alpha" {
		t.Fatalf("skip result = %+v", skipped)
	}
	items, err = service.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Fatalf("skip policy should import only non-conflicting entries: %+v", items)
	}

	var alphaBefore model.MCPDefinition
	for _, item := range items {
		if item.Name == "alpha" {
			alphaBefore = item
		}
	}
	updated, err := service.ApplyMCPImport(MCPImportApplyRequest{
		Content:        `{"mcpServers":{"alpha":{"url":"http://127.0.0.1:9020/mcp"}}}`,
		ConflictPolicy: MCPImportConflictUpdateByName,
		SelectedNames:  []string{"alpha"},
		DefaultInclude: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.Updated) != 1 || updated.Updated[0].ID != alphaBefore.ID || updated.Updated[0].Transport != MCPTransportStreamableHTTP || !updated.Updated[0].DefaultIncludeInEnv {
		t.Fatalf("update_by_name result = %+v alphaBefore=%+v", updated, alphaBefore)
	}

	invalidBatch := `{"mcpServers":{"delta":{"command":"delta-mcp"},"bad":{"transport":"sse","url":"http://127.0.0.1:9999/sse"}}}`
	if _, err := service.ApplyMCPImport(MCPImportApplyRequest{Content: invalidBatch}); err == nil || !strings.Contains(err.Error(), "invalid") {
		t.Fatalf("invalid selected batch must fail atomically, got %v", err)
	}
	items, err = service.List()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range items {
		if item.Name == "delta" || item.Name == "bad" {
			t.Fatalf("invalid batch persisted partial item: %+v", items)
		}
	}
}

func TestApplyMCPImportConvertsLiteralCredentialsWithoutPersistenceLeak(t *testing.T) {
	service := NewMCP(store.New(filepath.Join(t.TempDir(), "state.json")))
	secret := "plain-secret-token"
	applied, err := service.ApplyMCPImport(MCPImportApplyRequest{Content: `{
		"mcpServers": {
			"remote": {"url": "http://127.0.0.1:9010/mcp", "headers": {"Authorization": "Bearer plain-secret-token"}}
		}
	}`})
	if err != nil {
		t.Fatal(err)
	}
	if len(applied.Imported) != 1 || len(applied.ReferenceRequirements) != 1 {
		t.Fatalf("apply result = %+v", applied)
	}
	definition := applied.Imported[0]
	if definition.HeaderRefs["Authorization"] != "Bearer ${REMOTE_AUTHORIZATION}" {
		t.Fatalf("literal Authorization not converted to generated reference: %+v", definition.HeaderRefs)
	}
	encoded, err := json.Marshal(applied)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), secret) {
		t.Fatalf("apply result leaked secret: %s", encoded)
	}
	items, err := service.List()
	if err != nil {
		t.Fatal(err)
	}
	persisted, err := json.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(persisted), secret) {
		t.Fatalf("persisted MCP leaked secret: %s", persisted)
	}
}

func TestPreviewMCPImportRejectsUnsupportedSSE(t *testing.T) {
	preview, err := PreviewMCPImport(MCPImportPreviewRequest{
		Format:  MCPImportFormatCodexPlugin,
		Content: `{"mcpServers":{"legacy":{"transport":"sse","url":"http://127.0.0.1:9000/sse"}}}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.Candidates) != 1 || len(preview.Candidates[0].Errors) == 0 || preview.Candidates[0].Errors[0].Code != "unsupported_transport" {
		t.Fatalf("legacy sse should be candidate error: %+v", preview)
	}
}

func previewText(value any) string {
	return fmt.Sprintf("%+v", value)
}

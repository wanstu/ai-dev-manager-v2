package catalog

import (
	"fmt"
	"strings"
	"testing"
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

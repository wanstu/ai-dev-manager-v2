package catalog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"ai-dev-manager-v2/internal/identity"
	"ai-dev-manager-v2/internal/model"
)

const (
	MCPImportFormatAuto        = "auto"
	MCPImportFormatOpenCode    = "opencode"
	MCPImportFormatCodexPlugin = "codex-plugin"

	MCPImportConflictError        = "error"
	MCPImportConflictSkip         = "skip"
	MCPImportConflictUpdateByName = "update_by_name"
)

type MCPImportPreviewRequest struct {
	Format         string
	Content        string
	DefaultInclude bool
}

type MCPImportApplyRequest struct {
	Format         string
	Content        string
	DefaultInclude bool
	SelectedNames  []string
	ConflictPolicy string
}

type MCPImportApplyResult struct {
	Format                string                    `json:"format"`
	Imported              []model.MCPDefinition     `json:"imported"`
	Updated               []model.MCPDefinition     `json:"updated"`
	Skipped               []MCPImportSkipped        `json:"skipped"`
	ReferenceRequirements []MCPReferenceRequirement `json:"reference_requirements,omitempty"`
	Warnings              []MCPImportDiagnostic     `json:"warnings,omitempty"`
}

type MCPImportSkipped struct {
	Name       string `json:"name"`
	ExistingID string `json:"existing_id,omitempty"`
	Reason     string `json:"reason"`
}

type MCPImportPreview struct {
	Format     string                `json:"format"`
	Candidates []MCPImportCandidate  `json:"candidates"`
	Warnings   []MCPImportDiagnostic `json:"warnings"`
	Errors     []MCPImportDiagnostic `json:"errors"`
}

type MCPImportCandidate struct {
	Name                  string                    `json:"name"`
	SourcePath            string                    `json:"source_path"`
	Format                string                    `json:"format"`
	Definition            model.MCPDefinition       `json:"definition"`
	ReferenceRequirements []MCPReferenceRequirement `json:"reference_requirements,omitempty"`
	Warnings              []MCPImportDiagnostic     `json:"warnings,omitempty"`
	Errors                []MCPImportDiagnostic     `json:"errors,omitempty"`
}

type MCPReferenceRequirement struct {
	Name      string `json:"name"`
	FieldPath string `json:"field_path"`
	Reason    string `json:"reason"`
}

type MCPImportDiagnostic struct {
	Code      string `json:"code"`
	FieldPath string `json:"field_path,omitempty"`
	Message   string `json:"message"`
}

func PreviewMCPImport(request MCPImportPreviewRequest) (MCPImportPreview, error) {
	format := strings.TrimSpace(request.Format)
	if format == "" {
		format = MCPImportFormatAuto
	}
	root, err := parseMCPImportJSONC(request.Content)
	if err != nil {
		return MCPImportPreview{}, err
	}
	selected, err := detectMCPImportFormat(format, root)
	if err != nil {
		return MCPImportPreview{}, err
	}
	preview := MCPImportPreview{Format: selected}
	servers, sourcePath, err := importServerMap(selected, root)
	if err != nil {
		return MCPImportPreview{}, err
	}
	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		candidate := normalizeImportServer(selected, sourcePath+"."+name, name, servers[name], request.DefaultInclude)
		preview.Candidates = append(preview.Candidates, candidate)
	}
	return preview, nil
}

func (s *MCPService) ApplyMCPImport(request MCPImportApplyRequest) (MCPImportApplyResult, error) {
	preview, err := PreviewMCPImport(MCPImportPreviewRequest{Format: request.Format, Content: request.Content, DefaultInclude: request.DefaultInclude})
	if err != nil {
		return MCPImportApplyResult{}, err
	}
	policy, err := normalizeMCPImportConflictPolicy(request.ConflictPolicy)
	if err != nil {
		return MCPImportApplyResult{}, err
	}
	candidates, err := selectMCPImportCandidates(preview.Candidates, request.SelectedNames)
	if err != nil {
		return MCPImportApplyResult{}, err
	}
	if len(candidates) == 0 {
		return MCPImportApplyResult{}, fmt.Errorf("mcp import selected no candidates")
	}
	result := MCPImportApplyResult{Format: preview.Format, Imported: []model.MCPDefinition{}, Updated: []model.MCPDefinition{}, Skipped: []MCPImportSkipped{}}
	result.Warnings = append(result.Warnings, preview.Warnings...)
	for _, candidate := range candidates {
		if len(candidate.Errors) != 0 {
			return MCPImportApplyResult{}, fmt.Errorf("mcp import candidate %q is invalid: %s at %s", candidate.Name, candidate.Errors[0].Code, candidate.Errors[0].FieldPath)
		}
		result.Warnings = append(result.Warnings, candidate.Warnings...)
		result.ReferenceRequirements = append(result.ReferenceRequirements, candidate.ReferenceRequirements...)
	}
	err = s.store.Update(func(state *model.State) error {
		existingByName := map[string]int{}
		for i := range state.MCPs {
			existingByName[strings.ToLower(state.MCPs[i].Name)] = i
		}
		for _, candidate := range candidates {
			definition := cloneMCPDefinition(candidate.Definition)
			nameKey := strings.ToLower(definition.Name)
			if existingIndex, exists := existingByName[nameKey]; exists {
				existing := state.MCPs[existingIndex]
				switch policy {
				case MCPImportConflictError:
					return fmt.Errorf("mcp import conflict for %q: existing mcp %q", definition.Name, existing.ID)
				case MCPImportConflictSkip:
					result.Skipped = append(result.Skipped, MCPImportSkipped{Name: definition.Name, ExistingID: existing.ID, Reason: "name_conflict"})
					continue
				case MCPImportConflictUpdateByName:
					definition.ID = existing.ID
					state.MCPs[existingIndex] = definition
					result.Updated = append(result.Updated, cloneMCPDefinition(definition))
				}
				continue
			}
			id, err := identity.New("mcp")
			if err != nil {
				return err
			}
			definition.ID = id
			state.MCPs = append(state.MCPs, definition)
			existingByName[nameKey] = len(state.MCPs) - 1
			result.Imported = append(result.Imported, cloneMCPDefinition(definition))
		}
		sort.Slice(state.MCPs, func(i, j int) bool {
			return strings.ToLower(state.MCPs[i].Name) < strings.ToLower(state.MCPs[j].Name)
		})
		return nil
	})
	if err != nil {
		return MCPImportApplyResult{}, err
	}
	return result, nil
}

func normalizeMCPImportConflictPolicy(policy string) (string, error) {
	policy = strings.TrimSpace(policy)
	if policy == "" {
		policy = MCPImportConflictError
	}
	switch policy {
	case MCPImportConflictError, MCPImportConflictSkip, MCPImportConflictUpdateByName:
		return policy, nil
	default:
		return "", fmt.Errorf("unsupported mcp import conflict_policy %q", policy)
	}
}

func selectMCPImportCandidates(candidates []MCPImportCandidate, selectedNames []string) ([]MCPImportCandidate, error) {
	byName := map[string]MCPImportCandidate{}
	for _, candidate := range candidates {
		key := strings.ToLower(strings.TrimSpace(candidate.Name))
		if key == "" {
			return nil, fmt.Errorf("mcp import candidate name is required")
		}
		if _, exists := byName[key]; exists {
			return nil, fmt.Errorf("duplicate mcp import candidate name %q", candidate.Name)
		}
		byName[key] = candidate
	}
	if len(selectedNames) == 0 {
		return append([]MCPImportCandidate(nil), candidates...), nil
	}
	seen := map[string]struct{}{}
	selected := make([]MCPImportCandidate, 0, len(selectedNames))
	for _, rawName := range selectedNames {
		name := strings.TrimSpace(rawName)
		if name == "" {
			return nil, fmt.Errorf("mcp import selected_names contains an empty name")
		}
		key := strings.ToLower(name)
		if _, duplicate := seen[key]; duplicate {
			return nil, fmt.Errorf("mcp import selected_names contains duplicate %q", name)
		}
		seen[key] = struct{}{}
		candidate, ok := byName[key]
		if !ok {
			return nil, fmt.Errorf("mcp import selected candidate %q not found", name)
		}
		selected = append(selected, candidate)
	}
	return selected, nil
}

func parseMCPImportJSONC(content string) (map[string]any, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, fmt.Errorf("mcp import content is required")
	}
	withoutComments, err := stripJSONCComments(trimmed)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewBufferString(withoutComments))
	decoder.UseNumber()
	var root map[string]any
	if err := decoder.Decode(&root); err != nil {
		return nil, fmt.Errorf("invalid mcp import json/jsonc: %w", err)
	}
	if len(root) == 0 {
		return nil, fmt.Errorf("mcp import content must be a non-empty object")
	}
	return root, nil
}

func stripJSONCComments(input string) (string, error) {
	var out strings.Builder
	inString := false
	escaped := false
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inString {
			out.WriteByte(ch)
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}
		if ch == '"' {
			inString = true
			out.WriteByte(ch)
			continue
		}
		if ch == '/' && i+1 < len(input) {
			next := input[i+1]
			if next == '/' {
				for i < len(input) && input[i] != '\n' {
					i++
				}
				if i < len(input) {
					out.WriteByte(input[i])
				}
				continue
			}
			if next == '*' {
				i += 2
				closed := false
				for i < len(input)-1 {
					if input[i] == '*' && input[i+1] == '/' {
						closed = true
						i++
						break
					}
					if input[i] == '\n' {
						out.WriteByte('\n')
					}
					i++
				}
				if !closed {
					return "", fmt.Errorf("invalid jsonc: unterminated block comment")
				}
				continue
			}
		}
		out.WriteByte(ch)
	}
	return out.String(), nil
}

func detectMCPImportFormat(format string, root map[string]any) (string, error) {
	switch format {
	case MCPImportFormatOpenCode, MCPImportFormatCodexPlugin:
		return format, nil
	case MCPImportFormatAuto:
		matches := []string{}
		if hasOpenCodeShape(root) {
			matches = append(matches, MCPImportFormatOpenCode)
		}
		if hasCodexPluginShape(root) {
			matches = append(matches, MCPImportFormatCodexPlugin)
		}
		if len(matches) == 0 {
			return "", fmt.Errorf("unsupported mcp import format")
		}
		if len(matches) > 1 {
			return "", fmt.Errorf("ambiguous_format: %s", strings.Join(matches, ","))
		}
		return matches[0], nil
	default:
		return "", fmt.Errorf("unsupported mcp import format %q", format)
	}
}

func hasOpenCodeShape(root map[string]any) bool {
	mcpValue, ok := root["mcp"].(map[string]any)
	if !ok {
		return false
	}
	_, ok = mcpValue["servers"].(map[string]any)
	return ok
}

func hasCodexPluginShape(root map[string]any) bool {
	if _, ok := root["mcpServers"].(map[string]any); ok {
		return true
	}
	if _, ok := root["mcp"].(map[string]any); ok {
		return false
	}
	for _, value := range root {
		entry, ok := value.(map[string]any)
		if !ok {
			return false
		}
		if _, hasCommand := entry["command"]; hasCommand {
			continue
		}
		if _, hasURL := entry["url"]; hasURL {
			continue
		}
		return false
	}
	return len(root) > 0
}

func importServerMap(format string, root map[string]any) (map[string]any, string, error) {
	switch format {
	case MCPImportFormatOpenCode:
		mcpValue, ok := root["mcp"].(map[string]any)
		if !ok {
			return nil, "", fmt.Errorf("opencode import requires mcp.servers")
		}
		servers, ok := mcpValue["servers"].(map[string]any)
		if !ok || len(servers) == 0 {
			return nil, "", fmt.Errorf("opencode import requires non-empty mcp.servers")
		}
		return servers, "mcp.servers", nil
	case MCPImportFormatCodexPlugin:
		if servers, ok := root["mcpServers"].(map[string]any); ok {
			if len(servers) == 0 {
				return nil, "", fmt.Errorf("codex-plugin import requires non-empty mcpServers")
			}
			return servers, "mcpServers", nil
		}
		if !hasCodexPluginShape(root) {
			return nil, "", fmt.Errorf("codex-plugin import requires mcpServers or a direct server map")
		}
		return root, "", nil
	default:
		return nil, "", fmt.Errorf("unsupported mcp import format %q", format)
	}
}

func normalizeImportServer(format, sourcePath, name string, raw any, defaultInclude bool) MCPImportCandidate {
	candidate := MCPImportCandidate{Name: name, SourcePath: strings.TrimPrefix(sourcePath, "."), Format: format}
	entry, ok := raw.(map[string]any)
	if !ok {
		candidate.Errors = append(candidate.Errors, MCPImportDiagnostic{Code: "invalid_server", FieldPath: candidate.SourcePath, Message: "server entry must be an object"})
		return candidate
	}
	definition := model.MCPDefinition{Name: name, DefaultIncludeInEnv: defaultInclude, AuthMode: MCPAuthNone}
	if format == MCPImportFormatOpenCode {
		normalizeOpenCodeServer(&candidate, &definition, entry)
	} else {
		normalizeCodexPluginServer(&candidate, &definition, entry)
	}
	definition.HealthPolicy = normalizeMCPHealthPolicy(definition.HealthPolicy)
	if normalized, err := normalizeMCPDefinition(definition); err != nil {
		candidate.Errors = append(candidate.Errors, MCPImportDiagnostic{Code: "invalid_normalized_config", FieldPath: candidate.SourcePath, Message: err.Error()})
	} else {
		candidate.Definition = normalized
	}
	return candidate
}

func normalizeOpenCodeServer(candidate *MCPImportCandidate, definition *model.MCPDefinition, entry map[string]any) {
	typeValue := strings.TrimSpace(importString(entry["type"]))
	switch typeValue {
	case "remote":
		definition.Transport = MCPTransportStreamableHTTP
		definition.Endpoint = importString(entry["url"])
		definition.HeaderRefs = convertCredentialMap(candidate, entry["headers"], candidate.SourcePath+".headers")
		if len(definition.HeaderRefs) > 0 {
			definition.AuthMode = MCPAuthHeaders
		}
	case "local":
		definition.Transport = MCPTransportStdio
		command := importStringSlice(entry["command"])
		if len(command) == 0 {
			candidate.Errors = append(candidate.Errors, MCPImportDiagnostic{Code: "missing_command", FieldPath: candidate.SourcePath + ".command", Message: "local server requires command array"})
			return
		}
		definition.Executable = command[0]
		definition.Args = command[1:]
		definition.EnvRefs = convertCredentialMap(candidate, entry["environment"], candidate.SourcePath+".environment")
	default:
		candidate.Errors = append(candidate.Errors, MCPImportDiagnostic{Code: "unsupported_type", FieldPath: candidate.SourcePath + ".type", Message: "opencode server type must be local or remote"})
	}
	if _, ok := entry["disabled"]; ok {
		candidate.Warnings = append(candidate.Warnings, MCPImportDiagnostic{Code: "source_disabled_ignored", FieldPath: candidate.SourcePath + ".disabled", Message: "source disabled flag is preview-only and does not change ADM Environment selections"})
	}
}

func normalizeCodexPluginServer(candidate *MCPImportCandidate, definition *model.MCPDefinition, entry map[string]any) {
	if transport := strings.TrimSpace(importString(entry["transport"])); strings.EqualFold(transport, "sse") {
		candidate.Errors = append(candidate.Errors, MCPImportDiagnostic{Code: "unsupported_transport", FieldPath: candidate.SourcePath + ".transport", Message: "legacy SSE MCP transport is not imported"})
		return
	}
	if urlValue := strings.TrimSpace(importString(entry["url"])); urlValue != "" {
		definition.Transport = MCPTransportStreamableHTTP
		definition.Endpoint = urlValue
		definition.HeaderRefs = convertCredentialMap(candidate, entry["headers"], candidate.SourcePath+".headers")
		if len(definition.HeaderRefs) > 0 {
			definition.AuthMode = MCPAuthHeaders
		}
		return
	}
	definition.Transport = MCPTransportStdio
	definition.Executable = importString(entry["command"])
	definition.Args = importStringSlice(entry["args"])
	definition.EnvRefs = convertCredentialMap(candidate, entry["env"], candidate.SourcePath+".env")
}

func convertCredentialMap(candidate *MCPImportCandidate, raw any, fieldPath string) map[string]string {
	values, ok := raw.(map[string]any)
	if !ok || len(values) == 0 {
		return nil
	}
	result := map[string]string{}
	for rawKey, rawValue := range values {
		key := strings.TrimSpace(rawKey)
		if key == "" {
			candidate.Errors = append(candidate.Errors, MCPImportDiagnostic{Code: "empty_key", FieldPath: fieldPath, Message: "credential map contains an empty key"})
			continue
		}
		value := importString(rawValue)
		path := fieldPath + "." + key
		if converted, ok := recognizedSourceEnvReference(value); ok {
			result[key] = converted
			continue
		}
		if containsMCPEnvironmentReference(value) {
			result[key] = value
			continue
		}
		if mcpReferenceRequiresEnvironmentReference(fieldKindFromPath(fieldPath), key) {
			refName := generatedMCPImportReferenceName(candidate.Name, key)
			result[key] = referenceTemplateForLiteralCredential(key, value, refName)
			candidate.ReferenceRequirements = append(candidate.ReferenceRequirements, MCPReferenceRequirement{Name: refName, FieldPath: path, Reason: "literal credential converted to environment reference requirement"})
			continue
		}
		result[key] = value
	}
	return result
}

var sourceEnvReferencePattern = regexp.MustCompile(`\{env:([A-Za-z_][A-Za-z0-9_]*)\}`)

func recognizedSourceEnvReference(value string) (string, bool) {
	trimmed := strings.TrimSpace(value)
	converted := sourceEnvReferencePattern.ReplaceAllString(trimmed, `${$1}`)
	return converted, converted != trimmed
}

func referenceTemplateForLiteralCredential(key, value, refName string) string {
	reference := "${" + refName + "}"
	if strings.EqualFold(strings.TrimSpace(key), "Authorization") {
		fields := strings.Fields(value)
		if len(fields) >= 2 {
			scheme := fields[0]
			if strings.EqualFold(scheme, "Bearer") || strings.EqualFold(scheme, "Basic") {
				return scheme + " " + reference
			}
		}
	}
	return reference
}

func generatedMCPImportReferenceName(candidateName, key string) string {
	base := strings.ToUpper(candidateName + "_" + key)
	var b strings.Builder
	lastUnderscore := false
	for _, r := range base {
		valid := (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
		if !valid {
			r = '_'
		}
		if r == '_' {
			if lastUnderscore {
				continue
			}
			lastUnderscore = true
		} else {
			lastUnderscore = false
		}
		b.WriteRune(r)
	}
	result := strings.Trim(b.String(), "_")
	if result == "" || (result[0] >= '0' && result[0] <= '9') {
		result = "MCP_" + result
	}
	return result
}

func fieldKindFromPath(fieldPath string) string {
	if strings.Contains(fieldPath, "header") {
		return "header_refs"
	}
	return "env_refs"
}

func importString(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return v.String()
	default:
		return ""
	}
}

func importStringSlice(value any) []string {
	items, ok := value.([]any)
	if !ok {
		if single := importString(value); single != "" {
			return []string{single}
		}
		return nil
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if value := importString(item); value != "" {
			result = append(result, value)
		}
	}
	return result
}

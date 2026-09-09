package catalog

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"ai-dev-manager-v2/internal/identity"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/store"
)

const (
	MCPTransportStreamableHTTP = "streamable-http"
	MCPTransportStdio          = "stdio"

	MCPAuthNone    = "none"
	MCPAuthHeaders = "headers"

	defaultMCPCheckIntervalSeconds     int64 = 30
	defaultMCPProbeTimeoutSeconds      int64 = 10
	defaultMCPReconnectIntervalSeconds int64 = 30
)

var (
	bracedEnvReferencePattern = regexp.MustCompile(`\$\{[A-Za-z_][A-Za-z0-9_]*(:-[^}]*)?\}`)
	plainEnvReferencePattern  = regexp.MustCompile(`(^|[^$])\$[A-Za-z_][A-Za-z0-9_]*`)
)

type MCPConfig struct {
	Transport      string
	AuthMode       string
	Endpoint       string
	HeaderRefs     map[string]string
	Executable     string
	Args           []string
	EnvRefs        map[string]string
	HealthPolicy   model.MCPHealthPolicy
	DefaultInclude bool
}

type MCPService struct {
	store *store.Store
}

func NewMCP(s *store.Store) *MCPService { return &MCPService{store: s} }

// Add creates an intentionally incomplete MCP definition for internal state
// tests and management inspection. Product configuration surfaces should use
// AddMCP or AddMCPConfig so transport-local validation runs before persistence.
func (s *MCPService) Add(name string, defaultInclude bool) (model.MCPDefinition, error) {
	return s.add(model.MCPDefinition{
		Name:                name,
		DefaultIncludeInEnv: defaultInclude,
		Transport:           MCPTransportStreamableHTTP,
		AuthMode:            MCPAuthNone,
		HealthPolicy:        normalizeMCPHealthPolicy(model.MCPHealthPolicy{}),
	})
}

func (s *MCPService) AddMCP(name, endpoint string, defaultInclude bool) (model.MCPDefinition, error) {
	return s.AddMCPConfig(name, MCPConfig{
		Transport:      MCPTransportStreamableHTTP,
		AuthMode:       MCPAuthNone,
		Endpoint:       endpoint,
		DefaultInclude: defaultInclude,
	})
}

func (s *MCPService) AddMCPConfig(name string, config MCPConfig) (model.MCPDefinition, error) {
	definition, err := normalizeMCPDefinition(model.MCPDefinition{
		Name:                name,
		DefaultIncludeInEnv: config.DefaultInclude,
		Transport:           config.Transport,
		AuthMode:            config.AuthMode,
		Endpoint:            config.Endpoint,
		HeaderRefs:          cloneStringMap(config.HeaderRefs),
		Executable:          config.Executable,
		Args:                append([]string(nil), config.Args...),
		EnvRefs:             cloneStringMap(config.EnvRefs),
		HealthPolicy:        config.HealthPolicy,
	})
	if err != nil {
		return model.MCPDefinition{}, err
	}
	return s.add(definition)
}

func (s *MCPService) add(definition model.MCPDefinition) (model.MCPDefinition, error) {
	definition.Name = strings.TrimSpace(definition.Name)
	if definition.Name == "" {
		return model.MCPDefinition{}, fmt.Errorf("mcp name is required")
	}
	var result model.MCPDefinition
	err := s.store.Update(func(state *model.State) error {
		for _, existing := range state.MCPs {
			if strings.EqualFold(existing.Name, definition.Name) {
				return fmt.Errorf("mcp %q already exists", definition.Name)
			}
		}
		id, err := identity.New("mcp")
		if err != nil {
			return err
		}
		result = cloneMCPDefinition(definition)
		result.ID = id
		state.MCPs = append(state.MCPs, result)
		sort.Slice(state.MCPs, func(i, j int) bool {
			return strings.ToLower(state.MCPs[i].Name) < strings.ToLower(state.MCPs[j].Name)
		})
		return nil
	})
	return result, err
}

func (s *MCPService) List() ([]model.MCPDefinition, error) {
	state, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	result := make([]model.MCPDefinition, 0, len(state.MCPs))
	for _, definition := range state.MCPs {
		result = append(result, cloneMCPDefinition(definition))
	}
	return result, nil
}

func (s *MCPService) Get(id string) (model.MCPDefinition, error) {
	items, err := s.List()
	if err != nil {
		return model.MCPDefinition{}, err
	}
	for _, definition := range items {
		if definition.ID == id {
			return definition, nil
		}
	}
	return model.MCPDefinition{}, fmt.Errorf("mcp %q not found", id)
}

func (s *MCPService) SetDefault(id string, value bool) (model.MCPDefinition, error) {
	var result model.MCPDefinition
	err := s.store.Update(func(state *model.State) error {
		for i := range state.MCPs {
			if state.MCPs[i].ID != id {
				continue
			}
			state.MCPs[i].DefaultIncludeInEnv = value
			result = cloneMCPDefinition(state.MCPs[i])
			return nil
		}
		return fmt.Errorf("mcp %q not found", id)
	})
	return result, err
}

func (s *MCPService) Remove(id string) error {
	return s.store.Update(func(state *model.State) error {
		for i := range state.MCPs {
			if state.MCPs[i].ID != id {
				continue
			}
			state.MCPs = append(state.MCPs[:i], state.MCPs[i+1:]...)
			return nil
		}
		return fmt.Errorf("mcp %q not found", id)
	})
}

func (s *MCPService) Exists(id string) (bool, error) {
	items, err := s.List()
	if err != nil {
		return false, err
	}
	for _, definition := range items {
		if definition.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func normalizeMCPDefinition(definition model.MCPDefinition) (model.MCPDefinition, error) {
	definition.Name = strings.TrimSpace(definition.Name)
	definition.Transport = strings.TrimSpace(definition.Transport)
	if definition.Transport == "" {
		definition.Transport = MCPTransportStreamableHTTP
	}
	definition.AuthMode = strings.TrimSpace(definition.AuthMode)
	if definition.AuthMode == "" {
		definition.AuthMode = MCPAuthNone
	}
	definition.Endpoint = strings.TrimSpace(definition.Endpoint)
	definition.Executable = strings.TrimSpace(definition.Executable)
	definition.HealthPolicy = normalizeMCPHealthPolicy(definition.HealthPolicy)

	headerRefs, err := normalizeMCPReferenceMap("header_refs", definition.HeaderRefs)
	if err != nil {
		return model.MCPDefinition{}, err
	}
	definition.HeaderRefs = headerRefs
	envRefs, err := normalizeMCPReferenceMap("env_refs", definition.EnvRefs)
	if err != nil {
		return model.MCPDefinition{}, err
	}
	definition.EnvRefs = envRefs

	switch definition.Transport {
	case MCPTransportStreamableHTTP:
		parsed, err := url.Parse(definition.Endpoint)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return model.MCPDefinition{}, fmt.Errorf("mcp endpoint must be a valid http(s) URL")
		}
		if definition.Executable != "" || len(definition.Args) != 0 || len(definition.EnvRefs) != 0 {
			return model.MCPDefinition{}, fmt.Errorf("streamable-http mcp must not configure stdio fields")
		}
		switch definition.AuthMode {
		case MCPAuthNone:
			if len(definition.HeaderRefs) != 0 {
				return model.MCPDefinition{}, fmt.Errorf("mcp auth_mode %q must not configure header_refs", MCPAuthNone)
			}
		case MCPAuthHeaders:
			if len(definition.HeaderRefs) == 0 {
				return model.MCPDefinition{}, fmt.Errorf("mcp auth_mode %q requires header_refs", MCPAuthHeaders)
			}
		default:
			return model.MCPDefinition{}, fmt.Errorf("unsupported mcp auth_mode %q", definition.AuthMode)
		}
	case MCPTransportStdio:
		if definition.Executable == "" {
			return model.MCPDefinition{}, fmt.Errorf("stdio mcp executable is required")
		}
		if definition.Endpoint != "" || len(definition.HeaderRefs) != 0 {
			return model.MCPDefinition{}, fmt.Errorf("stdio mcp must not configure http fields")
		}
		if definition.AuthMode != MCPAuthNone {
			return model.MCPDefinition{}, fmt.Errorf("stdio mcp auth_mode must be %q", MCPAuthNone)
		}
	default:
		return model.MCPDefinition{}, fmt.Errorf("unsupported mcp transport %q", definition.Transport)
	}
	return definition, nil
}

func normalizeMCPReferenceMap(kind string, refs map[string]string) (map[string]string, error) {
	if len(refs) == 0 {
		return nil, nil
	}
	result := make(map[string]string, len(refs))
	for rawKey, rawValue := range refs {
		key := strings.TrimSpace(rawKey)
		if key == "" {
			return nil, fmt.Errorf("mcp %s contains an empty key", kind)
		}
		value := strings.TrimSpace(rawValue)
		if value == "" {
			return nil, fmt.Errorf("mcp %s for %q must use a non-empty environment reference", kind, key)
		}
		if !containsMCPEnvironmentReference(value) && mcpReferenceRequiresEnvironmentReference(kind, key) {
			return nil, fmt.Errorf("mcp %s for %q must use an environment reference; literal credential-bearing values are not persisted", kind, key)
		}
		result[key] = value
	}
	return result, nil
}

func containsMCPEnvironmentReference(value string) bool {
	return bracedEnvReferencePattern.MatchString(value) || plainEnvReferencePattern.MatchString(value)
}

func mcpReferenceRequiresEnvironmentReference(kind, key string) bool {
	if kind == "header_refs" {
		return true
	}
	key = strings.ToLower(strings.TrimSpace(key))
	for _, marker := range []string{"auth", "token", "secret", "password", "passwd", "credential", "api_key", "apikey", "private_key"} {
		if strings.Contains(key, marker) {
			return true
		}
	}
	return false
}

func normalizeMCPHealthPolicy(policy model.MCPHealthPolicy) model.MCPHealthPolicy {
	if policy.CheckIntervalSeconds <= 0 {
		policy.CheckIntervalSeconds = defaultMCPCheckIntervalSeconds
	}
	if policy.ProbeTimeoutSeconds <= 0 {
		policy.ProbeTimeoutSeconds = defaultMCPProbeTimeoutSeconds
	}
	if policy.ReconnectIntervalSeconds <= 0 {
		policy.ReconnectIntervalSeconds = defaultMCPReconnectIntervalSeconds
	}
	return policy
}

func cloneMCPDefinition(definition model.MCPDefinition) model.MCPDefinition {
	definition.HeaderRefs = cloneStringMap(definition.HeaderRefs)
	definition.Args = append([]string(nil), definition.Args...)
	definition.EnvRefs = cloneStringMap(definition.EnvRefs)
	return definition
}

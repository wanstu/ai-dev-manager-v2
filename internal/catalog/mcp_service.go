package catalog

import (
	"fmt"
	"net/url"
	"os"
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

func (s *MCPService) AddMCP(name, endpoint string, defaultInclude bool) (model.MCPDefinition, error) {
	return s.AddMCPConfig(name, MCPConfig{
		Transport:      MCPTransportStreamableHTTP,
		AuthMode:       MCPAuthNone,
		Endpoint:       endpoint,
		DefaultInclude: defaultInclude,
	})
}

func (s *MCPService) AddMCPConfig(name string, config MCPConfig) (model.MCPDefinition, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return model.MCPDefinition{}, fmt.Errorf("mcp name is required")
	}
	definition, err := validateMCPDefinition(model.MCPDefinition{
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

	var result model.MCPDefinition
	err = s.store.Update(func(state *model.State) error {
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
	return cloneMCPDefinition(result), err
}

func (s *MCPService) UpdateMCPConfig(id, name string, config MCPConfig) (model.MCPDefinition, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.MCPDefinition{}, fmt.Errorf("mcp id is required")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return model.MCPDefinition{}, fmt.Errorf("mcp name is required")
	}
	definition, err := validateMCPDefinition(model.MCPDefinition{
		ID:                  id,
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

	var result model.MCPDefinition
	err = s.store.Update(func(state *model.State) error {
		index := -1
		for i, existing := range state.MCPs {
			if existing.ID == id {
				index = i
				continue
			}
			if strings.EqualFold(existing.Name, definition.Name) {
				return fmt.Errorf("mcp %q already exists", definition.Name)
			}
		}
		if index < 0 {
			return fmt.Errorf("mcp %q not found", id)
		}
		result = cloneMCPDefinition(definition)
		state.MCPs[index] = result
		sort.Slice(state.MCPs, func(i, j int) bool {
			return strings.ToLower(state.MCPs[i].Name) < strings.ToLower(state.MCPs[j].Name)
		})
		return nil
	})
	return cloneMCPDefinition(result), err
}

func (s *MCPService) List() ([]model.MCPDefinition, error) {
	state, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	items := make([]model.MCPDefinition, len(state.MCPs))
	for i, item := range state.MCPs {
		items[i] = cloneMCPDefinition(item)
	}
	return items, nil
}

func (s *MCPService) Get(id string) (model.MCPDefinition, error) {
	items, err := s.List()
	if err != nil {
		return model.MCPDefinition{}, err
	}
	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}
	return model.MCPDefinition{}, fmt.Errorf("mcp %q not found", id)
}

func (s *MCPService) SetDefault(id string, value bool) (model.MCPDefinition, error) {
	var result model.MCPDefinition
	err := s.store.Update(func(state *model.State) error {
		for i := range state.MCPs {
			if state.MCPs[i].ID == id {
				state.MCPs[i].DefaultIncludeInEnv = value
				result = cloneMCPDefinition(state.MCPs[i])
				return nil
			}
		}
		return fmt.Errorf("mcp %q not found", id)
	})
	return result, err
}

func (s *MCPService) Remove(id string) error {
	return s.store.Update(func(state *model.State) error {
		for i := range state.MCPs {
			if state.MCPs[i].ID == id {
				state.MCPs = append(state.MCPs[:i], state.MCPs[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("mcp %q not found", id)
	})
}

func (s *MCPService) Exists(id string) (bool, error) {
	items, err := s.List()
	if err != nil {
		return false, err
	}
	for _, item := range items {
		if item.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func validateMCPDefinition(definition model.MCPDefinition) (model.MCPDefinition, error) {
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
	definition.HeaderRefs = trimStringMap(definition.HeaderRefs)
	definition.EnvRefs = trimStringMap(definition.EnvRefs)
	definition.Args = append([]string(nil), definition.Args...)

	if err := validateMCPHealthPolicy(definition.HealthPolicy); err != nil {
		return model.MCPDefinition{}, err
	}

	switch definition.Transport {
	case MCPTransportStreamableHTTP:
		if definition.Executable != "" || len(definition.Args) != 0 || len(definition.EnvRefs) != 0 {
			return model.MCPDefinition{}, fmt.Errorf("streamable-http mcp cannot configure stdio executable, args, or env_refs")
		}
		if !validHTTPURLTemplate(definition.Endpoint) {
			return model.MCPDefinition{}, fmt.Errorf("mcp endpoint must be a valid http(s) URL or environment-reference template")
		}
		switch definition.AuthMode {
		case MCPAuthNone:
			if len(definition.HeaderRefs) != 0 {
				return model.MCPDefinition{}, fmt.Errorf("mcp auth_mode %q cannot configure header_refs", MCPAuthNone)
			}
		case MCPAuthHeaders:
			if len(definition.HeaderRefs) == 0 {
				return model.MCPDefinition{}, fmt.Errorf("mcp auth_mode %q requires header_refs", MCPAuthHeaders)
			}
			for header, value := range definition.HeaderRefs {
				if !containsEnvironmentReference(value) {
					return model.MCPDefinition{}, fmt.Errorf("mcp header %q must use an environment reference", header)
				}
			}
		default:
			return model.MCPDefinition{}, fmt.Errorf("unsupported mcp auth_mode %q", definition.AuthMode)
		}
	case MCPTransportStdio:
		if definition.AuthMode != MCPAuthNone {
			return model.MCPDefinition{}, fmt.Errorf("stdio mcp auth_mode must be %q", MCPAuthNone)
		}
		if definition.Endpoint != "" || len(definition.HeaderRefs) != 0 {
			return model.MCPDefinition{}, fmt.Errorf("stdio mcp cannot configure endpoint or header_refs")
		}
		if definition.Executable == "" {
			return model.MCPDefinition{}, fmt.Errorf("stdio mcp executable is required")
		}
		for key, value := range definition.EnvRefs {
			if !containsEnvironmentReference(value) {
				return model.MCPDefinition{}, fmt.Errorf("stdio mcp env_refs[%q] must use an environment reference", key)
			}
		}
	default:
		return model.MCPDefinition{}, fmt.Errorf("unsupported mcp transport %q", definition.Transport)
	}
	return definition, nil
}

func validateMCPHealthPolicy(policy model.MCPHealthPolicy) error {
	if policy.CheckIntervalSeconds < 0 || policy.ProbeTimeoutSeconds < 0 || policy.ReconnectIntervalSeconds < 0 {
		return fmt.Errorf("mcp health policy intervals must be nonnegative")
	}
	if policy.HealthCheckEnabled {
		if policy.CheckIntervalSeconds <= 0 {
			return fmt.Errorf("mcp health_check_enabled requires positive check_interval_seconds")
		}
		if policy.ProbeTimeoutSeconds <= 0 {
			return fmt.Errorf("mcp health_check_enabled requires positive probe_timeout_seconds")
		}
	}
	if policy.AutoReconnect && policy.ReconnectIntervalSeconds <= 0 {
		return fmt.Errorf("mcp auto_reconnect requires positive reconnect_interval_seconds")
	}
	return nil
}

func validHTTPURLTemplate(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	expanded := os.Expand(value, func(string) string { return "adm-reference" })
	parsed, err := url.Parse(expanded)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

func containsEnvironmentReference(value string) bool {
	found := false
	_ = os.Expand(value, func(string) string {
		found = true
		return ""
	})
	return found
}

func trimStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	result := make(map[string]string, len(input))
	for key, value := range input {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key != "" {
			result[key] = value
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func cloneMCPDefinition(input model.MCPDefinition) model.MCPDefinition {
	input.HeaderRefs = cloneStringMap(input.HeaderRefs)
	input.EnvRefs = cloneStringMap(input.EnvRefs)
	input.Args = append([]string(nil), input.Args...)
	return input
}

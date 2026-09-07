package catalog

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"ai-dev-manager-v2/internal/identity"
	"ai-dev-manager-v2/internal/model"
	skillruntime "ai-dev-manager-v2/internal/skill"
	"ai-dev-manager-v2/internal/store"
)

type Kind string

const (
	KindMCP   Kind = "mcp"
	KindSkill Kind = "skill"

	MCPTransportStreamableHTTP = "streamable-http"
)

type MCPConfig struct {
	Endpoint       string
	Transport      string
	HeaderRefs     map[string]string
	DefaultInclude bool
}

type Service struct {
	store *store.Store
	kind  Kind
}

func New(s *store.Store, kind Kind) *Service { return &Service{store: s, kind: kind} }

// Add creates a metadata-only catalog entry. Product management surfaces must
// use AddMCP or AddSkillRoot for runtime-backed definitions.
func (s *Service) Add(name string, defaultInclude bool) (model.CatalogEntry, error) {
	return s.add(model.CatalogEntry{Name: name, DefaultIncludeInEnv: defaultInclude})
}

func (s *Service) AddMCP(name, endpoint string, defaultInclude bool) (model.CatalogEntry, error) {
	return s.AddMCPConfig(name, MCPConfig{
		Endpoint:       endpoint,
		Transport:      MCPTransportStreamableHTTP,
		DefaultInclude: defaultInclude,
	})
}

func (s *Service) AddMCPConfig(name string, config MCPConfig) (model.CatalogEntry, error) {
	if s.kind != KindMCP {
		return model.CatalogEntry{}, fmt.Errorf("catalog kind %q is not MCP", s.kind)
	}
	endpoint := strings.TrimSpace(config.Endpoint)
	parsed, err := url.Parse(endpoint)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return model.CatalogEntry{}, fmt.Errorf("mcp endpoint must be a valid http(s) URL")
	}
	transport := strings.TrimSpace(config.Transport)
	if transport == "" {
		transport = MCPTransportStreamableHTTP
	}
	if transport != MCPTransportStreamableHTTP {
		return model.CatalogEntry{}, fmt.Errorf("mcp transport must be %q", MCPTransportStreamableHTTP)
	}
	return s.add(model.CatalogEntry{
		Name:                name,
		DefaultIncludeInEnv: config.DefaultInclude,
		Endpoint:            endpoint,
		Transport:           transport,
		HeaderRefs:          cloneStringMap(config.HeaderRefs),
	})
}

// AddSkill is retained only for existing internal tests/dev-state inspection.
// New product surfaces configure real Skills by explicit discovery root.
func (s *Service) AddSkill(name, instructions string, defaultInclude bool) (model.CatalogEntry, error) {
	if s.kind != KindSkill {
		return model.CatalogEntry{}, fmt.Errorf("catalog kind %q is not Skill", s.kind)
	}
	instructions = strings.TrimSpace(instructions)
	if instructions == "" {
		return model.CatalogEntry{}, fmt.Errorf("skill instructions are required")
	}
	return s.add(model.CatalogEntry{Name: name, DefaultIncludeInEnv: defaultInclude, Instructions: instructions})
}

func (s *Service) AddSkillRoot(root string, supportRoots []string, defaultInclude bool) ([]model.CatalogEntry, error) {
	if s.kind != KindSkill {
		return nil, fmt.Errorf("catalog kind %q is not Skill", s.kind)
	}
	discovered, err := skillruntime.Discover(root, supportRoots, defaultInclude)
	if err != nil {
		return nil, err
	}
	err = s.store.Update(func(state *model.State) error {
		for _, incoming := range discovered {
			replaced := false
			for i := range state.Skills {
				if state.Skills[i].ID == incoming.ID {
					state.Skills[i] = incoming
					replaced = true
					break
				}
			}
			if !replaced {
				state.Skills = append(state.Skills, incoming)
			}
		}
		sort.Slice(state.Skills, func(i, j int) bool {
			return strings.ToLower(state.Skills[i].Name) < strings.ToLower(state.Skills[j].Name)
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return discovered, nil
}

func (s *Service) add(entry model.CatalogEntry) (model.CatalogEntry, error) {
	entry.Name = strings.TrimSpace(entry.Name)
	if entry.Name == "" {
		return model.CatalogEntry{}, fmt.Errorf("%s name is required", s.kind)
	}
	var result model.CatalogEntry
	err := s.store.Update(func(state *model.State) error {
		items := s.items(state)
		for _, existing := range *items {
			if strings.EqualFold(existing.Name, entry.Name) {
				return fmt.Errorf("%s %q already exists", s.kind, entry.Name)
			}
		}
		id, err := identity.New(string(s.kind))
		if err != nil {
			return err
		}
		result = entry
		result.ID = id
		*items = append(*items, result)
		sort.Slice(*items, func(i, j int) bool {
			return strings.ToLower((*items)[i].Name) < strings.ToLower((*items)[j].Name)
		})
		return nil
	})
	return result, err
}

func (s *Service) List() ([]model.CatalogEntry, error) {
	state, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	items := s.items(&state)
	result := append([]model.CatalogEntry(nil), (*items)...)
	for i := range result {
		result[i] = s.normalizeEntry(result[i])
	}
	return result, nil
}

func (s *Service) Get(id string) (model.CatalogEntry, error) {
	items, err := s.List()
	if err != nil {
		return model.CatalogEntry{}, err
	}
	for _, entry := range items {
		if entry.ID == id {
			return entry, nil
		}
	}
	return model.CatalogEntry{}, fmt.Errorf("%s %q not found", s.kind, id)
}

func (s *Service) SetDefault(id string, value bool) (model.CatalogEntry, error) {
	var result model.CatalogEntry
	err := s.store.Update(func(state *model.State) error {
		items := s.items(state)
		for i := range *items {
			if (*items)[i].ID == id {
				(*items)[i].DefaultIncludeInEnv = value
				result = s.normalizeEntry((*items)[i])
				return nil
			}
		}
		return fmt.Errorf("%s %q not found", s.kind, id)
	})
	return result, err
}

func (s *Service) Remove(id string) error {
	return s.store.Update(func(state *model.State) error {
		items := s.items(state)
		for i := range *items {
			if (*items)[i].ID == id {
				*items = append((*items)[:i], (*items)[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("%s %q not found", s.kind, id)
	})
}

func (s *Service) Exists(id string) (bool, error) {
	items, err := s.List()
	if err != nil {
		return false, err
	}
	for _, entry := range items {
		if entry.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) normalizeEntry(entry model.CatalogEntry) model.CatalogEntry {
	if s.kind == KindMCP && strings.TrimSpace(entry.Transport) == "" {
		entry.Transport = MCPTransportStreamableHTTP
	}
	return entry
}

func (s *Service) items(state *model.State) *[]model.CatalogEntry {
	if s.kind == KindSkill {
		return &state.Skills
	}
	return &state.MCPs
}

func cloneStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	result := make(map[string]string, len(input))
	for key, value := range input {
		result[key] = value
	}
	return result
}

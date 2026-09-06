package catalog

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"ai-dev-manager-v2/internal/identity"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/store"
)

type Kind string

const (
	KindMCP   Kind = "mcp"
	KindSkill Kind = "skill"
)

type Service struct {
	store *store.Store
	kind  Kind
}

func New(s *store.Store, kind Kind) *Service { return &Service{store: s, kind: kind} }

func (s *Service) Add(name string, defaultInclude bool) (model.CatalogEntry, error) {
	return s.add(name, "", "", defaultInclude)
}

func (s *Service) AddMCP(name, endpoint string, defaultInclude bool) (model.CatalogEntry, error) {
	if s.kind != KindMCP {
		return model.CatalogEntry{}, fmt.Errorf("catalog kind %q is not MCP", s.kind)
	}
	endpoint = strings.TrimSpace(endpoint)
	parsed, err := url.Parse(endpoint)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return model.CatalogEntry{}, fmt.Errorf("mcp endpoint must be a valid http(s) URL")
	}
	return s.add(name, endpoint, "", defaultInclude)
}

func (s *Service) AddSkill(name, instructions string, defaultInclude bool) (model.CatalogEntry, error) {
	if s.kind != KindSkill {
		return model.CatalogEntry{}, fmt.Errorf("catalog kind %q is not Skill", s.kind)
	}
	instructions = strings.TrimSpace(instructions)
	if instructions == "" {
		return model.CatalogEntry{}, fmt.Errorf("skill instructions are required")
	}
	return s.add(name, "", instructions, defaultInclude)
}

func (s *Service) add(name, endpoint, instructions string, defaultInclude bool) (model.CatalogEntry, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return model.CatalogEntry{}, fmt.Errorf("%s name is required", s.kind)
	}
	var result model.CatalogEntry
	err := s.store.Update(func(state *model.State) error {
		items := s.items(state)
		for _, entry := range *items {
			if strings.EqualFold(entry.Name, name) {
				return fmt.Errorf("%s %q already exists", s.kind, name)
			}
		}
		id, err := identity.New(string(s.kind))
		if err != nil {
			return err
		}
		result = model.CatalogEntry{ID: id, Name: name, DefaultIncludeInEnv: defaultInclude, Endpoint: endpoint, Instructions: instructions}
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
	return append([]model.CatalogEntry(nil), (*items)...), nil
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
				result = (*items)[i]
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

func (s *Service) items(state *model.State) *[]model.CatalogEntry {
	if s.kind == KindSkill {
		return &state.Skills
	}
	return &state.MCPs
}

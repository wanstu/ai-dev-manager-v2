package catalog

import (
	"fmt"
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
		result = model.CatalogEntry{ID: id, Name: name, DefaultIncludeInEnv: defaultInclude}
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

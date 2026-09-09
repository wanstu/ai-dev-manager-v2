package catalog

import (
	"fmt"
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
)

// Service owns the remaining generic catalog behavior used by Skills. MCPs use
// MCPService because their desired configuration is transport-specific.
type Service struct {
	store *store.Store
	kind  Kind
}

func New(s *store.Store, kind Kind) *Service { return &Service{store: s, kind: kind} }

func (s *Service) requireSkill() error {
	if s.kind != KindSkill {
		return fmt.Errorf("catalog kind %q is not Skill", s.kind)
	}
	return nil
}

func (s *Service) Add(name string, defaultInclude bool) (model.CatalogEntry, error) {
	if err := s.requireSkill(); err != nil {
		return model.CatalogEntry{}, err
	}
	return s.add(model.CatalogEntry{Name: name, DefaultIncludeInEnv: defaultInclude})
}

// AddSkill is retained only for existing internal tests/dev-state inspection.
// New product surfaces configure real Skills by explicit discovery root.
func (s *Service) AddSkill(name, instructions string, defaultInclude bool) (model.CatalogEntry, error) {
	if err := s.requireSkill(); err != nil {
		return model.CatalogEntry{}, err
	}
	instructions = strings.TrimSpace(instructions)
	if instructions == "" {
		return model.CatalogEntry{}, fmt.Errorf("skill instructions are required")
	}
	return s.add(model.CatalogEntry{Name: name, DefaultIncludeInEnv: defaultInclude, Instructions: instructions})
}

func (s *Service) AddSkillRoot(root string, supportRoots []string, defaultInclude bool) ([]model.CatalogEntry, error) {
	if err := s.requireSkill(); err != nil {
		return nil, err
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
		return model.CatalogEntry{}, fmt.Errorf("skill name is required")
	}
	var result model.CatalogEntry
	err := s.store.Update(func(state *model.State) error {
		for _, existing := range state.Skills {
			if strings.EqualFold(existing.Name, entry.Name) {
				return fmt.Errorf("skill %q already exists", entry.Name)
			}
		}
		id, err := identity.New("skill")
		if err != nil {
			return err
		}
		result = entry
		result.ID = id
		state.Skills = append(state.Skills, result)
		sort.Slice(state.Skills, func(i, j int) bool {
			return strings.ToLower(state.Skills[i].Name) < strings.ToLower(state.Skills[j].Name)
		})
		return nil
	})
	return result, err
}

func (s *Service) List() ([]model.CatalogEntry, error) {
	if err := s.requireSkill(); err != nil {
		return nil, err
	}
	state, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	return append([]model.CatalogEntry(nil), state.Skills...), nil
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
	return model.CatalogEntry{}, fmt.Errorf("skill %q not found", id)
}

func (s *Service) SetDefault(id string, value bool) (model.CatalogEntry, error) {
	if err := s.requireSkill(); err != nil {
		return model.CatalogEntry{}, err
	}
	var result model.CatalogEntry
	err := s.store.Update(func(state *model.State) error {
		for i := range state.Skills {
			if state.Skills[i].ID == id {
				state.Skills[i].DefaultIncludeInEnv = value
				result = state.Skills[i]
				return nil
			}
		}
		return fmt.Errorf("skill %q not found", id)
	})
	return result, err
}

func (s *Service) Remove(id string) error {
	if err := s.requireSkill(); err != nil {
		return err
	}
	return s.store.Update(func(state *model.State) error {
		for i := range state.Skills {
			if state.Skills[i].ID == id {
				state.Skills = append(state.Skills[:i], state.Skills[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("skill %q not found", id)
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

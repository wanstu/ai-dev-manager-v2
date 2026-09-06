package model

import "time"

type Workspace struct {
	ID        string    `json:"workspace_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

type WriterLease struct {
	Owner      string    `json:"owner"`
	AcquiredAt time.Time `json:"acquired_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type Environment struct {
	ID              string            `json:"environment_id"`
	WorkspaceID     string            `json:"workspace_id"`
	Name            string            `json:"name"`
	Root            string            `json:"root"`
	State           string            `json:"state"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	LastActivityAt  time.Time         `json:"last_activity_at"`
	Writer          *WriterLease      `json:"writer,omitempty"`
	EnabledMCPIDs   []string          `json:"enabled_mcp_ids,omitempty"`
	EnabledSkillIDs []string          `json:"enabled_skill_ids,omitempty"`
	PrivateMemory   map[string]string `json:"private_memory,omitempty"`
}

type CatalogEntry struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	DefaultIncludeInEnv bool     `json:"default_include_in_environment"`
	Endpoint            string   `json:"endpoint,omitempty"`
	Instructions        string   `json:"instructions,omitempty"`
	ArtifactPath        string   `json:"artifact_path,omitempty"`
	SourceRoot          string   `json:"source_root,omitempty"`
	SupportRoots        []string `json:"support_roots,omitempty"`
}

type State struct {
	Version            int               `json:"version"`
	Workspaces         []Workspace       `json:"workspaces"`
	Environments       []Environment     `json:"environments"`
	AllowedExecutables []string          `json:"allowed_executables,omitempty"`
	MCPs               []CatalogEntry    `json:"mcps,omitempty"`
	Skills             []CatalogEntry    `json:"skills,omitempty"`
	GlobalMemory       map[string]string `json:"global_memory,omitempty"`
}

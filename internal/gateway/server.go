package gateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/management"
	"ai-dev-manager-v2/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName         = "ai-dev-manager-v2"
	serverVersion      = "v0.1.0-dev"
	runtimeOwnerHeader = "X-ADM-Runtime-Owner"
)

type EmptyInput struct{}

type EnvironmentInput struct {
	EnvironmentID string `json:"environment_id"`
}

type EnvironmentRenameInput struct {
	EnvironmentID string `json:"environment_id"`
	Name          string `json:"name"`
}

type WorkspaceAddInput struct {
	Path string `json:"path"`
	Name string `json:"name,omitempty"`
}

type WorkspaceInput struct {
	WorkspaceID string `json:"workspace_id"`
}

type WorkspaceRenameInput struct {
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
}

type ExecutableInput struct {
	Executable string `json:"executable"`
}

type EnvironmentCreateInput struct {
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Root        string `json:"root,omitempty" jsonschema:"optional workspace-contained directory; defaults to workspace root"`
}

type EnvironmentWorktreeCreateInput struct {
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	BaseRef     string `json:"base_ref,omitempty" jsonschema:"optional Git ref; defaults to HEAD"`
}

type EnvironmentWorktreeDestroyInput struct {
	EnvironmentID string `json:"environment_id"`
	WriterOwner   string `json:"writer_owner"`
	Force         bool   `json:"force,omitempty"`
}

type WriterAcquireInput struct {
	EnvironmentID string `json:"environment_id"`
	Owner         string `json:"owner" jsonschema:"stable agent/session owner identifier"`
}

type MCPAddInput struct {
	Name           string                `json:"name"`
	Transport      string                `json:"transport"`
	AuthMode       string                `json:"auth_mode"`
	Endpoint       string                `json:"endpoint,omitempty"`
	HeaderRefs     map[string]string     `json:"header_refs,omitempty"`
	Executable     string                `json:"executable,omitempty"`
	Args           []string              `json:"args,omitempty"`
	EnvRefs        map[string]string     `json:"env_refs,omitempty"`
	HealthPolicy   model.MCPHealthPolicy `json:"health_policy"`
	DefaultInclude bool                  `json:"default_include_in_environment,omitempty"`
}

type MCPUpdateInput struct {
	ID             string                `json:"id"`
	Name           string                `json:"name"`
	Transport      string                `json:"transport"`
	AuthMode       string                `json:"auth_mode"`
	Endpoint       string                `json:"endpoint,omitempty"`
	HeaderRefs     map[string]string     `json:"header_refs,omitempty"`
	Executable     string                `json:"executable,omitempty"`
	Args           []string              `json:"args,omitempty"`
	EnvRefs        map[string]string     `json:"env_refs,omitempty"`
	HealthPolicy   model.MCPHealthPolicy `json:"health_policy"`
	DefaultInclude bool                  `json:"default_include_in_environment,omitempty"`
}

type CatalogAddInput struct {
	Name           string   `json:"name,omitempty"`
	Endpoint       string   `json:"endpoint,omitempty"`
	Instructions   string   `json:"instructions,omitempty"`
	Root           string   `json:"root,omitempty"`
	SupportRoots   []string `json:"support_roots,omitempty"`
	DefaultInclude bool     `json:"default_include_in_environment,omitempty"`
}

type CatalogIDInput struct {
	ID string `json:"id"`
}

type CatalogDefaultInput struct {
	ID             string `json:"id"`
	DefaultInclude bool   `json:"default_include_in_environment"`
}

type EnvironmentSelectionInput struct {
	EnvironmentID string `json:"environment_id"`
	ID            string `json:"id"`
	Enabled       bool   `json:"enabled"`
}

type EnvironmentMCPRuntimeInput struct {
	EnvironmentID string `json:"environment_id"`
	MCPID         string `json:"mcp_id"`
}

type EnvironmentMCPStatusInput struct {
	EnvironmentID string `json:"environment_id"`
	MCPID         string `json:"mcp_id"`
}

type EnvironmentMCPCallInput struct {
	EnvironmentID string         `json:"environment_id"`
	MCPID         string         `json:"mcp_id"`
	Tool          string         `json:"tool"`
	Arguments     map[string]any `json:"arguments,omitempty"`
}

type EnvironmentSkillReadInput struct {
	EnvironmentID string `json:"environment_id"`
	SkillID       string `json:"skill_id"`
	Path          string `json:"path,omitempty"`
	MaxBytes      int    `json:"max_bytes,omitempty"`
}

type EnvironmentSkillFilesInput struct {
	EnvironmentID string `json:"environment_id"`
	SkillID       string `json:"skill_id"`
	RootKind      string `json:"root_kind,omitempty" jsonschema:"artifact, support, or support:<index>; defaults to artifact"`
	MaxEntries    int    `json:"max_entries,omitempty"`
}

type EnvironmentVerifierRunInput struct {
	EnvironmentID  string `json:"environment_id"`
	WriterOwner    string `json:"writer_owner"`
	VerifierID     string `json:"verifier_id"`
	MaxOutputBytes int    `json:"max_output_bytes,omitempty"`
}

type EnvironmentVerifierAddInput struct {
	EnvironmentID string                   `json:"environment_id"`
	Definition    model.VerifierDefinition `json:"definition"`
}

type EnvironmentVerifierRemoveInput struct {
	EnvironmentID string `json:"environment_id"`
	VerifierID    string `json:"verifier_id"`
}

type MemoryKeyInput struct {
	Key string `json:"key"`
}

type MemoryWriteInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EnvironmentMemoryKeyInput struct {
	EnvironmentID string `json:"environment_id"`
	Key           string `json:"key"`
}

type EnvironmentMemoryWriteInput struct {
	EnvironmentID string `json:"environment_id"`
	Key           string `json:"key"`
	Value         string `json:"value"`
}

type WriterReleaseInput struct {
	EnvironmentID string `json:"environment_id"`
	Owner         string `json:"owner,omitempty"`
	Force         bool   `json:"force,omitempty"`
}

type TreeInput struct {
	EnvironmentID string `json:"environment_id"`
	Path          string `json:"path,omitempty"`
	MaxDepth      int    `json:"max_depth,omitempty"`
	MaxEntries    int    `json:"max_entries,omitempty"`
}

type ReadInput struct {
	EnvironmentID string `json:"environment_id"`
	Path          string `json:"path"`
	MaxBytes      int    `json:"max_bytes,omitempty"`
}

type SearchInput struct {
	EnvironmentID   string `json:"environment_id"`
	Path            string `json:"path,omitempty"`
	Query           string `json:"query"`
	MaxFiles        int    `json:"max_files,omitempty"`
	MaxMatches      int    `json:"max_matches,omitempty"`
	MaxBytesPerFile int    `json:"max_bytes_per_file,omitempty"`
}

type EndpointInvestigateInput struct {
	EnvironmentID   string `json:"environment_id"`
	Target          string `json:"target" jsonschema:"URL or path to investigate, for example /api/users/123"`
	Method          string `json:"method,omitempty" jsonschema:"optional HTTP method such as GET or POST"`
	Path            string `json:"path,omitempty" jsonschema:"optional environment-relative search root"`
	MaxFiles        int    `json:"max_files,omitempty"`
	MaxMatches      int    `json:"max_matches,omitempty"`
	MaxBytesPerFile int    `json:"max_bytes_per_file,omitempty"`
}

type WriteInput struct {
	EnvironmentID string `json:"environment_id"`
	WriterOwner   string `json:"writer_owner"`
	Path          string `json:"path"`
	Content       string `json:"content"`
	CreateParents bool   `json:"create_parents,omitempty"`
}

type EditInput struct {
	EnvironmentID        string `json:"environment_id"`
	WriterOwner          string `json:"writer_owner"`
	Path                 string `json:"path"`
	OldText              string `json:"old_text"`
	NewText              string `json:"new_text"`
	ExpectedReplacements int    `json:"expected_replacements,omitempty"`
}

type DeleteInput struct {
	EnvironmentID string `json:"environment_id"`
	WriterOwner   string `json:"writer_owner"`
	Path          string `json:"path"`
}

type ExecInput struct {
	EnvironmentID  string   `json:"environment_id"`
	WriterOwner    string   `json:"writer_owner"`
	Executable     string   `json:"executable"`
	Args           []string `json:"args,omitempty"`
	Cwd            string   `json:"cwd,omitempty"`
	TimeoutMS      int64    `json:"timeout_ms,omitempty"`
	MaxOutputBytes int      `json:"max_output_bytes,omitempty"`
}

type ProcessStartInput struct {
	EnvironmentID string   `json:"environment_id"`
	WriterOwner   string   `json:"writer_owner"`
	Executable    string   `json:"executable"`
	Args          []string `json:"args,omitempty"`
	Cwd           string   `json:"cwd,omitempty"`
	MaxLogBytes   int      `json:"max_log_bytes,omitempty"`
}

type ProcessInput struct {
	EnvironmentID string `json:"environment_id"`
	ProcessID     string `json:"process_id"`
}

type ProcessStopInput struct {
	EnvironmentID string `json:"environment_id"`
	WriterOwner   string `json:"writer_owner"`
	ProcessID     string `json:"process_id"`
}

type RunStartInput struct {
	EnvironmentID  string   `json:"environment_id"`
	WriterOwner    string   `json:"writer_owner"`
	Executable     string   `json:"executable"`
	Args           []string `json:"args,omitempty"`
	Cwd            string   `json:"cwd,omitempty"`
	TimeoutMS      int64    `json:"timeout_ms,omitempty"`
	MaxOutputBytes int      `json:"max_output_bytes,omitempty"`
}

type RunInput struct {
	EnvironmentID string `json:"environment_id"`
	RunID         string `json:"run_id"`
}

type RunCancelInput struct {
	EnvironmentID string `json:"environment_id"`
	WriterOwner   string `json:"writer_owner"`
	RunID         string `json:"run_id"`
}

type EnvironmentInfoOutput = app.EnvironmentInspection
type EnvironmentCapabilityReportOutput = model.CapabilityReport

type serverSurface string

const (
	serverSurfaceAgent serverSurface = "agent"
	serverSurfaceAdmin serverSurface = "admin"
)

func New(service *app.Service) *mcp.Server {
	return newServer(service, nil)
}

func NewAdmin(service *app.Service) *mcp.Server {
	return newServerForSurface(service, nil, serverSurfaceAdmin)
}

func newServer(service *app.Service, owner *runtimeOwner) *mcp.Server {
	return newServerForSurface(service, owner, serverSurfaceAgent)
}

func addScopedTool[In, Out any](server *mcp.Server, surface serverSurface, tool *mcp.Tool, handler mcp.ToolHandlerFor[In, Out]) {
	if surface == serverSurfaceAgent && isAdminOnlyTool(tool.Name) {
		return
	}
	mcp.AddTool(server, tool, handler)
}

func isAdminOnlyTool(name string) bool {
	switch name {
	case "management_snapshot",
		"workspace_add", "workspace_rename", "workspace_remove",
		"environment_create", "environment_rename", "environment_remove", "environment_verifier_add", "environment_verifier_remove",
		"exec_allow", "exec_allow_remove",
		"mcp_list", "mcp_add", "mcp_update", "mcp_remove", "mcp_set_default", "mcp_import_preview", "mcp_import_apply",
		"environment_mcp_set",
		"skill_list", "skill_add", "skill_remove", "skill_set_default", "skill_source_list", "skill_source_add", "skill_source_refresh", "skill_source_remove",
		"environment_skill_set",
		"memory_global_write", "memory_global_delete":
		return true
	default:
		return false
	}
}

func newServerForSurface(service *app.Service, owner *runtimeOwner, surface serverSurface) *mcp.Server {
	implementationName := serverName
	if surface == serverSurfaceAdmin {
		implementationName += "-admin"
	}
	server := mcp.NewServer(&mcp.Implementation{Name: implementationName, Version: serverVersion}, nil)

	addScopedTool(server, surface, &mcp.Tool{Name: "gateway_info", Description: "Describe the ADM V2 Agent Gateway and its core semantics."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			role := "local AI development gateway"
			if surface == serverSurfaceAdmin {
				role = "local ADM management gateway"
			}
			info := map[string]any{
				"name":        implementationName,
				"api_version": "v2-dev",
				"surface":     string(surface),
				"role":        role,
				"notes": []string{
					"Workspace is a registered local directory; Git is optional.",
					"Environment is a persistent development context; worktree is not a prerequisite.",
					"MCP/Skill catalogs and Memory are optional development-context capabilities.",
				},
			}
			if owner != nil {
				info["runtime_owner"] = owner.Info()
			}
			return toolResult(info, nil)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "management_snapshot", Description: "Return a safe Desktop/CLI management overview without exposing Global or Environment-private Memory values."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, management.Snapshot, error) {
			snapshot, err := management.New(service).Snapshot()
			return nil, snapshot, err
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "workspace_list", Description: "List local directories explicitly registered as ADM Workspaces."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Workspaces.List()
			return toolResult(items, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "workspace_add", Description: "Register an existing local directory as an ADM Workspace. Git is not required."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WorkspaceAddInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Workspaces.Add(in.Path, in.Name)
			return toolResult(item, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "workspace_inspect", Description: "Inspect one registered ADM Workspace by stable ID."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WorkspaceInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Workspaces.Get(in.WorkspaceID)
			return toolResult(item, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "workspace_rename", Description: "Change one Workspace display name without moving or renaming its directory."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WorkspaceRenameInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Workspaces.Rename(in.WorkspaceID, in.Name)
			return toolResult(item, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "workspace_remove", Description: "Remove one ADM Workspace record without deleting project files. Existing Environment references block removal."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WorkspaceInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Workspaces.Remove(in.WorkspaceID)
			return toolResult(map[string]any{"removed": item}, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "exec_allow", Description: "Allow one executable for Environment command execution."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ExecutableInput) (*mcp.CallToolResult, any, error) {
			err := service.AllowExecutable(in.Executable)
			if err != nil {
				return toolResult(nil, err)
			}
			items, err := service.AllowedExecutables()
			return toolResult(items, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "exec_allow_remove", Description: "Remove one executable from the Environment command execution allowlist."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ExecutableInput) (*mcp.CallToolResult, any, error) {
			err := service.RemoveAllowedExecutable(in.Executable)
			if err != nil {
				return toolResult(nil, err)
			}
			items, err := service.AllowedExecutables()
			return toolResult(items, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "exec_allow_list", Description: "List executables allowed for Environment command execution."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.AllowedExecutables()
			return toolResult(items, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_list", Description: "List lightweight Environment summaries. Private Memory values are omitted; only the entry count is exposed."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.EnvironmentSummaries()
			return toolResult(items, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_create", Description: "Create a development Environment for a registered Workspace. Git, worktree and verifier are not prerequisites."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentCreateInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.Create(in.WorkspaceID, in.Name, in.Root)
			return toolResult(env, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_worktree_create", Description: "Create an optional managed Git worktree Environment under the ADM-owned worktree root. The source Workspace must be a Git top-level; caller does not choose filesystem destination or branch name."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentWorktreeCreateInput) (*mcp.CallToolResult, any, error) {
			value, err := service.CreateManagedWorktree(ctx, in.WorkspaceID, in.Name, in.BaseRef)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_worktree_list", Description: "List persisted ADM-managed Git worktree records. Ordinary Environments are not included."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.ManagedWorktrees()
			return toolResult(items, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_worktree_destroy", Description: "Destroy one ADM-managed Git worktree Environment. Requires the matching writer_owner. Dirty or unpublished work is refused unless force=true; the managed branch is always retained."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentWorktreeDestroyInput) (*mcp.CallToolResult, any, error) {
			value, err := service.DestroyManagedWorktree(ctx, in.EnvironmentID, in.WriterOwner, in.Force)
			if err == nil && owner != nil {
				owner.DropEnvironment(in.EnvironmentID)
			}
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_inspect", Description: "Inspect Workspace relation, structured capability facts, resolved/unresolved MCP and Skill selections, and private Memory entry count without exposing Memory values."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, EnvironmentInfoOutput, error) {
			info, err := service.InspectEnvironment(ctx, in.EnvironmentID)
			if err != nil {
				return nil, EnvironmentInfoOutput{}, err
			}
			return nil, info, nil
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_capability_report", Description: "Return the canonical side-effect-free Environment capability report. With a Gateway runtime owner, owner-local MCP/process/run observations enrich the same CapabilityFact schema without reconnecting, probing, calling tools, or running verifiers."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, EnvironmentCapabilityReportOutput, error) {
			var (
				report model.CapabilityReport
				err    error
			)
			if owner != nil {
				report, err = owner.CapabilityReport(ctx, in.EnvironmentID)
			} else {
				report, err = service.EnvironmentCapabilityReport(ctx, in.EnvironmentID)
			}
			if err != nil {
				return nil, EnvironmentCapabilityReportOutput{}, err
			}
			return nil, report, nil
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_rename", Description: "Rename one Environment in ADM metadata only. The root directory, selections, private memory, writer state, and project files are unchanged."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentRenameInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.Rename(in.EnvironmentID, in.Name)
			return toolResult(env, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_remove", Description: "Remove one ADM Environment record without deleting its root directory or project files. Active writers block removal."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.Remove(in.EnvironmentID)
			if err == nil && owner != nil {
				owner.DropEnvironment(in.EnvironmentID)
			}
			return toolResult(map[string]any{"removed": env}, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_writer_acquire", Description: "Acquire or renew the single writer lease for an Environment physical root."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WriterAcquireInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.AcquireWriter(in.EnvironmentID, in.Owner)
			return toolResult(env, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_writer_heartbeat", Description: "Renew an active writer lease without performing a file mutation."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WriterAcquireInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.HeartbeatWriter(in.EnvironmentID, in.Owner)
			return toolResult(env, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_writer_release", Description: "Release an Environment writer lease."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WriterReleaseInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.ReleaseWriter(in.EnvironmentID, in.Owner, in.Force)
			return toolResult(env, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_verifier_add", Description: "Add one structured verifier definition for an Environment. Admin surface only; this does not modify the executable allowlist."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentVerifierAddInput) (*mcp.CallToolResult, any, error) {
			definition, err := service.AddVerifier(in.EnvironmentID, in.Definition)
			return toolResult(definition, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_verifier_remove", Description: "Remove one structured verifier definition from an Environment. Admin surface only."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentVerifierRemoveInput) (*mcp.CallToolResult, any, error) {
			definition, err := service.RemoveVerifier(in.EnvironmentID, in.VerifierID)
			return toolResult(map[string]any{"removed": definition}, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_verifier_list", Description: "List structured verifier definitions configured for one Environment. No writer is required."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			items, err := service.ListVerifiers(in.EnvironmentID)
			return toolResult(items, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_verifier_run", Description: "Run one configured Environment verifier through the existing Runtime execution policy. Requires the matching writer_owner."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentVerifierRunInput) (*mcp.CallToolResult, any, error) {
			result, err := service.RunVerifier(ctx, in.EnvironmentID, in.WriterOwner, in.VerifierID, in.MaxOutputBytes)
			return toolResult(result, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "mcp_list", Description: "List global MCP catalog entries."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.MCPs.List()
			return toolResult(items, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "mcp_add", Description: "Add one typed global MCP definition using streamable-http or stdio configuration."},
		func(_ context.Context, _ *mcp.CallToolRequest, in MCPAddInput) (*mcp.CallToolResult, any, error) {
			item, err := service.MCPs.AddMCPConfig(in.Name, catalog.MCPConfig{
				Transport:      in.Transport,
				AuthMode:       in.AuthMode,
				Endpoint:       in.Endpoint,
				HeaderRefs:     in.HeaderRefs,
				Executable:     in.Executable,
				Args:           in.Args,
				EnvRefs:        in.EnvRefs,
				HealthPolicy:   in.HealthPolicy,
				DefaultInclude: in.DefaultInclude,
			})
			return toolResult(item, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "mcp_update", Description: "Atomically update one global MCP definition while preserving its stable ID; owned runtime state is invalidated so new configuration takes effect without a Gateway restart."},
		func(_ context.Context, _ *mcp.CallToolRequest, in MCPUpdateInput) (*mcp.CallToolResult, any, error) {
			item, err := service.MCPs.UpdateMCPConfig(in.ID, in.Name, catalog.MCPConfig{
				Transport:      in.Transport,
				AuthMode:       in.AuthMode,
				Endpoint:       in.Endpoint,
				HeaderRefs:     in.HeaderRefs,
				Executable:     in.Executable,
				Args:           in.Args,
				EnvRefs:        in.EnvRefs,
				HealthPolicy:   in.HealthPolicy,
				DefaultInclude: in.DefaultInclude,
			})
			if err == nil && owner != nil {
				owner.DropMCP(in.ID)
			}
			return toolResult(item, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "mcp_import_preview", Description: "Preview sanitized MCP JSON/JSONC import candidates without persisting source content or changing Environment selections."},
		func(_ context.Context, _ *mcp.CallToolRequest, in app.MCPImportInput) (*mcp.CallToolResult, any, error) {
			preview, err := service.PreviewMCPImport(in)
			return toolResult(preview, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "mcp_import_apply", Description: "Reparse and atomically apply selected MCP JSON/JSONC candidates to the global MCP catalog using explicit conflict policy."},
		func(_ context.Context, _ *mcp.CallToolRequest, in app.MCPImportInput) (*mcp.CallToolResult, any, error) {
			result, err := service.ApplyMCPImport(in)
			if err == nil && owner != nil {
				for _, mutation := range result.Result.Mutations {
					if mutation.Action == catalog.MCPConflictUpdateByName {
						owner.DropMCP(mutation.Definition.ID)
					}
				}
			}
			return toolResult(result, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "mcp_remove", Description: "Remove one global MCP catalog entry. Existing Environment ID references are not silently rewritten."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogIDInput) (*mcp.CallToolResult, any, error) {
			err := service.MCPs.Remove(in.ID)
			if err == nil && owner != nil {
				owner.DropMCP(in.ID)
			}
			return toolResult(map[string]any{"removed": in.ID}, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "mcp_set_default", Description: "Change whether a global MCP is selected by newly created Environments."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogDefaultInput) (*mcp.CallToolResult, any, error) {
			item, err := service.MCPs.SetDefault(in.ID, in.DefaultInclude)
			return toolResult(item, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "environment_mcp_set", Description: "Enable or disable one global MCP ID for one Environment only."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentSelectionInput) (*mcp.CallToolResult, any, error) {
			env, err := service.SetEnvironmentMCP(in.EnvironmentID, in.ID, in.Enabled)
			if err == nil && owner != nil && !in.Enabled {
				owner.Drop(in.EnvironmentID, in.ID)
			}
			return toolResult(env, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_mcp_status", Description: "Probe one MCP selected for an Environment and return configured, disabled, healthy, or error status. No writer is required."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentMCPStatusInput) (*mcp.CallToolResult, app.MCPHealthStatus, error) {
			var (
				status app.MCPHealthStatus
				err    error
			)
			if owner != nil {
				status, err = owner.Status(ctx, in.EnvironmentID, in.MCPID)
			} else {
				status, err = service.ProbeMCPHealth(ctx, in.EnvironmentID, in.MCPID)
			}
			if err != nil {
				return nil, app.MCPHealthStatus{}, err
			}
			return nil, status, nil
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_mcp_inspect", Description: "Inspect sanitized desired MCP configuration and owner-local runtime health, recovery, and tool inventory evidence."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentMCPRuntimeInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			inspection, err := owner.Inspect(in.EnvironmentID, in.MCPID)
			return toolResult(inspection, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_mcp_refresh", Description: "Explicitly discard one MCP session and observation, reconnect, Ping, and refresh its bounded public tool inventory."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentMCPRuntimeInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			observation, err := owner.Refresh(ctx, in.EnvironmentID, in.MCPID)
			return toolResult(observation, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_mcp_tools", Description: "List tools from one external MCP that is enabled and healthy for the selected Environment."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentMCPRuntimeInput) (*mcp.CallToolResult, any, error) {
			if owner != nil {
				tools, err := owner.ListTools(ctx, in.EnvironmentID, in.MCPID)
				return toolResult(map[string]any{"mcp_id": in.MCPID, "tools": tools}, err)
			}
			if err := requireHealthyMCP(ctx, service, in.EnvironmentID, in.MCPID); err != nil {
				return toolResult(nil, err)
			}
			connection, err := enabledMCPConnection(ctx, service, in.EnvironmentID, in.MCPID)
			if err != nil {
				return toolResult(nil, err)
			}
			session, err := connectExternalMCP(ctx, in.MCPID, connection.endpoint, connection.headers)
			if err != nil {
				return toolResult(nil, err)
			}
			defer session.Close()
			params := &mcp.ListToolsParams{}
			var tools []*mcp.Tool
			for {
				page, err := session.ListTools(ctx, params)
				if err != nil {
					kind := app.ClassifyMCPError(err)
					if kind == "connection_failed" {
						kind = "tool_list_failed"
					}
					return toolResult(nil, &app.MCPError{MCPID: in.MCPID, ErrorKind: kind, Message: "external MCP tool listing failed"})
				}
				tools = append(tools, page.Tools...)
				if page.NextCursor == "" {
					break
				}
				params.Cursor = page.NextCursor
			}
			return toolResult(map[string]any{"mcp_id": in.MCPID, "tools": tools}, nil)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "environment_mcp_call", Description: "Call a tool on one external MCP that is enabled and healthy for the selected Environment."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentMCPCallInput) (*mcp.CallToolResult, any, error) {
			if strings.TrimSpace(in.Tool) == "" {
				return toolResult(nil, &app.MCPError{MCPID: in.MCPID, ErrorKind: "missing_tool_name", Message: "external MCP tool name is required"})
			}
			if owner != nil {
				result, err := owner.CallTool(ctx, in.EnvironmentID, in.MCPID, in.Tool, in.Arguments)
				return toolResult(result, err)
			}
			if err := requireHealthyMCP(ctx, service, in.EnvironmentID, in.MCPID); err != nil {
				return toolResult(nil, err)
			}
			connection, err := enabledMCPConnection(ctx, service, in.EnvironmentID, in.MCPID)
			if err != nil {
				return toolResult(nil, err)
			}
			session, err := connectExternalMCP(ctx, in.MCPID, connection.endpoint, connection.headers)
			if err != nil {
				return toolResult(nil, err)
			}
			defer session.Close()
			result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: in.Tool, Arguments: in.Arguments})
			if err != nil {
				kind := app.ClassifyMCPError(err)
				if kind == "connection_failed" {
					kind = "tool_call_failed"
				}
				return toolResult(nil, &app.MCPError{MCPID: in.MCPID, ErrorKind: kind, Message: "external MCP tool call failed"})
			}
			return toolResult(result, nil)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "skill_list", Description: "List global Skill catalog entries."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Skills.List()
			return toolResult(items, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "skill_add", Description: "Compatibility helper: register one explicit Skill source and refresh it once. Optional support roots authorize Skill-owned supporting files."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogAddInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Skills.AddSkillRoot(in.Root, in.SupportRoots, in.DefaultInclude)
			return toolResult(items, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "skill_source_list", Description: "List explicit persisted Skill sources without scanning host paths."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Skills.ListSkillSources()
			return toolResult(items, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "skill_source_add", Description: "Register one explicit Skill source root and support roots without refreshing or scanning outside the source."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogAddInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Skills.AddSkillSource(in.Root, in.SupportRoots, in.DefaultInclude)
			return toolResult(item, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "skill_source_refresh", Description: "Atomically refresh one Skill source and replace only that source-owned discovered Skill snapshot."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogIDInput) (*mcp.CallToolResult, any, error) {
			result, err := service.Skills.RefreshSkillSource(in.ID)
			return toolResult(result, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "skill_source_remove", Description: "Remove one Skill source and its source-owned discovered Skills while preserving unresolved Environment selections."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogIDInput) (*mcp.CallToolResult, any, error) {
			result, err := service.Skills.RemoveSkillSource(in.ID)
			return toolResult(result, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "skill_remove", Description: "Remove one global Skill catalog entry. Existing Environment ID references are not silently rewritten."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogIDInput) (*mcp.CallToolResult, any, error) {
			err := service.Skills.Remove(in.ID)
			return toolResult(map[string]any{"removed": in.ID}, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "skill_set_default", Description: "Change whether a global Skill is selected by newly created Environments."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogDefaultInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Skills.SetDefault(in.ID, in.DefaultInclude)
			return toolResult(item, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "environment_skill_set", Description: "Enable or disable one global Skill ID for one Environment only."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentSelectionInput) (*mcp.CallToolResult, any, error) {
			env, err := service.SetEnvironmentSkill(in.EnvironmentID, in.ID, in.Enabled)
			return toolResult(env, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_skill_list", Description: "List Environment-specific Skill availability for enabled selections and known disabled Skills without interpreting Skill instructions."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			availability, err := service.EnvironmentSkillAvailabilities(in.EnvironmentID)
			return toolResult(availability, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_skill_inspect", Description: "Inspect one Skill's Environment-specific availability, source/artifact facts and structured unavailability reason."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentSkillReadInput) (*mcp.CallToolResult, any, error) {
			availability, err := service.InspectEnvironmentSkill(in.EnvironmentID, in.SkillID)
			return toolResult(availability, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_skill_files", Description: "List a bounded inventory under an enabled Skill's artifact directory or one explicit support root."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentSkillFilesInput) (*mcp.CallToolResult, any, error) {
			inventory, err := service.EnvironmentSkillFiles(in.EnvironmentID, in.SkillID, in.RootKind, in.MaxEntries)
			return toolResult(inventory, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "environment_skill_read", Description: "Read an enabled Skill's SKILL.md or a file contained by its explicitly configured artifact/support roots."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentSkillReadInput) (*mcp.CallToolResult, any, error) {
			content, err := service.ReadEnvironmentSkill(in.EnvironmentID, in.SkillID, in.Path, in.MaxBytes)
			return toolResult(content, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "memory_global_list", Description: "List global durable memory entries."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Memory.GlobalList()
			return toolResult(items, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "memory_global_read", Description: "Read one global durable memory entry."},
		func(_ context.Context, _ *mcp.CallToolRequest, in MemoryKeyInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Memory.GlobalRead(in.Key)
			return toolResult(item, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "memory_global_write", Description: "Write one global durable memory entry; scope is explicit."},
		func(_ context.Context, _ *mcp.CallToolRequest, in MemoryWriteInput) (*mcp.CallToolResult, any, error) {
			err := service.Memory.GlobalWrite(in.Key, in.Value)
			if err != nil {
				return toolResult(nil, err)
			}
			item, err := service.Memory.GlobalRead(in.Key)
			return toolResult(item, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "memory_global_delete", Description: "Delete one global durable memory entry."},
		func(_ context.Context, _ *mcp.CallToolRequest, in MemoryKeyInput) (*mcp.CallToolResult, any, error) {
			err := service.Memory.GlobalDelete(in.Key)
			return toolResult(map[string]any{"deleted": in.Key}, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "memory_environment_list", Description: "List private memory entries for one Environment."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Memory.EnvironmentList(in.EnvironmentID)
			return toolResult(items, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "memory_environment_read", Description: "Read one private memory entry from one Environment only."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentMemoryKeyInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Memory.EnvironmentRead(in.EnvironmentID, in.Key)
			return toolResult(item, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "memory_environment_write", Description: "Write one Environment-private memory entry; scope is explicit."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentMemoryWriteInput) (*mcp.CallToolResult, any, error) {
			err := service.Memory.EnvironmentWrite(in.EnvironmentID, in.Key, in.Value)
			if err != nil {
				return toolResult(nil, err)
			}
			item, err := service.Memory.EnvironmentRead(in.EnvironmentID, in.Key)
			return toolResult(item, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "memory_environment_delete", Description: "Delete one private memory entry from one Environment only."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentMemoryKeyInput) (*mcp.CallToolResult, any, error) {
			err := service.Memory.EnvironmentDelete(in.EnvironmentID, in.Key)
			return toolResult(map[string]any{"deleted": in.Key}, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "tree", Description: "List files under an Environment root. Works in ordinary non-Git directories."},
		func(_ context.Context, _ *mcp.CallToolRequest, in TreeInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Tree(in.EnvironmentID, in.Path, in.MaxDepth, in.MaxEntries)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "read", Description: "Read a text file under an Environment root."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ReadInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Read(in.EnvironmentID, in.Path, in.MaxBytes)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "search", Description: "Search literal text under an Environment root."},
		func(_ context.Context, _ *mcp.CallToolRequest, in SearchInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Search(in.EnvironmentID, in.Path, in.Query, in.MaxFiles, in.MaxMatches, in.MaxBytesPerFile)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "investigate_endpoint", Description: "Resolve one URL/path to bounded static route evidence with confidence and uncertainties. Does not execute project code, call the endpoint, run verifiers, or mutate state."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EndpointInvestigateInput) (*mcp.CallToolResult, any, error) {
			report, err := service.InvestigateEndpoint(in.EnvironmentID, model.EndpointInvestigationRequest{
				Target:          in.Target,
				Method:          in.Method,
				Path:            in.Path,
				MaxFiles:        in.MaxFiles,
				MaxMatches:      in.MaxMatches,
				MaxBytesPerFile: in.MaxBytesPerFile,
			})
			return toolResult(report, err)
		})
	addScopedTool(server, surface, &mcp.Tool{Name: "write", Description: "Write a text file under an Environment root. Requires the matching writer_owner."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WriteInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Write(in.EnvironmentID, in.WriterOwner, in.Path, in.Content, in.CreateParents)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "edit", Description: "Apply an exact text replacement. Requires the matching writer_owner."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EditInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Edit(in.EnvironmentID, in.WriterOwner, in.Path, in.OldText, in.NewText, in.ExpectedReplacements)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "delete", Description: "Delete one file. Requires the matching writer_owner; recursive directory deletion is not exposed."},
		func(_ context.Context, _ *mcp.CallToolRequest, in DeleteInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Delete(in.EnvironmentID, in.WriterOwner, in.Path)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "exec", Description: "Run one explicitly allowlisted executable inside an Environment root. Requires the matching writer_owner."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in ExecInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Exec(ctx, in.EnvironmentID, in.WriterOwner, in.Executable, in.Args, in.Cwd, in.TimeoutMS, in.MaxOutputBytes)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "process_start", Description: "Start one allowlisted long-running development process owned by this Gateway. Requires the matching Environment writer_owner."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ProcessStartInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			value, err := owner.StartDevProcess(in.EnvironmentID, in.WriterOwner, in.Executable, in.Args, in.Cwd, in.MaxLogBytes)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "process_list", Description: "List development processes owned by this Gateway for one Environment."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			value, err := owner.ListDevProcesses(in.EnvironmentID)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "process_status", Description: "Inspect one Gateway-owned development process by stable ADM process identity."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ProcessInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			value, err := owner.DevProcessStatus(in.EnvironmentID, in.ProcessID)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "process_logs", Description: "Read bounded stdout/stderr tails from one Gateway-owned development process."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ProcessInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			value, err := owner.DevProcessLogs(in.EnvironmentID, in.ProcessID)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "process_stop", Description: "Stop one Gateway-owned development process. Requires the matching Environment writer_owner; arbitrary OS PIDs are not accepted."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ProcessStopInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			value, err := owner.StopDevProcess(in.EnvironmentID, in.WriterOwner, in.ProcessID)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "run_start", Description: "Start one asynchronous single-command Agent Run owned by this Gateway. Requires the matching Environment writer_owner and reuses the existing Runtime allowlist/cwd policy."},
		func(_ context.Context, _ *mcp.CallToolRequest, in RunStartInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			value, err := owner.StartAgentRun(in.EnvironmentID, in.WriterOwner, in.Executable, in.Args, in.Cwd, in.TimeoutMS, in.MaxOutputBytes)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "run_list", Description: "List Agent Runs owned by this Gateway for one Environment. Run observations are owner-local and are not persisted across restart."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			value, err := owner.ListAgentRuns(in.EnvironmentID)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "run_status", Description: "Inspect one Gateway-owned Agent Run by stable ADM run identity."},
		func(_ context.Context, _ *mcp.CallToolRequest, in RunInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			value, err := owner.AgentRunStatus(in.EnvironmentID, in.RunID)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "run_cancel", Description: "Cancel one running Gateway-owned Agent Run. Requires the matching Environment writer_owner and stable run identity."},
		func(_ context.Context, _ *mcp.CallToolRequest, in RunCancelInput) (*mcp.CallToolResult, any, error) {
			if owner == nil {
				return toolResult(nil, fmt.Errorf("persistent runtime owner is unavailable"))
			}
			value, err := owner.CancelAgentRun(in.EnvironmentID, in.WriterOwner, in.RunID)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "git_status", Description: "Optional Git status tool. Fails locally when the Environment root is not a Git repository."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			value, err := service.GitStatus(ctx, in.EnvironmentID)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "git_diff", Description: "Optional Git diff tool. Git is never required for Environment creation or file development."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			value, err := service.GitDiff(ctx, in.EnvironmentID)
			return toolResult(value, err)
		})

	addScopedTool(server, surface, &mcp.Tool{Name: "git_branch", Description: "Optional Git branch tool."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			value, err := service.GitBranch(ctx, in.EnvironmentID)
			return toolResult(value, err)
		})

	return server
}

type externalMCPConnection struct {
	endpoint string
	headers  map[string]string
}

func requireHealthyMCP(ctx context.Context, service *app.Service, environmentID, mcpID string) error {
	status, err := service.ProbeMCPHealth(ctx, environmentID, mcpID)
	if err != nil {
		return err
	}
	if status.State == app.MCPHealthHealthy {
		return nil
	}
	kind := status.ErrorKind
	if kind == "" {
		if status.State == app.MCPHealthDisabled {
			kind = "not_enabled"
		} else {
			kind = string(status.State)
		}
	}
	message := status.Message
	if strings.TrimSpace(message) == "" {
		message = fmt.Sprintf("external MCP is %s", status.State)
	}
	return &app.MCPError{MCPID: mcpID, ErrorKind: kind, Message: message}
}

func enabledMCPConnection(_ context.Context, service *app.Service, environmentID, mcpID string) (externalMCPConnection, error) {
	activation, status, err := service.ResolveMCPActivation(environmentID, mcpID)
	if err != nil {
		return externalMCPConnection{}, err
	}
	if activation == nil {
		return externalMCPConnection{}, statusAsMCPError(status)
	}
	return externalMCPConnection{endpoint: activation.Endpoint, headers: activation.Headers}, nil
}

func connectExternalMCP(ctx context.Context, mcpID, endpoint string, headers map[string]string) (*mcp.ClientSession, error) {
	client := mcp.NewClient(&mcp.Implementation{Name: serverName + "-proxy", Version: serverVersion}, nil)
	transport := &mcp.StreamableClientTransport{
		Endpoint:             endpoint,
		MaxRetries:           -1,
		DisableStandaloneSSE: true,
	}
	if len(headers) != 0 {
		transport.HTTPClient = &http.Client{Transport: externalMCPHeaderRoundTripper{base: http.DefaultTransport, headers: headers}}
	}
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, &app.MCPError{MCPID: mcpID, ErrorKind: app.ClassifyMCPError(err), Message: "external MCP connection failed"}
	}
	return session, nil
}

func connectStdioMCP(ctx, commandCtx context.Context, service *app.Service, environmentID string, activation *app.MCPActivation) (*mcp.ClientSession, error) {
	cmd, err := service.MCPCommand(commandCtx, environmentID, activation)
	if err != nil {
		return nil, &app.MCPError{MCPID: activation.MCPID, ErrorKind: "executable_not_allowed", Message: "stdio MCP executable is unavailable under Environment authority"}
	}
	client := mcp.NewClient(&mcp.Implementation{Name: serverName + "-proxy", Version: serverVersion}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		return nil, &app.MCPError{MCPID: activation.MCPID, ErrorKind: app.ClassifyMCPError(err), Message: "external MCP stdio connection failed"}
	}
	return session, nil
}

type externalMCPHeaderRoundTripper struct {
	base    http.RoundTripper
	headers map[string]string
}

func (r externalMCPHeaderRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	clone := request.Clone(request.Context())
	clone.Header = request.Header.Clone()
	for key, value := range r.headers {
		clone.Header.Set(key, value)
	}
	return r.base.RoundTrip(clone)
}

func RunStdio(ctx context.Context, service *app.Service) error {
	owner := newRuntimeOwner(service)
	ownerCtx, cancelOwner := context.WithCancel(ctx)
	defer owner.Close()
	defer cancelOwner()
	go owner.Reconcile(ownerCtx)
	return newServer(service, owner).Run(ctx, &mcp.StdioTransport{})
}

func NewHTTPHandler(service *app.Service) http.Handler {
	return newHTTPHandler(service, nil)
}

func newHTTPHandler(service *app.Service, owner *runtimeOwner) http.Handler {
	return newHTTPHandlerWithShutdown(service, owner, nil)
}

func newHTTPHandlerWithShutdown(service *app.Service, owner *runtimeOwner, shutdown func()) http.Handler {
	agentServer := newServer(service, owner)
	adminServer := newServerForSurface(service, owner, serverSurfaceAdmin)
	newSurfaceHandler := func(server *mcp.Server) http.Handler {
		base := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, &mcp.StreamableHTTPOptions{Stateless: true, DisableLocalhostProtection: true})
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !allowedGatewayHost(r.Host) {
				http.Error(w, "Forbidden: invalid Host header", http.StatusForbidden)
				return
			}
			base.ServeHTTP(w, r)
		})
	}
	mux := http.NewServeMux()
	mux.Handle("/mcp", newSurfaceHandler(agentServer))
	mux.Handle("/admin/mcp", newSurfaceHandler(adminServer))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		ownerID := ""
		if owner != nil {
			ownerID = owner.Info().ID
		}
		_, _ = fmt.Fprintf(w, `{"name":%q,"version":%q,"status":"ok","pid":%d,"transport":"http","owner_id":%q}`, serverName, serverVersion, os.Getpid(), ownerID)
	})
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		if owner == nil || shutdown == nil {
			http.NotFound(w, r)
			return
		}
		expectedOwnerID := owner.Info().ID
		if providedOwnerID := strings.TrimSpace(r.Header.Get(runtimeOwnerHeader)); providedOwnerID == "" || providedOwnerID != expectedOwnerID {
			http.Error(w, "runtime owner mismatch", http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		go shutdown()
	})
	return mux
}

func allowedGatewayHost(authority string) bool {
	host := strings.TrimSpace(authority)
	if parsed, _, err := net.SplitHostPort(host); err == nil {
		host = parsed
	}
	host = strings.Trim(strings.ToLower(host), "[]")
	switch host {
	case "localhost", "127.0.0.1", "::1", "host.docker.internal":
		return true
	default:
		return false
	}
}

func RunHTTP(ctx context.Context, service *app.Service, listen string) error {
	listen = strings.TrimSpace(listen)
	if listen == "" {
		return fmt.Errorf("gateway listen address is required")
	}
	host, _, err := net.SplitHostPort(listen)
	if err != nil {
		return fmt.Errorf("invalid gateway listen address %q: %w", listen, err)
	}
	if host != "localhost" {
		ip := net.ParseIP(host)
		if ip == nil || !ip.IsLoopback() {
			return fmt.Errorf("gateway HTTP listen must be loopback; got %q", listen)
		}
	}

	runCtx, cancelRun := context.WithCancel(ctx)
	defer cancelRun()
	owner := newRuntimeOwner(service)
	ownerCtx, cancelOwner := context.WithCancel(runCtx)
	defer owner.Close()
	defer cancelOwner()
	go owner.Reconcile(ownerCtx)

	httpServer := &http.Server{
		Addr:              listen,
		Handler:           newHTTPHandlerWithShutdown(service, owner, cancelRun),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		<-runCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
	}()
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func toolResult(value any, err error) (*mcp.CallToolResult, any, error) {
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}
	return nil, map[string]any{"result": value}, nil
}

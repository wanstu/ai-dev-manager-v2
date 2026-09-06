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
	"ai-dev-manager-v2/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "ai-dev-manager-v2"
	serverVersion = "v0.1.0-dev"
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

type WriterAcquireInput struct {
	EnvironmentID string `json:"environment_id"`
	Owner         string `json:"owner" jsonschema:"stable agent/session owner identifier"`
}

type CatalogAddInput struct {
	Name           string `json:"name"`
	DefaultInclude bool   `json:"default_include_in_environment,omitempty"`
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

type EnvironmentInfoOutput struct {
	Environment  model.Environment `json:"environment"`
	Capabilities []string          `json:"capabilities"`
}

func New(service *app.Service) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: serverName, Version: serverVersion}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "gateway_info", Description: "Describe the ADM V2 Agent Gateway and its core semantics."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			return toolResult(map[string]any{
				"name":        serverName,
				"api_version": "v2-dev",
				"role":        "local AI development gateway",
				"notes": []string{
					"Workspace is a registered local directory; Git is optional.",
					"Environment is a persistent development context; worktree is not a prerequisite.",
					"MCP/Skill catalogs and Memory are optional development-context capabilities.",
				},
			}, nil)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "workspace_list", Description: "List local directories explicitly registered as ADM Workspaces."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Workspaces.List()
			return toolResult(items, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "workspace_add", Description: "Register an existing local directory as an ADM Workspace. Git is not required."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WorkspaceAddInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Workspaces.Add(in.Path, in.Name)
			return toolResult(item, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "workspace_inspect", Description: "Inspect one registered ADM Workspace by stable ID."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WorkspaceInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Workspaces.Get(in.WorkspaceID)
			return toolResult(item, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "workspace_rename", Description: "Change one Workspace display name without moving or renaming its directory."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WorkspaceRenameInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Workspaces.Rename(in.WorkspaceID, in.Name)
			return toolResult(item, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "workspace_remove", Description: "Remove one ADM Workspace record without deleting project files. Existing Environment references block removal."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WorkspaceInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Workspaces.Remove(in.WorkspaceID)
			return toolResult(map[string]any{"removed": item}, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "exec_allow", Description: "Allow one executable for Environment command execution."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ExecutableInput) (*mcp.CallToolResult, any, error) {
			err := service.AllowExecutable(in.Executable)
			if err != nil {
				return toolResult(nil, err)
			}
			items, err := service.AllowedExecutables()
			return toolResult(items, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "exec_allow_remove", Description: "Remove one executable from the Environment command execution allowlist."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ExecutableInput) (*mcp.CallToolResult, any, error) {
			err := service.RemoveAllowedExecutable(in.Executable)
			if err != nil {
				return toolResult(nil, err)
			}
			items, err := service.AllowedExecutables()
			return toolResult(items, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "exec_allow_list", Description: "List executables allowed for Environment command execution."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.AllowedExecutables()
			return toolResult(items, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "environment_list", Description: "List persistent ADM development contexts. Environments do not require Git."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Environments.List()
			return toolResult(items, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "environment_create", Description: "Create a development Environment for a registered Workspace. Git, worktree and verifier are not prerequisites."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentCreateInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.Create(in.WorkspaceID, in.Name, in.Root)
			return toolResult(env, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "environment_inspect", Description: "Inspect an Environment and the capabilities currently available at its root."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, EnvironmentInfoOutput, error) {
			env, err := service.Environments.Get(in.EnvironmentID)
			if err != nil {
				return nil, EnvironmentInfoOutput{}, err
			}
			caps, err := service.Capabilities(ctx, in.EnvironmentID)
			if err != nil {
				return nil, EnvironmentInfoOutput{}, err
			}
			return nil, EnvironmentInfoOutput{Environment: env, Capabilities: caps}, nil
		})

	mcp.AddTool(server, &mcp.Tool{Name: "environment_rename", Description: "Rename one Environment in ADM metadata only. The root directory, selections, private memory, writer state, and project files are unchanged."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentRenameInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.Rename(in.EnvironmentID, in.Name)
			return toolResult(env, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "environment_remove", Description: "Remove one ADM Environment record without deleting its root directory or project files. Active writers block removal."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.Remove(in.EnvironmentID)
			return toolResult(map[string]any{"removed": env}, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "environment_writer_acquire", Description: "Acquire or renew the single writer lease for an Environment physical root."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WriterAcquireInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.AcquireWriter(in.EnvironmentID, in.Owner)
			return toolResult(env, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "environment_writer_heartbeat", Description: "Renew an active writer lease without performing a file mutation."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WriterAcquireInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.HeartbeatWriter(in.EnvironmentID, in.Owner)
			return toolResult(env, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "environment_writer_release", Description: "Release an Environment writer lease."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WriterReleaseInput) (*mcp.CallToolResult, any, error) {
			env, err := service.Environments.ReleaseWriter(in.EnvironmentID, in.Owner, in.Force)
			return toolResult(env, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "mcp_list", Description: "List global MCP catalog entries."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.MCPs.List()
			return toolResult(items, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "mcp_add", Description: "Add one global MCP catalog entry."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogAddInput) (*mcp.CallToolResult, any, error) {
			item, err := service.MCPs.Add(in.Name, in.DefaultInclude)
			return toolResult(item, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "mcp_remove", Description: "Remove one global MCP catalog entry. Existing Environment ID references are not silently rewritten."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogIDInput) (*mcp.CallToolResult, any, error) {
			err := service.MCPs.Remove(in.ID)
			return toolResult(map[string]any{"removed": in.ID}, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "mcp_set_default", Description: "Change whether a global MCP is selected by newly created Environments."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogDefaultInput) (*mcp.CallToolResult, any, error) {
			item, err := service.MCPs.SetDefault(in.ID, in.DefaultInclude)
			return toolResult(item, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "environment_mcp_set", Description: "Enable or disable one global MCP ID for one Environment only."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentSelectionInput) (*mcp.CallToolResult, any, error) {
			env, err := service.SetEnvironmentMCP(in.EnvironmentID, in.ID, in.Enabled)
			return toolResult(env, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "skill_list", Description: "List global Skill catalog entries."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Skills.List()
			return toolResult(items, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "skill_add", Description: "Add one global Skill catalog entry."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogAddInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Skills.Add(in.Name, in.DefaultInclude)
			return toolResult(item, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "skill_remove", Description: "Remove one global Skill catalog entry. Existing Environment ID references are not silently rewritten."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogIDInput) (*mcp.CallToolResult, any, error) {
			err := service.Skills.Remove(in.ID)
			return toolResult(map[string]any{"removed": in.ID}, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "skill_set_default", Description: "Change whether a global Skill is selected by newly created Environments."},
		func(_ context.Context, _ *mcp.CallToolRequest, in CatalogDefaultInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Skills.SetDefault(in.ID, in.DefaultInclude)
			return toolResult(item, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "environment_skill_set", Description: "Enable or disable one global Skill ID for one Environment only."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentSelectionInput) (*mcp.CallToolResult, any, error) {
			env, err := service.SetEnvironmentSkill(in.EnvironmentID, in.ID, in.Enabled)
			return toolResult(env, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "memory_global_list", Description: "List global durable memory entries."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Memory.GlobalList()
			return toolResult(items, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "memory_global_read", Description: "Read one global durable memory entry."},
		func(_ context.Context, _ *mcp.CallToolRequest, in MemoryKeyInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Memory.GlobalRead(in.Key)
			return toolResult(item, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "memory_global_write", Description: "Write one global durable memory entry; scope is explicit."},
		func(_ context.Context, _ *mcp.CallToolRequest, in MemoryWriteInput) (*mcp.CallToolResult, any, error) {
			err := service.Memory.GlobalWrite(in.Key, in.Value)
			if err != nil {
				return toolResult(nil, err)
			}
			item, err := service.Memory.GlobalRead(in.Key)
			return toolResult(item, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "memory_global_delete", Description: "Delete one global durable memory entry."},
		func(_ context.Context, _ *mcp.CallToolRequest, in MemoryKeyInput) (*mcp.CallToolResult, any, error) {
			err := service.Memory.GlobalDelete(in.Key)
			return toolResult(map[string]any{"deleted": in.Key}, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "memory_environment_list", Description: "List private memory entries for one Environment."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			items, err := service.Memory.EnvironmentList(in.EnvironmentID)
			return toolResult(items, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "memory_environment_read", Description: "Read one private memory entry from one Environment only."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentMemoryKeyInput) (*mcp.CallToolResult, any, error) {
			item, err := service.Memory.EnvironmentRead(in.EnvironmentID, in.Key)
			return toolResult(item, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "memory_environment_write", Description: "Write one Environment-private memory entry; scope is explicit."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentMemoryWriteInput) (*mcp.CallToolResult, any, error) {
			err := service.Memory.EnvironmentWrite(in.EnvironmentID, in.Key, in.Value)
			if err != nil {
				return toolResult(nil, err)
			}
			item, err := service.Memory.EnvironmentRead(in.EnvironmentID, in.Key)
			return toolResult(item, err)
		})
	mcp.AddTool(server, &mcp.Tool{Name: "memory_environment_delete", Description: "Delete one private memory entry from one Environment only."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EnvironmentMemoryKeyInput) (*mcp.CallToolResult, any, error) {
			err := service.Memory.EnvironmentDelete(in.EnvironmentID, in.Key)
			return toolResult(map[string]any{"deleted": in.Key}, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "tree", Description: "List files under an Environment root. Works in ordinary non-Git directories."},
		func(_ context.Context, _ *mcp.CallToolRequest, in TreeInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Tree(in.EnvironmentID, in.Path, in.MaxDepth, in.MaxEntries)
			return toolResult(value, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "read", Description: "Read a text file under an Environment root."},
		func(_ context.Context, _ *mcp.CallToolRequest, in ReadInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Read(in.EnvironmentID, in.Path, in.MaxBytes)
			return toolResult(value, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "search", Description: "Search literal text under an Environment root."},
		func(_ context.Context, _ *mcp.CallToolRequest, in SearchInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Search(in.EnvironmentID, in.Path, in.Query, in.MaxFiles, in.MaxMatches, in.MaxBytesPerFile)
			return toolResult(value, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "write", Description: "Write a text file under an Environment root. Requires the matching writer_owner."},
		func(_ context.Context, _ *mcp.CallToolRequest, in WriteInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Write(in.EnvironmentID, in.WriterOwner, in.Path, in.Content, in.CreateParents)
			return toolResult(value, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "edit", Description: "Apply an exact text replacement. Requires the matching writer_owner."},
		func(_ context.Context, _ *mcp.CallToolRequest, in EditInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Edit(in.EnvironmentID, in.WriterOwner, in.Path, in.OldText, in.NewText, in.ExpectedReplacements)
			return toolResult(value, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "delete", Description: "Delete one file. Requires the matching writer_owner; recursive directory deletion is not exposed."},
		func(_ context.Context, _ *mcp.CallToolRequest, in DeleteInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Delete(in.EnvironmentID, in.WriterOwner, in.Path)
			return toolResult(value, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "exec", Description: "Run one explicitly allowlisted executable inside an Environment root. Requires the matching writer_owner."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in ExecInput) (*mcp.CallToolResult, any, error) {
			value, err := service.Exec(ctx, in.EnvironmentID, in.WriterOwner, in.Executable, in.Args, in.Cwd, in.TimeoutMS, in.MaxOutputBytes)
			return toolResult(value, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "git_status", Description: "Optional Git status tool. Fails locally when the Environment root is not a Git repository."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			value, err := service.GitStatus(ctx, in.EnvironmentID)
			return toolResult(value, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "git_diff", Description: "Optional Git diff tool. Git is never required for Environment creation or file development."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			value, err := service.GitDiff(ctx, in.EnvironmentID)
			return toolResult(value, err)
		})

	mcp.AddTool(server, &mcp.Tool{Name: "git_branch", Description: "Optional Git branch tool."},
		func(ctx context.Context, _ *mcp.CallToolRequest, in EnvironmentInput) (*mcp.CallToolResult, any, error) {
			value, err := service.GitBranch(ctx, in.EnvironmentID)
			return toolResult(value, err)
		})

	return server
}

func RunStdio(ctx context.Context, service *app.Service) error {
	return New(service).Run(ctx, &mcp.StdioTransport{})
}

func NewHTTPHandler(service *app.Service) http.Handler {
	server := New(service)
	base := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{Stateless: true, DisableLocalhostProtection: true})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !allowedGatewayHost(r.Host) {
			http.Error(w, "Forbidden: invalid Host header", http.StatusForbidden)
			return
		}
		base.ServeHTTP(w, r)
	})
	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"name":%q,"version":%q,"status":"ok","pid":%d,"transport":"http"}`, serverName, serverVersion, os.Getpid())
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

	httpServer := &http.Server{
		Addr:              listen,
		Handler:           NewHTTPHandler(service),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		<-ctx.Done()
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

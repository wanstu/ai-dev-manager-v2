package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/gateway"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/store"
)

const defaultGatewayListen = gateway.DefaultHTTPListen

type gatewayHealth struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	PID       int    `json:"pid"`
	Transport string `json:"transport"`
	OwnerID   string `json:"owner_id,omitempty"`
}

type incompatibleGatewayError struct{ detail string }

func (e *incompatibleGatewayError) Error() string { return e.detail }

var (
	matchesADMExecutable = sameADMExecutable
	startGatewayDetached = startDetachedHTTPGateway
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "错误：", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	statePath, err := store.DefaultPath()
	if err != nil {
		return err
	}
	service := app.New(statePath)
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "workspace":
		return runWorkspace(service, args[1:])
	case "environment", "env":
		return runEnvironment(service, args[1:])
	case "exec":
		return runExec(service, args[1:])
	case "mcp":
		return runCatalog("mcp", service, service.MCPs, args[1:])
	case "skill":
		return runCatalog("skill", service, service.Skills, args[1:])
	case "memory":
		return runMemory(service, args[1:])
	case "gateway":
		return runGateway(service, args[1:])
	case "doctor":
		return runDoctor(service, statePath, args[1:])
	case "state":
		return runState(statePath, args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("未知命令 %q；运行 ai-dev-manager-v2 -h 查看帮助", args[0])
	}
}

func runWorkspace(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printWorkspaceHelp()
		return nil
	}
	switch args[0] {
	case "add":
		fs := newFlagSet("workspace add", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 workspace add --path PATH [--name NAME]")
			fmt.Fprintln(os.Stdout, "\n登记一个现有本地目录为 Workspace；不要求 Git。")
		})
		path := fs.String("path", "", "Workspace 目录路径")
		name := fs.String("name", "", "Workspace 显示名称")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*path) == "" {
			return fmt.Errorf("缺少 --path；运行 ai-dev-manager-v2 workspace add -h 查看帮助")
		}
		ws, err := service.Workspaces.Add(*path, *name)
		if err != nil {
			return err
		}
		return writeJSON(ws)
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("workspace list 不接受参数")
		}
		items, err := service.Workspaces.List()
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "inspect":
		fs := newFlagSet("workspace inspect", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 workspace inspect --workspace-id WS_ID")
		})
		workspaceID := fs.String("workspace-id", "", "Workspace ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*workspaceID) == "" {
			return fmt.Errorf("缺少 --workspace-id；运行 ai-dev-manager-v2 workspace inspect -h 查看帮助")
		}
		ws, err := service.Workspaces.Get(*workspaceID)
		if err != nil {
			return err
		}
		return writeJSON(ws)
	case "rename":
		fs := newFlagSet("workspace rename", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 workspace rename --workspace-id WS_ID --name NAME")
			fmt.Fprintln(os.Stdout, "\n只修改 ADM 中的显示名称，不移动或重命名项目目录。")
		})
		workspaceID := fs.String("workspace-id", "", "Workspace ID")
		name := fs.String("name", "", "新的 Workspace 显示名称")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*workspaceID) == "" || strings.TrimSpace(*name) == "" {
			return fmt.Errorf("必须提供 --workspace-id 和 --name；运行 ai-dev-manager-v2 workspace rename -h 查看帮助")
		}
		ws, err := service.Workspaces.Rename(*workspaceID, *name)
		if err != nil {
			return err
		}
		return writeJSON(ws)
	case "remove":
		fs := newFlagSet("workspace remove", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 workspace remove --workspace-id WS_ID")
			fmt.Fprintln(os.Stdout, "\n只删除 ADM 中的 Workspace 记录，不会删除项目目录或文件；仍有 Environment 引用时禁止删除。")
		})
		workspaceID := fs.String("workspace-id", "", "Workspace ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*workspaceID) == "" {
			return fmt.Errorf("缺少 --workspace-id；运行 ai-dev-manager-v2 workspace remove -h 查看帮助")
		}
		removed, err := service.Workspaces.Remove(*workspaceID)
		if err != nil {
			return err
		}
		return writeJSON(map[string]any{"removed": removed})
	default:
		return fmt.Errorf("未知 workspace 命令 %q；运行 ai-dev-manager-v2 workspace -h 查看帮助", args[0])
	}
}

func runEnvironment(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printEnvironmentHelp()
		return nil
	}
	switch args[0] {
	case "create":
		fs := newFlagSet("environment create", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 environment create --workspace-id WS_ID --name NAME [--root PATH]")
			fmt.Fprintln(os.Stdout, "\n创建持久开发上下文；不写 --root 时默认使用整个 Workspace。")
		})
		workspaceID := fs.String("workspace-id", "", "Workspace ID")
		name := fs.String("name", "", "Environment 名称")
		root := fs.String("root", "", "Workspace 内的可选根目录；默认使用整个 Workspace")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*workspaceID) == "" || strings.TrimSpace(*name) == "" {
			return fmt.Errorf("必须提供 --workspace-id 和 --name；运行 ai-dev-manager-v2 environment create -h 查看帮助")
		}
		env, err := service.Environments.Create(*workspaceID, *name, *root)
		if err != nil {
			return err
		}
		return writeJSON(env)
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("environment list 不接受参数")
		}
		items, err := service.EnvironmentSummaries()
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "inspect":
		fs := newFlagSet("environment inspect", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 environment inspect --environment-id ENV_ID")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" {
			return fmt.Errorf("缺少 --environment-id；运行 ai-dev-manager-v2 environment inspect -h 查看帮助")
		}
		info, err := service.InspectEnvironment(context.Background(), *environmentID)
		if err != nil {
			return err
		}
		return writeJSON(info)
	case "rename":
		fs := newFlagSet("environment rename", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 environment rename --environment-id ENV_ID --name NAME")
			fmt.Fprintln(os.Stdout, "\n只修改 ADM 中的 Environment 显示名称，不移动根目录、不修改选择或 Memory，也不触碰项目文件。")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		name := fs.String("name", "", "新的 Environment 显示名称")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" || strings.TrimSpace(*name) == "" {
			return fmt.Errorf("必须提供 --environment-id 和 --name；运行 ai-dev-manager-v2 environment rename -h 查看帮助")
		}
		env, err := service.Environments.Rename(*environmentID, *name)
		if err != nil {
			return err
		}
		return writeJSON(env)
	case "remove":
		fs := newFlagSet("environment remove", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 environment remove --environment-id ENV_ID")
			fmt.Fprintln(os.Stdout, "\n只删除 ADM 中的 Environment 记录，不会删除根目录或项目文件；存在有效 Writer 时禁止删除。")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" {
			return fmt.Errorf("缺少 --environment-id；运行 ai-dev-manager-v2 environment remove -h 查看帮助")
		}
		removed, err := service.Environments.Remove(*environmentID)
		if err != nil {
			return err
		}
		return writeJSON(map[string]any{"removed": removed})
	case "mcp", "skill":
		return runEnvironmentSelection(args[0], service, args[1:])
	case "verifier":
		return runEnvironmentVerifier(service, args[1:])
	case "writer":
		return runWriter(service, args[1:])
	default:
		return fmt.Errorf("未知 environment 命令 %q；运行 ai-dev-manager-v2 environment -h 查看帮助", args[0])
	}
}

func runEnvironmentVerifier(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printEnvironmentVerifierHelp()
		return nil
	}
	switch args[0] {
	case "add":
		fs := newFlagSet("environment verifier add", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 environment verifier add --environment-id ENV_ID --kind test|lint|build|custom --executable NAME_OR_PATH [--name NAME] [--arg ARG ...] [--cwd RELATIVE_PATH] [--timeout-seconds N] [--enabled=true|false]")
			fmt.Fprintln(os.Stdout, "\n只声明 Environment-scoped verifier；不会把 executable 加入执行白名单。")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		name := fs.String("name", "", "可选的人类可读名称")
		kind := fs.String("kind", "", "verifier 类型：test、lint、build、custom")
		executable := fs.String("executable", "", "程序名或绝对路径；运行时仍必须在全局执行白名单中")
		cwd := fs.String("cwd", "", "Environment 内的相对工作目录；空值表示 Environment root")
		timeoutSeconds := fs.Int64("timeout-seconds", 0, "超时秒数；0 使用默认值")
		enabled := fs.Bool("enabled", true, "是否启用；默认 true")
		var verifierArgs []string
		fs.Func("arg", "传给 verifier 的一个参数；可重复", func(value string) error {
			verifierArgs = append(verifierArgs, value)
			return nil
		})
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" || strings.TrimSpace(*kind) == "" || strings.TrimSpace(*executable) == "" {
			return fmt.Errorf("必须提供 --environment-id、--kind 和 --executable；运行 ai-dev-manager-v2 environment verifier add -h 查看帮助")
		}
		definition, err := service.AddVerifier(*environmentID, model.VerifierDefinition{
			Name:           *name,
			Kind:           *kind,
			Enabled:        *enabled,
			Executable:     *executable,
			Args:           verifierArgs,
			Cwd:            *cwd,
			TimeoutSeconds: *timeoutSeconds,
		})
		if err != nil {
			return err
		}
		return writeJSON(definition)
	case "list":
		fs := newFlagSet("environment verifier list", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 environment verifier list --environment-id ENV_ID")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" {
			return fmt.Errorf("缺少 --environment-id；运行 ai-dev-manager-v2 environment verifier list -h 查看帮助")
		}
		items, err := service.ListVerifiers(*environmentID)
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "remove":
		fs := newFlagSet("environment verifier remove", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 environment verifier remove --environment-id ENV_ID --verifier-id VF_ID")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		verifierID := fs.String("verifier-id", "", "Verifier ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" || strings.TrimSpace(*verifierID) == "" {
			return fmt.Errorf("必须提供 --environment-id 和 --verifier-id；运行 ai-dev-manager-v2 environment verifier remove -h 查看帮助")
		}
		removed, err := service.RemoveVerifier(*environmentID, *verifierID)
		if err != nil {
			return err
		}
		return writeJSON(map[string]any{"removed": removed})
	default:
		return fmt.Errorf("未知 environment verifier 命令 %q；运行 ai-dev-manager-v2 environment verifier -h 查看帮助", args[0])
	}
}

func runEnvironmentSelection(kind string, service *app.Service, args []string) error {
	if wantsHelp(args) {
		printEnvironmentSelectionHelp(kind)
		return nil
	}
	if kind != "mcp" && kind != "skill" {
		return fmt.Errorf("unsupported Environment selection kind %q", kind)
	}
	action := args[0]
	if action != "enable" && action != "disable" {
		return fmt.Errorf("未知 environment %s 命令 %q；运行 ai-dev-manager-v2 environment %s -h 查看帮助", kind, action, kind)
	}
	fs := newFlagSet("environment "+kind+" "+action, func() {
		fmt.Fprintf(os.Stdout, "用法：ai-dev-manager-v2 environment %s %s --environment-id ENV_ID --%s-id ID\n", kind, action, kind)
	})
	environmentID := fs.String("environment-id", "", "Environment ID")
	entryID := fs.String(kind+"-id", "", strings.ToUpper(kind)+" catalog ID")
	if err := fs.Parse(args[1:]); err != nil {
		return flagError(err)
	}
	if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" || strings.TrimSpace(*entryID) == "" {
		return fmt.Errorf("必须提供 --environment-id 和 --%s-id；运行 ai-dev-manager-v2 environment %s %s -h 查看帮助", kind, kind, action)
	}
	enabled := action == "enable"
	if kind == "mcp" {
		env, err := service.SetEnvironmentMCP(*environmentID, *entryID, enabled)
		if err != nil {
			return err
		}
		return writeJSON(env)
	}
	env, err := service.SetEnvironmentSkill(*environmentID, *entryID, enabled)
	if err != nil {
		return err
	}
	return writeJSON(env)
}

func runWriter(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printWriterHelp()
		return nil
	}
	switch args[0] {
	case "acquire", "heartbeat":
		action := args[0]
		fs := newFlagSet("environment writer "+action, func() {
			fmt.Fprintf(os.Stdout, "用法：ai-dev-manager-v2 environment writer %s --environment-id ENV_ID --owner OWNER\n", action)
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		owner := fs.String("owner", "", "稳定的 Agent/会话 Writer 标识")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" || strings.TrimSpace(*owner) == "" {
			return fmt.Errorf("必须提供 --environment-id 和 --owner；运行 ai-dev-manager-v2 environment writer %s -h 查看帮助", action)
		}
		var (
			env any
			err error
		)
		if action == "acquire" {
			env, err = service.Environments.AcquireWriter(*environmentID, *owner)
		} else {
			env, err = service.Environments.HeartbeatWriter(*environmentID, *owner)
		}
		if err != nil {
			return err
		}
		return writeJSON(env)
	case "release":
		fs := newFlagSet("environment writer release", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 environment writer release --environment-id ENV_ID (--owner OWNER | --force)")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		owner := fs.String("owner", "", "当前 Writer owner")
		force := fs.Bool("force", false, "忽略 owner，强制释放 Writer")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" || (!*force && strings.TrimSpace(*owner) == "") {
			return fmt.Errorf("必须提供 --environment-id，并提供 --owner 或 --force；运行 ai-dev-manager-v2 environment writer release -h 查看帮助")
		}
		env, err := service.Environments.ReleaseWriter(*environmentID, *owner, *force)
		if err != nil {
			return err
		}
		return writeJSON(env)
	default:
		return fmt.Errorf("未知 writer 命令 %q；运行 ai-dev-manager-v2 environment writer -h 查看帮助", args[0])
	}
}

func runExec(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printExecHelp()
		return nil
	}
	switch args[0] {
	case "allow":
		fs := newFlagSet("exec allow", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 exec allow --executable NAME_OR_PATH")
		})
		executable := fs.String("executable", "", "程序名或绝对路径")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*executable) == "" {
			return fmt.Errorf("缺少 --executable；运行 ai-dev-manager-v2 exec allow -h 查看帮助")
		}
		if err := service.AllowExecutable(*executable); err != nil {
			return err
		}
		items, err := service.AllowedExecutables()
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "remove":
		fs := newFlagSet("exec remove", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 exec remove --executable NAME_OR_PATH")
			fmt.Fprintln(os.Stdout, "\n从执行白名单移除一个程序；后续 Environment exec 将立即按新的白名单判断。")
		})
		executable := fs.String("executable", "", "程序名或绝对路径")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*executable) == "" {
			return fmt.Errorf("缺少 --executable；运行 ai-dev-manager-v2 exec remove -h 查看帮助")
		}
		if err := service.RemoveAllowedExecutable(*executable); err != nil {
			return err
		}
		items, err := service.AllowedExecutables()
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("exec list 不接受参数")
		}
		items, err := service.AllowedExecutables()
		if err != nil {
			return err
		}
		return writeJSON(items)
	default:
		return fmt.Errorf("未知 exec 命令 %q；运行 ai-dev-manager-v2 exec -h 查看帮助", args[0])
	}
}

func runCatalog(kind string, application *app.Service, service any, args []string) error {
	label := "MCP"
	if kind == "skill" {
		label = "Skill"
	} else if kind != "mcp" {
		return fmt.Errorf("unknown catalog kind %q", kind)
	}
	if wantsHelp(args) {
		printCatalogHelp(kind)
		return nil
	}
	mcpService, _ := service.(*catalog.MCPService)
	skillService, _ := service.(*catalog.Service)
	switch args[0] {
	case "add":
		if kind == "mcp" {
			fs := newFlagSet("mcp add", func() {
				fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 mcp add --name NAME [--transport streamable-http --endpoint URL | --transport stdio --executable PATH] [选项]")
				fmt.Fprintln(os.Stdout, "\n添加 typed MCP 定义。args/header/env refs 使用 JSON；secret 值必须写成环境变量引用。")
			})
			name := fs.String("name", "", "MCP 名称")
			transport := fs.String("transport", catalog.MCPTransportStreamableHTTP, "streamable-http 或 stdio")
			authMode := fs.String("auth-mode", catalog.MCPAuthNone, "none 或 headers")
			endpoint := fs.String("endpoint", "", "MCP Streamable HTTP endpoint")
			executable := fs.String("executable", "", "stdio MCP executable")
			argsJSON := fs.String("args-json", "", "stdio 参数 JSON 数组")
			headerRefsJSON := fs.String("header-refs-json", "", "HTTP header 环境引用 JSON 对象")
			envRefsJSON := fs.String("env-refs-json", "", "stdio 环境引用 JSON 对象")
			healthEnabled := fs.Bool("health-check-enabled", false, "启用定期健康检查")
			checkInterval := fs.Int64("check-interval-seconds", 0, "健康检查间隔")
			probeTimeout := fs.Int64("probe-timeout-seconds", 0, "健康探测超时")
			autoReconnect := fs.Bool("auto-reconnect", false, "启用固定间隔自动重连")
			reconnectInterval := fs.Int64("reconnect-interval-seconds", 0, "自动重连间隔")
			defaultInclude := fs.Bool("default", false, "新建 Environment 时默认启用")
			if err := fs.Parse(args[1:]); err != nil {
				return flagError(err)
			}
			if fs.NArg() != 0 || strings.TrimSpace(*name) == "" {
				return fmt.Errorf("必须提供 --name；运行 ai-dev-manager-v2 mcp add -h 查看帮助")
			}
			var args []string
			var headerRefs, envRefs map[string]string
			if err := decodeOptionalJSON(*argsJSON, &args, "--args-json"); err != nil {
				return err
			}
			if err := decodeOptionalJSON(*headerRefsJSON, &headerRefs, "--header-refs-json"); err != nil {
				return err
			}
			if err := decodeOptionalJSON(*envRefsJSON, &envRefs, "--env-refs-json"); err != nil {
				return err
			}
			item, err := mcpService.AddMCPConfig(*name, catalog.MCPConfig{
				Transport:  *transport,
				AuthMode:   *authMode,
				Endpoint:   *endpoint,
				HeaderRefs: headerRefs,
				Executable: *executable,
				Args:       args,
				EnvRefs:    envRefs,
				HealthPolicy: model.MCPHealthPolicy{
					HealthCheckEnabled:       *healthEnabled,
					CheckIntervalSeconds:     *checkInterval,
					ProbeTimeoutSeconds:      *probeTimeout,
					AutoReconnect:            *autoReconnect,
					ReconnectIntervalSeconds: *reconnectInterval,
				},
				DefaultInclude: *defaultInclude,
			})
			if err != nil {
				return err
			}
			return writeJSON(item)
		}

		fs := newFlagSet("skill add", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 skill add --root PATH [--support-root PATH] [--default]")
			fmt.Fprintln(os.Stdout, "\n从一个显式全局 Skill root 发现真实 SKILL.md；可选 support root 只授权该 Skill 所需的共享支持文件。")
		})
		root := fs.String("root", "", "包含 Skill 目录/SKILL.md 的显式 discovery root")
		supportRoot := fs.String("support-root", "", "可选的 Skill supporting-files root")
		defaultInclude := fs.Bool("default", false, "发现的 Skill 在新建 Environment 中默认启用")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*root) == "" {
			return fmt.Errorf("缺少 --root；运行 ai-dev-manager-v2 skill add -h 查看帮助")
		}
		supportRoots := []string{}
		if value := strings.TrimSpace(*supportRoot); value != "" {
			supportRoots = append(supportRoots, value)
		}
		items, err := skillService.AddSkillRoot(*root, supportRoots, *defaultInclude)
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "import-preview":
		if kind != "mcp" {
			return fmt.Errorf("未知 %s 命令 %q；运行 ai-dev-manager-v2 %s -h 查看帮助", kind, args[0], kind)
		}
		fs := newFlagSet("mcp import-preview", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 mcp import-preview --json-or-jsonc CONTENT [--format auto|opencode|workbuddy|codex-plugin|claude-code|mcphub] [--source-scope SCOPE] [--default]")
			fmt.Fprintln(os.Stdout, "\n解析并脱敏预览外部 MCP JSON/JSONC；不写入 catalog，也不修改 Environment 选择。")
		})
		format := fs.String("format", app.MCPImportAuto, "导入格式；默认 auto")
		content := fs.String("json-or-jsonc", "", "JSON/JSONC 内容")
		sourceScope := fs.String("source-scope", "", "Claude Code project scope 等显式来源 scope")
		defaultInclude := fs.Bool("default", false, "导入后供新建 Environment 默认选择")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*content) == "" {
			return fmt.Errorf("必须提供 --json-or-jsonc；运行 ai-dev-manager-v2 mcp import-preview -h 查看帮助")
		}
		if application == nil {
			return fmt.Errorf("MCP import service is not initialized")
		}
		preview, err := application.PreviewMCPImport(app.MCPImportInput{Format: *format, Content: *content, SourceScope: *sourceScope, DefaultInclude: *defaultInclude})
		if err != nil {
			return err
		}
		return writeJSON(preview)
	case "import-apply":
		if kind != "mcp" {
			return fmt.Errorf("未知 %s 命令 %q；运行 ai-dev-manager-v2 %s -h 查看帮助", kind, args[0], kind)
		}
		fs := newFlagSet("mcp import-apply", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 mcp import-apply --json-or-jsonc CONTENT [--format FORMAT] [--selected-names A,B] [--conflict-policy error|skip|update_by_name] [--source-scope SCOPE] [--default]")
			fmt.Fprintln(os.Stdout, "\n重新解析并原子写入选中的全局 MCP 定义；不会启用任何已有 Environment。")
		})
		format := fs.String("format", app.MCPImportAuto, "导入格式；默认 auto")
		content := fs.String("json-or-jsonc", "", "JSON/JSONC 内容")
		selectedText := fs.String("selected-names", "", "逗号分隔的 MCP 名称；空值表示全部候选")
		conflictPolicy := fs.String("conflict-policy", catalog.MCPConflictError, "error、skip 或 update_by_name")
		sourceScope := fs.String("source-scope", "", "Claude Code project scope 等显式来源 scope")
		defaultInclude := fs.Bool("default", false, "导入后供新建 Environment 默认选择")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*content) == "" {
			return fmt.Errorf("必须提供 --json-or-jsonc；运行 ai-dev-manager-v2 mcp import-apply -h 查看帮助")
		}
		if application == nil {
			return fmt.Errorf("MCP import service is not initialized")
		}
		selectedNames := []string{}
		for _, value := range strings.Split(*selectedText, ",") {
			if value = strings.TrimSpace(value); value != "" {
				selectedNames = append(selectedNames, value)
			}
		}
		result, err := application.ApplyMCPImport(app.MCPImportInput{Format: *format, Content: *content, SelectedNames: selectedNames, ConflictPolicy: *conflictPolicy, SourceScope: *sourceScope, DefaultInclude: *defaultInclude})
		if err != nil {
			return err
		}
		return writeJSON(result)
	case "source-list":
		if kind != "skill" {
			return fmt.Errorf("unknown %s command %q; run ai-dev-manager-v2 %s -h for help", kind, args[0], kind)
		}
		if len(args) != 1 {
			return fmt.Errorf("skill source-list does not accept arguments")
		}
		sources, err := skillService.ListSkillSources()
		if err != nil {
			return err
		}
		return writeJSON(sources)
	case "source-add":
		if kind != "skill" {
			return fmt.Errorf("unknown %s command %q; run ai-dev-manager-v2 %s -h for help", kind, args[0], kind)
		}
		fs := newFlagSet("skill source-add", func() {
			fmt.Fprintln(os.Stdout, "Usage: ai-dev-manager-v2 skill source-add --root PATH [--support-root PATH] [--default]")
			fmt.Fprintln(os.Stdout, "\\nRegister one explicit Skill source without refreshing it.")
		})
		root := fs.String("root", "", "Skill source discovery root")
		supportRoot := fs.String("support-root", "", "Optional Skill support root")
		defaultInclude := fs.Bool("default", false, "Default-enable Skills refreshed from this source for newly created Environments")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*root) == "" {
			return fmt.Errorf("must provide --root; run ai-dev-manager-v2 skill source-add -h for help")
		}
		supportRoots := []string{}
		if value := strings.TrimSpace(*supportRoot); value != "" {
			supportRoots = append(supportRoots, value)
		}
		source, err := skillService.AddSkillSource(*root, supportRoots, *defaultInclude)
		if err != nil {
			return err
		}
		return writeJSON(source)
	case "source-refresh":
		if kind != "skill" {
			return fmt.Errorf("unknown %s command %q; run ai-dev-manager-v2 %s -h for help", kind, args[0], kind)
		}
		fs := newFlagSet("skill source-refresh", func() {
			fmt.Fprintln(os.Stdout, "Usage: ai-dev-manager-v2 skill source-refresh --id SOURCE_ID")
			fmt.Fprintln(os.Stdout, "\\nAtomically refresh one Skill source snapshot.")
		})
		id := fs.String("id", "", "Skill source ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*id) == "" {
			return fmt.Errorf("must provide --id; run ai-dev-manager-v2 skill source-refresh -h for help")
		}
		result, err := skillService.RefreshSkillSource(*id)
		if err != nil {
			return err
		}
		return writeJSON(result)
	case "source-remove":
		if kind != "skill" {
			return fmt.Errorf("unknown %s command %q; run ai-dev-manager-v2 %s -h for help", kind, args[0], kind)
		}
		fs := newFlagSet("skill source-remove", func() {
			fmt.Fprintln(os.Stdout, "Usage: ai-dev-manager-v2 skill source-remove --id SOURCE_ID")
			fmt.Fprintln(os.Stdout, "\\nRemove one Skill source and its source-owned Skills; existing Environment selections become unresolved.")
		})
		id := fs.String("id", "", "Skill source ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*id) == "" {
			return fmt.Errorf("must provide --id; run ai-dev-manager-v2 skill source-remove -h for help")
		}
		result, err := skillService.RemoveSkillSource(*id)
		if err != nil {
			return err
		}
		return writeJSON(result)
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("%s list 不接受参数", kind)
		}
		if kind == "mcp" {
			items, err := mcpService.List()
			if err != nil {
				return err
			}
			return writeJSON(items)
		}
		items, err := skillService.List()
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "status":
		if kind != "mcp" {
			return fmt.Errorf("未知 %s 命令 %q；运行 ai-dev-manager-v2 %s -h 查看帮助", kind, args[0], kind)
		}
		fs := newFlagSet("mcp status", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 mcp status --id MCP_ID --environment-id ENV_ID")
			fmt.Fprintln(os.Stdout, "\n按 Environment 选择策略对一个 MCP 做即时健康检查，并输出 configured / disabled / healthy / error JSON 状态。")
		})
		id := fs.String("id", "", "MCP ID")
		environmentID := fs.String("environment-id", "", "Environment ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		mcpID := strings.TrimSpace(*id)
		envID := strings.TrimSpace(*environmentID)
		if fs.NArg() != 0 || mcpID == "" || envID == "" {
			return fmt.Errorf("必须提供 --id 和 --environment-id；运行 ai-dev-manager-v2 mcp status -h 查看帮助")
		}
		if application == nil {
			return fmt.Errorf("MCP health service is not initialized")
		}
		status, err := application.ProbeMCPHealth(context.Background(), envID, mcpID)
		if err != nil {
			return err
		}
		return writeJSON(status)
	case "remove":
		fs := newFlagSet(kind+" remove", func() {
			fmt.Fprintf(os.Stdout, "用法：ai-dev-manager-v2 %s remove --id ID\n", kind)
		})
		id := fs.String("id", "", label+" ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		value := strings.TrimSpace(*id)
		if fs.NArg() != 0 || value == "" {
			return fmt.Errorf("缺少 --id；运行 ai-dev-manager-v2 %s remove -h 查看帮助", kind)
		}
		if kind == "mcp" {
			if err := mcpService.Remove(value); err != nil {
				return err
			}
		} else if err := skillService.Remove(value); err != nil {
			return err
		}
		return writeJSON(map[string]any{"removed": value})
	case "set-default":
		fs := newFlagSet(kind+" set-default", func() {
			fmt.Fprintf(os.Stdout, "用法：ai-dev-manager-v2 %s set-default --id ID --enabled true|false\n", kind)
			fmt.Fprintln(os.Stdout, "\n只影响之后新建的 Environment，不重写已有 Environment 的选择。")
		})
		id := fs.String("id", "", label+" ID")
		enabledText := fs.String("enabled", "", "是否默认启用：true 或 false")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		value := strings.TrimSpace(*id)
		if fs.NArg() != 0 || value == "" || strings.TrimSpace(*enabledText) == "" {
			return fmt.Errorf("必须提供 --id 和 --enabled；运行 ai-dev-manager-v2 %s set-default -h 查看帮助", kind)
		}
		enabled, err := strconv.ParseBool(strings.TrimSpace(*enabledText))
		if err != nil {
			return fmt.Errorf("--enabled 必须是 true 或 false")
		}
		if kind == "mcp" {
			item, err := mcpService.SetDefault(value, enabled)
			if err != nil {
				return err
			}
			return writeJSON(item)
		}
		item, err := skillService.SetDefault(value, enabled)
		if err != nil {
			return err
		}
		return writeJSON(item)
	default:
		return fmt.Errorf("未知 %s 命令 %q；运行 ai-dev-manager-v2 %s -h 查看帮助", kind, args[0], kind)
	}
}

func runMemory(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printMemoryHelp()
		return nil
	}
	switch args[0] {
	case "global":
		return runGlobalMemory(service, args[1:])
	case "environment":
		return runEnvironmentMemory(service, args[1:])
	default:
		return fmt.Errorf("未知 memory 命令 %q；运行 ai-dev-manager-v2 memory -h 查看帮助", args[0])
	}
}

func runGlobalMemory(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printGlobalMemoryHelp()
		return nil
	}
	switch args[0] {
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("memory global list 不接受参数")
		}
		items, err := service.Memory.GlobalList()
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "read":
		fs := newFlagSet("memory global read", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 memory global read --key KEY")
		})
		key := fs.String("key", "", "Global Memory key")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		value := strings.TrimSpace(*key)
		if fs.NArg() != 0 || value == "" {
			return fmt.Errorf("缺少 --key；运行 ai-dev-manager-v2 memory global read -h 查看帮助")
		}
		item, err := service.Memory.GlobalRead(value)
		if err != nil {
			return err
		}
		return writeJSON(item)
	case "write":
		fs := newFlagSet("memory global write", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 memory global write --key KEY --value VALUE")
			fmt.Fprintln(os.Stdout, "\n显式写入 Global Memory；VALUE 可以是空字符串，但必须提供 --value。")
		})
		key := fs.String("key", "", "Global Memory key")
		value := fs.String("value", "", "Global Memory value")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		keyValue := strings.TrimSpace(*key)
		if fs.NArg() != 0 || keyValue == "" || !flagWasSet(fs, "value") {
			return fmt.Errorf("必须提供 --key 和 --value；运行 ai-dev-manager-v2 memory global write -h 查看帮助")
		}
		if err := service.Memory.GlobalWrite(keyValue, *value); err != nil {
			return err
		}
		item, err := service.Memory.GlobalRead(keyValue)
		if err != nil {
			return err
		}
		return writeJSON(item)
	case "delete":
		fs := newFlagSet("memory global delete", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 memory global delete --key KEY")
		})
		key := fs.String("key", "", "Global Memory key")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		value := strings.TrimSpace(*key)
		if fs.NArg() != 0 || value == "" {
			return fmt.Errorf("缺少 --key；运行 ai-dev-manager-v2 memory global delete -h 查看帮助")
		}
		if err := service.Memory.GlobalDelete(value); err != nil {
			return err
		}
		return writeJSON(map[string]any{"deleted": value})
	default:
		return fmt.Errorf("未知 memory global 命令 %q；运行 ai-dev-manager-v2 memory global -h 查看帮助", args[0])
	}
}

func runEnvironmentMemory(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printEnvironmentMemoryHelp()
		return nil
	}
	switch args[0] {
	case "list":
		fs := newFlagSet("memory environment list", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 memory environment list --environment-id ENV_ID")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" {
			return fmt.Errorf("缺少 --environment-id；运行 ai-dev-manager-v2 memory environment list -h 查看帮助")
		}
		items, err := service.Memory.EnvironmentList(*environmentID)
		if err != nil {
			return err
		}
		return writeJSON(items)
	case "read":
		fs := newFlagSet("memory environment read", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 memory environment read --environment-id ENV_ID --key KEY")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		key := fs.String("key", "", "Environment-private Memory key")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		keyValue := strings.TrimSpace(*key)
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" || keyValue == "" {
			return fmt.Errorf("必须提供 --environment-id 和 --key；运行 ai-dev-manager-v2 memory environment read -h 查看帮助")
		}
		item, err := service.Memory.EnvironmentRead(*environmentID, keyValue)
		if err != nil {
			return err
		}
		return writeJSON(item)
	case "write":
		fs := newFlagSet("memory environment write", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 memory environment write --environment-id ENV_ID --key KEY --value VALUE")
			fmt.Fprintln(os.Stdout, "\n显式写入指定 Environment 的 private Memory；VALUE 可以是空字符串，但必须提供 --value。")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		key := fs.String("key", "", "Environment-private Memory key")
		value := fs.String("value", "", "Environment-private Memory value")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		keyValue := strings.TrimSpace(*key)
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" || keyValue == "" || !flagWasSet(fs, "value") {
			return fmt.Errorf("必须提供 --environment-id、--key 和 --value；运行 ai-dev-manager-v2 memory environment write -h 查看帮助")
		}
		if err := service.Memory.EnvironmentWrite(*environmentID, keyValue, *value); err != nil {
			return err
		}
		item, err := service.Memory.EnvironmentRead(*environmentID, keyValue)
		if err != nil {
			return err
		}
		return writeJSON(item)
	case "delete":
		fs := newFlagSet("memory environment delete", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 memory environment delete --environment-id ENV_ID --key KEY")
		})
		environmentID := fs.String("environment-id", "", "Environment ID")
		key := fs.String("key", "", "Environment-private Memory key")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		keyValue := strings.TrimSpace(*key)
		if fs.NArg() != 0 || strings.TrimSpace(*environmentID) == "" || keyValue == "" {
			return fmt.Errorf("必须提供 --environment-id 和 --key；运行 ai-dev-manager-v2 memory environment delete -h 查看帮助")
		}
		if err := service.Memory.EnvironmentDelete(*environmentID, keyValue); err != nil {
			return err
		}
		return writeJSON(map[string]any{"environment_id": *environmentID, "deleted": keyValue})
	default:
		return fmt.Errorf("未知 memory environment 命令 %q；运行 ai-dev-manager-v2 memory environment -h 查看帮助", args[0])
	}
}

func runGateway(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printGatewayHelp()
		return nil
	}

	switch args[0] {
	case "start", "http":
		fs := newFlagSet("gateway start", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 gateway start [--listen 127.0.0.1:41137] [-d|--detach]")
			fmt.Fprintln(os.Stdout, "\n在当前终端前台启动 HTTP MCP Gateway；按 Ctrl+C 停止。")
			fmt.Fprintln(os.Stdout, "加 -d 或 --detach 可脱离当前终端运行，健康检查通过后命令返回。")
		})
		listen := fs.String("listen", defaultGatewayListen, "本机回环监听地址")
		var detach bool
		fs.BoolVar(&detach, "detach", false, "脱离当前终端运行，并在健康检查通过后返回")
		fs.BoolVar(&detach, "d", false, "--detach 的简写")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("gateway start 只接受 --flag 参数；运行 ai-dev-manager-v2 gateway start -h 查看帮助")
		}
		if detach {
			return startGatewayDetached(*listen)
		}
		return startHTTPGateway(service, *listen)
	case "status":
		fs := newFlagSet("gateway status", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 gateway status [--listen 127.0.0.1:41137]")
		})
		listen := fs.String("listen", defaultGatewayListen, "Gateway 监听地址")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("gateway status 只接受 --flag 参数")
		}
		return printGatewayStatus(*listen)
	case "stop":
		fs := newFlagSet("gateway stop", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 gateway stop [--listen 127.0.0.1:41137]")
			fmt.Fprintln(os.Stdout, "\n停止 HTTP Gateway；也支持安全识别并停止同一路径启动的旧版 ADM V2 Gateway。")
		})
		listen := fs.String("listen", defaultGatewayListen, "Gateway 监听地址")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("gateway stop 只接受 --flag 参数")
		}
		return stopHTTPGateway(*listen)
	case "restart":
		fs := newFlagSet("gateway restart", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 gateway restart [--listen 127.0.0.1:41137]")
			fmt.Fprintln(os.Stdout, "\n停止当前 HTTP Gateway，然后在这个终端启动新 Gateway。")
		})
		listen := fs.String("listen", defaultGatewayListen, "Gateway 监听地址")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("gateway restart 只接受 --flag 参数")
		}
		if err := stopHTTPGateway(*listen); err != nil {
			return err
		}
		return startHTTPGateway(service, *listen)
	case "stdio":
		if len(args) != 1 {
			return fmt.Errorf("gateway stdio 不接受参数；运行 ai-dev-manager-v2 gateway -h 查看帮助")
		}
		if stdinIsTerminal() {
			return fmt.Errorf("gateway stdio 是给 MCP 客户端使用的协议通道，不是人工终端命令；人工启动 HTTP Gateway 请运行 ai-dev-manager-v2 gateway start")
		}
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		return gateway.RunStdio(ctx, service)
	default:
		return fmt.Errorf("未知 gateway 命令 %q；运行 ai-dev-manager-v2 gateway -h 查看帮助", args[0])
	}
}

func runDoctor(service *app.Service, statePath string, args []string) error {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 doctor")
		fmt.Fprintln(os.Stdout, "\n一次查看当前 ADM 程序、状态文件、Gateway、Workspace、Environment、Writer 和执行白名单。")
		return nil
	}
	if len(args) != 0 {
		return fmt.Errorf("doctor 不接受参数")
	}

	workspaces, err := service.Workspaces.List()
	if err != nil {
		return fmt.Errorf("读取 Workspace 失败: %w", err)
	}
	environments, err := service.Environments.List()
	if err != nil {
		return fmt.Errorf("读取 Environment 失败: %w", err)
	}
	allowed, err := service.AllowedExecutables()
	if err != nil {
		return fmt.Errorf("读取执行白名单失败: %w", err)
	}
	executable, executableErr := os.Executable()
	if executableErr != nil {
		executable = "（无法获取：" + executableErr.Error() + "）"
	}

	fmt.Println("ADM V2 诊断")
	fmt.Println("当前程序：", executable)
	fmt.Println("状态文件：", statePath)
	fmt.Println()

	health, running, healthErr := fetchGatewayHealth(defaultGatewayListen)
	fmt.Println("Gateway")
	switch {
	case healthErr != nil:
		var incompatible *incompatibleGatewayError
		if errors.As(healthErr, &incompatible) {
			fmt.Println("  状态：    版本不兼容")
			fmt.Println("  MCP 地址：http://127.0.0.1:41137/mcp")
			fmt.Println("  详情：   ", incompatible)
			fmt.Println("  处理：    ai-dev-manager-v2 gateway restart")
		} else {
			fmt.Println("  状态：    未知")
			fmt.Println("  MCP 地址：http://127.0.0.1:41137/mcp")
			fmt.Println("  详情：   ", healthErr)
		}
	case running:
		fmt.Println("  状态：    运行中")
		fmt.Println("  MCP 地址：http://127.0.0.1:41137/mcp")
		if health.PID > 0 {
			fmt.Println("  PID：    ", health.PID)
		} else {
			fmt.Println("  PID：     旧版 Gateway 未提供")
		}
		fmt.Println("  版本：   ", health.Version)
		if health.OwnerID != "" {
			fmt.Println("  Runtime Owner：", health.OwnerID)
		}
	default:
		fmt.Println("  状态：    已停止")
		fmt.Println("  MCP 地址：http://127.0.0.1:41137/mcp")
		fmt.Println("  启动：    ai-dev-manager-v2 gateway start")
	}
	fmt.Println()

	fmt.Printf("Workspace（%d）\n", len(workspaces))
	if len(workspaces) == 0 {
		fmt.Println("  暂无；添加：ai-dev-manager-v2 workspace add --path PATH --name NAME")
	} else {
		for _, ws := range workspaces {
			fmt.Printf("  %s  %s  %s\n", ws.ID, ws.Name, ws.Path)
		}
	}
	fmt.Println()

	fmt.Printf("Environment（%d）\n", len(environments))
	if len(environments) == 0 {
		fmt.Println("  暂无；创建：ai-dev-manager-v2 environment create --workspace-id WS_ID --name NAME")
	} else {
		for _, env := range environments {
			writer := "无 Writer"
			if env.Writer != nil {
				writer = "Writer=" + env.Writer.Owner
			}
			fmt.Printf("  %s  %s  root=%s  %s\n", env.ID, env.Name, env.Root, writer)
		}
	}
	fmt.Println()

	fmt.Printf("执行白名单（%d）\n", len(allowed))
	if len(allowed) == 0 {
		fmt.Println("  暂无")
	} else {
		for _, executable := range allowed {
			fmt.Println(" ", executable)
		}
	}
	return nil
}

func runState(statePath string, args []string) error {
	if wantsHelp(args) {
		fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 state path")
		fmt.Fprintln(os.Stdout, "\n打印 ADM V2 持久状态文件路径。")
		return nil
	}
	if len(args) == 1 && args[0] == "path" {
		fmt.Println(statePath)
		return nil
	}
	return fmt.Errorf("未知 state 命令；运行 ai-dev-manager-v2 state -h 查看帮助")
}

func startHTTPGateway(service *app.Service, listen string) error {
	listen = strings.TrimSpace(listen)
	baseURL, err := gatewayBaseURL(listen)
	if err != nil {
		return err
	}
	fmt.Println("ADM V2 HTTP Gateway")
	fmt.Println("状态：    正在启动")
	fmt.Println("MCP 地址：", baseURL+"/mcp")
	fmt.Println("健康检查：", baseURL+"/healthz")
	fmt.Println("停止方式：当前终端按 Ctrl+C，或另开终端运行 ai-dev-manager-v2 gateway stop")
	fmt.Println()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := gateway.RunHTTP(ctx, service, listen); err != nil {
		return fmt.Errorf("在 %s 启动 HTTP Gateway 失败: %w；如果端口可能已被占用，请运行 ai-dev-manager-v2 gateway status", listen, err)
	}
	return nil
}

func startDetachedHTTPGateway(listen string) error {
	listen = strings.TrimSpace(listen)
	baseURL, err := gatewayBaseURL(listen)
	if err != nil {
		return err
	}
	health, running, err := fetchGatewayHealth(listen)
	if err != nil {
		return err
	}
	if running {
		fmt.Println("ADM V2 HTTP Gateway")
		fmt.Println("状态：    已在运行")
		fmt.Println("MCP 地址：", baseURL+"/mcp")
		fmt.Println("PID：    ", health.PID)
		return nil
	}

	process, err := startDetachedGatewayProcess(listen)
	if err != nil {
		return fmt.Errorf("后台启动 Gateway 失败: %w", err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		health, running, healthErr := fetchGatewayHealth(listen)
		if healthErr != nil {
			_ = process.Kill()
			_ = process.Release()
			return fmt.Errorf("后台 Gateway 启动失败: %w", healthErr)
		}
		if running {
			_ = process.Release()
			fmt.Println("ADM V2 HTTP Gateway")
			fmt.Println("状态：    已在后台运行")
			fmt.Println("MCP 地址：", baseURL+"/mcp")
			fmt.Println("PID：    ", health.PID)
			fmt.Println("停止：    ai-dev-manager-v2 gateway stop")
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	_ = process.Kill()
	_ = process.Release()
	return fmt.Errorf("后台 Gateway 未能在 5 秒内通过健康检查: %s/healthz", baseURL)
}
func printGatewayStatus(listen string) error {
	baseURL, err := gatewayBaseURL(listen)
	if err != nil {
		return err
	}
	health, running, err := fetchGatewayHealth(listen)
	if err != nil {
		fmt.Println("ADM V2 HTTP Gateway")
		fmt.Println("MCP 地址：", baseURL+"/mcp")
		var incompatible *incompatibleGatewayError
		if errors.As(err, &incompatible) {
			fmt.Println("状态：    版本不兼容")
			fmt.Println("详情：   ", incompatible)
			fmt.Println("处理：    运行 ai-dev-manager-v2 gateway restart；新版 CLI 会尝试安全停止同一路径的旧版 ADM Gateway")
			return nil
		}
		fmt.Println("状态：    未知")
		return err
	}
	fmt.Println("ADM V2 HTTP Gateway")
	if !running {
		fmt.Println("状态：    已停止")
		fmt.Println("MCP 地址：", baseURL+"/mcp")
		fmt.Println("启动：    ai-dev-manager-v2 gateway start")
		return nil
	}
	fmt.Println("状态：    运行中")
	fmt.Println("MCP 地址：", baseURL+"/mcp")
	fmt.Println("PID：    ", health.PID)
	fmt.Println("版本：   ", health.Version)
	if health.OwnerID != "" {
		fmt.Println("Runtime Owner：", health.OwnerID)
	}
	fmt.Println("停止：    ai-dev-manager-v2 gateway stop")
	return nil
}

func stopHTTPGateway(listen string) error {
	baseURL, err := gatewayBaseURL(listen)
	if err != nil {
		return err
	}
	health, running, err := fetchGatewayHealth(listen)
	if err != nil {
		var incompatible *incompatibleGatewayError
		if errors.As(err, &incompatible) {
			pid, processPath, lookupErr := findListeningProcess(listen)
			if lookupErr != nil {
				return fmt.Errorf("检测到旧版或不兼容的 Gateway，但无法自动定位监听进程: %w", lookupErr)
			}
			currentExecutable, executableErr := os.Executable()
			if executableErr != nil {
				return fmt.Errorf("检测到监听进程 PID %d，但无法确认当前 ADM 可执行文件路径: %w", pid, executableErr)
			}
			if !matchesADMExecutable(processPath, currentExecutable) {
				return fmt.Errorf("%s 被其他程序占用（PID %d，%s）；为了避免误杀，ADM 不会自动停止它", listen, pid, processPath)
			}
			fmt.Printf("检测到旧版 ADM V2 Gateway（PID %d），正在停止以完成升级重启。\n", pid)
			return terminateGatewayProcess(pid, listen, baseURL)
		}
		return err
	}
	if !running {
		fmt.Println("ADM V2 HTTP Gateway 已经停止：", baseURL+"/mcp")
		return nil
	}
	if health.PID <= 0 {
		return fmt.Errorf("Gateway %s 没有提供可用 PID，无法自动停止", baseURL)
	}
	if health.OwnerID != "" {
		if _, err := gateway.StopHTTP(listen); err != nil {
			return err
		}
		fmt.Printf("ADM V2 HTTP Gateway 已停止（PID %d）。\n", health.PID)
		return nil
	}
	return terminateGatewayProcess(health.PID, listen, baseURL)
}

func sameADMExecutable(targetPath, currentPath string) bool {
	if !strings.EqualFold(filepath.Clean(filepath.Dir(targetPath)), filepath.Clean(filepath.Dir(currentPath))) {
		return false
	}
	targetName := strings.ToLower(filepath.Base(targetPath))
	currentName := strings.ToLower(filepath.Base(currentPath))
	if targetName != "ai-dev-manager-v2.exe" {
		return false
	}
	return currentName == "ai-dev-manager-v2.exe" || strings.HasPrefix(currentName, "ai-dev-manager-v2.")
}

func terminateGatewayProcess(pid int, listen, _ string) error {
	if err := gateway.TerminateHTTPProcess(pid, listen); err != nil {
		return err
	}
	fmt.Printf("ADM V2 HTTP Gateway 已停止（PID %d）。\n", pid)
	return nil
}

func fetchGatewayHealth(listen string) (gatewayHealth, bool, error) {
	status, err := gateway.InspectHTTP(listen)
	if err != nil {
		return gatewayHealth{}, false, err
	}
	switch status.State {
	case gateway.HTTPStateStopped:
		return gatewayHealth{}, false, nil
	case gateway.HTTPStateIncompatible:
		return gatewayHealth{}, false, &incompatibleGatewayError{detail: status.Detail}
	case gateway.HTTPStateRunning:
		return gatewayHealth{
			Name:      "ai-dev-manager-v2",
			Version:   status.Version,
			Status:    "ok",
			PID:       status.PID,
			Transport: "http",
			OwnerID:   status.OwnerID,
		}, true, nil
	default:
		return gatewayHealth{}, false, fmt.Errorf("unknown Gateway state %q", status.State)
	}
}

func gatewayBaseURL(listen string) (string, error) {
	return gateway.HTTPBaseURL(listen)
}
func stdinIsTerminal() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func newFlagSet(name string, usage func()) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = usage
	return fs
}

func decodeOptionalJSON(value string, target any, flagName string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(value), target); err != nil {
		return fmt.Errorf("%s 必须是有效 JSON: %w", flagName, err)
	}
	return nil
}

func flagWasSet(fs *flag.FlagSet, name string) bool {
	found := false
	fs.Visit(func(item *flag.Flag) {
		if item.Name == name {
			found = true
		}
	})
	return found
}

func flagError(err error) error {
	if err == flag.ErrHelp {
		return nil
	}
	return err
}

func wantsHelp(args []string) bool {
	return len(args) == 0 || (len(args) == 1 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help"))
}

func writeJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func printUsage() {
	fmt.Fprintln(os.Stdout, `AI Dev Manager V2

给 AI Agent 使用的本地开发网关。
Workspace 和 Environment 都只是配置/状态对象；真正运行中的服务只有 Gateway。

快速开始（HTTP Gateway）：
  ai-dev-manager-v2 workspace add --path D:\projects --name projects
  ai-dev-manager-v2 environment create --workspace-id WS_ID --name main
  ai-dev-manager-v2 gateway start
  ai-dev-manager-v2 gateway status

主要命令：
  workspace      登记、查看、重命名、移除允许 ADM 使用的本地目录
  environment    创建、查看、检查、删除开发上下文（也可以简写为 env）
  exec           管理 Agent 可以执行的程序白名单
  mcp            管理全局 MCP catalog
  skill          管理全局 Skill catalog
  memory         管理显式作用域的持久 Memory
  gateway        启动、查看、停止、重启 MCP Gateway
  doctor         一次查看 ADM 本机整体状态
  state          查看 ADM 状态文件位置

Gateway 常用命令：
  gateway start      启动 HTTP Gateway；加 -d / --detach 脱离终端运行（默认 127.0.0.1:41137）
  gateway status     查看 HTTP Gateway 是否运行、PID 和版本
  gateway stop       停止正在运行的 HTTP Gateway
  gateway restart    停止旧 Gateway，然后在当前终端启动新的 Gateway
  gateway stdio      仅供 MCP 客户端使用；不要在普通终端里手动运行

查看子命令帮助：
  ai-dev-manager-v2 workspace -h
  ai-dev-manager-v2 environment -h
  ai-dev-manager-v2 environment writer -h
  ai-dev-manager-v2 exec -h
  ai-dev-manager-v2 mcp -h
  ai-dev-manager-v2 skill -h
  ai-dev-manager-v2 memory -h
  ai-dev-manager-v2 gateway -h
  ai-dev-manager-v2 doctor`)
}

func printWorkspaceHelp() {
	fmt.Fprintln(os.Stdout, `Workspace = ADM 被允许操作的本地目录。它不是服务，也不要求 Git。

命令：
  ai-dev-manager-v2 workspace add --path PATH [--name NAME]
      登记一个本地目录。

  ai-dev-manager-v2 workspace list
      查看所有 Workspace。

  ai-dev-manager-v2 workspace inspect --workspace-id WS_ID
      按稳定 ID 查看一个 Workspace。

  ai-dev-manager-v2 workspace rename --workspace-id WS_ID --name NAME
      只修改显示名称，不移动或重命名项目目录。

  ai-dev-manager-v2 workspace remove --workspace-id WS_ID
      只移除 ADM 记录，不删除项目目录或文件；仍有 Environment 引用时拒绝移除。`)
}

func printEnvironmentHelp() {
	fmt.Fprintln(os.Stdout, `Environment = 位于 Workspace 内的持久开发上下文。它不是运行中的服务。

命令：
  ai-dev-manager-v2 environment create --workspace-id WS_ID --name NAME [--root PATH]
      创建 Environment；不写 --root 时默认使用整个 Workspace。

  ai-dev-manager-v2 environment list
      查看所有 Environment。

  ai-dev-manager-v2 environment inspect --environment-id ENV_ID
      查看 Workspace 关系、当前能力、已解析/未解析 MCP/Skill 选择和 private Memory 条目数；不展开 Memory 值。

  ai-dev-manager-v2 environment rename --environment-id ENV_ID --name NAME
      只修改显示名称，不移动根目录、不修改选择或 Memory，也不触碰项目文件。

  ai-dev-manager-v2 environment remove --environment-id ENV_ID
      只删除 ADM 中的 Environment 记录，不会删除项目目录或文件。

  ai-dev-manager-v2 environment mcp -h
      管理这个 Environment 启用的全局 MCP ID。

  ai-dev-manager-v2 environment skill -h
      管理这个 Environment 启用的全局 Skill ID。

  ai-dev-manager-v2 environment verifier -h
      声明、查看、删除这个 Environment 的 structured verifier 定义。

  ai-dev-manager-v2 environment writer -h
      查看 Writer 租约相关命令。`)
}

func printEnvironmentSelectionHelp(kind string) {
	label := strings.ToUpper(kind)
	fmt.Fprintf(os.Stdout, `Environment %s selection = 只修改一个 Environment 启用的全局 %s ID，不修改 catalog 默认值或其他 Environment。

命令：
  ai-dev-manager-v2 environment %s enable --environment-id ENV_ID --%s-id ID
      为一个 Environment 启用全局 %s。

  ai-dev-manager-v2 environment %s disable --environment-id ENV_ID --%s-id ID
      为一个 Environment 禁用全局 %s。
`, label, label, kind, kind, label, kind, kind, label)
}

func printEnvironmentVerifierHelp() {
	fmt.Fprintln(os.Stdout, `Environment verifier = 这个 Environment 的结构化 test/lint/build/custom 验证定义。定义本身不会授予执行权限；运行时仍受全局 exec 白名单和 Environment cwd 约束。

命令：
  ai-dev-manager-v2 environment verifier add --environment-id ENV_ID --kind test|lint|build|custom --executable NAME_OR_PATH [--name NAME] [--arg ARG ...] [--cwd RELATIVE_PATH] [--timeout-seconds N] [--enabled=true|false]
      添加一个 Environment-scoped verifier 定义；--arg 可重复。

  ai-dev-manager-v2 environment verifier list --environment-id ENV_ID
      查看这个 Environment 的 verifier 定义。

  ai-dev-manager-v2 environment verifier remove --environment-id ENV_ID --verifier-id VF_ID
      删除一个 verifier 定义。`)
}

func printWriterHelp() {
	fmt.Fprintln(os.Stdout, `Writer = 对同一个物理目录进行修改时使用的单写入租约。

命令：
  ai-dev-manager-v2 environment writer acquire --environment-id ENV_ID --owner OWNER
      获取或续租 Writer。

  ai-dev-manager-v2 environment writer heartbeat --environment-id ENV_ID --owner OWNER
      只续租，不执行文件修改。

  ai-dev-manager-v2 environment writer release --environment-id ENV_ID --owner OWNER
      正常释放自己的 Writer。

  ai-dev-manager-v2 environment writer release --environment-id ENV_ID --force
      强制释放，用于人工恢复。`)
}

func printExecHelp() {
	fmt.Fprintln(os.Stdout, `Exec 白名单决定 Agent 可以运行哪些本地程序。

命令：
  ai-dev-manager-v2 exec allow --executable NAME_OR_PATH
      加入一个允许执行的程序。

  ai-dev-manager-v2 exec remove --executable NAME_OR_PATH
      从白名单移除一个程序；后续 exec 立即按新的白名单判断。

  ai-dev-manager-v2 exec list
      查看当前白名单。`)
}

func printCatalogHelp(kind string) {
	if kind == "skill" {
		fmt.Fprintln(os.Stdout, `Skill catalog = 从显式配置的全局 Skill root 发现真实 SKILL.md；Environment 只保存启用的稳定 Skill ID。

命令：
  ai-dev-manager-v2 skill add --root PATH [--support-root PATH] [--default]
      扫描一个显式 discovery root。support root 只用于授权 Skill 需要读取的共享支持文件。

  ai-dev-manager-v2 skill list
      查看已发现的真实 Skill artifact/source 信息。

  ai-dev-manager-v2 skill remove --id ID
      删除一个 catalog 条目；已有 Environment 中的 ID 引用不会被静默改写。

  ai-dev-manager-v2 skill set-default --id ID --enabled true|false
      修改新建 Environment 的默认选择，不重写已有 Environment。`)
		return
	}
	fmt.Fprintln(os.Stdout, `MCP catalog = 全局定义；Environment 只保存启用的 ID。

命令：
	  ai-dev-manager-v2 mcp add --name NAME --transport streamable-http --endpoint URL [--auth-mode headers --header-refs-json JSON] [--default]
	  ai-dev-manager-v2 mcp add --name NAME --transport stdio --executable PATH [--args-json JSON] [--env-refs-json JSON] [--default]
	      添加 typed Streamable HTTP 或 stdio MCP 定义；secret 值使用环境变量引用。

  ai-dev-manager-v2 mcp list
      查看所有全局条目。

  ai-dev-manager-v2 mcp import-preview --json-or-jsonc CONTENT [--format FORMAT] [--source-scope SCOPE]
      脱敏预览 OpenCode / WorkBuddy / Codex plugin / Claude Code / MCPHub JSON/JSONC，不写入 catalog。

  ai-dev-manager-v2 mcp import-apply --json-or-jsonc CONTENT [--selected-names A,B] [--conflict-policy error|skip|update_by_name]
      原子写入选中的全局 MCP 定义；不会修改已有 Environment 选择。

  ai-dev-manager-v2 mcp status --id MCP_ID --environment-id ENV_ID
      即时检查一个 MCP 在指定 Environment 中的 configured / disabled / healthy / error 状态。

  ai-dev-manager-v2 mcp remove --id ID
      删除一个全局条目；已有 Environment 中的 ID 引用不会被静默改写。

  ai-dev-manager-v2 mcp set-default --id ID --enabled true|false
      修改新建 Environment 的默认选择，不重写已有 Environment。`)
}

func printMemoryHelp() {
	fmt.Fprintln(os.Stdout, `Memory = ADM 持久开发上下文。写入时必须显式选择作用域。

命令：
  ai-dev-manager-v2 memory global -h
      管理跨 Environment 共享的 Global Memory。

  ai-dev-manager-v2 memory environment -h
      按显式 Environment ID 管理 Environment-private Memory。`)
}

func printGlobalMemoryHelp() {
	fmt.Fprintln(os.Stdout, `Global Memory = 跨 Environment 共享的持久上下文。

命令：
  ai-dev-manager-v2 memory global list
      查看所有 Global Memory 条目。

  ai-dev-manager-v2 memory global read --key KEY
      读取一个条目。

  ai-dev-manager-v2 memory global write --key KEY --value VALUE
      显式写入 Global Memory。

  ai-dev-manager-v2 memory global delete --key KEY
      删除一个 Global Memory 条目。`)
}

func printEnvironmentMemoryHelp() {
	fmt.Fprintln(os.Stdout, `Environment-private Memory = 只属于一个显式 Environment 的持久上下文。

命令：
  ai-dev-manager-v2 memory environment list --environment-id ENV_ID
      查看一个 Environment 的 private Memory。

  ai-dev-manager-v2 memory environment read --environment-id ENV_ID --key KEY
      读取一个 Environment-private Memory 条目。

  ai-dev-manager-v2 memory environment write --environment-id ENV_ID --key KEY --value VALUE
      显式写入一个 Environment 的 private Memory。

  ai-dev-manager-v2 memory environment delete --environment-id ENV_ID --key KEY
      删除一个 Environment-private Memory 条目。`)
}

func printGatewayHelp() {
	fmt.Fprintln(os.Stdout, `Gateway = 真正运行中的 MCP 服务进程。

人工使用的 HTTP Gateway：
  ai-dev-manager-v2 gateway start [--listen 127.0.0.1:41137] [-d|--detach]
      在当前终端前台启动。终端会被占用，按 Ctrl+C 停止。
      加 -d 或 --detach 后脱离当前终端运行，健康检查通过后命令立即返回。

  ai-dev-manager-v2 gateway status [--listen 127.0.0.1:41137]
      查看运行状态、MCP 地址、PID 和版本。

  ai-dev-manager-v2 gateway stop [--listen 127.0.0.1:41137]
      从另一个终端停止正在运行的 HTTP Gateway。
      也能识别并停止同一路径启动的旧版 ADM V2 Gateway。

  ai-dev-manager-v2 gateway restart [--listen 127.0.0.1:41137]
      停止旧 Gateway，然后在当前终端启动新 Gateway。

仅供 MCP 客户端使用：
  ai-dev-manager-v2 gateway stdio
      stdin/stdout 是 MCP 协议通道。通常由 MCP 客户端自动启动，人不要手动运行。`)
}

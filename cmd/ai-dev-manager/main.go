package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
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
	"ai-dev-manager-v2/internal/store"
)

const defaultGatewayListen = "127.0.0.1:41137"

type gatewayHealth struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Status    string `json:"status"`
	PID       int    `json:"pid"`
	Transport string `json:"transport"`
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
		return runCatalog("mcp", service.MCPs, args[1:])
	case "skill":
		return runCatalog("skill", service.Skills, args[1:])
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
		items, err := service.Environments.List()
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
		env, err := service.Environments.Get(*environmentID)
		if err != nil {
			return err
		}
		caps, err := service.Capabilities(context.Background(), env.ID)
		if err != nil {
			return err
		}
		return writeJSON(map[string]any{"environment": env, "capabilities": caps})
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
	case "writer":
		return runWriter(service, args[1:])
	default:
		return fmt.Errorf("未知 environment 命令 %q；运行 ai-dev-manager-v2 environment -h 查看帮助", args[0])
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

func runCatalog(kind string, service *catalog.Service, args []string) error {
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
	switch args[0] {
	case "add":
		fs := newFlagSet(kind+" add", func() {
			fmt.Fprintf(os.Stdout, "用法：ai-dev-manager-v2 %s add --name NAME [--default]\n", kind)
			fmt.Fprintf(os.Stdout, "\n添加一个全局 %s catalog 条目；--default 表示新建 Environment 时默认启用。\n", label)
		})
		name := fs.String("name", "", label+" 名称")
		defaultInclude := fs.Bool("default", false, "新建 Environment 时默认启用")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 || strings.TrimSpace(*name) == "" {
			return fmt.Errorf("缺少 --name；运行 ai-dev-manager-v2 %s add -h 查看帮助", kind)
		}
		item, err := service.Add(*name, *defaultInclude)
		if err != nil {
			return err
		}
		return writeJSON(item)
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("%s list 不接受参数", kind)
		}
		items, err := service.List()
		if err != nil {
			return err
		}
		return writeJSON(items)
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
		if err := service.Remove(value); err != nil {
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
		item, err := service.SetDefault(value, enabled)
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

func terminateGatewayProcess(pid int, listen, baseURL string) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("查找 Gateway 进程 %d 失败: %w", pid, err)
	}
	if err := process.Kill(); err != nil {
		return fmt.Errorf("停止 Gateway 进程 %d 失败: %w", pid, err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		connection, dialErr := net.DialTimeout("tcp", listen, 200*time.Millisecond)
		if dialErr != nil {
			fmt.Printf("ADM V2 HTTP Gateway 已停止（PID %d）。\n", pid)
			return nil
		}
		_ = connection.Close()
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("已终止 Gateway 进程 %d，但端点 %s 仍然有响应", pid, baseURL)
}

func fetchGatewayHealth(listen string) (gatewayHealth, bool, error) {
	baseURL, err := gatewayBaseURL(listen)
	if err != nil {
		return gatewayHealth{}, false, err
	}
	client := &http.Client{Timeout: 1200 * time.Millisecond}
	response, err := client.Get(baseURL + "/healthz")
	if err != nil {
		if isConnectionFailure(err) {
			return gatewayHealth{}, false, nil
		}
		return gatewayHealth{}, false, fmt.Errorf("检查 Gateway %s 失败: %w", baseURL, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return gatewayHealth{}, false, &incompatibleGatewayError{detail: fmt.Sprintf("端口 %s 有程序响应，但它不是当前版本可识别的 ADM V2 Gateway（健康检查返回 %s）", listen, response.Status)}
	}
	var health gatewayHealth
	if err := json.NewDecoder(response.Body).Decode(&health); err != nil {
		return gatewayHealth{}, false, &incompatibleGatewayError{detail: fmt.Sprintf("端口 %s 有程序响应，但健康检查不是有效的 ADM V2 JSON: %v", listen, err)}
	}
	if health.Name != "ai-dev-manager-v2" || health.Status != "ok" {
		return gatewayHealth{}, false, &incompatibleGatewayError{detail: fmt.Sprintf("端口 %s 有程序响应，但它不是预期的 ADM V2 Gateway", listen)}
	}
	return health, true, nil
}

func gatewayBaseURL(listen string) (string, error) {
	listen = strings.TrimSpace(listen)
	host, port, err := net.SplitHostPort(listen)
	if err != nil {
		return "", fmt.Errorf("invalid Gateway listen address %q; expected host:port", listen)
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return "http://" + net.JoinHostPort(host, port), nil
}

func isConnectionFailure(err error) bool {
	var netErr net.Error
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "connection refused") || strings.Contains(message, "actively refused") || (errors.As(err, &netErr) && netErr.Timeout())
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
      查看一个 Environment 及其当前能力。

  ai-dev-manager-v2 environment remove --environment-id ENV_ID
      只删除 ADM 中的 Environment 记录，不会删除项目目录或文件。

  ai-dev-manager-v2 environment mcp -h
      管理这个 Environment 启用的全局 MCP ID。

  ai-dev-manager-v2 environment skill -h
      管理这个 Environment 启用的全局 Skill ID。

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
	label := "MCP"
	if kind == "skill" {
		label = "Skill"
	}
	fmt.Fprintf(os.Stdout, `%s catalog = 全局定义；Environment 只保存启用的 ID。

命令：
  ai-dev-manager-v2 %s add --name NAME [--default]
      添加全局条目；--default 表示新建 Environment 时默认启用。

  ai-dev-manager-v2 %s list
      查看所有全局条目。

  ai-dev-manager-v2 %s remove --id ID
      删除一个全局条目；已有 Environment 中的 ID 引用不会被静默改写。

  ai-dev-manager-v2 %s set-default --id ID --enabled true|false
      修改新建 Environment 的默认选择，不重写已有 Environment。
`, label, kind, kind, kind, kind)
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

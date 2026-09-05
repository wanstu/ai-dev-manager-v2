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
	"strings"
	"syscall"
	"time"

	"ai-dev-manager-v2/internal/app"
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

var matchesADMExecutable = sameADMExecutable

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
	case "writer":
		return runWriter(service, args[1:])
	default:
		return fmt.Errorf("未知 environment 命令 %q；运行 ai-dev-manager-v2 environment -h 查看帮助", args[0])
	}
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

func runGateway(service *app.Service, args []string) error {
	if wantsHelp(args) {
		printGatewayHelp()
		return nil
	}

	switch args[0] {
	case "start", "http":
		fs := newFlagSet("gateway start", func() {
			fmt.Fprintln(os.Stdout, "用法：ai-dev-manager-v2 gateway start [--listen 127.0.0.1:41137]")
			fmt.Fprintln(os.Stdout, "\n在当前终端前台启动 HTTP MCP Gateway；按 Ctrl+C 停止。")
		})
		listen := fs.String("listen", defaultGatewayListen, "本机回环监听地址")
		if err := fs.Parse(args[1:]); err != nil {
			return flagError(err)
		}
		if fs.NArg() != 0 {
			return fmt.Errorf("gateway start 只接受 --flag 参数；运行 ai-dev-manager-v2 gateway start -h 查看帮助")
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
  workspace      登记、查看允许 ADM 使用的本地目录
  environment    创建、查看、检查、删除开发上下文（也可以简写为 env）
  exec           管理 Agent 可以执行的程序白名单
  gateway        启动、查看、停止、重启 MCP Gateway
  doctor         一次查看 ADM 本机整体状态
  state          查看 ADM 状态文件位置

Gateway 常用命令：
  gateway start      在当前终端启动 HTTP Gateway（默认 127.0.0.1:41137）
  gateway status     查看 HTTP Gateway 是否运行、PID 和版本
  gateway stop       停止正在运行的 HTTP Gateway
  gateway restart    停止旧 Gateway，然后在当前终端启动新的 Gateway
  gateway stdio      仅供 MCP 客户端使用；不要在普通终端里手动运行

查看子命令帮助：
  ai-dev-manager-v2 workspace -h
  ai-dev-manager-v2 environment -h
  ai-dev-manager-v2 environment writer -h
  ai-dev-manager-v2 exec -h
  ai-dev-manager-v2 gateway -h
  ai-dev-manager-v2 doctor`)
}

func printWorkspaceHelp() {
	fmt.Fprintln(os.Stdout, `Workspace = ADM 被允许操作的本地目录。它不是服务，也不要求 Git。

命令：
  ai-dev-manager-v2 workspace add --path PATH [--name NAME]
      登记一个本地目录。

  ai-dev-manager-v2 workspace list
      查看所有 Workspace。`)
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

  ai-dev-manager-v2 environment writer -h
      查看 Writer 租约相关命令。`)
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

  ai-dev-manager-v2 exec list
      查看当前白名单。`)
}

func printGatewayHelp() {
	fmt.Fprintln(os.Stdout, `Gateway = 真正运行中的 MCP 服务进程。

人工使用的 HTTP Gateway：
  ai-dev-manager-v2 gateway start [--listen 127.0.0.1:41137]
      在当前终端前台启动。终端会被占用，按 Ctrl+C 停止。

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

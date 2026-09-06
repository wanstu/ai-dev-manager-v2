# AI Dev Manager V2

ADM V2 是本地 AI 开发 Gateway。它让 Agent 在你明确登记的目录中读写文件、搜索和执行允许的命令。

**Workspace 和 Environment 都不是后台服务。Gateway 才是运行中的 MCP 服务进程。Git 不是前置条件。**

## 第一次使用

在项目目录编译：

```powershell
go build -o ai-dev-manager-v2.exe ./cmd/ai-dev-manager
```

先看帮助。如果不知道当前 ADM 到底是什么状态，直接跑 `doctor`：

```powershell
.\ai-dev-manager-v2.exe -h
.\ai-dev-manager-v2.exe doctor
```

## Desktop Manager

桌面端是独立入口，和 CLI / MCP Gateway 共用同一份 ADM state，不需要先启动 Gateway：

```powershell
go build -o ai-dev-manager-v2-desktop.exe ./cmd/ai-dev-manager-desktop
.\ai-dev-manager-v2-desktop.exe
```

桌面端使用 Wails v2 + 内嵌 HTML/CSS/JavaScript，不需要 npm、Vite 或 Node 构建链。目前已经可以管理 Workspace / Environment 生命周期、查看 Environment detail、维护 exec allowlist、MCP / Skill catalog 和每个 Environment 的选择。

Global Memory 和 Environment-private Memory 也可以在桌面端显式读取、写入和删除，但 Memory 值不会进入普通 Snapshot 或 Environment 总览；只有点击对应的“加载 Memory”后才会读取值。

Gateway 生命周期控制仍在后续桌面 slice 中。

登记 `D:\projects`：

```powershell
.\ai-dev-manager-v2.exe workspace add --path D:\projects --name projects
.\ai-dev-manager-v2.exe workspace list
```

记下返回的 `workspace_id`，然后创建 Environment：

```powershell
.\ai-dev-manager-v2.exe environment create --workspace-id ws_xxx --name main
.\ai-dev-manager-v2.exe environment list
```

同一个 Workspace、同一个 name、同一个 root 再次执行 `environment create` 会复用已有 Environment，不会继续制造重复记录。

Environment 可以只修改显示名称，不移动 root、不修改 MCP/Skill 选择或 Memory，也不触碰项目文件：

```powershell
.\ai-dev-manager-v2.exe environment rename --environment-id env_xxx --name review
```

`environment list` 只返回轻量摘要和 `private_memory_count`，不会展开 private Memory 值。需要看完整管理上下文时使用：

```powershell
.\ai-dev-manager-v2.exe environment inspect --environment-id env_xxx
```

`inspect` 会补充 Workspace 关系、当前 capabilities、已解析的 MCP/Skill catalog 条目以及仍然存在于 Environment 选择中的 unresolved IDs；private Memory 值仍然只能通过显式 `memory environment ...` 命令读取。

删除 Environment 只会删除 ADM 里的上下文记录，不会删除 root 或项目文件：

```powershell
.\ai-dev-manager-v2.exe environment remove --environment-id env_xxx
```

有 active writer 时 remove 会被拒绝。

如果不传 `--root`，Environment root 就是 Workspace 本身。因此一个 `D:\projects` Environment 可以直接操作它下面的多个项目。

## MCP / Skill catalog

全局 MCP 和 Skill 定义现在也可以直接从 CLI 管理；Environment 只保存自己启用的全局 ID。

```powershell
.\ai-dev-manager-v2.exe mcp add --name filesystem --default
.\ai-dev-manager-v2.exe mcp list
.\ai-dev-manager-v2.exe mcp set-default --id mcp_xxx --enabled false
.\ai-dev-manager-v2.exe mcp remove --id mcp_xxx

.\ai-dev-manager-v2.exe skill add --name go-project
.\ai-dev-manager-v2.exe skill list
.\ai-dev-manager-v2.exe skill set-default --id skill_xxx --enabled true
.\ai-dev-manager-v2.exe skill remove --id skill_xxx
```

`set-default` 只影响之后新建的 Environment，不会重写已有 Environment 的 MCP / Skill 选择。

已有 Environment 可以独立启用或禁用 catalog 中的 ID：

```powershell
.\ai-dev-manager-v2.exe environment mcp enable --environment-id ENV_ID --mcp-id mcp_xxx
.\ai-dev-manager-v2.exe environment mcp disable --environment-id ENV_ID --mcp-id mcp_xxx
.\ai-dev-manager-v2.exe environment skill enable --environment-id ENV_ID --skill-id skill_xxx
.\ai-dev-manager-v2.exe environment skill disable --environment-id ENV_ID --skill-id skill_xxx
```

这些命令只修改指定 Environment 的选择，不修改 catalog 默认值，也不影响其他 Environment。

## Global Memory

Global Memory 是跨 Environment 共享的持久上下文。CLI 要求显式写出 `global`，避免以后和 Environment-private Memory 混淆作用域。

```powershell
.\ai-dev-manager-v2.exe memory global list
.\ai-dev-manager-v2.exe memory global read --key machine
.\ai-dev-manager-v2.exe memory global write --key machine --value windows
.\ai-dev-manager-v2.exe memory global delete --key machine
```

Environment-private Memory 必须显式指定 Environment ID：

```powershell
.\ai-dev-manager-v2.exe memory environment list --environment-id ENV_ID
.\ai-dev-manager-v2.exe memory environment read --environment-id ENV_ID --key task
.\ai-dev-manager-v2.exe memory environment write --environment-id ENV_ID --key task --value "private context"
.\ai-dev-manager-v2.exe memory environment delete --environment-id ENV_ID --key task
```

Environment-private Memory 不会自动写入 Global Memory，也不会通过另一个 Environment ID 读取。

## Gateway：真正需要启动的服务

人工使用 HTTP Gateway：

```powershell
.\ai-dev-manager-v2.exe gateway start
```

默认监听：

```text
http://127.0.0.1:41137/mcp
```

`gateway start` 是**前台常驻进程**。启动它的终端会一直被占用，按 `Ctrl+C` 停止。

从另一个终端查看、停止或重启：

```powershell
.\ai-dev-manager-v2.exe gateway status
.\ai-dev-manager-v2.exe gateway stop
.\ai-dev-manager-v2.exe gateway restart
```

如果 41137 已经被旧版本或其他进程占用，`gateway status` / `doctor` 会明确显示 `INCOMPATIBLE`，而不是让你猜发生了什么。

### `gateway stdio` 是什么？

```powershell
.\ai-dev-manager-v2.exe gateway stdio
```

这是给 MCP 客户端通过 stdin/stdout 拉起的 transport，**不是给人手动在终端里运行的服务模式**。交互终端误跑时 CLI 会直接解释并拒绝启动。

## Writer

Agent 修改文件或执行命令前需要持有对应 physical root 的 Writer lease：

```powershell
.\ai-dev-manager-v2.exe environment writer acquire --environment-id env_xxx --owner chatgpt
```

显式续租：

```powershell
.\ai-dev-manager-v2.exe environment writer heartbeat --environment-id env_xxx --owner chatgpt
```

释放：

```powershell
.\ai-dev-manager-v2.exe environment writer release --environment-id env_xxx --owner chatgpt
```

恢复场景可强制释放：

```powershell
.\ai-dev-manager-v2.exe environment writer release --environment-id env_xxx --force
```

Writer 默认 TTL 为 5 分钟。成功的写入、编辑、删除和命令执行会续租；长时间 `exec` 会 heartbeat；异常退出后 lease 到期自动失效。

## Exec allowlist

Agent 只能执行显式允许的 executable：

```powershell
.\ai-dev-manager-v2.exe exec allow --executable go
.\ai-dev-manager-v2.exe exec allow --executable git
.\ai-dev-manager-v2.exe exec remove --executable git
.\ai-dev-manager-v2.exe exec list
```

## 一个 Environment 能不能管整个 `D:\projects`？

可以，而且这是合法的正常用法：

```text
D:\projects
├── ai-dev-manager-v2
├── project-a
├── project-b
└── mcphub

Workspace root   = D:\projects
Environment root = D:\projects
```

文件操作和 `exec --cwd` 都可以进入 Environment root 下的子目录。不需要每个项目创建一个 Environment。

多个 Environment 主要用于需要不同 Writer、Memory、MCP、Skill 或隔离上下文时。

## Git

Workspace / Environment 不要求 Git。这些能力可以在普通目录工作：

```text
tree
read
search
write
edit
delete
exec
```

当前专用 `git_status` / `git_diff` / `git_branch` 仍然针对 Environment root 本身。如果 Environment root 是 `D:\projects` 而 Git 仓库在 `D:\projects\ai-dev-manager-v2`，专用 Git 工具当前不会自动切到该子目录；Agent 仍可通过 `exec` 的 `cwd` 在子仓库运行 Git。

这是 Git capability 的已知限制，不应该通过拆分 Workspace/Environment 来规避。

## 状态与诊断

一条命令看整体状态：

```powershell
.\ai-dev-manager-v2.exe doctor
```

它会显示当前 executable、state 文件、Gateway、Workspace、Environment、Writer 和 exec allowlist。

查看 state 文件路径：

```powershell
.\ai-dev-manager-v2.exe state path
```

详细产品语义见 `docs/PRODUCT_CONTRACT.md`。

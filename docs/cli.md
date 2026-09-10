# CLI 命令

本地 CLI 发布名为 `adm`。Windows 构建产物默认叫：

```text
dist\adm-windows-amd64.exe
```

示例中用变量简化命令：

```powershell
$adm = '.\dist\adm-windows-amd64.exe'
```

## Workspace

```powershell
& $adm workspace add --path D:\projects --name projects
& $adm workspace list
& $adm workspace inspect --workspace-id WS_ID
& $adm workspace rename --workspace-id WS_ID --name projects-main
& $adm workspace remove --workspace-id WS_ID
```

Workspace 记录的是 ADM 可以访问的本地目录。移除 Workspace 只删除 ADM 记录，不删除项目文件；仍有 Environment 引用时会拒绝移除。

## Environment

```powershell
& $adm environment create --workspace-id WS_ID --name main
& $adm environment list
& $adm environment inspect --environment-id ENV_ID
& $adm environment rename --environment-id ENV_ID --name review
& $adm environment remove --environment-id ENV_ID
```

Environment 可以简写成 `env`。移除 Environment 只删除 ADM 上下文记录，不删除 root 或项目文件。有 active writer 时 remove 会被拒绝。

## Writer

Agent 修改文件或执行命令前需要持有对应 physical root 的 Writer lease：

```powershell
& $adm environment writer acquire --environment-id ENV_ID --owner chatgpt
& $adm environment writer heartbeat --environment-id ENV_ID --owner chatgpt
& $adm environment writer release --environment-id ENV_ID --owner chatgpt
```

恢复场景可以强制释放：

```powershell
& $adm environment writer release --environment-id ENV_ID --force
```

Writer 默认 TTL 为 5 分钟。成功的写入、编辑、删除和命令执行会续租；长时间 `exec` 会 heartbeat；异常退出后 lease 到期自动失效。

## Exec allowlist

Agent 只能执行显式允许的 executable：

```powershell
& $adm exec allow --executable go
& $adm exec allow --executable git
& $adm exec list
& $adm exec remove --executable git
```

## Gateway

```powershell
& $adm gateway start
& $adm gateway start --detach
& $adm gateway status
& $adm gateway stop
& $adm gateway restart
```

`gateway start` 默认是前台常驻进程。`--detach` 会在后台启动本机 HTTP Gateway。

`gateway stdio` 仅供 MCP 客户端通过 stdin/stdout 拉起，不是给人手动在普通终端里运行的服务模式：

```powershell
& $adm gateway stdio
```

## 状态文件

```powershell
& $adm state path
& $adm doctor
```

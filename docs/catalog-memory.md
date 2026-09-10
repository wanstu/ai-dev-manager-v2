# MCP、Skill 与 Memory

## MCP catalog

全局 MCP 定义由 CLI 或 Desktop 管理。Environment 只保存自己启用的全局 MCP ID。

```powershell
$adm = '.\dist\adm-windows-amd64.exe'

& $adm mcp add --name filesystem --transport streamable-http --endpoint http://127.0.0.1:9000/mcp --default
& $adm mcp add --name local-helper --transport stdio --executable mcp-helper --args-json '["serve","--stdio"]' --env-refs-json '{"API_TOKEN":"${MCP_API_TOKEN}"}'
& $adm mcp list
& $adm mcp set-default --id mcp_xxx --enabled false
& $adm mcp remove --id mcp_xxx
```

HTTP MCP 支持 `none` 或 `headers` 认证模式。header 与 stdio 环境变量只保存 `${ENV_NAME}` 引用，并在 Environment 激活时解析。stdio executable 必须先加入 ADM exec allowlist。

Gateway 会在进程内保存 MCP observation：健康状态、失败阶段、连续失败次数、下一次固定间隔重连时间，以及最多 256 个已发现工具的公开清单。这些运行时事实不会写入 ADM state。

## Skill catalog

Skill 不再是手填一段 instructions。ADM 只扫描你显式配置的 discovery root，发现真实 `SKILL.md` 并记录 artifact/source。

```powershell
& $adm skill add --root C:\Users\you\.config\opencode\skills --support-root C:\Users\you\.config\opencode\gsd-core
& $adm skill list
& $adm skill set-default --id skill_xxx --enabled true
& $adm skill remove --id skill_xxx
```

如果 Skill 引用 discovery root 之外的共享支持文件，需要显式配置 `--support-root`。

## Environment 选择

`set-default` 只影响之后新建的 Environment，不会重写已有 Environment 的 MCP / Skill 选择。已有 Environment 可以独立启用或禁用 catalog 中的 ID：

```powershell
& $adm environment mcp enable --environment-id ENV_ID --mcp-id mcp_xxx
& $adm environment mcp disable --environment-id ENV_ID --mcp-id mcp_xxx
& $adm environment skill enable --environment-id ENV_ID --skill-id skill_xxx
& $adm environment skill disable --environment-id ENV_ID --skill-id skill_xxx
```

这些命令只修改指定 Environment 的选择，不修改 catalog 默认值，也不影响其他 Environment。

## Global Memory

Global Memory 是跨 Environment 共享的持久上下文。CLI 要求显式写出 `global`，避免和 Environment-private Memory 混淆作用域。

```powershell
& $adm memory global list
& $adm memory global read --key machine
& $adm memory global write --key machine --value windows
& $adm memory global delete --key machine
```

## Environment-private Memory

Environment-private Memory 必须显式指定 Environment ID：

```powershell
& $adm memory environment list --environment-id ENV_ID
& $adm memory environment read --environment-id ENV_ID --key task
& $adm memory environment write --environment-id ENV_ID --key task --value "private context"
& $adm memory environment delete --environment-id ENV_ID --key task
```

Environment-private Memory 不会自动写入 Global Memory，也不会通过另一个 Environment ID 读取。

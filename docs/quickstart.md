# 快速开始

## 编译 CLI

```powershell
New-Item -ItemType Directory -Force -Path dist | Out-Null
go build -trimpath -o dist\adm-windows-amd64.exe ./cmd/ai-dev-manager
```

后续示例都使用 Windows 产物名：

```powershell
$adm = '.\dist\adm-windows-amd64.exe'
```

## 查看状态

```powershell
& $adm -h
& $adm doctor
```

`doctor` 会显示当前 executable、ADM state 文件、Gateway、Workspace、Environment、Writer 和 exec allowlist 状态。

## 启动 Gateway

Gateway 是真正运行的 MCP/HTTP 服务。Workspace 和 Environment 都只是持久记录，不是后台服务。

```powershell
& $adm gateway start --detach
& $adm gateway status
```

默认 ADM Base URL 是：

```text
http://127.0.0.1:43137
```

也可以显式选择管理目标：

```powershell
& $adm --adm-url http://127.0.0.1:8001 gateway status
$env:ADM_V2_URL = 'http://127.0.0.1:8001'
& $adm gateway status
```

`--adm-url` / `ADM_V2_URL` 只选择管理目标；连接失败不会回退到本地 `state.json`。

## 登记 Workspace

```powershell
& $adm workspace add --path D:\projects --name projects
& $adm workspace list
```

Workspace 是 ADM 被允许访问的本地目录。Git 不是前置条件。

## 创建 Environment

```powershell
& $adm environment create --workspace-id WS_ID --name main
& $adm environment list
```

不传 `--root` 时，Environment root 就是 Workspace root。一个 `D:\projects` Environment 可以操作其下多个项目。

```text
D:\projects
├── ai-dev-manager-v2
├── project-a
└── project-b

Workspace root   = D:\projects
Environment root = D:\projects
```

# adm

adm 是本地 AI 开发 Gateway。它让 Agent 在你明确登记的目录中读写文件、搜索文件，并只执行允许的命令。

**Workspace 和 Environment 不是后台服务；真正运行的是 Gateway。Git 不是前置条件。**

## 目录

- [快速开始](#快速开始)
- [Desktop](#desktop)
- [发布产物](#发布产物)
- [核心概念](#核心概念)
- [文档](#文档)

## 快速开始

在项目目录编译 CLI：

```powershell
New-Item -ItemType Directory -Force -Path dist | Out-Null
go build -trimpath -o dist\adm-windows-amd64.exe ./cmd/ai-dev-manager
```

查看帮助和本机状态：

```powershell
.\dist\adm-windows-amd64.exe -h
.\dist\adm-windows-amd64.exe doctor
```

启动本机 HTTP Gateway：

```powershell
.\dist\adm-windows-amd64.exe gateway start --detach
.\dist\adm-windows-amd64.exe gateway status
```

登记一个目录并创建 Environment：

```powershell
.\dist\adm-windows-amd64.exe workspace add --path D:\projects --name projects
.\dist\adm-windows-amd64.exe environment create --workspace-id WS_ID --name main
```

## Desktop

Desktop 管理端叫 `adm-desktop`。它通过所选 ADM Base URL 的 `/admin/mcp` 管理 ADM，不会在连接失败时直接回退读写本地 `state.json`。

Desktop 必须用 Wails 构建：

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build-desktop.ps1 -clean -trimpath
.\dist\adm-desktop-windows-amd64.exe
```

不要把 `go build ./cmd/ai-dev-manager-desktop` 当作可运行 Desktop 产物；普通 `go build` 缺少 Wails 运行所需的 build tags。

Desktop 支持保存多个 ADM 连接、编辑和切换连接。连接配置保存在：

```text
~/.config/adm/desktop-connections.json
```

## 发布产物

本地 RC 打包统一写入 repo 根目录的 `dist/`：

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build-rc.ps1 -Version v1.0.0-rc.4
```

产物命名：

```text
dist\adm-v1.0.0-rc.4-windows-amd64.exe
dist\adm-desktop-v1.0.0-rc.4-windows-amd64.exe
dist\SHA256SUMS-v1.0.0-rc.4.txt
```

GitHub Actions 在 tag 构建时也会把 tag 写进文件名和 artifact name，例如 `adm-v1.0.0-rc.4-windows-amd64` 与 `adm-desktop-v1.0.0-rc.4-windows-amd64`。

## 核心概念

- **Workspace**：ADM 被允许访问的本地目录。
- **Environment**：位于 Workspace 内的持久开发上下文；保存 Writer、Memory、MCP、Skill 等选择。
- **Gateway**：真正运行的本地 MCP/HTTP 服务。
- **Writer lease**：Agent 写文件或执行命令前必须持有的单写者租约。
- **Exec allowlist**：Agent 只能运行显式允许的 executable。
- **MCP / Skill catalog**：全局定义；Environment 只保存启用选择。
- **Memory**：分为 Global Memory 和 Environment-private Memory，读取值必须显式请求。

## 文档

完整文档入口见 [docs/index.md](docs/index.md)。常用主题：

- [快速开始](docs/quickstart.md)
- [Desktop 管理端](docs/desktop.md)
- [CLI 命令](docs/cli.md)
- [MCP、Skill 与 Memory](docs/catalog-memory.md)
- [打包与 GitHub Actions](docs/packaging.md)
- [产品语义合同](docs/PRODUCT_CONTRACT.md)

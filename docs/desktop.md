# adm-desktop

`adm-desktop` 是 ADM 的 Wails 桌面管理端。它通过所选 ADM Base URL 的 `/admin/mcp` 管理 ADM，不会在连接失败时直接读写本地 `state.json`。

## 构建

Desktop 必须通过 Wails 构建：

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build-desktop.ps1 -clean -trimpath
```

默认产物在统一目录：

```text
dist\adm-desktop-windows-amd64.exe
```

不要把普通 `go build ./cmd/ai-dev-manager-desktop` 当作可运行 Desktop artifact；它缺少 Wails 运行所需 build tags。

## 连接配置

Desktop 可以保存多个 ADM 连接，支持添加、编辑、删除和切换。配置文件位于：

```text
~/.config/adm/desktop-connections.json
```

切换连接时，Desktop 会清空旧连接的 Workspace、Environment、MCP health、Runtime 等界面状态，并等待旧请求完成，避免两边数据混在一起。

## 本地 Gateway 控制

Desktop 会先检查所选 Base URL 的 `/healthz`，再通过 `/admin/mcp` 管理。只有 loopback HTTP root URL 可以使用 Desktop 的本地启动/停止按钮。

常见 Base URL：

```text
http://127.0.0.1:43137
http://127.0.0.1:8001
```

## 托盘和开机启动

Windows 上 Desktop 支持系统托盘：

- 关闭窗口时隐藏到托盘，不直接退出。
- 托盘菜单可以显示主窗口、隐藏主窗口、切换开机启动、退出。
- `--autostart` 启动时默认隐藏窗口。
- 开机启动使用 HKCU Run，value name 为 `adm-desktop`，command 为当前 Desktop exe 加 `--autostart`。

托盘图标来自 `assets/icons/ai-dev-manager-tray.png`，构建时会裁掉透明留白并生成 `cmd/ai-dev-manager-desktop/assets/tray.png`，避免 Windows 托盘里显示过小。

## 图标

- `assets/icons/ai-dev-manager-app.png`：Windows exe / Wails app icon。
- `assets/icons/ai-dev-manager-tray.png`：托盘源图。
- `assets/icons/ai-dev-manager-window.png`：Desktop 顶部品牌图。

`cmd/ai-dev-manager-desktop/wails.json` 的 `frontend:build` 会执行 `scripts/prepare-desktop-icons.ps1`，所以直接执行 `wails build` 和仓库脚本都会刷新图标。

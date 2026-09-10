# 打包与 GitHub Actions

## 命名

用户可见应用名统一为：

- CLI：`adm`
- Desktop：`adm-desktop`

Go module path 暂时仍是 `ai-dev-manager-v2`，这是源码内部 import 路径，不作为用户可见应用名。

## 统一产物目录

本地打包产物统一写入 repo 根目录的 `dist/`。

构建 Desktop：

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build-desktop.ps1 -clean -trimpath
```

默认输出：

```text
dist\adm-desktop-windows-amd64.exe
```

构建 RC：

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build-rc.ps1 -Version v1.0.0-rc.4
```

默认输出：

```text
dist\adm-v1.0.0-rc.4-windows-amd64.exe
dist\adm-desktop-v1.0.0-rc.4-windows-amd64.exe
dist\SHA256SUMS-v1.0.0-rc.4.txt
```

也可以指定输出目录：

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\build-rc.ps1 -Version v1.0.0-rc.4 -OutputDir .\dist\rc
```

## GitHub Actions artifact 规则

CI 的普通 branch / PR 构建使用无 tag 名称：

```text
adm-windows-amd64.exe
adm-desktop-windows-amd64.exe
adm-linux-amd64
adm-darwin-amd64
```

tag 构建会把 tag 写入文件名和 artifact name：

```text
adm-v1.0.0-rc.4-windows-amd64.exe
adm-desktop-v1.0.0-rc.4-windows-amd64.exe
adm-v1.0.0-rc.4-linux-amd64
adm-v1.0.0-rc.4-darwin-amd64
```

当前 CI 负责测试、vet、构建并上传 artifact；正式 GitHub Release 发布仍需单独流程或人工发布。

## Desktop icon 构建

`cmd/ai-dev-manager-desktop/wails.json` 的 `frontend:build` 会运行：

```powershell
scripts\prepare-desktop-icons.ps1
```

该脚本会：

- 复制 `assets/icons/ai-dev-manager-app.png` 到 Wails `build/appicon.png`。
- 复制 `assets/icons/ai-dev-manager-window.png` 到 Desktop 前端资源目录。
- 裁掉 `assets/icons/ai-dev-manager-tray.png` 的透明留白，生成 `cmd/ai-dev-manager-desktop/assets/tray.png`。
- 删除旧的 Wails generated `build/windows/icon.ico`，确保新 app icon 生效。

# adm 文档索引

这里是 adm / adm-desktop 的文档入口。README 只保留快速说明，详细操作放在本目录。

## 入门

- [快速开始](quickstart.md)：编译 CLI、启动 Gateway、登记 Workspace、创建 Environment。
- [Desktop 管理端](desktop.md)：Wails 构建、多连接、托盘、开机启动和 Desktop 注意事项。
- [CLI 命令](cli.md)：Workspace、Environment、Writer、Gateway、exec allowlist 的常用命令。

## 管理能力

- [MCP、Skill 与 Memory](catalog-memory.md)：全局 catalog、Environment 选择、引用变量和显式 Memory 读取。
- [产品语义合同](PRODUCT_CONTRACT.md)：ADM 的完整产品边界、状态语义和安全约束。
- [当前状态](STATUS.md)：阶段性实现状态和已知限制。

## 构建与发布

- [打包与 GitHub Actions](packaging.md)：`dist/` 产物目录、`adm` / `adm-desktop` 命名、RC 打包和 tag artifact 规则。

## 历史阶段文档

这些文档记录早期实现阶段，主要用于追溯：

- [Phase 01：Git-independent core](PHASE_01_GIT_INDEPENDENT_CORE.md)
- [Phase 02：Development context](PHASE_02_DEVELOPMENT_CONTEXT.md)
- [Phase 03：Management CLI](PHASE_03_MANAGEMENT_CLI.md)
- [Phase 04：Environment management UX](PHASE_04_ENVIRONMENT_MANAGEMENT_UX.md)
- [Phase 05：Management boundary](PHASE_05_MANAGEMENT_BOUNDARY.md)
- [Phase 06：Desktop UI](PHASE_06_DESKTOP_UI.md)
- [Phase 07：Desktop release readiness](PHASE_07_DESKTOP_RELEASE_READINESS.md)

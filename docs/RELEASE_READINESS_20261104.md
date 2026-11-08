# Thoughts 发行状态与验收 — 2026-11-04

## 当前公开发行状态：preview.21 未签名预览版

**`0.3.0-preview.21` 已于 2026-11-04 11:00（Asia/Shanghai）公开发布为未签名实验预览版（Pre-release）。它不是已签名稳定版，也不代表 SignPath 已批准申请。**

[打开实际 Release 和下载文件](https://github.com/xuli82221-pixel/helloThoughts/releases/tag/v0.3.0-preview.21)。该 Release 已发布，非草稿，包含七个文件：

| 文件 | 用途 |
| --- | --- |
| `thoughts-0.3.0-preview.21-windows-amd64-setup.exe` | Windows x64 离线、当前用户安装器 |
| `thoughts-0.3.0-preview.21-windows-amd64.zip` | 便携包，含 `thoughts.exe` 和客户端配置示例 |
| `thoughts-claude-desktop-0.3.0-preview.21-windows-amd64.mcpb` | Claude Desktop 扩展，使用相同核心 |
| `SHA256SUMS.txt` | 发布文件的 SHA256 校验值 |
| `BUILDINFO.json` | 公开源码 commit、构建记录与未签名状态 |
| `ACCEPTANCE.json` | 实际 stdio MCP 和隔离安装测试摘要及验收边界 |
| `UNSIGNED_NOTICE.txt` | 未签名与实验预览声明 |

构建来源为公开 commit [`0b155f8a4d0381315f742bb6661e2e0eb6c738c4`](https://github.com/xuli82221-pixel/helloThoughts/commit/0b155f8a4d0381315f742bb6661e2e0eb6c738c4)，版本标签 `v0.3.0-preview.21` 指向该 commit。[Windows 构建和验收记录](https://github.com/xuli82221-pixel/helloThoughts/actions/runs/41762640507)已成功完成。发布后的文档修订可以位于 `main` 的后续提交，不改变该标签或已发布文件。

preview.21 核心 SHA256：`6f045833ae29839955edee44fed9f36b0b76b84566951b47c3615ce5f4c9a021`。
preview.21 安装器 SHA256：`869ee203c41eda411c6f337390ed2a10327e69ba6b708f244f45dd14c8202d47`。
ZIP 与 MCPB 的核心字节一致；七个上传文件的 GitHub 服务端 SHA256 均已与下载后的 CI 产物核对一致。完整包校验值以该 Release 的 `SHA256SUMS.txt` 为准。

## preview.21 已完成的检查

| 项目 | 实际结果与范围 |
| --- | --- |
| 源码 CI | Windows/Linux 检查通过；不是 Linux 产品发行 |
| 前端 | 25 项测试、TypeScript/Vite 构建通过 |
| 后端与预检 | Windows 原生 Go 测试、Go vet、文档与源码预检、发行边界和八语言声明测试通过 |
| 包结构与完整性 | Setup.exe/ZIP/MCPB 的结构、架构、清单、许可证及校验通过 |
| 真实 stdio MCP | 使用虚构资料完成初始化、读写、搜索、图片、冲突拒绝、移动、回收站、恢复及服务重开测试 |
| 隔离安装 | GitHub 托管 Windows runner 上通过安装、实际开始菜单快捷方式、重装保留数据、错误哈希拒绝及卸载保留共享资料测试 |
| 签名检查 | `thoughts.exe`、安装器、注册辅助脚本为 `NotSigned`；没有可信发布者签名 |
| 公开发布 | 人工授权后上传并公开为 Pre-release；构建工作流本身不发布 Release |

测试原始目录、会话、凭据和机器日志未作为 Release 文件上传。GitHub runner 的测试不能替代干净 SAC 环境或真实 AI 客户端验收。

## 下载与安装

从上述实际 Release 下载所需包和 `SHA256SUMS.txt`，按 [INSTALL.md](../INSTALL.md) 核对完整性与未签名状态。GitHub 的 Code → Download ZIP 是源码，不是应用。

当前文件未签名，Windows SAC、SmartScreen 或组织策略可能拦截。**若被拦截，请停止并保持防护开启；不要关闭 SAC、修改执行策略、移除下载来源标记或绕过警告。** SHA256 不证明发布者身份。预编译包运行时不需要 Go、Node.js、Python、Docker 或 Thoughts 账号；AI 客户端有自己的要求。

安装为当前用户创建程序入口和快捷方式，资料默认在 Windows 用户的 Thoughts 主目录中。卸载保留共享笔记、图片、配置和受管理核心供其他客户端使用。旧开发版迁移需审阅安装说明并将 MCP 入口重新注册到 `thoughts.exe`；本次发行不会自动替换既有安装。软件回退不代替笔记备份。

## 已发布预览版仍未完成的工作

1. **可信签名与 SignPath 审批。** 申请材料已准备，申请未由本次工作提交或获批；签名 CI、MFA 核实、受信来源绑定和产品／版本元数据约束仍待完成。见 [Code signing policy](CODE_SIGNING_POLICY.md)。
2. **干净 Windows／SAC 验收。** GitHub runner 安装测试通过不代表所有 SAC 开启的干净系统可安装或运行。
3. **三个真实 AI 客户端验收。** 真实 stdio 协议测试不等于 Claude Desktop、Claude Code、OpenAI Codex 均已完成最终宿主验收。
4. **签名稳定版发行。** 仍需签名、安装器／卸载器验证、客户环境测试和单独人工批准。当前未签名预览的发布不使这些门槛自动完成。
5. **继续保护资料和私有历史。** 旧版本、备份和原私有 Git 历史保留；公开仓库使用独立的公开历史。不得把原私有历史整体推送。

<!-- historical-records:start -->

## 历史记录：preview.20 本地准备与 SAC 拦截

本节仅记录 preview.20 及更早安装器的结果，不作为 preview.21 的安装、哈希或发布状态。preview.20 阶段完成了源码公开准备；当时的授权范围限于源码，二进制文件保留为本地审查候选。维护者后来单独授权 preview.21 未签名预览发布，已按上文完成。

preview.20 本地核心 SHA256 为 `df959a3a34d051d8587f586f9418d2f9b6c13df17a2c1af28b3b9d791678bfda`；该阶段原生 Windows 后端测试记录为 59 项通过、0 失败、0 跳过，前端为 25 项测试通过。该阶段测试记录未包含在公开文档中。这些数量和哈希属于旧版本，不能套用到 preview.21。

preview.20 最终本地安装器 SHA256 为 `441cbabb4747fd0ac29d4c4d34aed37175041abf4435f25d4bfa08c282e5a6c2`，Authenticode 为 `NotSigned`；在维护者机器上启动返回 Windows 错误 4551。2026-11-04 06:19:59（本地时间）的 Code Integrity 事件 3033／3077 记录了签名级别或代码完整性策略拦截，后续快捷方式验收因此未完成。该历史拦截结果没有被改成通过，也没有通过修改安全策略规避。

更早通过独立安装流程测试的安装器 SHA256 为 `f1c25d9e8da9c4de708e1ebc6baac9bdd548eefaae7f1b046810e3c94f18c578`。它与 preview.20 最终安装器和 preview.21 发行安装器均不是同一文件，其测试不能证明另外两者通过。

旧文件、原始测试资料和备份继续保留；本次文档修订不替换已发布文件、不重打版本标签，也不改变任何签名或验收结果。

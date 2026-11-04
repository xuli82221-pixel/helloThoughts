# preview.2／preview.4 Windows 安全评估（历史记录）

本文件仅保留所标版本的历史测试结果。当时的本地状态、测试数量、哈希和验收缺口不代表当前版本；当前 preview.21 已发布未签名预览，见[发行状态](RELEASE_READINESS_20261104.md)。

<!-- historical-records:start -->

## Historical records / 历史记录

2026-11-01 补充：当时本地界面候选为 `0.3.0-preview.4`。新 exe 和实际服务仍为 NotSigned，核心 SHA256 `8c2a7dd55cf67d64c3bebee51ef9d5646c5811a8345060e78dd21c827686c96b`；同机 SAC On，单文件 Defender 扫描完成且匹配威胁记录 0，安全设置未改动。新证据在 `preview4-security-evidence.json`、`preview4-defender-scan.json`，详见 [界面验收](UI_TABS_IMAGES_ACCEPTANCE.md)。下面保留 .2 的原始测量，不能把其包哈希或扫描结果移用到新版。干净 SAC On/Off 两组验收仍缺失。

评估日期：2026-11-01（Asia/Shanghai）。候选版本：`0.3.0-preview.2`。仅本地准备，未公开推送、未发布。

**结论：签名检查及本机运行测试已完成；干净 Windows 11 的 SAC On / Off 两组正式安装验收均未完成。Windows 安全兼容性不能判定通过，M0 仍未通过。**

## 实测环境与边界

| 项目 | 实际值 |
| --- | --- |
| 系统 | Windows 11 Pro x64，10.0.26300 |
| SAC | On；只读检查，未改动 |
| Defender | AntivirusEnabled / RealTimeProtectionEnabled / IsTamperProtected 全部 true |
| Defender 引擎产品版本 | 4.18.26100.4 |
| 签名库 | 1.463.58.0 |
| 测试机性质 | 已使用的开发电脑，不是干净环境 |
| 本地文件来源 | 本机编译与打包；ZIP、MCPB 和源目录 exe 均未带 Zone.Identifier |
| 虚拟机能力 | 未找到 Get-VM、Windows Sandbox、VMware 或 VirtualBox 测试入口，也未找到已运行的 Windows 测试 VM |

存在 vmcompute 服务不等于已有可用 Windows VM。本次没有启用虚拟化功能、修改 SAC、修改 Defender、添加排除项、安装证书或尝试绕过拦截。没有使用用户真实资料库，也没有替换现有客户端扩展。源码和测试包均留在本地。

## 两个包与真正运行的程序

| 文件 | SHA256 |
| --- | --- |
| thoughts-0.3.0-preview.2-windows-amd64.zip | a21d5ee8e196e63ffc29afa9c7f441ed14ae457ff2d5796f7d361e301ee7c7b6 |
| thoughts-claude-desktop-0.3.0-preview.2-windows-amd64.mcpb | 2dde36daaedd656320054b14bdf9bf99025fc8af6ad68be9a5d03e15bb92dcdf |
| 两包内的 memo.exe | c92bbf64b056637f6881c2ceaf1acf09a111e37b6bf68541441981989685a7ed |

1. ZIP 的 `memo.exe` 与 MCPB 的 `server/memo.exe` 字节一致，均为 Windows AMD64 PE32+。
2. `Get-AuthenticodeSignature` 均为 **NotSigned**，无签名证书；PE Security Directory 的证书偏移及大小均为 **0**。注册用 PowerShell 脚本也未签名。
3. 官方 `@anthropic-ai/mcpb@2.1.2` 清单验证通过；`mcpb verify` 返回 **ERROR: Extension is not signed**。MCPB 的包签名与内部 exe 的 Windows Authenticode 是不同层级；不能互相替代。
4. 对实际运行的隔离服务进程重新读取 exe SHA256，和上述候选 exe 完全一致，签名仍为 NotSigned。这次运行的是候选程序，而不是旧服务或前端开发服务器。
5. 两个归档均通过 CRC、路径、逐文件 SHA256、清单、许可文件与同核心校验。SHA256 证明内容一致，**不能代替可信发布者签名**。

## 本机实际结果

| 测试 | 预期与实际结果 | 范围 |
| --- | --- | --- |
| Windows x64 编译 | 完成，内嵌 React 网页 | 开发机 |
| 前端 | 12 项测试通过，TypeScript/Vite 构建通过；仍有大 chunk 提示 | 开发机 |
| Go | 全套通过，internal/memo 用时约 10.9 秒 | 开发机、隔离临时目录 |
| 本地 stdio MCP | 初始化、14 个工具、重复初始化、真实 Markdown 读写搜索、旧 revision 拒绝、图片读写、改名移动、回收站恢复、路径隔离通过 | 直接协议测试，不等于三个宿主验收 |
| 生命周期 | 关闭 MCP 后服务继续；终止经过身份核验的测试服务并重开后，文件保留 | 隔离库 |
| 内置浏览器 | 实际程序网页可访问；默认预览、明确编辑、保存后核对磁盘一致、重载仍预览 | Codex 内置浏览器 |
| 外部文件刷新 | 本次观测 424 ms | 单次观测，不是性能保证 |
| Defender 自定义扫描 | 对候选 exe 扫描完成，约 5.79 秒；匹配该文件的威胁记录 0 条 | 当时引擎/签名库；不代表未来或其他机器零误报 |
| Claude Desktop 实际导入此版本 | 未测试 | 清单有效不能代替实际安装 |
| Claude Code / Codex 实际宿主工具调用此版本 | 未测试 | 直接 stdio 成功不能代替宿主调用 |

Windows CodeIntegrity 最近 7 天、有界读取最多 2000 条中，找到 14 条与 Thoughts/Memo 名称有关的历史事件，日期为 10 月 27 日和 29 日，含 3077、3033。3077 记录的策略为 `VerifiedAndReputableDesktop`，状态 `0xc0e90002`；提取的历史文件平面 SHA256 与本次候选 exe 不同。本次候选实际成功启动，不能把旧记录说成此次被拦截，也不能因此断言没有其他策略。

Defender 同一有界查询读取 604 条事件，名称匹配 0 条；这只是查询范围内没有匹配记录。没有查看或导出与项目无关的威胁正文、用户笔记、令牌或完整本机路径。

## 三种安全机制要分开判断

| 机制 | 可能看到什么 | 本候选的判断 |
| --- | --- | --- |
| SmartScreen 下载信誉 | 未知应用/发布者警告；企业策略也可能禁止继续 | 下载路径尚未实测；本机构建无 MOTW，不覆盖网络下载情景 |
| Smart App Control | 执行控制直接拒绝启动，常伴 CodeIntegrity 事件；不是普通可继续的警告 | 未知、未签名程序有被阻止风险；本机 On 可运行不能推广到干净 On 环境 |
| Defender 防病毒 | 检测、隔离或阻止可疑文件 | 本次单文件扫描未发现匹配威胁记录；其他版本/引擎仍需复测 |

SmartScreen 依据发布者和文件信誉，新文件即使有效签名也可能先出现警告；自签名不解决信誉问题，EV 也不保证立即免警告。[Microsoft SmartScreen 文档](https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/smartscreen-reputation)

SAC 对无法识别为可信且没有受信任签名的代码可能直接阻止；Evaluation 与 On enforcement 不同。不能把本次开发机成功当成所有 SAC On 设备可用。[Microsoft SAC 文档](https://learn.microsoft.com/en-us/windows/apps/develop/smart-app-control/overview)

## 正式兼容性矩阵

| 环境 | 安装 | 启动 | 判定 |
| --- | --- | --- | --- |
| 当前开发机，SAC On、Defender On、无 MOTW 的本机构建 | 只做归档/清单验证，未导入 MCPB | CLI、stdio、服务、内置网页成功 | 仅该机器与该哈希通过运行检查 |
| 干净 Windows 11 x64，SAC On、Defender On，真实下载包 | 未测试 | 未测试 | **未通过验收；可能硬拦截，不能保证直接运行** |
| 干净 Windows 11 x64，SAC Off、SmartScreen/Defender 保持正常，真实下载包 | 未测试 | 未测试 | **未通过验收；可能有 SmartScreen 警告或其他策略阻止** |
| 企业 WDAC / AppLocker 策略设备 | 未测试 | 未测试 | 不承诺支持；SAC Off 不等于所有执行控制关闭 |

目前不能诚实列出“所有可直接运行的 Windows 环境”。缺的是两份独立、已合法配置好模式的干净 Windows 11 测试环境，以及三个客户端的真实导入/调用结果。**本次未要求关闭当前电脑任何安全防护。**

## 复测步骤（仅在已有隔离环境执行）

准备两个独立 Windows 11 x64 快照，一个已经处于 SAC On，另一个已经处于 SAC Off，Defender 和 SmartScreen 保持正常。不要为测试去切换用户真实电脑的模式。不能使用普通共享宿主临时目录代替隔离环境，也不把 Windows Sandbox 自动当成 SAC On 的替身。

1. 记录系统、三个客户端版本、SAC 实际状态、Defender 版本与签名库、网络及企业策略状态。
2. 通过真实浏览器下载同一哈希包，保留文件来源标记。记录下载、解压、MCPB 导入、实际启动 exe 的来源标记与警告截图；不要执行 Unblock-File 或点击绕过安全警告。
3. 运行只读 `scripts/probe-windows-security.ps1 -Candidate <解包后的候选构建目录> -Output <证据.json> -ExpectedSAC On`（另一快照为 Off）。Candidate 指构建输出目录，包含 index.json、两个包和 sources；单独下载的 ZIP 不含这些检查材料。脚本只读安全设置并写报告，模式不符立即失败。
4. 在 Claude Desktop 导入 MCPB；在 Claude Code 和 Codex 按 INSTALL.md 用显式隔离 THOUGHTS_HOME 注册同一 exe。安装命令和 SHA256 先检查再执行，不覆盖既有配置。
5. 每个宿主验证实际进程路径/哈希、工具调用、初始化与重复安装、真实创建/读取/修改/搜索/移动/图片/恢复/冲突、网页打开及重启数据保留。
6. 出现警告或硬拦截立即保留时间、错误码、相关事件 ID 和截图，停止该路径。没有打开本机安全设置、添加排除、改策略或装证书的复测步骤。
7. 记录每步耗时、预期、实际结果、失败日志和复现步骤。模式未经核验、未走真实下载或只完成 CLI，都不能记为整体验收通过。

## 解决方向

- 当前包继续标为本地未签名预览，暂不宣称任何 Windows 用户安装即用。
- 正式分发前申请可信 CA 代码签名或 Microsoft Artifact Signing，签署实际 exe 及分发脚本并加可信时间戳，再重新打包/校验/执行两组隔离验收。MCPB 签名另行处理，不代替内部 exe 签名。[Microsoft SAC 签名文档](https://learn.microsoft.com/en-us/windows/apps/develop/smart-app-control/code-signing-for-smart-app-control)
- 签名不保证 SmartScreen 立刻建立信誉；每个新哈希保留可复现构建、SHA256、签名验证与测试证据。若确认误报，走 Microsoft 官方提交渠道复核。
- 不采用自签名、不安装不可信根证书、不关闭保护、不换壳或反复重编译来规避拦截；开源本身不解决签名或信誉问题。

本地证据位于忽略的 `.artifacts` 中：`preview2-build.log`、`preview2-mcp-smoke.log`、`preview2-security-evidence.json`、`preview2-codeintegrity-metadata.json`、`preview2-defender-scan.json`、`preview2-mcpb-signature.log`、`preview2-browser/`。这些不是公共上传清单，不包含于发行包。

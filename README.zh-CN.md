[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Deutsch](README.de.md) | [Français](README.fr.md) | [Español](README.es.md)

# Thoughts

[安装](INSTALL.md) · [文档](docs/README.md) · [捐赠](DONATE.md) · [代码签名政策 / Code signing policy](docs/CODE_SIGNING_POLICY.md)

免费的开源 Windows x64 本地优先笔记工具。让 AI 助手管理 Markdown 文件，在独立的本地网页中阅读和编辑。

**源码预览 — 2026-11-04：**本仓库以全新的公开历史发布经过审查的 `0.3.0-preview.20` 源码，采用 MIT 许可证。59 项后端测试已在原生 Windows 上通过。**目前没有获准下载的 Windows 安装程序：**SAC 拦截了最终的未签名 Setup.exe（错误 4551）。仓库包含安装器源码，二进制发行仍需可信签名和最终验收。参见[发行准备](docs/RELEASE_READINESS_20261104.md)、[安装说明](INSTALL.md)和 [Windows 测试](docs/WINDOWS_NATIVE_TESTS_20261104.md)。

## 笔记就在你的电脑上

Thoughts 不需要账号，没有云端数据库、笔记上传服务、遥测、内置 AI、付费主题、订阅、支付集成或浏览器扩展。Markdown 和本地图片属于你，可以用其他编辑器打开。

Go 可执行文件内嵌 React 网页。**预编译版本运行时不需要 Node.js、Go、Python 或 Docker。**你选择的 AI 客户端有自己的安装和账号要求。

## 将来如何安装

经过审查的发行版本可用后，可从 [GitHub Releases](https://github.com/xuli82221-pixel/hello-Thoughts-/releases) 下载 **`thoughts-<version>-windows-amd64-setup.exe`**，核对校验值和发布者签名，再双击并选择 **Install**。安装器无需管理员权限，会初始化共享主目录并创建桌面和开始菜单快捷方式。已有数据会保留；卸载移除入口和快捷方式，保留共享主目录供其他客户端使用。届时也计划提供便携 ZIP 供手动安装；目前这两种包都不是已获准的下载。

Claude Desktop 用户届时也可导入发行版本的 `.mcpb` 扩展。Codex 和 Claude Code 用户将 ZIP 内的可执行文件注册为本地 stdio MCP 服务。网页也能在没有 AI 客户端时运行。具体命令见[普通用户快速入门](INSTALL.md#ordinary-windows-users)。GitHub 的 **Code → Download ZIP** 下载的是源码，而非预编译应用。

AI 客户端启动同一个本地标准 stdio MCP 适配器，再由它调用经过身份验证的本地核心；网页仅监听 `127.0.0.1`。当你要求在线 AI 客户端读取笔记时，内容会交给该客户端，并受其数据政策约束。本地优先不代表主动调用 AI 得到的工具结果绝不离开电脑。Thoughts 不加密文件。

## 目标客户端

| 客户端 | 安装方式 |
| --- | --- |
| Claude Desktop | 官方 `.mcpb` 扩展格式 |
| Windows 上的 Claude Code | 本地 stdio MCP，便携 Windows ZIP |
| Windows 上的 OpenAI Codex | 本地 stdio MCP，同一个便携 ZIP |

这些是目标集成；实际验收状态以发行审查为准。本阶段不支持普通 ChatGPT 网页插件、macOS、Linux、移动端、远程 MCP 或云同步。

## 功能

- 在本地创建、读取、编辑、搜索、重命名和移动 Markdown。
- 紧凑的项目树、快速笔记、项目归档与取消归档，以及可恢复的笔记和项目回收站。
- 已有笔记默认以渲染预览打开，编辑需主动切换。
- 本地 PNG/JPEG/WebP 图片、预览和大图查看。
- 免费主题及经过校验的本地主题导入。
- 英语、韩语、日语、简体中文、繁体中文、德语、法语和西班牙语界面，以及对应的本地隐私与使用说明。
- 外部文件变更通知和 SHA256 冲突保护。
- 所有客户端使用同一个已配置资料库；重复初始化保留已有选择。
- 经校验值验证的本地核心升级和软件回退，并保护未保存的编辑。
- 服务端保存笔记标签页和阅读位置，服务端口变化后仍可恢复。

冲突保护会拒绝过时写入；尚未实现自动修订历史。请单独备份资料库。打开的标签页可在冲突或断线时保留草稿，但关闭标签页或重启 Windows 后不保证恢复。

![Thoughts 0.3.0-preview.20 与虚构示例笔记](docs/assets/preview-0.3.0-preview.20.jpg)

这是 `0.3.0-preview.20` 的真实界面，于 2026-11-04 在全新的独立临时资料库中拍摄，内容均为虚构笔记。画面没有个人笔记、真实图片、本机路径或凭据。[旧版界面截图](docs/assets/preview.png)作为历史资料保留，不代表当前版本。截图核验不等于安装器或客户端验收。

## 安装

SHA256 核验、先检查后执行的命令和三个客户端步骤见 [INSTALL.md](INSTALL.md)。从 preview.8 起，新安装在所有客户端中默认使用同一个 Windows 用户 Thoughts 目录，不受安装顺序影响。高级主目录选项保持 `auto`。只有刻意自定义目录时才需设置绝对路径 `THOUGHTS_HOME`：

```text
<Thoughts home>/
  library/     Markdown, images and recoverable Trash
  config/      Configuration, local sessions, preferences and logs
  app/         Managed executable versions (optional)
  backups/     Previous runtime pins; not a backup of note bodies
```

一般安装请使用共享默认目录。旧版显式 `MEMO_CONFIG_DIR` 仍可使用；本版不会自动迁移或合并已有资料库。若主动自定义目录，需要共享内容的客户端应使用同一个物理文件夹。

安装和更新必须保留该主目录。阅读状态保存路径和滚动位置，不保存未提交的正文。核心更新目前需要明确提供经过审查的本地可执行文件及 SHA256，没有后台联网更新器。客户端扩展和共享核心分别更新，参见[升级说明](INSTALL.md#managed-core-updates-and-rollback)。

未签名的 Windows 程序或脚本可能被 SmartScreen、Smart App Control 或组织应用控制策略拦截。SHA256 只能与校验值比对完整性，不能证明发布者身份。不要为了运行预览版关闭防护或绕过警告。公开源码本身无法解决签名问题。

## 工作区

参见[工作区布局](docs/WORKSPACE_LAYOUT.md)和[文档索引](docs/README.md)。旧需求和报告保留在本地归档。内部 Memo 路径是兼容名称，不代表额外产品，也不是可以随意删除的文件。

## 构建

构建需要 Windows x64、与 go.mod 匹配的 Go、Node.js 24、Python 3.12+ 和有数字签名的官方 [Inno Setup 7 编译器](https://jrsoftware.org/isdl.php)。最终用户不需要这些工具。参见 [CONTRIBUTING.md](CONTRIBUTING.md)。

```powershell
Set-Location web
npm ci
Set-Location ..
pwsh -File scripts/build-desktop-adapters.ps1 -Python 'C:\Tools\Python\python.exe' -Go 'C:\Tools\Go\bin\go.exe' -InnoCompiler 'C:\Tools\InnoSetup7\ISCC.exe'
```

请将工具路径换成实际安装位置。构建流程测试并构建前端和 Go 核心、生成声明、创建 Setup.exe/ZIP/MCPB 及校验值，并在新输出目录中验证包。不会覆盖已有输出，也不会自动上传。公开候选模式要求可信签名和真实发布者元数据；签名不能代替干净电脑和客户端验收。

## 安全与许可证

参见 [SECURITY.md](SECURITY.md)、[PRIVACY.md](PRIVACY.md)、[CHANGELOG.md](CHANGELOG.md)、[第三方声明](THIRD_PARTY_NOTICES.txt)和[依赖审查](docs/DEPENDENCIES.md)。项目代码采用 MIT 许可证；捆绑依赖保留各自许可证，包括 OFL 字体。

不要提交真实笔记、个人截图、凭据、会话文件或日志。漏洞请私下报告至 hello@spxipo.com。设置中提供维护者 Base 主网地址的可选捐赠入口。二维码仅包含地址；请在钱包中单独选择 Base。Thoughts 不连接钱包、不索取私钥、不指定币种或金额，也不处理支付。所有功能始终免费。

## 代码签名政策

**待批准 — SignPath Foundation 申请尚未获批。**所有当前 Windows 程序仍未签名。若获批准，免费代码签名由 SignPath.io 提供，证书由 SignPath Foundation 提供。具体范围、真实团队职责和待办条件见[代码签名政策 / Code signing policy](docs/CODE_SIGNING_POLICY.md)及[申请准备](docs/SIGNPATH_APPLICATION.md)。本次准备未提交申请。

## 捐赠

支持完全自愿，不解锁任何功能。独立[捐赠页面](DONATE.md)提供维护者已有的 Base 地址、仅含地址的二维码和隐私边界。你也可以通过反馈问题、翻译或文档帮助项目。Thoughts 不连接钱包，也不处理支付。

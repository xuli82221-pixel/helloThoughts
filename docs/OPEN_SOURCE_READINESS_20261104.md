# preview.20 开源展示准备验收（历史记录）

**当前公开版本为 `0.3.0-preview.21`，已公开发布为未签名预览版。** 详见[实际 Release](https://github.com/xuli82221-pixel/hello-Thoughts-/releases/tag/v0.3.0-preview.21)及[当前发行状态](RELEASE_READINESS_20261104.md)。签名与稳定版验收仍未完成。

本文件保留较早 preview.20 文档准备阶段的文件范围、31 份 Markdown／340 个链接等当时检查结果，不作为 preview.21 或后续文档修订的当前统计。维护者后来分别授权并完成文档推送、命名统一和未签名预览发布。

<!-- historical-records:start -->

## 历史记录：preview.20 文档准备阶段

preview.20 文档准备阶段仅修改公开文档、截图和文档导出／检查工具。该阶段的产品版本为 `0.3.0-preview.20`，目标平台为 Windows x64。该阶段已有公开源码，安装包仍为本地候选；preview.21 的后续二进制发布见上方链接。仓库原有大量未提交改动已保留，验收以本次开始时的文件哈希为比较基准。

### 后续提交与推送授权

维护者在本报告所记录的准备工作完成后，明确授权提交并推送本次改动。提交从公开仓库的 `main` 建立独立工作区，仅包含本报告列出的 22 个文档、截图及配套工具文件，不合并原私有历史或其他产品代码改动。下文“未提交／未推送”等表述记录的是原准备阶段；最新公开提交以仓库历史为准。签名申请仍未提交，Windows 程序仍未签名，该次文档授权仅涉及源码；维护者后来另行授权并完成了 preview.21 未签名预览发布。

### 该阶段已完成

- 保留英文 `README.md` 为内容基准，新增完整的简体中文、繁体中文、日文、韩文、德文、法文、西班牙文 README。
- 八份 README 使用相同的语言导航，并覆盖本地存储、未来安装方式、目标客户端、十项功能、安装目录、工作区、构建、安全与许可、签名政策和捐赠。代码示例、路径和链接目标保持一致。
- 新增 [Code signing policy](CODE_SIGNING_POLICY.md) 和 [SignPath 申请准备](SIGNPATH_APPLICATION.md)，并从全部 README、INSTALL、SECURITY、CONTRIBUTING 和文档索引连接到政策。
- 新增独立 [Donate 页面](../DONATE.md)，全部 README 顶部和独立章节均提供入口；复用维护者已配置的 Base 地址及地址二维码，没有编造捐赠平台或新增支付集成。
- GitHub CLI 已认证为实际管理员 `xuli82221-pixel`；仓库 Description 已更新为 `Free, local-first Markdown workspace for AI assistants on Windows.`，九个指定 Topics 已写入并经 API 回读确认。Homepage 保持为空。
- 更新为真实的 [preview.20 界面截图](assets/preview-0.3.0-preview.20.jpg)。使用全新临时资料库和虚构笔记，不读取个人资料库；画面无个人笔记、真实图片、本机绝对路径或凭据。旧截图保留不变。见[截图来源记录](assets/README.md)。
- 更新源码导出白名单，避免后续导出遗漏七份译文、Donate、政策、申请草稿、验收报告及新截图。没有生成或发布新的二进制包。

### 八语验收

| 语言 | 文件 | 完整性与链接检查 |
| --- | --- | --- |
| English | [README.md](../README.md) | 通过：完整 10 节；导航、链接与锚点正常 |
| 简体中文 | [README.zh-CN.md](../README.zh-CN.md) | 通过：完整 10 节；导航、链接与锚点正常 |
| 繁體中文 | [README.zh-TW.md](../README.zh-TW.md) | 通过：完整 10 节；导航、链接与锚点正常 |
| 日本語 | [README.ja.md](../README.ja.md) | 通过：完整 10 节；导航、链接与锚点正常 |
| 한국어 | [README.ko.md](../README.ko.md) | 通过：完整 10 节；导航、链接与锚点正常 |
| Deutsch | [README.de.md](../README.de.md) | 通过：完整 10 节；导航、链接与锚点正常 |
| Français | [README.fr.md](../README.fr.md) | 通过：完整 10 节；导航、链接与锚点正常 |
| Español | [README.es.md](../README.es.md) | 通过：完整 10 节；导航、链接与锚点正常 |

翻译逐节人工核对英文基准；自动检查只验证结构、目标链接、锚点与代码块，不构成母语专家校审。所有版本保留在线 AI 客户端的数据边界、文件不加密、用户自行备份、无内置 AI／付费功能、当时未签名、仅有源码公开的状态的说明。

### 该阶段的 SignPath 评估

已具备的项目基础：公开仓库、项目代码 MIT 授权、近期维护和成功的源码 CI、清晰的产品／未来下载说明、本地隐私说明、真实单维护者职责和签名政策草稿。仓库控制权由认证后的管理员权限核实。

尚未满足或无法证实的条件：

1. **公开发行形态：**preview.20 准备阶段检查时，GitHub Releases 和版本标签列表为空；本地测试包与源码 ZIP 不能冒充已公开发行的可执行程序。
2. **可验证声誉：**检查时仓库 0 stars、0 forks、0 subscribers，创建时间很近；没有可引用的用户数量或独立使用证据。不能认定声誉充分，也不编造星数门槛。
3. **签名来源验证：**现有 CI 只检查源码，没有可信 Windows 发行产物、SignPath 来源绑定、版本标签、受保护的人工签名审批或 PE 产品／版本元数据约束。
4. **账号与授权：**GitHub／SignPath MFA、SignPath 团队与角色尚未核实或配置。不能将 GitHub 管理权限等同于 MFA 或签名资格。
5. **签名范围与许可证：**该阶段拟申请第一方核心（当时文件名 `memo.exe`；preview.21 已统一为 `thoughts.exe`）。Inno Setup 安装器／卸载器、PowerShell 辅助脚本的签名路径及混合组件／字体许可证需另行审查。仅签核心不能通过现有公开候选构建门禁。
6. **发行验收：**最终安装器曾被 SAC 拦截，签名、干净 Windows、安装器／卸载器及三个真实客户端的最终验收仍未完成。

申请条件依据 [SignPath 官方条件](https://signpath.org/terms.html)和[官方申请页](https://signpath.org/apply.html)，读取日期为 2026-11-04。申请草稿列明证据及局限，最终资格由 Foundation 判断。

### 该阶段列出的后续工作

- 审阅本次文件差异，选择要公开的该阶段源码快照，再手动提交和推送。原始私有历史和本机资料不得一并推送。该阶段完成文件准备时尚未推送；随后这些文档、Donate 入口和截图已提交至公开仓库，preview.21 也已另行发布。
- 核实联系信息、GitHub MFA，并在建立 SignPath 团队后启用其 MFA；检查真实协作者名单。
- 决定先就“源码预览／尚无发行／声誉不足”向 Foundation 咨询，还是等待条件成熟，再手动提交申请。没有自动发送邮件或表单。
- 批准签名 CI 的具体设计和凭据权限后再接入；本轮没有新增 Secrets。签名后重新校验哈希、签名与最终包，再完成验收并单独批准发行。

### 需要 SignPath 审批

- 项目是否符合发行、声誉、许可证及安全要求。
- 初始第一方核心范围、安装器运行时处理方式及单维护者角色安排。
- 受信构建集成、签名配置和每次签名请求。无批准保证，不以文档或徽章暗示已认证。

### 该阶段检查结果

最终检查全部通过：31 份公开 Markdown、340 个相对链接（含锚点及公开导出覆盖）；八份 README 的 10 个对应章节、10 项功能、相同代码块与相同链接目标全部通过。源码预检覆盖 184 个文件，无发现；三项发行边界测试通过，0 失败。`git diff --check` 未发现空白错误。本轮哈希比较确认只修改了下列 8 个已有文件，并新增 14 个文件；无意外修改，旧截图哈希不变。截图临时服务已正常关闭，临时测试资料保留在忽略目录中。检查命令为 `python scripts/check-docs.py`、`python scripts/release-preflight.py`、`python -m unittest discover -s scripts -p "test_release_preflight.py"`。文档检查不会访问捐赠钱包、提交申请或发布文件。

相对链接和锚点可自动验证；外部 GitHub 和 SignPath 关键页面已实际读取。其他第三方链接不作永久可用保证。该阶段准备的文档随后已推送，远端入口现已生效。

### 文件范围

新增：七份语言 README、`DONATE.md`、`docs/CODE_SIGNING_POLICY.md`、`docs/SIGNPATH_APPLICATION.md`、本报告、`docs/assets/README.md`、新版 JPG 截图、`scripts/check-docs.py`。

修改：英文 README、INSTALL、SECURITY、CONTRIBUTING、CHANGELOG、文档索引、源码文件清单、`scripts/prepare-public-source.py`。PRIVACY 已核对，不需要改写；其本轮开始前已有的改动保留。本轮哈希核对已确认以上范围。

仅就该文档准备阶段而言，未改变产品逻辑、真实资料、已安装核心、安全设置或旧版本，且未在该阶段执行签名、申请提交或二进制发布。后续经授权的源码提交、Thoughts 命名统一及 preview.21 发布另有独立记录。

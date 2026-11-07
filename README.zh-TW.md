[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Deutsch](README.de.md) | [Français](README.fr.md) | [Español](README.es.md)

# Thoughts

[安裝](INSTALL.md) · [文件](docs/README.md) · [捐款](DONATE.md) · [程式碼簽章政策 / Code signing policy](docs/CODE_SIGNING_POLICY.md)

免費、開放原始碼的 Windows x64 本機優先筆記工具。讓 AI 助理管理 Markdown 檔案，並在獨立的本機網頁中閱讀與編輯。

**未簽章預覽 — 2026-11-04：**`0.3.0-preview.21` 將原始碼、執行檔與用戶端設定統一命名為 Thoughts。維護者已授權發布[未簽章實驗預覽版](https://github.com/xuli82221-pixel/hello-Thoughts-/releases/tag/v0.3.0-preview.21)，檔案已在公開 Windows 建置與檢查通過後提供。**不代表已有可信發行者簽章或穩定版驗收。**SignPath Foundation 尚未核准申請。Windows SAC/SmartScreen 可能封鎖；遇到封鎖請停止並保持防護啟用。請見[發行準備](docs/RELEASE_READINESS_20261104.md)、[安裝說明](INSTALL.md)。

## 筆記留在你的電腦上

Thoughts 不需要帳號，沒有雲端資料庫、筆記上傳服務、遙測、內建 AI、付費佈景主題、訂閱、付款整合或瀏覽器擴充功能。Markdown 與本機圖片屬於你，也能用其他編輯器開啟。

Go 執行檔內嵌 React 網頁。**預先編譯版本執行時不需要 Node.js、Go、Python 或 Docker。**你選用的 AI 用戶端另有自己的安裝與帳號要求。

## 下載與安裝

本實驗預覽版透過 [GitHub Releases](https://github.com/xuli82221-pixel/hello-Thoughts-/releases) 提供 **`thoughts-<version>-windows-amd64-setup.exe`**、可攜式 ZIP 與 Claude Desktop MCPB，內部程式名稱為 **`thoughts.exe`**。使用前請核對 SHA256；目前檔案**沒有 Authenticode 簽章**，雜湊相符無法證明發行者身分。Windows 可能拒絕執行，遇到封鎖請停止。安裝程式按使用者安裝，無需系統管理員權限，會初始化共用主目錄並建立桌面與開始功能表捷徑。既有資料保留；解除安裝移除入口與捷徑，保留共用主目錄供其他用戶端使用。已簽章穩定版仍待完成。

Claude Desktop 使用者屆時也可匯入該版本的 `.mcpb` 擴充功能。Codex 與 Claude Code 使用者將 ZIP 中的執行檔註冊為本機 stdio MCP 伺服器。網頁也可在沒有 AI 用戶端時運作。具體指令請見[一般使用者快速入門](INSTALL.md#ordinary-windows-users)。GitHub 的 **Code → Download ZIP** 下載的是原始碼，而非預先編譯的應用程式。

AI 用戶端啟動相同的本機標準 stdio MCP 配接器，再呼叫經過驗證的本機核心；網頁僅監聽 `127.0.0.1`。當你要求線上 AI 用戶端讀取筆記時，內容會交給該用戶端，受其資料政策規範。本機優先不代表主動呼叫 AI 所得的工具結果絕不離開電腦。Thoughts 不會加密檔案。

## 目標用戶端

| 用戶端 | 安裝方式 |
| --- | --- |
| Claude Desktop | 官方 `.mcpb` 擴充功能格式 |
| Windows 上的 Claude Code | 本機 stdio MCP，可攜式 Windows ZIP |
| Windows 上的 OpenAI Codex | 本機 stdio MCP，相同的可攜式 ZIP |

以上是目標整合；實際驗收狀態請見[發行狀態說明](docs/RELEASE_READINESS_20261104.md)。本階段不支援一般 ChatGPT 網頁外掛、macOS、Linux、行動裝置、遠端 MCP 或雲端同步。

## 功能

- 在本機建立、讀取、編輯、搜尋、重新命名與移動 Markdown。
- 精簡的專案樹、快速筆記、專案封存與取消封存，以及可還原的筆記與專案資源回收筒。
- 既有筆記預設以呈現後的預覽開啟，需主動切換才能編輯。
- 本機 PNG/JPEG/WebP 圖片、預覽與大圖檢視。
- 免費佈景主題，以及經過驗證的本機佈景主題匯入。
- 英文、韓文、日文、簡體中文、繁體中文、德文、法文與西班牙文介面，以及對應的本機隱私與使用說明。
- 外部檔案變更通知及 SHA256 衝突保護。
- 所有用戶端使用相同的已設定資料庫；重複初始化會保留既有選擇。
- 經校驗值驗證的本機核心升級與軟體回復，並保護尚未儲存的編輯內容。
- 由伺服器儲存筆記分頁與閱讀位置，服務連接埠改變後仍可還原。

衝突保護會拒絕過時寫入；尚未實作自動修訂歷程。請另外備份資料庫。開啟的分頁可在衝突或斷線時保留草稿，但關閉分頁或重新啟動 Windows 後不保證能還原。

![Thoughts 0.3.0-preview.20 與虛構示範筆記](docs/assets/preview-0.3.0-preview.20.jpg)

這是 `0.3.0-preview.20` 的真實介面，於 2026-11-04 使用全新、獨立的暫存資料庫與虛構筆記擷取。畫面沒有個人筆記、真實圖片、本機路徑或憑證。[舊版介面截圖](docs/assets/preview.png)保留作為歷史資料，不代表目前版本。截圖核對不等同安裝程式或用戶端驗收。

## 安裝

SHA256 核對、先檢查後執行的指令與三種用戶端步驟請見 [INSTALL.md](INSTALL.md)。從 preview.8 起，新安裝在所有用戶端中預設使用同一個 Windows 使用者 Thoughts 目錄，不受安裝順序影響。進階主目錄設定請維持 `auto`。只有刻意自訂主目錄時才需要絕對路徑 `THOUGHTS_HOME`：

```text
<Thoughts home>/
  library/     Markdown, images and recoverable Trash
  config/      Configuration, local sessions, preferences and logs
  app/         Managed executable versions (optional)
  backups/     Previous runtime pins; not a backup of note bodies
```

一般安裝請使用共用預設目錄。舊版的明確設定仍可使用；本版不會自動遷移或合併既有資料庫。若主動自訂目錄，需要共用內容的用戶端應使用同一個實體資料夾。

安裝與更新必須保留此主目錄。閱讀狀態只記錄路徑與捲動位置，不記錄尚未儲存的文字。核心更新目前需明確提供經過審查的本機執行檔與 SHA256，沒有背景連網更新程式。用戶端擴充功能與共用核心須分別更新，請見[升級說明](INSTALL.md#managed-core-updates-and-rollback)。

未簽章的 Windows 程式或指令碼可能被 SmartScreen、Smart App Control 或組織的應用程式控制原則封鎖。SHA256 只能比對校驗值以確認完整性，無法證明發行者身分。請勿為了執行預覽版關閉防護或略過警告。公開原始碼本身無法解決簽章問題。

## 工作區

請見[工作區配置](docs/WORKSPACE_LAYOUT.md)與[文件索引](docs/README.md)。歷史開發紀錄不再納入公開文件。原始碼與執行檔已統一使用 Thoughts 名稱，舊版資料識別保留相容。

## 建置

建置需要 Windows x64、與 go.mod 相符的 Go、Node.js 24、Python 3.12+，以及具有數位簽章的官方 [Inno Setup 7 編譯器](https://jrsoftware.org/isdl.php)。使用者不需要這些工具。請見 [CONTRIBUTING.md](CONTRIBUTING.md)。

```powershell
Set-Location web
npm ci
Set-Location ..
pwsh -File scripts/build-desktop-adapters.ps1 -Python 'C:\Tools\Python\python.exe' -Go 'C:\Tools\Go\bin\go.exe' -InnoCompiler 'C:\Tools\InnoSetup7\ISCC.exe'
```

請將工具路徑換成實際安裝位置。建置流程測試並建置前端與 Go 核心、產生聲明、建立 Setup.exe/ZIP/MCPB 與校驗值，並在新的輸出目錄中驗證套件。不會覆寫既有輸出，也不會自動上傳。公開候選模式需要可信簽章與真實發行者中繼資料；簽章不能代替乾淨電腦與用戶端驗收。

## 安全性與授權

請見 [SECURITY.md](SECURITY.md)、[PRIVACY.md](PRIVACY.md)、[CHANGELOG.md](CHANGELOG.md)、[第三方聲明](THIRD_PARTY_NOTICES.txt)與[相依套件審查](docs/DEPENDENCIES.md)。專案程式碼採 MIT 授權；隨附相依套件保留各自授權，包括 OFL 字型。

請勿提交真實筆記、個人截圖、憑證、工作階段檔案或記錄。漏洞請私下回報至 hello@spxipo.com。設定中提供維護者 Base 主網地址作為自願捐款入口。QR 碼僅包含地址；請在錢包中另行選擇 Base。Thoughts 不連接錢包、不索取私鑰、不指定資產或金額，也不處理付款。所有功能持續免費。

## 程式碼簽章政策

**待核准 — SignPath Foundation 申請尚未獲准。**所有目前的 Windows 程式仍未簽章。若獲核准，免費程式碼簽章由 SignPath.io 提供，憑證由 SignPath Foundation 提供。程式碼簽章範圍、目前狀態與待完成事項請見[程式碼簽章政策 / Code signing policy](docs/CODE_SIGNING_POLICY.md)。SignPath Foundation 申請尚未提交。

## 捐款

支持完全自願，不會解鎖任何功能。獨立[捐款頁面](DONATE.md)提供維護者既有的 Base 地址、僅含地址的 QR 碼與隱私邊界。你也可以透過問題回報、翻譯與文件幫助專案。Thoughts 不連接錢包，也不處理付款。

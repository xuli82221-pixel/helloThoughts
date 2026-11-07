; Offline, current-user installer. Never remove the shared Thoughts home.
#ifndef CoreDir
  #error CoreDir must contain the reviewed Windows executable
#endif
#ifndef ProjectRoot
  #error ProjectRoot is required
#endif
#ifndef CoreSHA256
  #error CoreSHA256 is required
#endif
#ifndef AppVersion
  #error AppVersion is required
#endif
#ifndef OutputFolder
  #error OutputFolder is required
#endif
#ifndef InstallerId
  #define InstallerId "Thoughts.Desktop.Windows"
#endif
#define LauncherFolder "launchers\" + AppVersion + "-" + Copy(CoreSHA256, 1, 12)

[Setup]
AppId={#InstallerId}
AppName=Thoughts
AppVersion={#AppVersion}
AppPublisher=xuli82221-pixel
AppPublisherURL=https://github.com/xuli82221-pixel/helloThoughts
AppSupportURL=https://github.com/xuli82221-pixel/helloThoughts/issues
DefaultDirName={localappdata}\Programs\Thoughts
DefaultGroupName=Thoughts
PrivilegesRequired=lowest
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
MinVersion=10.0.19044
DisableWelcomePage=yes
DisableDirPage=yes
DisableProgramGroupPage=no
DisableReadyPage=no
ShowLanguageDialog=no
AllowNoIcons=yes
UsePreviousTasks=yes
CloseApplications=no
RestartApplications=no
RestartIfNeededByRun=no
UninstallDisplayIcon={app}\{#LauncherFolder}\thoughts.exe
UninstallDisplayName=Thoughts (program and shortcuts)
Compression=lzma2/fast
SolidCompression=yes
WizardStyle=modern
OutputDir={#OutputFolder}
OutputBaseFilename=thoughts-{#AppVersion}-windows-amd64-setup
SetupLogging=yes
#ifdef EnableSigning
SignTool=thoughts
SignedUninstaller=yes
#endif

[Languages]
Name: "en"; MessagesFile: "compiler:Default.isl"
Name: "zhHans"; MessagesFile: "compiler:Languages\ChineseSimplified.isl"
Name: "zhHant"; MessagesFile: "compiler:Languages\ChineseTraditional.isl"
Name: "ja"; MessagesFile: "compiler:Languages\Japanese.isl"
Name: "ko"; MessagesFile: "compiler:Languages\Korean.isl"
Name: "de"; MessagesFile: "compiler:Languages\German.isl"
Name: "fr"; MessagesFile: "compiler:Languages\French.isl"
Name: "es"; MessagesFile: "compiler:Languages\Spanish.isl"

[Messages]
en.ReadyLabel2a=Click Install to continue.
zhHans.ReadyLabel2a=点击“安装”即可开始。
zhHant.ReadyLabel2a=按一下「安裝」即可開始。
ja.ReadyLabel2a=「インストール」をクリックして開始します。
ko.ReadyLabel2a=설치를 시작하려면 설치를 클릭하세요.
de.ReadyLabel2a=Klicken Sie auf Installieren, um zu beginnen.
fr.ReadyLabel2a=Cliquez sur Installer pour commencer.
es.ReadyLabel2a=Haga clic en Instalar para comenzar.
en.ReadyLabel1=Install Thoughts for your Windows account. Existing notes, images and settings are kept.
zhHans.ReadyLabel1=为当前 Windows 用户安装随心记。已有笔记、图片和设置会保留。
zhHant.ReadyLabel1=為目前 Windows 使用者安裝隨心記。現有筆記、圖片和設定會保留。
ja.ReadyLabel1=現在の Windows ユーザーに Thoughts をインストールします。既存のノート、画像、設定は保持されます。
ko.ReadyLabel1=현재 Windows 사용자용 Thoughts를 설치합니다. 기존 노트, 이미지와 설정은 유지됩니다.
de.ReadyLabel1=Thoughts für dieses Windows-Konto installieren. Vorhandene Notizen, Bilder und Einstellungen bleiben erhalten.
fr.ReadyLabel1=Installer Thoughts pour ce compte Windows. Les notes, images et paramètres existants sont conservés.
es.ReadyLabel1=Instalar Thoughts para esta cuenta de Windows. Se conservan las notas, imágenes y opciones existentes.

[Files]
Source: "{#CoreDir}\thoughts.exe"; DestDir: "{app}\{#LauncherFolder}"; Flags: ignoreversion
Source: "{#ProjectRoot}\LICENSE"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#ProjectRoot}\PRIVACY.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#ProjectRoot}\INSTALL.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#ProjectRoot}\SECURITY.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#ProjectRoot}\THIRD_PARTY_NOTICES.txt"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#ProjectRoot}\licenses\InnoSetup-LICENSE.txt"; DestDir: "{app}\licenses"; Flags: ignoreversion

[Icons]
Name: "{group}\Thoughts"; Filename: "{code:ManagedExecutable}"; Parameters: "open --browser system --json"; WorkingDir: "{app}"; Flags: runminimized; Check: CreateShortcuts
Name: "{autodesktop}\Thoughts"; Filename: "{code:ManagedExecutable}"; Parameters: "open --browser system --json"; WorkingDir: "{app}"; Flags: runminimized; Check: CreateDesktopShortcuts

[Run]
Filename: "{code:ManagedExecutable}"; Parameters: "open --browser system --json"; Flags: postinstall nowait skipifsilent runhidden

[CustomMessages]
en.OfflineSummary=Install offline for this Windows account, without administrator privileges or developer tools.
en.KeepSummary=Desktop and Start menu shortcuts will be created. Existing notes, images and settings are retained.
en.AISummary=AI clients are optional and can be connected separately after installation.
zhHans.OfflineSummary=离线安装到当前 Windows 用户，无需管理员权限或开发工具。
zhHans.KeepSummary=自动创建桌面和开始菜单快捷方式，保留已有笔记、图片和设置。
zhHans.AISummary=不需要 AI 客户端也能使用；需要 AI 功能时，可在安装后单独接入。
zhHant.OfflineSummary=離線安裝至目前 Windows 使用者，無需管理員權限或開發工具。
zhHant.KeepSummary=自動建立桌面和開始功能表捷徑，保留現有筆記、圖片和設定。
zhHant.AISummary=無需 AI 用戶端也能使用；如需 AI 功能，可於安裝後另外連接。
ja.OfflineSummary=管理者権限や開発ツールを使わず、現在の Windows ユーザーにオフラインでインストールします。
ja.KeepSummary=デスクトップとスタートメニューにショートカットを作成し、既存のノート、画像、設定を保持します。
ja.AISummary=AI クライアントは任意です。必要に応じてインストール後に接続できます。
ko.OfflineSummary=관리자 권한이나 개발 도구 없이 현재 Windows 사용자용으로 오프라인 설치합니다.
ko.KeepSummary=바탕 화면과 시작 메뉴에 바로 가기를 만들고 기존 노트, 이미지와 설정을 유지합니다.
ko.AISummary=AI 클라이언트는 선택 사항이며 설치 후 별도로 연결할 수 있습니다.
de.OfflineSummary=Offline-Installation für dieses Windows-Konto, ohne Administratorrechte oder Entwicklungswerkzeuge.
de.KeepSummary=Verknüpfungen werden auf dem Desktop und im Startmenü erstellt. Notizen, Bilder und Einstellungen bleiben erhalten.
de.AISummary=KI-Clients sind optional und können nach der Installation separat verbunden werden.
fr.OfflineSummary=Installation hors ligne pour ce compte Windows, sans droits administrateur ni outils de développement.
fr.KeepSummary=Des raccourcis sont créés sur le Bureau et dans le menu Démarrer. Les notes, images et paramètres sont conservés.
fr.AISummary=Les clients IA sont facultatifs et peuvent être connectés séparément après l’installation.
es.OfflineSummary=Instalación sin conexión para esta cuenta de Windows, sin permisos de administrador ni herramientas de desarrollo.
es.KeepSummary=Se crean accesos directos en el escritorio y en el menú Inicio. Se conservan notas, imágenes y opciones.
es.AISummary=Los clientes de IA son opcionales y se pueden conectar por separado después de la instalación.

; No UninstallDelete entry: the shared home, managed versions and user files
; belong to the user and may still be in use by other MCP clients.
[Code]
var
  CorePrepared: Boolean;
  WindowsAccountProfile: String;

function InitializeSetup: Boolean;
begin
  WindowsAccountProfile := GetShellFolderByCSIDL($0028, False);
  Result := (WindowsAccountProfile <> '') and DirExists(WindowsAccountProfile);
  if not Result then
    MsgBox('Windows could not locate your account profile. No installation was performed.', mbError, MB_OK);
end;

function ShouldSkipPage(PageID: Integer): Boolean;
begin
  // Keep /GROUP available for an explicitly chosen folder, without another UI step.
  Result := PageID = wpSelectProgramGroup;
end;

function UpdateReadyMemo(Space, NewLine, MemoUserInfoInfo, MemoDirInfo,
  MemoTypeInfo, MemoComponentsInfo, MemoGroupInfo, MemoTasksInfo: String): String;
begin
  Result := CustomMessage('OfflineSummary') + NewLine + NewLine +
    CustomMessage('KeepSummary') + NewLine + NewLine + CustomMessage('AISummary');
end;

function CreateShortcuts: Boolean;
begin
  Result := not WizardNoIcons;
end;

function CreateDesktopShortcuts: Boolean;
begin
  Result := CreateShortcuts and
    (CompareText(ExpandConstant('{param:NoDesktop|false}'), 'true') <> 0);
end;

function ManagedExecutable(Param: String): String;
var
  Home: String;
begin
  Home := GetEnv('THOUGHTS_HOME');
  if (Home = '') or (CompareText(Home, 'auto') = 0) or
      (Home = '${HOME}/Thoughts') or (Home = '${HOME}\Thoughts') then
    Home := AddBackslash(WindowsAccountProfile) + 'Thoughts';
  Result := AddBackslash(Home) + 'app\versions\{#AppVersion}-' +
    Copy('{#CoreSHA256}', 1, 12) + '\thoughts.exe';
end;

function RunCore(const ExePath, Arguments: String): String;
var
  ExitCode, I: Integer;
  Captured: TExecOutput;
begin
  Result := '';
  // A stable working directory prevents the service retaining Setup's temp folder.
  if not ExecAndCaptureOutput(ExePath, Arguments, ExpandConstant('{localappdata}'),
      SW_SHOWNORMAL, ewWaitUntilTerminated, ExitCode, Captured) then
  begin
    Result := 'Windows could not run Thoughts. Your current runtime is kept. ' +
      'Do not disable Windows protection. Error: ' + SysErrorMessage(ExitCode);
    Exit;
  end;
  if (ExitCode <> 0) or Captured.Error then
  begin
    Result := 'Thoughts setup did not complete. Your existing library is kept.';
    for I := 0 to GetArrayLength(Captured.StdErr) - 1 do
      Result := Result + #13#10 + Captured.StdErr[I];
    // Only init/runtime install call this helper; errors on stdout are useful
    // here, but successful output and browser-opening commands are never logged.
    for I := 0 to GetArrayLength(Captured.StdOut) - 1 do
      Result := Result + #13#10 + Copy(Captured.StdOut[I], 1, 2048);
  end;
end;

function PrepareToInstall(var NeedsRestart: Boolean): String;
var
  Core: String;
  Zone: AnsiString;
begin
  Result := '';
  if CorePrepared then Exit;
  if DirExists(ExpandConstant('{app}\Inbox')) or
      DirExists(ExpandConstant('{app}\Projects')) or
      DirExists(ExpandConstant('{app}\library')) or
      FileExists(ExpandConstant('{app}\config\config.json')) then
  begin
    Result := 'Choose a program directory outside your existing library and Thoughts home.';
    Exit;
  end;
  ExtractTemporaryFile('thoughts.exe');
  Core := ExpandConstant('{tmp}\thoughts.exe');
  if CompareText(GetSHA256OfFile(Core), '{#CoreSHA256}') <> 0 then
  begin
    Result := 'The packaged executable checksum does not match. Installation stopped.';
    Exit;
  end;
  // Preserve readable download provenance; never strip Mark of the Web.
  if LoadStringFromFile(ExpandConstant('{srcexe}') + ':Zone.Identifier', Zone) then
    if not SaveStringToFile(Core + ':Zone.Identifier', Zone, False) then
    begin
      Result := 'Could not preserve download provenance. Installation stopped.';
      Exit;
    end;
  Result := RunCore(Core, 'init --json');
  if Result <> '' then Exit;
  // The runtime command validates identity/hash, backs up the previous pin,
  // refuses busy editors, health-checks the new core and restores on failure.
  Result := RunCore(Core, 'runtime install --source "' + Core +
    '" --sha256 {#CoreSHA256} --json');
  if Result <> '' then Exit;
  CorePrepared := True;
end;

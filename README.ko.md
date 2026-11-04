[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Deutsch](README.de.md) | [Français](README.fr.md) | [Español](README.es.md)

# Thoughts

[설치](INSTALL.md) · [문서](docs/README.md) · [후원](DONATE.md) · [코드 서명 정책 / Code signing policy](docs/CODE_SIGNING_POLICY.md)

Windows x64용 무료 오픈 소스 로컬 우선 노트 도구입니다. AI 어시스턴트로 Markdown 파일을 관리하고, 독립적인 로컬 웹페이지에서 읽고 편집할 수 있습니다.

**미서명 미리보기 — 2026-11-04:** `0.3.0-preview.21`은 소스, 실행 파일, 클라이언트 설정의 이름을 Thoughts로 통일합니다. 관리자가 [미서명 실험용 사전 릴리스](https://github.com/xuli82221-pixel/hello-Thoughts-/releases/tag/v0.3.0-preview.21) 공개를 승인했으며, 공개 Windows 빌드와 검증을 통과한 파일을 제공합니다. **신뢰할 수 있는 게시자 서명이나 안정 버전 검증을 뜻하지 않습니다.** SignPath Foundation 승인은 아직 받지 않았습니다. Windows SAC/SmartScreen이 차단하면 실행을 멈추고 보호를 유지하세요. [배포 준비](docs/RELEASE_READINESS_20261104.md), [설치 안내](INSTALL.md), [이전 네이티브 Windows 테스트](docs/WINDOWS_NATIVE_TESTS_20261104.md)를 확인하세요.

## 노트는 내 컴퓨터에

Thoughts에는 계정, 클라우드 데이터베이스, 노트 업로드 서비스, 원격 측정, 내장 AI, 유료 테마, 구독, 결제 연동 또는 브라우저 확장 기능이 없습니다. Markdown과 로컬 이미지는 사용자 소유이며 다른 편집기로도 열 수 있습니다.

Go 실행 파일에 React 웹페이지가 내장되어 있습니다. **미리 빌드된 배포 파일은 실행 시 Node.js, Go, Python 또는 Docker가 필요하지 않습니다.** 선택한 AI 클라이언트에는 해당 클라이언트의 설치 및 계정 요건이 별도로 적용됩니다.

## 다운로드 및 설치

이 실험용 미리보기는 [GitHub Releases](https://github.com/xuli82221-pixel/hello-Thoughts-/releases)에서 **`thoughts-<version>-windows-amd64-setup.exe`**, 포터블 ZIP, Claude Desktop MCPB를 제공합니다. 내부 프로그램 이름은 **`thoughts.exe`**입니다. 사용 전 SHA256을 확인하세요. 현재 파일에는 **Authenticode 서명이 없으며**, 해시 일치가 게시자 신원을 증명하지는 않습니다. Windows가 차단하면 실행을 멈추세요. 설치는 사용자별로 이루어지며 관리자 권한 없이 공유 홈과 바탕 화면·시작 메뉴 바로가기를 만듭니다. 기존 자료는 보존하고 제거 후에도 다른 클라이언트를 위해 공유 홈을 남깁니다. 서명된 안정 버전은 아직 제공하지 않습니다.

Claude Desktop 사용자는 해당 릴리스의 `.mcpb` 확장을 가져오는 방법도 사용할 수 있습니다. Codex와 Claude Code 사용자는 ZIP의 실행 파일을 로컬 stdio MCP 서버로 등록합니다. 웹페이지는 AI 클라이언트 없이도 실행할 수 있습니다. 명령은 [일반 사용자 빠른 시작](INSTALL.md#ordinary-windows-users)을 확인하세요. GitHub의 **Code → Download ZIP**은 소스 코드를 받는 기능이며, 미리 빌드된 앱을 받는 기능이 아닙니다.

AI 클라이언트는 동일한 로컬 표준 stdio MCP 어댑터를 시작하고 인증된 로컬 코어를 호출합니다. 웹페이지는 `127.0.0.1`에서만 수신합니다. 온라인 AI 클라이언트에 노트 읽기를 요청하면 내용이 해당 클라이언트로 전달되어 그 데이터 정책을 따릅니다. 로컬 우선이라는 말이 사용자가 요청한 AI 도구 결과까지 컴퓨터 밖으로 나가지 않는다는 뜻은 아닙니다. Thoughts는 파일을 암호화하지 않습니다.

## 대상 클라이언트

| 클라이언트 | 설치 방식 |
| --- | --- |
| Claude Desktop | 공식 `.mcpb` 확장 형식 |
| Windows용 Claude Code | 로컬 stdio MCP, Windows 포터블 ZIP |
| Windows용 OpenAI Codex | 로컬 stdio MCP, 동일한 포터블 ZIP |

대상으로 삼는 연동이며, 실제 검증 상태는 배포 검토 문서에 기록되어 있습니다. 일반 ChatGPT 웹 플러그인, macOS, Linux, 모바일, 원격 MCP 및 클라우드 동기화는 이번 단계의 지원 범위에 포함되지 않습니다.

## 기능

- 로컬 Markdown 생성, 읽기, 편집, 검색, 이름 변경 및 이동.
- 간결한 프로젝트 트리, 빠른 노트, 프로젝트 보관 및 보관 해제, 복원 가능한 노트·프로젝트 휴지통.
- 기존 노트는 렌더링된 미리보기로 열리며, 편집은 직접 전환.
- 로컬 PNG/JPEG/WebP 이미지, 미리보기 및 확대 보기.
- 무료 테마와 검증을 거치는 로컬 테마 가져오기.
- 영어, 한국어, 일본어, 중국어 간체·번체, 독일어, 프랑스어, 스페인어 UI 및 각 언어의 로컬 개인정보·사용 안내.
- 외부 파일 변경 알림과 SHA256 충돌 보호.
- 모든 클라이언트가 동일하게 설정된 라이브러리를 사용하며, 재초기화해도 기존 선택 유지.
- 체크섬으로 검증하는 로컬 코어 업그레이드 및 소프트웨어 롤백과 미저장 편집 보호.
- 서비스 포트가 바뀌어도 복원되는 서버 저장 노트 탭 및 읽기 위치.

충돌 보호는 오래된 버전을 바탕으로 한 쓰기를 거부합니다. 자동 수정 이력은 구현되지 않았으므로 라이브러리를 별도로 백업하세요. 충돌이나 연결 끊김 중 열린 탭에 초안이 남을 수 있지만, 탭을 닫거나 Windows를 재시작한 뒤 복구되는 것은 보장하지 않습니다.

![가상 예제 노트를 표시한 Thoughts 0.3.0-preview.20](docs/assets/preview-0.3.0-preview.20.jpg)

2026-11-04에 새로 만든 독립적인 임시 라이브러리와 가상 노트로 촬영한 실제 `0.3.0-preview.20` UI입니다. 개인 노트, 실제 이미지, 컴퓨터의 로컬 경로 또는 자격 증명이 표시되지 않습니다. [이전 레이아웃 스크린샷](docs/assets/preview.png)은 기록으로 보존하며 현재 버전으로 제시하지 않습니다. 스크린샷 확인은 설치 프로그램이나 클라이언트 검증을 대신하지 않습니다.

## 설치

SHA256 확인, 실행 전 점검할 명령 및 세 클라이언트의 절차는 [INSTALL.md](INSTALL.md)를 확인하세요. preview.8부터 새 설치는 순서에 관계없이 모든 클라이언트에서 같은 Windows 사용자 Thoughts 디렉터리를 기본으로 사용합니다. 고급 홈 설정은 `auto`로 유지하세요. 의도적으로 홈을 바꿀 때만 절대 경로 `THOUGHTS_HOME`을 지정합니다.

```text
<Thoughts home>/
  library/     Markdown, images and recoverable Trash
  config/      Configuration, local sessions, preferences and logs
  app/         Managed executable versions (optional)
  backups/     Previous runtime pins; not a backup of note bodies
```

일반 설치에는 공유 기본 홈을 사용하세요. 명시적으로 지정한 기존 설정도 계속 지원하지만, 이 버전은 기존 라이브러리를 자동으로 이전하거나 합치지 않습니다. 홈을 직접 바꾸면 공유하려는 클라이언트 모두에서 동일한 물리 폴더를 사용하세요.

설치와 업데이트는 홈을 보존해야 합니다. 읽기 상태에는 경로와 스크롤 위치가 저장되며 미저장 본문은 저장되지 않습니다. 코어 업데이트에는 검토된 로컬 실행 파일과 SHA256을 명시적으로 제공해야 하며, 백그라운드 네트워크 업데이트 기능은 없습니다. 클라이언트 확장과 공유 코어는 별도로 업데이트합니다. [업그레이드 안내](INSTALL.md#managed-core-updates-and-rollback)를 확인하세요.

미서명 Windows 프로그램이나 스크립트는 SmartScreen, Smart App Control 또는 조직의 앱 제어 정책에 의해 차단될 수 있습니다. SHA256은 체크섬과의 무결성을 확인하며 게시자 신원을 증명하지 않습니다. 미리보기를 실행하려고 보호를 끄거나 경고를 우회하지 마세요. 소스 공개만으로 서명 문제가 해결되지는 않습니다.

## 작업 공간

[작업 공간 구성](docs/WORKSPACE_LAYOUT.md)과 [문서 목록](docs/README.md)을 확인하세요. 이전 요구사항과 보고서는 로컬 보관소에 보존합니다. 소스와 실행 파일 이름은 Thoughts로 통일하며 기존 데이터 식별자는 호환성을 위해 유지합니다.

## 빌드

Windows x64, go.mod에 맞는 Go, Node.js 24, Python 3.12+, 디지털 서명된 공식 [Inno Setup 7 컴파일러](https://jrsoftware.org/isdl.php)가 필요합니다. 최종 사용자에게는 이 도구들이 필요하지 않습니다. [CONTRIBUTING.md](CONTRIBUTING.md)를 확인하세요.

```powershell
Set-Location web
npm ci
Set-Location ..
pwsh -File scripts/build-desktop-adapters.ps1 -Python 'C:\Tools\Python\python.exe' -Go 'C:\Tools\Go\bin\go.exe' -InnoCompiler 'C:\Tools\InnoSetup7\ISCC.exe'
```

도구 경로를 실제 설치 위치로 바꾸세요. 빌드는 프런트엔드와 Go 코어를 테스트·빌드하고, 고지문과 Setup.exe/ZIP/MCPB 및 체크섬을 만들고, 새 출력 디렉터리에서 패키지를 검증합니다. 기존 출력을 덮어쓰거나 자동 업로드하지 않습니다. 공개 후보 모드에는 신뢰할 수 있는 서명과 실제 게시자 메타데이터가 필요합니다. 서명은 깨끗한 Windows 환경이나 클라이언트 검증을 대신하지 않습니다.

## 보안 및 라이선스

[SECURITY.md](SECURITY.md), [PRIVACY.md](PRIVACY.md), [CHANGELOG.md](CHANGELOG.md), [제3자 고지](THIRD_PARTY_NOTICES.txt), [의존성 검토](docs/DEPENDENCIES.md)를 확인하세요. 프로젝트 코드는 MIT 라이선스이며, OFL 글꼴을 포함한 번들 의존성에는 각자의 라이선스가 적용됩니다.

실제 노트, 개인 스크린샷, 자격 증명, 세션 파일 또는 로그를 커밋하지 마세요. 취약점은 hello@spxipo.com으로 비공개 보고하세요. 설정에는 소유자의 Base 메인넷 주소를 이용한 선택적 후원 안내가 있습니다. QR에는 주소만 포함되므로 지갑에서 Base를 별도로 선택하세요. Thoughts는 지갑 연결, 개인 키 요청, 자산·금액 지정 또는 결제 처리를 하지 않습니다. 모든 기능은 무료입니다.

## 코드 서명 정책

**승인 대기 — SignPath Foundation 신청은 아직 승인되지 않았습니다.** 현재 Windows 프로그램은 모두 미서명 상태입니다. 승인되면 무료 코드 서명은 SignPath.io가, 인증서는 SignPath Foundation이 제공합니다. 범위, 실제 담당자 및 미완료 요건은 [코드 서명 정책 / Code signing policy](docs/CODE_SIGNING_POLICY.md)과 [신청 준비](docs/SIGNPATH_APPLICATION.md)에 기록되어 있습니다. 이번 준비 작업은 신청서를 제출하지 않았습니다.

## 후원

후원은 자발적이며 기능을 해제하지 않습니다. 별도의 [후원 페이지](DONATE.md)에 소유자의 기존 Base 주소, 주소만 포함된 QR 및 개인정보 처리 경계가 있습니다. 문제 보고, 번역 또는 문서 개선으로도 도울 수 있습니다. Thoughts는 지갑에 연결하거나 결제를 처리하지 않습니다.

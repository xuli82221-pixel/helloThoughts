# Thoughts · Privacy and use / 隐私与使用说明

2026-11-03

## English

Updated: 2026-11-03

### Privacy notice

Thoughts is a free, open-source, local-first Markdown tool for Windows. You do not need a Thoughts account. This notice describes Thoughts itself and the separate services you may choose to use.

#### Files and local settings

Notes and images are files in your authorized local library. Clients configured to the same library read and modify those same files. Configuration stores the library path, language, themes and other interface preferences.

Local reading state contains note paths, open tabs and scroll positions, not draft bodies. Browser storage may cache interface settings. Runtime files contain service identity and authorization information. Diagnostic logs are local and are designed not to record note bodies or session credentials.

#### What Thoughts does not collect

Thoughts does not operate a cloud note database, account service, analytics, advertising, telemetry or automatic note-upload service. It does not automatically collect conversations from your AI clients. It has no built-in AI model or subscription service.

#### AI clients are a separate data boundary

When you authorize Codex, Claude Desktop, Claude Code or another client to read a note, returned text and tool results are passed to that client. An online client may send them to its AI provider under its own policies and your account settings.

Local storage does not mean content you send to an online AI stays offline. Review the client's permissions and data controls before sharing sensitive material. Thoughts cannot delete copies retained by another service.

#### Local sessions and external links

The reader service listens on 127.0.0.1. Local authorization uses a one-use opening ticket, a session cookie and request checks. These authorize local access, not a Thoughts account.

The reader does not load remote note images, web fonts or analytics scripts. External links open when you choose them; the destination can receive normal browser network information. Thoughts does not append note bodies, library paths or local credentials to support links.

#### Retention, deletion and backups

Your files remain until you change or remove them. Moving a note or project to local Trash is recoverable and is not permanent erasure. Shared images, old copies and backups may remain. Uninstalling an integration does not delete your library.

You can open, copy, export and back up the Markdown and image files yourself. To remove all copies, review Trash, images and any backups or synchronized folders you use. A remote service or operating-system backup follows its own retention rules.

#### Security limits

Thoughts does not encrypt your files. Local session checks do not prevent other programs or users with filesystem permissions from reading or changing them. Use operating-system access controls and suitable backups.

Revision checks protect supported writes against stale updates; they are not a guarantee against arbitrary external edits, device failure or data loss. Do not share local authorization tickets or diagnostic files containing credentials.

#### Information you choose to send for support

Publisher: xuli82221-pixel. Support email: hello@spxipo.com. Thoughts does not automatically send support email or diagnostic reports.

If you contact support, you choose what to disclose. Email contents are handled by the email providers involved. Omit private notes and credentials unless necessary, and contact the support address with questions or requests about information you supplied.

#### Changes to this notice

This notice is bundled with the application and carries an update date. Review the notice supplied with a new version. A local notice does not imply that a public policy page, marketplace listing or publisher verification has been completed.

### Use, license and support

#### Free software and license

Thoughts is free and open source under the MIT License. The original license below governs the software; this translated description is a summary. Third-party components retain their own licenses, listed in the bundled notices.

MIT permits use, modification and redistribution subject to retaining its required notices, and supplies the software without warranty as stated in the license. Keep backups and verify important changes.

#### Your files and client permissions

Choose a stable library location and configure clients to the same library when you want shared content. Updates and integration installation are separate from your Markdown files. Give AI clients only the file access and actions you intend to authorize.

#### Optional support

Donations are voluntary, unlock no features and are not a subscription or purchase. Thoughts does not require wallet connection or private keys and does not submit transactions for you.

If configured, the local support page shows a Base address and an address-only QR code. Sending funds is a separate action in your wallet. Wallet services and public blockchain activity have their own data exposure and policies; Thoughts does not automatically collect your transaction details.

## 한국어

최종 수정: 2026-11-03

### 개인정보 처리 안내

Thoughts는 Windows용 무료 오픈 소스 로컬 우선 Markdown 도구이며 Thoughts 계정이 필요하지 않습니다. 이 안내는 Thoughts 자체와 사용자가 선택하는 별도 서비스의 경계를 설명합니다.

#### 파일과 로컬 설정

노트와 이미지는 사용자가 허용한 로컬 라이브러리에 파일로 저장됩니다. 같은 라이브러리로 설정한 클라이언트는 동일한 파일을 읽고 수정합니다. 설정에는 라이브러리 경로, 언어, 테마 및 화면 환경설정이 포함됩니다.

로컬 읽기 상태에는 노트 경로, 열린 탭, 스크롤 위치가 포함되며 초안 본문은 포함되지 않습니다. 브라우저는 화면 설정을 캐시할 수 있습니다. 실행 정보에는 서비스 식별 및 인증 정보가 저장됩니다. 진단 로그는 로컬에 보관되며 노트 본문이나 세션 자격 증명을 기록하지 않도록 설계되었습니다.

#### Thoughts가 수집하지 않는 정보

Thoughts는 클라우드 노트 데이터베이스, 계정 서비스, 분석·광고·원격 측정 또는 자동 노트 업로드 서비스를 운영하지 않습니다. AI 클라이언트의 대화를 자동 수집하지 않으며 내장 AI 모델이나 구독 서비스가 없습니다.

#### AI 클라이언트의 별도 데이터 처리

Codex, Claude Desktop, Claude Code 또는 다른 클라이언트에 노트 읽기를 허용하면 본문과 도구 결과가 해당 클라이언트로 전달됩니다. 온라인 클라이언트는 자체 정책 및 계정 설정에 따라 이를 AI 제공업체에 보낼 수 있습니다.

로컬 저장은 온라인 AI에 보낸 내용까지 오프라인에 머문다는 뜻이 아닙니다. 민감한 정보를 공유하기 전에 클라이언트의 권한과 데이터 설정을 확인하세요. Thoughts는 다른 서비스에 남은 사본을 삭제할 수 없습니다.

#### 로컬 세션과 외부 링크

리더 서비스는 127.0.0.1에서 수신합니다. 일회용 열기 티켓, 세션 쿠키 및 요청 검사를 사용해 로컬 접근을 인증하며 Thoughts 계정 인증은 아닙니다.

리더는 노트의 원격 이미지, 외부 웹 폰트 또는 분석 스크립트를 불러오지 않습니다. 외부 링크는 사용자가 선택할 때 열리며 방문 사이트는 일반적인 브라우저 네트워크 정보를 받을 수 있습니다. 지원 링크에 노트 본문, 라이브러리 경로 또는 로컬 자격 증명을 자동 첨부하지 않습니다.

#### 보관·삭제·백업

파일은 사용자가 변경하거나 제거할 때까지 유지됩니다. 노트나 프로젝트를 로컬 휴지통으로 이동하는 것은 복원 가능한 작업이며 영구 삭제가 아닙니다. 공유 이미지, 이전 사본 및 백업이 남을 수 있습니다. 연동 기능을 제거해도 라이브러리는 삭제되지 않습니다.

Markdown과 이미지 파일을 직접 열고 복사·내보내기·백업할 수 있습니다. 모든 사본을 제거하려면 휴지통, 이미지, 백업 및 동기화 폴더를 확인하세요. 원격 서비스와 운영체제 백업에는 각자의 보관 규칙이 적용됩니다.

#### 보안의 한계

Thoughts는 파일을 암호화하지 않습니다. 로컬 세션 검사는 파일 접근 권한이 있는 다른 프로그램이나 사용자의 읽기·변경을 막지 못합니다. 운영체제 접근 제어와 적절한 백업을 사용하세요.

버전 검사는 지원되는 쓰기 작업에서 오래된 버전의 덮어쓰기를 방지하지만 외부 편집, 장치 고장 또는 데이터 손실을 완전히 막는 보장은 아닙니다. 로컬 인증 티켓이나 자격 증명이 담긴 진단 파일을 공유하지 마세요.

#### 지원을 위해 자발적으로 보내는 정보

게시자: xuli82221-pixel. 지원 이메일: hello@spxipo.com. Thoughts는 지원 이메일이나 진단 보고서를 자동 전송하지 않습니다.

지원 문의 시 공개할 정보를 직접 선택합니다. 이메일 내용은 관련 이메일 제공업체에서 처리합니다. 필요한 경우가 아니라면 개인 노트와 자격 증명을 제외하고, 제공한 정보에 관한 질문이나 요청은 지원 주소로 보내세요.

#### 안내 변경

이 안내는 앱에 포함되며 수정 날짜를 표시합니다. 새 버전의 안내를 확인하세요. 로컬 안내가 제공된다는 사실만으로 공개 정책 페이지, 마켓 등록 또는 게시자 확인이 완료되었다는 뜻은 아닙니다.

### 사용·라이선스·지원 안내

#### 무료 소프트웨어와 라이선스

Thoughts는 MIT License로 제공되는 무료 오픈 소스입니다. 아래 영어 원문 라이선스가 소프트웨어에 적용되며 번역된 설명은 요약입니다. 타사 구성요소에는 각자의 라이선스가 적용되며 포함된 고지에서 확인할 수 있습니다.

MIT는 필수 고지를 유지하는 조건으로 사용·수정·재배포를 허용하며 라이선스에 명시된 대로 보증 없이 제공됩니다. 백업을 유지하고 중요한 변경을 확인하세요.

#### 파일과 클라이언트 권한

안정적인 라이브러리 경로를 선택하고 내용을 공유하려면 클라이언트를 같은 라이브러리로 설정하세요. 프로그램 업데이트 및 연동 설치는 Markdown 파일과 별개입니다. AI 클라이언트에는 의도한 파일 접근과 작업만 허용하세요.

#### 선택적 후원

후원은 자발적이며 기능을 해제하지 않고 구독이나 구매도 아닙니다. Thoughts는 지갑 연결이나 개인 키를 요구하지 않으며 거래를 대신 제출하지 않습니다.

설정된 경우 로컬 지원 페이지에 Base 주소와 주소만 담긴 QR 코드가 표시됩니다. 송금은 사용자의 지갑에서 별도로 수행합니다. 지갑 서비스와 공개 블록체인 활동에는 별도 정보 노출 및 정책이 적용되며 Thoughts는 거래 내역을 자동 수집하지 않습니다.

## 日本語

更新日: 2026-11-03

### プライバシーについて

Thoughts は Windows 向けの無料・オープンソースのローカル優先 Markdown ツールです。Thoughts アカウントは不要です。この文書は Thoughts 自体と、利用者が選ぶ外部サービスとの境界を説明します。

#### ファイルとローカル設定

ノートと画像は、許可したローカルライブラリ内のファイルとして保存されます。同じライブラリを設定したクライアントは同じファイルを読み書きします。設定にはライブラリの場所、言語、テーマ、画面の環境設定が含まれます。

ローカルの閲覧状態にはノートのパス、開いたタブ、スクロール位置が含まれ、下書き本文は含まれません。ブラウザーは画面設定をキャッシュする場合があります。実行情報にはサービス識別情報と認証情報が保存されます。診断ログはローカルに保存され、ノート本文やセッションの認証情報を記録しない設計です。

#### Thoughts が収集しない情報

Thoughts はクラウドのノートデータベース、アカウントサービス、分析、広告、テレメトリ、自動ノートアップロードを運営しません。AI クライアントの会話を自動収集せず、内蔵 AI モデルやサブスクリプションもありません。

#### AI クライアントによる別のデータ処理

Codex、Claude Desktop、Claude Code などにノートの読み取りを許可すると、本文とツール結果がそのクライアントへ渡されます。オンラインのクライアントは、自身のポリシーとアカウント設定に従って AI 提供元へ送信する場合があります。

ローカル保存は、オンライン AI に渡した内容までオフラインに保たれることを意味しません。機密情報を共有する前にクライアントの権限とデータ設定を確認してください。Thoughts は他サービスに保持されたコピーを削除できません。

#### ローカルセッションと外部リンク

リーダーサービスは 127.0.0.1 で待ち受けます。一回限りの起動チケット、セッション Cookie、リクエスト検査でローカルアクセスを認証します。Thoughts アカウントの認証ではありません。

リーダーはノートのリモート画像、外部 Web フォント、分析スクリプトを読み込みません。外部リンクは選択したときに開き、リンク先には通常のブラウザーネットワーク情報が伝わる場合があります。サポートリンクに本文、ライブラリのパス、ローカル認証情報を自動付加しません。

#### 保存・削除・バックアップ

ファイルは利用者が変更・削除するまで残ります。ノートやプロジェクトのローカルごみ箱への移動は復元可能で、完全消去ではありません。共有画像、古いコピー、バックアップが残る場合があります。連携機能をアンインストールしてもライブラリは削除されません。

Markdown と画像は自分で開き、コピー、書き出し、バックアップできます。すべてのコピーを消す場合は、ごみ箱、画像、バックアップ、同期フォルダーも確認してください。外部サービスや OS のバックアップにはそれぞれの保存ルールが適用されます。

#### セキュリティの限界

Thoughts はファイルを暗号化しません。ローカルセッションの検査は、ファイル権限を持つ他のプログラムや利用者による読み取り・変更を防ぐものではありません。OS のアクセス制御と適切なバックアップを利用してください。

バージョン検査は対応する書き込みで古い内容の上書きを防ぎますが、外部編集、機器故障、データ損失を完全に防ぐ保証ではありません。ローカル認証チケットや認証情報を含む診断ファイルを共有しないでください。

#### サポートへ自分で送る情報

公開者: xuli82221-pixel。サポートメール: hello@spxipo.com。Thoughts はサポートメールや診断報告を自動送信しません。

問い合わせで開示する情報は自分で選べます。メール本文は関係するメール提供元によって処理されます。必要がなければ個人のノートや認証情報を含めず、提供した情報に関する質問や依頼はサポートアドレスへ連絡してください。

#### この文書の変更

この文書はアプリに同梱され、更新日を表示します。新しいバージョンに含まれる文書を確認してください。ローカル文書の提供は、公開ポリシーページ、マーケット掲載、公開者確認の完了を意味しません。

### 利用・ライセンス・サポート

#### 無料ソフトウェアとライセンス

Thoughts は MIT License の無料オープンソースです。ソフトウェアには以下の英語原文ライセンスが適用され、翻訳された説明は要約です。第三者の構成要素には各ライセンスが適用され、同梱の通知で確認できます。

MIT は必要な通知の保持を条件に使用・変更・再配布を認め、ライセンスに記載のとおり無保証で提供します。バックアップを取り、重要な変更を確認してください。

#### ファイルとクライアントの権限

安定したライブラリの場所を選び、内容を共有したいクライアントには同じライブラリを設定してください。更新や連携のインストールは Markdown ファイルとは別です。AI クライアントには意図したファイルアクセスと操作だけを許可してください。

#### 任意の支援

寄付は任意で、機能の解除、購入、サブスクリプションではありません。Thoughts はウォレット接続や秘密鍵を要求せず、取引を代行送信しません。

設定済みの場合、ローカルの支援ページに Base アドレスとアドレスのみの QR コードを表示します。送金はウォレットで別途行う操作です。ウォレットサービスと公開ブロックチェーンの活動には独自の情報公開やポリシーがあり、Thoughts は取引情報を自動収集しません。

## 简体中文

更新日期: 2026-11-03

### 隐私说明

Thoughts 是面向 Windows 的免费、开源、本地优先 Markdown 工具，不需要 Thoughts 账号。本说明区分 Thoughts 自身与用户主动选择的其他服务。

#### 文件与本地设置

笔记和图片以真实文件保存在你授权的本地资料库。配置到同一资料库的客户端读写同一份文件。配置保存资料库路径、语言、主题和界面偏好。

本地阅读状态包含笔记路径、已打开的标签和滚动位置，不包含草稿正文；浏览器可能缓存界面设置。运行文件包含服务身份和授权信息。诊断日志保存在本机，设计上不记录笔记正文或会话凭据。

#### Thoughts 不收集什么

Thoughts 不运营云端笔记数据库、账号服务、分析、广告、遥测或自动笔记上传服务，不自动抓取 AI 客户端中的聊天，没有内置 AI 模型和订阅服务。

#### AI 客户端是另一条数据边界

当你授权 Codex、Claude Desktop、Claude Code 或其他客户端读取笔记时，正文和工具结果会交给该客户端。在线客户端可能依照自身政策及你的账号设置，把这些内容发送给其 AI 服务商。

资料保存在本地，不等于主动交给在线 AI 的内容也只在本地处理。分享敏感资料前请核对客户端权限与数据设置。Thoughts 无法删除其他服务保留的副本。

#### 本地会话与外部链接

阅读器服务监听 127.0.0.1，通过一次性打开凭证、会话 Cookie 和请求校验授权本地访问。这不是 Thoughts 账号。

阅读器不加载笔记中的远程图片、外部网络字体或分析脚本。外部链接在你主动选择时打开，目标网站可能收到正常浏览器访问产生的网络信息。支持链接不自动附带笔记正文、资料库路径或本地凭据。

#### 保留、删除与备份

文件会保留到你修改或移除为止。把笔记或项目移入本地回收站是可恢复操作，不是永久擦除；共享图片、旧副本和备份可能仍在。卸载集成不会删除资料库。

你可以直接打开、复制、导出和备份 Markdown 与图片。需要清除所有副本时，还应检查回收站、图片以及你使用的备份和同步文件夹。远程服务或操作系统备份按各自规则保留资料。

#### 安全边界

Thoughts 不对文件加密。本地会话校验无法阻止其他具有文件权限的程序或用户读写这些文件。请使用操作系统访问控制并保留合适的备份。

版本校验会在支持的写入流程中拒绝过期覆盖，但不保证阻止任意外部编辑、设备故障或数据丢失。不要分享本地授权凭证或包含凭据的诊断文件。

#### 你主动发送给支持方的信息

发布者：xuli82221-pixel。支持邮箱：hello@spxipo.com。Thoughts 不自动发送支持邮件或诊断报告。

联系支持时，由你决定提供哪些信息；邮件内容由相关邮件服务商处理。非必要时不要附带私人笔记或凭据。对你提供的信息有疑问或请求，可联系支持邮箱。

#### 本说明的变化

本说明随程序本地提供并标注更新日期，使用新版时请查看其附带说明。本地提供说明不代表公开隐私网页、插件市场上架或发布者验证已经完成。

### 使用、许可与支持说明

#### 免费软件与许可

Thoughts 根据 MIT License 免费开源。软件适用下方英文许可原文，本地化说明是摘要。第三方组件保留各自许可，可在随包的第三方声明中查看。

MIT 允许在保留必要版权和许可声明的条件下使用、修改与再分发，并按许可原文不提供保证。请保留备份并核对重要修改。

#### 文件与客户端权限

请选择稳定的资料库位置；希望共享内容时，把各客户端配置到同一资料库。程序升级及集成安装与 Markdown 文件分开管理。只授予 AI 客户端你确实希望授权的文件访问和操作。

#### 自愿支持

打赏完全自愿，不解锁功能，不是订阅或购买。Thoughts 不要求连接钱包或提供私钥，也不代你提交交易。

配置后，本地支持页会显示 Base 地址和只包含地址的二维码。转账是你在钱包中的独立操作。钱包服务和公开链上活动有各自的信息公开范围与政策；Thoughts 不自动收集你的交易信息。

## 繁體中文

更新日期: 2026-11-03

### 隱私說明

Thoughts 是適用於 Windows 的免費、開源、本機優先 Markdown 工具，不需要 Thoughts 帳號。本說明區分 Thoughts 本身與使用者主動選擇的其他服務。

#### 檔案與本機設定

筆記和圖片以實際檔案儲存在你授權的本機資料庫。設定為同一資料庫的用戶端會讀寫同一份檔案。設定保存資料庫路徑、語言、主題與介面偏好。

本機閱讀狀態包含筆記路徑、開啟的分頁與捲動位置，不包含草稿內文；瀏覽器可能快取介面設定。執行檔案包含服務身分與授權資訊。診斷記錄保存在本機，設計上不記錄筆記內文或工作階段憑證。

#### Thoughts 不蒐集什麼

Thoughts 不營運雲端筆記資料庫、帳號服務、分析、廣告、遙測或自動筆記上傳服務，不自動擷取 AI 用戶端中的對話，沒有內建 AI 模型與訂閱服務。

#### AI 用戶端是另一個資料處理邊界

當你授權 Codex、Claude Desktop、Claude Code 或其他用戶端讀取筆記時，內文與工具結果會交給該用戶端。線上用戶端可能依自身政策及你的帳號設定，將這些內容傳送給其 AI 服務商。

資料儲存在本機，不表示主動交給線上 AI 的內容也僅在本機處理。分享敏感資料前請確認用戶端權限與資料設定。Thoughts 無法刪除其他服務保留的副本。

#### 本機工作階段與外部連結

閱讀器服務監聽 127.0.0.1，透過一次性開啟憑證、工作階段 Cookie 與請求檢查授權本機存取。這不是 Thoughts 帳號。

閱讀器不載入筆記中的遠端圖片、外部網路字型或分析指令碼。外部連結在你主動選擇時開啟，目標網站可能收到一般瀏覽器存取產生的網路資訊。支援連結不自動附帶筆記內文、資料庫路徑或本機憑證。

#### 保留、刪除與備份

檔案會保留到你修改或移除為止。將筆記或專案移至本機資源回收筒是可還原操作，不是永久抹除；共用圖片、舊副本與備份可能仍在。解除安裝整合功能不會刪除資料庫。

你可以直接開啟、複製、匯出及備份 Markdown 與圖片。需要清除所有副本時，也應檢查資源回收筒、圖片與你使用的備份及同步資料夾。遠端服務或作業系統備份依各自規則保留資料。

#### 安全邊界

Thoughts 不加密檔案。本機工作階段檢查無法阻止其他具有檔案權限的程式或使用者讀寫這些檔案。請使用作業系統存取控制並保留適當備份。

版本檢查會在支援的寫入流程中拒絕過期覆寫，但不保證能阻止任意外部編輯、裝置故障或資料遺失。不要分享本機授權憑證或包含憑證的診斷檔案。

#### 你主動傳送給支援方的資訊

發布者：xuli82221-pixel。支援信箱：hello@spxipo.com。Thoughts 不自動傳送支援郵件或診斷報告。

聯絡支援時，由你決定提供哪些資訊；郵件內容由相關電子郵件服務商處理。非必要時不要附帶私人筆記或憑證。對你提供的資訊有疑問或請求，可聯絡支援信箱。

#### 本說明的變更

本說明隨程式在本機提供並標示更新日期，使用新版時請查看其附帶說明。本機提供說明不表示公開隱私網頁、外掛市集上架或發布者驗證已完成。

### 使用、授權與支援說明

#### 免費軟體與授權

Thoughts 依 MIT License 免費開源。軟體適用下方英文授權原文，本地化說明是摘要。第三方元件保留各自授權，可在隨附的第三方聲明中查看。

MIT 允許在保留必要版權與授權聲明的條件下使用、修改及再散布，並依授權原文不提供保證。請保留備份並確認重要修改。

#### 檔案與用戶端權限

請選擇穩定的資料庫位置；希望共用內容時，將各用戶端設定為同一資料庫。程式升級與整合功能安裝和 Markdown 檔案分開管理。僅授予 AI 用戶端你確實希望授權的檔案存取與操作。

#### 自願支持

打賞完全自願，不解鎖功能，不是訂閱或購買。Thoughts 不要求連接錢包或提供私鑰，也不代你提交交易。

設定後，本機支援頁會顯示 Base 地址與只包含地址的 QR 碼。轉帳是你在錢包中的獨立操作。錢包服務與公開鏈上活動有各自的資訊公開範圍與政策；Thoughts 不自動蒐集你的交易資訊。

## Deutsch

Aktualisiert: 2026-11-03

### Datenschutzhinweise

Thoughts ist ein kostenloses, quelloffenes Markdown-Werkzeug für Windows, das vorrangig lokal arbeitet. Ein Thoughts-Konto ist nicht erforderlich. Diese Hinweise erläutern Thoughts selbst und die davon getrennten Dienste, die Sie gegebenenfalls nutzen möchten.

#### Dateien und lokale Einstellungen

Notizen und Bilder liegen als Dateien in Ihrer autorisierten lokalen Bibliothek. Clients, die für dieselbe Bibliothek eingerichtet sind, lesen und ändern dieselben Dateien. Die Konfiguration speichert den Bibliothekspfad, die Sprache, Designs und weitere Einstellungen der Oberfläche.

Der lokale Lesestatus enthält Notizpfade, geöffnete Tabs und Scrollpositionen, aber keine Entwurfstexte. Der Browser kann Oberflächeneinstellungen zwischenspeichern. Laufzeitdateien enthalten Angaben zur Dienstidentität und Autorisierung. Diagnoseprotokolle bleiben lokal und sind so gestaltet, dass sie weder Notiztexte noch Sitzungszugangsdaten aufzeichnen.

#### Was Thoughts nicht erhebt

Thoughts betreibt keine Cloud-Notizdatenbank, keinen Kontodienst und keine Dienste für Analyse, Werbung, Telemetrie oder das automatische Hochladen von Notizen. Gespräche aus Ihren KI-Clients werden nicht automatisch erfasst. Es gibt weder ein integriertes KI-Modell noch einen Abonnementdienst.

#### KI-Clients verarbeiten Daten eigenständig

Wenn Sie Codex, Claude Desktop, Claude Code oder einem anderen Client das Lesen einer Notiz erlauben, werden der zurückgegebene Text und die Werkzeugergebnisse an diesen Client übermittelt. Ein Online-Client kann diese nach seinen eigenen Richtlinien und Ihren Kontoeinstellungen an seinen KI-Anbieter senden.

Lokale Speicherung bedeutet nicht, dass Inhalte, die Sie an eine Online-KI senden, offline bleiben. Prüfen Sie vor dem Teilen sensibler Inhalte die Berechtigungen und Dateneinstellungen des Clients. Thoughts kann keine Kopien löschen, die ein anderer Dienst aufbewahrt.

#### Lokale Sitzungen und externe Links

Der Lesedienst lauscht auf 127.0.0.1. Die lokale Autorisierung verwendet ein einmaliges Öffnungsticket, ein Sitzungs-Cookie und Anfrageprüfungen. Diese erlauben den lokalen Zugriff und sind kein Thoughts-Konto.

Der Leser lädt keine entfernten Notizbilder, externen Webschriften oder Analyseskripte. Externe Links werden geöffnet, wenn Sie sie auswählen; die Zielseite kann dabei übliche Netzwerkdaten eines Browserzugriffs erhalten. Thoughts ergänzt Supportlinks nicht um Notiztexte, Bibliothekspfade oder lokale Zugangsdaten.

#### Aufbewahrung, Löschung und Sicherungen

Ihre Dateien bleiben erhalten, bis Sie sie ändern oder entfernen. Das Verschieben einer Notiz oder eines Projekts in den lokalen Papierkorb ist rückgängig zu machen und keine endgültige Löschung. Gemeinsam genutzte Bilder, alte Kopien und Sicherungen können bestehen bleiben. Das Deinstallieren einer Integration löscht Ihre Bibliothek nicht.

Sie können die Markdown- und Bilddateien selbst öffnen, kopieren, exportieren und sichern. Um alle Kopien zu entfernen, prüfen Sie auch den Papierkorb, Bilder sowie verwendete Sicherungen und synchronisierte Ordner. Externe Dienste und Betriebssystem-Sicherungen folgen ihren eigenen Aufbewahrungsregeln.

#### Grenzen der Sicherheit

Thoughts verschlüsselt Ihre Dateien nicht. Lokale Sitzungsprüfungen verhindern nicht, dass andere Programme oder Benutzer mit Dateisystemberechtigungen die Dateien lesen oder ändern. Nutzen Sie die Zugriffskontrollen des Betriebssystems und geeignete Sicherungen.

Versionsprüfungen schützen unterstützte Schreibvorgänge vor veralteten Aktualisierungen. Sie bieten keine Garantie gegen beliebige externe Änderungen, Geräteausfälle oder Datenverlust. Teilen Sie keine lokalen Öffnungstickets oder Diagnosedateien, die Zugangsdaten enthalten.

#### Informationen, die Sie freiwillig an den Support senden

Herausgeber: xuli82221-pixel. Support-E-Mail: hello@spxipo.com. Thoughts versendet weder Support-E-Mails noch Diagnoseberichte automatisch.

Bei einer Supportanfrage entscheiden Sie selbst, welche Angaben Sie übermitteln. E-Mail-Inhalte werden von den beteiligten E-Mail-Anbietern verarbeitet. Lassen Sie private Notizen und Zugangsdaten weg, sofern diese nicht erforderlich sind. Fragen und Anliegen zu übermittelten Informationen können Sie an die Supportadresse richten.

#### Änderungen dieser Hinweise

Diese Hinweise werden mit der Anwendung bereitgestellt und tragen ein Aktualisierungsdatum. Prüfen Sie die Hinweise einer neuen Version. Lokal bereitgestellte Hinweise bedeuten nicht, dass eine öffentliche Datenschutzseite, ein Marktplatzeintrag oder eine Herausgeberprüfung bereits abgeschlossen ist.

### Nutzung, Lizenz und Unterstützung

#### Freie Software und Lizenz

Thoughts ist kostenlos und quelloffen unter der MIT License. Für die Software gilt der untenstehende englische Originaltext; diese übersetzte Beschreibung ist eine Zusammenfassung. Komponenten Dritter behalten ihre eigenen Lizenzen, die in den mitgelieferten Hinweisen aufgeführt sind.

Die MIT-Lizenz erlaubt Nutzung, Änderung und Weiterverteilung unter Beibehaltung der vorgeschriebenen Hinweise und stellt die Software ohne Gewährleistung gemäß dem Lizenztext bereit. Erstellen Sie Sicherungen und prüfen Sie wichtige Änderungen.

#### Ihre Dateien und die Berechtigungen der Clients

Wählen Sie einen dauerhaften Bibliotheksort und richten Sie Clients für dieselbe Bibliothek ein, wenn Sie Inhalte gemeinsam nutzen möchten. Aktualisierungen und die Installation von Integrationen sind von Ihren Markdown-Dateien getrennt. Gewähren Sie KI-Clients nur den Dateizugriff und die Aktionen, die Sie tatsächlich erlauben möchten.

#### Freiwillige Unterstützung

Spenden sind freiwillig, schalten keine Funktionen frei und sind weder ein Abonnement noch ein Kauf. Thoughts verlangt weder eine Wallet-Verbindung noch private Schlüssel und führt keine Transaktionen für Sie aus.

Wenn dies eingerichtet ist, zeigt die lokale Supportseite eine Base-Adresse und einen QR-Code, der nur die Adresse enthält. Das Senden von Geld ist eine getrennte Aktion in Ihrer Wallet. Wallet-Dienste und öffentliche Blockchain-Aktivitäten haben eigene Regeln und Offenlegungen von Daten; Thoughts erfasst Ihre Transaktionsdaten nicht automatisch.

## Français

Mise à jour: 2026-11-03

### Notice de confidentialité

Thoughts est un outil Markdown gratuit, open source et privilégiant le stockage local, conçu pour Windows. Aucun compte Thoughts n’est nécessaire. Cette notice décrit Thoughts lui-même ainsi que les services distincts que vous pouvez choisir d’utiliser.

#### Fichiers et paramètres locaux

Les notes et les images sont des fichiers conservés dans votre bibliothèque locale autorisée. Les clients configurés avec la même bibliothèque lisent et modifient ces mêmes fichiers. La configuration conserve le chemin de la bibliothèque, la langue, les thèmes et les autres préférences de l’interface.

L’état de lecture local contient les chemins des notes, les onglets ouverts et les positions de défilement, mais pas le texte des brouillons. Le navigateur peut mettre en cache des paramètres d’interface. Les fichiers d’exécution contiennent des informations d’identité du service et d’autorisation. Les journaux de diagnostic restent locaux et sont conçus pour ne pas enregistrer le texte des notes ni les identifiants de session.

#### Ce que Thoughts ne collecte pas

Thoughts n’exploite ni base de notes dans le cloud, ni service de comptes, d’analyse, de publicité, de télémétrie ou d’envoi automatique de notes. Il ne collecte pas automatiquement les conversations de vos clients IA. Il ne comporte ni modèle d’IA intégré ni service d’abonnement.

#### Les clients IA traitent les données séparément

Lorsque vous autorisez Codex, Claude Desktop, Claude Code ou un autre client à lire une note, le texte renvoyé et les résultats des outils sont transmis à ce client. Un client en ligne peut les envoyer à son fournisseur d’IA selon ses propres politiques et les paramètres de votre compte.

Le stockage local ne signifie pas que le contenu transmis à une IA en ligne reste hors ligne. Examinez les autorisations et les contrôles des données du client avant de partager des informations sensibles. Thoughts ne peut pas supprimer les copies conservées par un autre service.

#### Sessions locales et liens externes

Le service de lecture écoute sur 127.0.0.1. L’autorisation locale repose sur un ticket d’ouverture à usage unique, un cookie de session et des contrôles des requêtes. Ils autorisent l’accès local et ne constituent pas un compte Thoughts.

Le lecteur ne charge pas les images distantes des notes, les polices web externes ni les scripts d’analyse. Les liens externes s’ouvrent lorsque vous les choisissez ; le site de destination peut recevoir les informations réseau habituelles d’une visite par navigateur. Thoughts n’ajoute pas le texte des notes, les chemins de bibliothèque ou les identifiants locaux aux liens d’assistance.

#### Conservation, suppression et sauvegardes

Vos fichiers restent présents jusqu’à ce que vous les modifiiez ou les supprimiez. Le déplacement d’une note ou d’un projet vers la corbeille locale est réversible et ne constitue pas un effacement définitif. Des images partagées, d’anciennes copies et des sauvegardes peuvent subsister. Désinstaller une intégration ne supprime pas votre bibliothèque.

Vous pouvez ouvrir, copier, exporter et sauvegarder vous-même les fichiers Markdown et les images. Pour supprimer toutes les copies, vérifiez aussi la corbeille, les images, les sauvegardes et les dossiers synchronisés que vous utilisez. Les services distants et les sauvegardes du système d’exploitation appliquent leurs propres règles de conservation.

#### Limites de sécurité

Thoughts ne chiffre pas vos fichiers. Les contrôles de session locale n’empêchent pas d’autres programmes ou utilisateurs disposant de permissions sur les fichiers de les lire ou de les modifier. Utilisez les contrôles d’accès du système d’exploitation et des sauvegardes adaptées.

Les contrôles de version protègent les écritures prises en charge contre les mises à jour obsolètes. Ils ne garantissent pas une protection contre toute modification externe, panne matérielle ou perte de données. Ne partagez pas les tickets d’autorisation locale ni les fichiers de diagnostic contenant des identifiants.

#### Informations que vous choisissez d’envoyer à l’assistance

Éditeur : xuli82221-pixel. Adresse d’assistance : hello@spxipo.com. Thoughts n’envoie pas automatiquement de courriels d’assistance ni de rapports de diagnostic.

Si vous contactez l’assistance, vous choisissez les informations à communiquer. Les courriels sont traités par les fournisseurs de messagerie concernés. Sauf nécessité, n’incluez ni notes privées ni identifiants. Pour toute question ou demande concernant les informations fournies, contactez l’adresse d’assistance.

#### Modification de cette notice

Cette notice est fournie avec l’application et indique une date de mise à jour. Consultez la notice accompagnant chaque nouvelle version. La présence d’une notice locale ne signifie pas qu’une page publique de confidentialité, une publication sur une place de marché ou une vérification de l’éditeur a été finalisée.

### Utilisation, licence et assistance

#### Logiciel libre et licence

Thoughts est gratuit et open source sous licence MIT. Le texte original en anglais ci-dessous régit le logiciel ; cette description traduite est un résumé. Les composants tiers conservent leurs propres licences, indiquées dans les notices fournies.

La licence MIT autorise l’utilisation, la modification et la redistribution sous réserve de conserver les mentions requises. Le logiciel est fourni sans garantie, conformément au texte de la licence. Conservez des sauvegardes et vérifiez les modifications importantes.

#### Vos fichiers et les autorisations des clients

Choisissez un emplacement stable pour la bibliothèque et configurez les clients avec la même bibliothèque si vous souhaitez partager le contenu. Les mises à jour et l’installation des intégrations sont distinctes de vos fichiers Markdown. N’accordez aux clients IA que les accès aux fichiers et les actions que vous souhaitez autoriser.

#### Soutien facultatif

Les dons sont volontaires, ne débloquent aucune fonction et ne constituent ni un abonnement ni un achat. Thoughts ne demande ni connexion à un portefeuille ni clés privées et ne soumet pas de transactions à votre place.

Si cette fonction est configurée, la page locale de soutien affiche une adresse Base et un code QR contenant uniquement l’adresse. L’envoi de fonds est une action distincte dans votre portefeuille. Les services de portefeuille et l’activité publique sur la blockchain impliquent leurs propres divulgations de données et politiques ; Thoughts ne collecte pas automatiquement les détails de vos transactions.

## Español

Actualización: 2026-11-03

### Aviso de privacidad

Thoughts es una herramienta Markdown gratuita, de código abierto y que prioriza el almacenamiento local, diseñada para Windows. No necesitas una cuenta de Thoughts. Este aviso describe Thoughts y los servicios independientes que puedes decidir utilizar.

#### Archivos y ajustes locales

Las notas y las imágenes son archivos en tu biblioteca local autorizada. Los clientes configurados con la misma biblioteca leen y modifican esos mismos archivos. La configuración guarda la ruta de la biblioteca, el idioma, los temas y otras preferencias de la interfaz.

El estado de lectura local contiene rutas de notas, pestañas abiertas y posiciones de desplazamiento, pero no el texto de los borradores. El navegador puede almacenar ajustes de interfaz en caché. Los archivos de ejecución contienen información de identidad del servicio y de autorización. Los registros de diagnóstico son locales y están diseñados para no registrar el contenido de las notas ni las credenciales de sesión.

#### Qué no recopila Thoughts

Thoughts no ofrece una base de datos de notas en la nube ni servicios de cuentas, análisis, publicidad, telemetría o carga automática de notas. No recopila automáticamente conversaciones de tus clientes de IA. No tiene un modelo de IA integrado ni un servicio de suscripción.

#### Los clientes de IA procesan datos por separado

Cuando autorizas a Codex, Claude Desktop, Claude Code u otro cliente a leer una nota, el texto devuelto y los resultados de las herramientas se entregan a ese cliente. Un cliente en línea puede enviarlos a su proveedor de IA según sus propias políticas y los ajustes de tu cuenta.

El almacenamiento local no significa que el contenido enviado a una IA en línea permanezca sin conexión. Revisa los permisos y controles de datos del cliente antes de compartir información sensible. Thoughts no puede eliminar las copias que conserve otro servicio.

#### Sesiones locales y enlaces externos

El servicio de lectura escucha en 127.0.0.1. La autorización local utiliza un ticket de apertura de un solo uso, una cookie de sesión y comprobaciones de solicitudes. Estos mecanismos autorizan el acceso local; no constituyen una cuenta de Thoughts.

El lector no carga imágenes remotas de las notas, fuentes web externas ni scripts de análisis. Los enlaces externos se abren cuando los eliges; el destino puede recibir la información de red habitual de una visita del navegador. Thoughts no añade contenido de notas, rutas de biblioteca ni credenciales locales a los enlaces de soporte.

#### Conservación, eliminación y copias de seguridad

Tus archivos permanecen hasta que los cambies o elimines. Mover una nota o un proyecto a la papelera local es reversible y no equivale a borrarlo definitivamente. Pueden permanecer imágenes compartidas, copias antiguas y copias de seguridad. Desinstalar una integración no elimina tu biblioteca.

Puedes abrir, copiar, exportar y respaldar los archivos Markdown y de imagen por tu cuenta. Para eliminar todas las copias, revisa también la papelera, las imágenes, las copias de seguridad y las carpetas sincronizadas que utilices. Los servicios remotos y las copias del sistema operativo siguen sus propias reglas de conservación.

#### Límites de seguridad

Thoughts no cifra tus archivos. Las comprobaciones de las sesiones locales no impiden que otros programas o usuarios con permisos sobre el sistema de archivos los lean o modifiquen. Usa los controles de acceso del sistema operativo y copias de seguridad adecuadas.

Las comprobaciones de revisión protegen las escrituras compatibles frente a actualizaciones obsoletas; no garantizan la protección ante cualquier edición externa, fallo del dispositivo o pérdida de datos. No compartas tickets de autorización local ni archivos de diagnóstico que contengan credenciales.

#### Información que decides enviar al soporte

Editor: xuli82221-pixel. Correo de soporte: hello@spxipo.com. Thoughts no envía automáticamente correos de soporte ni informes de diagnóstico.

Si contactas con soporte, tú decides qué información compartir. Los proveedores de correo implicados procesan el contenido de los mensajes. Salvo que sea necesario, omite las notas privadas y las credenciales. Si tienes preguntas o solicitudes sobre la información que proporcionaste, contacta con la dirección de soporte.

#### Cambios en este aviso

Este aviso se incluye con la aplicación e indica una fecha de actualización. Revisa el aviso que acompaña a cada nueva versión. Disponer de un aviso local no implica que se hayan completado una página pública de privacidad, una publicación en un mercado de aplicaciones o la verificación del editor.

### Uso, licencia y soporte

#### Software libre y licencia

Thoughts es gratuito y de código abierto bajo la licencia MIT. El texto original en inglés que aparece más abajo rige el software; esta descripción traducida es un resumen. Los componentes de terceros conservan sus propias licencias, indicadas en los avisos incluidos.

MIT permite el uso, la modificación y la redistribución siempre que se conserven los avisos requeridos, y proporciona el software sin garantías según los términos de la licencia. Mantén copias de seguridad y verifica los cambios importantes.

#### Tus archivos y los permisos de los clientes

Elige una ubicación estable para la biblioteca y configura los clientes con la misma biblioteca cuando quieras compartir contenido. Las actualizaciones y la instalación de integraciones son independientes de tus archivos Markdown. Concede a los clientes de IA solo el acceso a archivos y las acciones que quieras autorizar.

#### Apoyo opcional

Las donaciones son voluntarias, no desbloquean funciones y no son una suscripción ni una compra. Thoughts no requiere conectar una cartera ni proporcionar claves privadas, y no envía transacciones en tu nombre.

Si está configurada, la página local de apoyo muestra una dirección de Base y un código QR que solo contiene la dirección. Enviar fondos es una acción independiente en tu cartera. Los servicios de cartera y la actividad pública de la cadena de bloques tienen sus propias políticas y exposición de datos; Thoughts no recopila automáticamente los detalles de tus transacciones.

## MIT License — English original

```text
MIT License

Copyright (c) 2026 xuli82221-pixel

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

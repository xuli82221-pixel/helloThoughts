[English](README.md) | [简体中文](README.zh-CN.md) | [繁體中文](README.zh-TW.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Deutsch](README.de.md) | [Français](README.fr.md) | [Español](README.es.md)

# Thoughts

[インストール](INSTALL.md) · [ドキュメント](docs/README.md) · [寄付](DONATE.md) · [コード署名ポリシー / Code signing policy](docs/CODE_SIGNING_POLICY.md)

Windows x64 向けの無料・オープンソースのローカル優先ノートツールです。AI アシスタントで Markdown ファイルを管理し、独立したローカル Web ページで閲覧・編集できます。

**未署名プレビュー — 2026-11-04:** `0.3.0-preview.21` ではソース、実行ファイル、クライアント設定の名称を Thoughts に統一します。管理者は[未署名の実験的プレリリース](https://github.com/xuli82221-pixel/hello-Thoughts-/releases/tag/v0.3.0-preview.21)の公開を承認しており、公開 Windows ビルドと検証が成功し、ファイルを公開しました。**信頼された発行元署名や安定版の受け入れ完了を意味しません。** SignPath Foundation の承認は未取得です。Windows SAC/SmartScreen がブロックした場合は実行を中止し、保護を維持してください。[公開準備](docs/RELEASE_READINESS_20261104.md)、[インストール](INSTALL.md)をご覧ください。

## ノートは自分のコンピューターに

Thoughts にはアカウント、クラウドデータベース、ノートのアップロードサービス、テレメトリ、内蔵 AI、有料テーマ、サブスクリプション、決済連携、ブラウザー拡張はありません。Markdown とローカル画像は利用者のもので、他のエディターでも開けます。

Go 実行ファイルに React の Web ページを内蔵しています。**ビルド済みの配布物は、実行時に Node.js、Go、Python、Docker を必要としません。**利用する AI クライアントには、別途そのクライアントのインストール・アカウント要件が適用されます。

## ダウンロードとインストール

この実験的プレビューでは、[GitHub Releases](https://github.com/xuli82221-pixel/hello-Thoughts-/releases) から **`thoughts-<version>-windows-amd64-setup.exe`**、ポータブル ZIP、Claude Desktop MCPB を提供します。プログラム名は **`thoughts.exe`** です。使用前に SHA256 を確認してください。現行ファイルに **Authenticode 署名はなく**、ハッシュの一致は発行元の証明にはなりません。Windows にブロックされた場合は中止してください。セットアップはユーザー単位で、管理者権限なしに共有ホームとデスクトップ・スタートメニューのショートカットを作成します。既存データは保持し、アンインストール後も他のクライアント用に共有ホームを残します。署名済み安定版は未提供です。

Claude Desktop では、そのリリースの `.mcpb` 拡張をインポートする方法もあります。Codex と Claude Code では ZIP 内の実行ファイルをローカル stdio MCP サーバーとして登録します。Web ページは AI クライアントなしでも使用できます。コマンドは[一般利用者向けクイックスタート](INSTALL.md#ordinary-windows-users)をご覧ください。GitHub の **Code → Download ZIP** はソースの取得であり、ビルド済みアプリの取得ではありません。

AI クライアントは同じローカルの標準 stdio MCP アダプターを起動し、認証済みのローカルコアを呼び出します。Web ページは `127.0.0.1` のみで待ち受けます。オンライン AI クライアントにノートを読むよう依頼すると、その内容はクライアントに渡り、クライアント自身のデータポリシーの対象になります。ローカル優先でも、依頼した AI ツールの結果が必ずコンピューター内に留まるわけではありません。Thoughts はファイルを暗号化しません。

## 対象クライアント

| クライアント | 導入方法 |
| --- | --- |
| Claude Desktop | 公式の `.mcpb` 拡張形式 |
| Windows 版 Claude Code | ローカル stdio MCP、Windows ポータブル ZIP |
| Windows 版 OpenAI Codex | ローカル stdio MCP、同じポータブル ZIP |

これらは対象とする連携です。実際の受け入れ状況は[公開リリース状況](docs/RELEASE_READINESS_20261104.md)に記載しています。通常の ChatGPT Web プラグイン、macOS、Linux、モバイル、リモート MCP、クラウド同期は今回の対象外です。

## 機能

- ローカル Markdown の作成、閲覧、編集、検索、名前変更、移動。
- コンパクトなプロジェクトツリー、クイックノート、プロジェクトのアーカイブと解除、復元可能なノート・プロジェクトのごみ箱。
- 既存ノートは整形済みプレビューで開き、編集は明示的に切り替え。
- ローカル PNG/JPEG/WebP 画像、プレビュー、拡大表示。
- 無料テーマと、検証付きのローカルテーマインポート。
- 英語、韓国語、日本語、簡体字中国語、繁体字中国語、ドイツ語、フランス語、スペイン語の UI と、対応するローカルのプライバシー・利用案内。
- 外部でのファイル変更通知と SHA256 による競合保護。
- 全クライアントで同じ設定済みライブラリを使用し、再初期化時も既存の選択を保持。
- チェックサムで検証したローカルコアの更新とソフトウェアのロールバック、未保存編集の保護。
- サービスのポートが変わっても復元できる、サーバー側保存のノートタブと閲覧位置。

競合保護は古い版に基づく書き込みを拒否します。自動の変更履歴は未実装なので、ライブラリは別途バックアップしてください。競合や切断時に開いたタブが下書きを保持することはありますが、タブを閉じたり Windows を再起動した後の復元は保証されません。

![架空のサンプルノートを表示した Thoughts 0.3.0-preview.20](docs/assets/preview-0.3.0-preview.20.jpg)

2026-11-04 に、新しい独立した一時ライブラリと架空のノートで撮影した `0.3.0-preview.20` の実際の UI です。個人のノート、実画像、端末のパス、認証情報は写っていません。[旧レイアウトの画像](docs/assets/preview.png)は履歴資料として保持し、現行版としては使用しません。画像の確認はインストーラーやクライアントの受け入れ確認を意味しません。

## インストール

SHA256 の確認、実行前に内容を確認するコマンド、3 クライアントの手順は [INSTALL.md](INSTALL.md) にあります。preview.8 以降の新規導入では、導入順序に関係なく全クライアントが同じ Windows ユーザーの Thoughts ディレクトリを既定で使用します。詳細設定のホームは `auto` のままにします。意図的にホームを変更する場合のみ、絶対パスの `THOUGHTS_HOME` を指定します。

```text
<Thoughts home>/
  library/     Markdown, images and recoverable Trash
  config/      Configuration, local sessions, preferences and logs
  app/         Managed executable versions (optional)
  backups/     Previous runtime pins; not a backup of note bodies
```

通常は共有の既定ホームを使ってください。従来の明示的な設定も利用できますが、既存ライブラリの自動移行や統合は行いません。ホームを変更する場合、共有する全クライアントで同じ物理フォルダーを指定します。

インストールや更新時にはホームを保持する必要があります。閲覧状態はパスとスクロール位置を保存し、未保存の本文は保存しません。コア更新には、レビュー済みのローカル実行ファイルと SHA256 の明示的な指定が必要で、バックグラウンドのネットワーク更新機能はありません。クライアント拡張と共有コアは別々に更新します。[更新手順](INSTALL.md#managed-core-updates-and-rollback)をご覧ください。

未署名の Windows プログラムやスクリプトは、SmartScreen、Smart App Control、組織のアプリ制御ポリシーでブロックされることがあります。SHA256 はチェックサムとの整合性を示すもので、発行元を証明しません。プレビューを実行するために保護を無効化したり警告を回避したりしないでください。ソース公開だけでは署名の問題は解決しません。

## ワークスペース

[構成](docs/WORKSPACE_LAYOUT.md)と[ドキュメント一覧](docs/README.md)をご覧ください。過去の開発記録は公開ドキュメントには含めていません。ソースと実行ファイルは Thoughts 名に統一し、旧データの識別子は互換性のため保持します。

## ビルド

Windows x64、go.mod に対応する Go、Node.js 24、Python 3.12+、署名済みの公式 [Inno Setup 7 コンパイラー](https://jrsoftware.org/isdl.php)が必要です。一般利用者にこれらのツールは不要です。[CONTRIBUTING.md](CONTRIBUTING.md)をご覧ください。

```powershell
Set-Location web
npm ci
Set-Location ..
pwsh -File scripts/build-desktop-adapters.ps1 -Python 'C:\Tools\Python\python.exe' -Go 'C:\Tools\Go\bin\go.exe' -InnoCompiler 'C:\Tools\InnoSetup7\ISCC.exe'
```

ツールのパスは実際の導入先に置き換えます。フロントエンドと Go コアのテスト・ビルド、通知の生成、Setup.exe/ZIP/MCPB とチェックサムの作成、パッケージ検証を新しい出力先で行います。既存出力の上書きや自動アップロードはしません。公開候補モードには信頼された署名と実在する発行元メタデータが必要です。署名はクリーン環境やクライアントの受け入れ確認の代わりにはなりません。

## セキュリティとライセンス

[SECURITY.md](SECURITY.md)、[PRIVACY.md](PRIVACY.md)、[CHANGELOG.md](CHANGELOG.md)、[第三者の通知](THIRD_PARTY_NOTICES.txt)、[依存関係のレビュー](docs/DEPENDENCIES.md)をご覧ください。プロジェクトのコードは MIT License、同梱する依存物には OFL フォントを含め、それぞれのライセンスが適用されます。

実際のノート、個人のスクリーンショット、認証情報、セッションファイル、ログはコミットしないでください。脆弱性は hello@spxipo.com へ非公開で報告してください。設定には、所有者の Base メインネットアドレスによる任意の支援案内があります。QR に含まれるのは送金先アドレスのみです。ウォレットでは Base を別途選択してください。Thoughts はウォレット接続、秘密鍵の要求、資産や金額の指定、決済処理を行いません。全機能は無料です。

## コード署名ポリシー

**審査待ち — SignPath Foundation の申請はまだ承認されていません。**現行の Windows プログラムはすべて未署名です。承認された場合、無料のコード署名は SignPath.io、証明書は SignPath Foundation が提供します。[コード署名ポリシー / Code signing policy](docs/CODE_SIGNING_POLICY.md)に、署名の対象範囲、現在の状態、未完了の要件を記載しています。SignPath Foundation への申請はまだ送信していません。

## 寄付

支援は任意で、機能の解除には結び付きません。専用の[寄付ページ](DONATE.md)に、所有者の既存の Base アドレス、アドレスのみの QR、プライバシーの境界を記載しています。問題報告、翻訳、ドキュメントでも協力できます。Thoughts はウォレットに接続せず、決済を処理しません。

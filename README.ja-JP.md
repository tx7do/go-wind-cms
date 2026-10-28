<div align="center">

<img src="docs/brand/vortex-tile.svg" width="120" alt="GoWind Content Hub" />

# GoWind Content Hub

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org/)
[![React](https://img.shields.io/badge/React-19.x-61DAFB?logo=react&logoColor=white)](https://react.dev/)
[![Flutter](https://img.shields.io/badge/Flutter-3.x-02569B?logo=flutter&logoColor=white)](https://flutter.dev/)
[![Kratos](https://img.shields.io/badge/Kratos-2.9-00ADD8?logo=go&logoColor=white)](https://go-kratos.dev/)

[English](./README.en-US.md) | [中文](./README.md) | **日本語**

</div>

---

GoWind Content Hub（風行）は、すぐに使える企業向けの Golang フルスタック Headless コンテンツプラットフォームであり、柔軟で拡張性の高いコンテンツ管理と配信ソリューションを企業に提供します。

**主な特徴：**

- **API ファースト** — 完全な RESTful / gRPC デュアルプロトコルインターフェース、OpenAPI ドキュメント自動生成
- **マルチプラットフォーム** — Vue、React、Taro（ミニプログラム）、Flutter の 4 つのフロントエンドをサポート
- **マルチテナント** — テナントデータの分離、ロール・管理者の自動初期化
- **きめ細かい権限管理** — メニュー権限、API権限、データ権限の 3 レベル制御、ボタンレベルの制御に対応
- **ネイティブ i18n** — コンテンツ、メニュー、UIテキストの統一翻訳管理
- **マイクロサービス** — go-kratos ベース、サービスディスカバリ、分散トレーシング、分散トランザクション対応

## デモ

| デモタイプ | アクセスアドレス |
|:---|:---|
| 管理画面 | [https://admin.cms.gowind.cloud](https://admin.cms.gowind.cloud) |
| 管理画面 API Swagger | [https://api.admin.cms.gowind.cloud/docs/](https://api.admin.cms.gowind.cloud/docs/) |
| フロント API Swagger | [https://api.cms.gowind.cloud/docs/](https://api.cms.gowind.cloud/docs/) |
| Vue フロントエンド | [https://cms.gowind.cloud](https://cms.gowind.cloud) |
| React フロントエンド | [https://react.cms.gowind.cloud](https://react.cms.gowind.cloud) |
| Taro フロントエンド | [https://taro.cms.gowind.cloud](https://taro.cms.gowind.cloud) |
| Flutter フロントエンド | [https://flutter.cms.gowind.cloud](https://flutter.cms.gowind.cloud) |

> すべてのデモサイトの共通アカウント：`admin` / `admin`

## テクノロジースタック

### バックエンド

| レイヤー | 技術 | 説明 |
|:---|:---|:---|
| 言語 | [Go 1.25+](https://go.dev/) | 高性能コンパイル言語 |
| フレームワーク | [go-kratos](https://go-kratos.dev/) | Bilibili オープンソースマイクロサービスフレームワーク |
| DI | 手書き組み立て(`wiring.go`) | DIフレームワークなし、`make register`で新モジュールを登録 |
| ORM | [Ent](https://entgo.io/) | Go エンティティフレームワーク |
| データベース | [PostgreSQL](https://www.postgresql.org/) / [MySQL](https://www.mysql.com/) | リレーショナルデータベース |
| キャッシュ | [Redis](https://redis.io/) | インメモリデータストア |
| オブジェクトストレージ | [MinIO](https://min.io/) | S3 互換オブジェクトストレージ |
| 検索エンジン | [OpenSearch](https://opensearch.org/) | 全文検索 |
| サービスレジストリ | [Etcd](https://etcd.io/) / [Consul](https://www.consul.io/) | サービスディスカバリ & 設定管理 |
| トレーシング | [Jaeger](https://www.jaegertracing.io/) + [OpenTelemetry](https://opentelemetry.io/) | 分散オブザーバビリティ |
| API 定義 | [Protobuf](https://protobuf.dev/) + [buf.build](https://buf.build/) | コントラクトファースト API 設計 |
| 認可エンジン | [Casbin](https://casbin.org/) / [OPA](https://www.openpolicyagent.org/) | ポリシーベース認可 |

### 管理画面フロントエンド

| 技術 | 説明 |
|:---|:---|
| [React 19](https://react.dev/) | プログレッシブフロントエンドフレームワーク |
| [TypeScript](https://www.typescriptlang.org/) | 型安全開発 |
| [Ant Design](https://ant.design/) | エンタープライズ UI コンポーネントライブラリ |
| [ProComponents](https://procomponents.ant.design/) | エンタープライズ管理画面コンポーネント |

### アプリケーションフロントエンド

**まず「フロントエンド」の 2 種類を整理**：本プロジェクトにはユーザー向けインターフェースが 2 種類あります。**管理画面**は運用チーム向け——PC ブラウザで開き、記事の作成・コメント管理・サイト設定を行います。**アプリケーションフロントエンド**（下記 4 種）は読者・顧客向け——公開されたコンテンツを閲覧します。バックエンドは 1 つだけで、公開後のコンテンツは 4 つのフロントエンドすべてが同じデータを参照します。**読者の接点に合わせて 1 つ（または複数）を選べばよく、すべてをデプロイする必要はありません。**

| バージョン | 技術スタック | 動作環境 | 典型的な用途 |
|:---|:---|:---|:---|
| Vue | [Nuxt.js](https://nuxt.com/) + [shadcn-vue](https://www.shadcn-vue.com/) | PC / スマホの**ブラウザ（Web ページ）** | 企業サイト、コンテンツポータル、ブログ、ニュースサイト。SEO に有利 |
| React | [Next.js](https://nextjs.org/) + [shadcn/ui](https://ui.shadcn.com/) | PC / スマホの**ブラウザ（Web ページ）** | Vue 版と機能的に同等。チームの技術スタックに合わせてどちらかを選択 |
| Taro | [Taro](https://docs.taro.zone/en/docs/) + React + [shadcn/ui](https://ui.shadcn.com/) | **WeChat ミニプログラム**、H5 モバイル Web にもビルド可 | WeChat 内での運用：QR コードや検索で即開け、インストール不要 |
| Flutter | [Flutter](https://flutter.dev/) + [BLoC](https://bloclibrary.dev/) | **iOS / Android ネイティブアプリ**（アプリストアからインストール） | 独立アプリ：プッシュ通知・オフライン利用に対応。1 つのコードベースで両プラットフォームを生成 |

#### 用語解説（フロントエンドに不慣れな方向け）

- **Web ページ**：ブラウザで URL を開くだけでアクセスでき、PC・スマホ両対応、インストール不要。React 版と Vue 版がこれに該当し、どちらも **SSR**（サーバーサイドレンダリング：ページをサーバー側で組み立ててからブラウザへ送る方式。検索エンジンにクロールされやすく、**SEO** に有利）に対応しています。
- **H5**：WeChat やモバイルブラウザ内で開くモバイル向け Web ページ。実態は Web ページで、インストール不要、タップして即閲覧できます。
- **WeChat ミニプログラム**：WeChat 内で動く「ライトアプリ」。QR コードのスキャンや WeChat 内検索で開け、ダウンロード不要・ホーム画面も占有しません。公式アカウント導線・グループ共有・オフライン QR など、中国の WeChat エコシステムで最も一般的な接点です。
- **ネイティブアプリ**：App Store / Android アプリストアからダウンロードし、スマホのホーム画面にインストールするアプリ。動作が最も滑らかで、プッシュ通知・カメラ・オフライン保存など端末機能を利用できます。Flutter の価値は**1 つのコードベースから iOS と Android の両アプリを生成できる**ことです。

#### どれを選ぶべきか？

- 読者がブラウザでアクセスする（公式サイト / ポータル / ブログ）→ **React** または **Vue** を選択。機能は同等なので、チームの得意なスタックでどうぞ。
- ターゲットが WeChat ユーザー（公式アカウント、グループ共有、オフライン QR）→ **Taro**（WeChat ミニプログラム）。
- アプリストアに自社アプリを出したい、プッシュ通知など端末機能が必要 → **Flutter**。
- 組み合わせも可能：例えば公式サイトに React、ミニプログラムに Taro——両方とも同じコンテンツと設定を参照するため、二重入力は不要です。

## コア機能

### 組織 & 権限

| 機能 | 説明 |
|:---|:---|
| マルチテナント管理 | テナントの追加、有効化/無効化、パッケージ設定、データ分離。新規テナントのロール、管理者を自動初期化 |
| ユーザー管理 | ライフサイクル全体を管理、複数ロール/部署のバインディング、高度な検索 |
| ロール管理 | ロール & ロールグループ管理、メニュー/API/データ権限のきめ細かい設定 |
| 権限管理 | 権限グループ、メニューノード、権限ポイント、ボタンレベルの制御に対応 |
| メニュー管理 | ディレクトリ/ページ/ボタンの 3 レベル設定、権限に基づく動的レンダリング |
| OpenAPI 認証情報 | テナント単位の AK/SK マシン認証情報管理、トークン交換で短期 JWT 発行、IP+AK 総当たり制限 |
| テナントクォータ | ユーザー/ストレージ/API 呼び出し/AI トークンの 4 種クォータをテナント別に設定、ハードゲート + 透過アラート、セルフサービス使用量パネル |
| ログインポリシー | IP/時間/デバイス/地域の 4 軸ブラックリスト/ホワイトリスト、グローバル+ターゲット 2 段階執行；パスワード有効期限切れはログイン拒否 |
| 部門管理 | 複数レベルの部門ツリー管理、ユーザーとの連携バインディング |

### コンテンツ & サイト

| 機能 | 説明 |
|:---|:---|
| コンテンツモデリング | テキスト、数値、リッチテキスト、画像、ファイル、関連フィールドを持つカスタムコンテンツモデル |
| コンテンツ管理 | CRUD、公開/非公開、ピン留め、並べ替え、ゴミ箱、一括操作。リッチテキスト / Markdown 編集 |
| カテゴリ管理 | 複数レベルのカテゴリツリー、コンテンツモデルのバインディング、フロントエンドフィルタリング |
| タグ管理 | タグ CRUD、コンテンツ関連付け、検索 & 集約表示 |
| コメント管理 | 審査、削除、返信、ブロック。コンテンツ/ユーザー/時間でフィルタリング |
| 多言語管理 | 言語管理、コンテンツ/メニュー/UIテキストの統一翻訳 |
| サイト管理 | 複数サイトの独立設定、ドメイン、タイトル、ロゴ、SEO、スタイル |
| サイト設定 | 基本情報、SEO、アップロード制限、メール/SMS グローバル設定 |
| ページビルダー | セクションベースのページエディター（14 種類のブロック）、3 つのフロントエンドで描画 |
| 承認ワークフロー | 記事のレビュー申請→承認/却下、結果を作者に通知（サイト内メッセージ+メール）；権限ゲートは API ロールバインディング |
| コンテンツ実験 | A/B ランダム分割とオーディエンスターゲティング（locale/デバイス/ログイン/パラメータ）、バリアント丸ごと置換、参加統計 |
| 共同編集ロック | 記事/ページエディタで他ユーザー編集中のバナー表示（ハートビート付き Redis ソフトロック）、同時上書きを防止 |
| リアルタイム共同編集 | リッチテキストの複数人同時編集を競合なくマージ、カーソル共有（ygo エンジン内蔵、コンテンツ+言語ごとのルーム）、保存は明示操作のまま |
| コンテンツバージョン | 保存時に自動スナップショット、バージョン一覧・行差分・ワンクリック復元 |
| 下書きプレビュー | エディターから新しいタブで下書きプレビュー（短命署名トークン） |
| リダイレクト管理 | 301/302 リダイレクトルール、nginx map 設定のエクスポート |
| フォーム送信 | ページ FORM セクションへの訪問者送信を収集、管理一覧/詳細/CSV エクスポート |

### システム & 運用

| 機能 | 説明 |
|:---|:---|
| ファイル管理 | 画像/ドキュメント/動画の管理、ローカルまたは OSS ストレージ、プレビュー、ダウンロード、グループ管理 |
| 辞書管理 | データ辞書カテゴリ & アイテム、連携検索、並べ替え、インポート/エクスポート |
| API 管理 | バックエンド API の保守 & 自動同期、ツリー表示 |
| タスクスケジューリング | Cron ジョブ管理、開始/一時停止/即時実行、実行履歴 & ログ |
| メッセージ通知 | 複数レベルのメッセージカテゴリ、対象ユーザーへの送信、既読ステータス追跡 |
| 站内メッセージ | 個人受信トレイ、表示/削除、個別/一括既読 |
| ログインログ | ログイン成功/失敗ログ、アカウント、IP、所在地、デバイス、時刻 |
| 操作ログ | 全リンク監査証跡、操作者、IP、パラメータ & 結果 |
| パーソナルセンター | プロフィール編集、アバター、パスワードリセット、ログイン履歴 |
| Headless API | コンテンツ CRUD の完全な OpenAPI、マルチターミナル対応；加えて読み取り専用 GraphQL コンテンツ API（記事/ページ/カテゴリ/タグ、`:6702`） |
| Webhook 配信 | コンテンツイベント（公開/更新/削除）を HMAC 署名付きで非同期配信、配信履歴とテスト送信対応 |
| メール送信 | システム通知を非同期 SMTP で配信、SMTP 可視化設定（認証情報暗号化）+ テスト送信 |
| API 同期 | OpenAPI 整合同期（空でないテーブルでも安全、ID 安定でロールバインディングを保護） |
| バックアップ・リストア | 中核 ID/権限/組織 24 テーブルのワンクリック非同期バックアップ（MinIO バケット + プロジェクト接頭辞分離）、空 DB 一括リストアで ID 保持、バケットライフサイクルクリーンアップ |
| システムパラメータ | プラットフォーム全体の動的 KV 設定、パラメータ管理ページで視覚的に調整、プロセス内キャッシュ + Redis マルチインスタンス無効化ブロードキャスト；内蔵パスワードポリシーしきい値 |
| オンラインセッション | 全オンラインセッション一覧と強制ログアウト、ユーザー名/IP フィルタ対応 |
| サーバーモニター | Go ランタイム/DB 接続プール/ホスト情報の読み取り専用ダッシュボード（10 秒自動更新） |
| Redis キャッシュモニター | INFO セクションツリー、DBSIZE とスローログの読み取り専用パネル |

### ソーシャルディストリビューション

| 機能 | 説明 |
|:---|:---|
| ソーシャルアカウント | 10 プラットフォームに接続 — Bluesky / Mastodon / DEV.to / Telegram / WordPress / WeChat公式アカウント / Nostr / Weibo / LinkedIn / cnblogs（MetaWeblog）（多くは審査不要）；保存時に認証情報を検証、AES-GCM で暗号化保存、読み取りはマスク表示のみ |
| コンテンツ配信 | プラットフォーム別のテキスト・タイトル上書き、予約/即時投稿；ステートマシンで全体追跡、ワンクリック再試行 |
| 投稿カレンダー | カレンダー/リストの切り替え、月次スケジュール表示、日付クリックで素早く予約 |
| エンゲージメント統計 | いいね/返信/リポストのスナップショット取得、手動即時更新 + 定期自動同期 |

### AI アシスタント

| 機能 | 説明 |
|:--------|:------------|
| モデルプロバイダー | OpenAI 互換エンドポイント管理（クラウド API / ローカル Ollama）、API キー暗号化、利用量記録 |
| チャットアシスタント | SSE ストリーミング対話、CMS 読み取り専用ツール（概览/投稿検索/コメント統計） |
| ナレッジベース RAG | ドキュメント（txt/md/docx/pdf）自動チャンク＆pgvector へベクトル化、チャット検索注入 |
| AI クエリ | 自然言語 → 読み取り専用 SQL → 構造化結果、テーブルホワイトリスト＆テナント fail-closed |
| セマンティック検索 | 投稿タイトル/要約をグローバル検索にベクトル化、増分＆全再構築 |
| 安全インサイト | 24h 監査異常パターンのルール検出＋LLM サマリー（プラットフォーム管理者） |
| AI コンテンツ生成 | 投稿タイトル/要約/本文、ページタイトル、カテゴリ/タグ説明、コメント返信、ロール説明のワンクリック生成 |
| 利用統計 | 対話/生成/embedding 全呼び出し記録、月次サマリー＆サーバーサイドエクスポート |
| AI Agent オーケストレーション | 5 種類ノード（開始/知識検索/LLM/条件/終了）のビジュアルキャンバスオーケストレーション、同期テストランでノード毎トレース + SSE リアルタイム軌跡、実行台帳と使用量記録 |

## プロジェクト構造

```
go-wind-cms/
├── backend/                        # バックエンドサービス
│   ├── api/                        # Protobuf API 定義 & 生成コード
│   │   ├── protos/                 # .proto ソースファイル
│   │   └── gen/                    # 生成コード (Go / TypeScript / OpenAPI)
│   ├── app/
│   │   ├── admin/service/          # 管理サービス (HTTP/gRPC)
│   │   ├── app/service/            # アプリサービス (HTTP/gRPC)
│   │   └── core/service/           # コアサービス (ビジネスロジック + データ層)
│   ├── pkg/                        # 共通パッケージ (認可/暗号化/イベントバス/JWT/ミドルウェア/OSS...)
│   └── scripts/                    # デプロイスクリプト (Docker/PM2/環境設定)
├── frontend/
│   ├── admin/                      # 管理画面フロントエンド (React 19 + Ant Design 6 + ProComponents)
│   └── app/                        # アプリケーションフロントエンド
│       ├── react/                  # React アプリ (Next.js)
│       ├── vue/                    # Vue アプリ (Nuxt.js)
│       ├── taro/                   # Taro アプリ (WeChat ミニプログラム / H5)
│       └── flutter_app/            # Flutter アプリ (クロスプラットフォームネイティブ)
└── docs/                           # ドキュメント & スクリーンショット
```

## クイックスタート

### 前提条件

- Go 1.25+
- Docker & Docker Compose
- Node.js 18+ & pnpm
- buf (Protobuf ツールチェーン)

### 1. 依存サービスの起動

```bash
cd backend

# Windows
.\scripts\docker\libs_only.ps1

# Linux / macOS
./scripts/docker/libs_only.sh
```

### 2. バックエンドの起動

```bash
# 推奨：gow CLI を使用
gow run admin

# または IDE で直接デバッグ
```

### 3. フロントエンドの起動

```bash
# 管理画面
cd frontend/admin
pnpm install
pnpm dev

# React アプリ
cd frontend/app/react
pnpm install
pnpm dev
```

### よく使うコマンド

```bash
cd backend

# Protobuf API コードの生成
make api

# OpenAPI ドキュメントの生成
make openapi

# TypeScript コードの生成
make ts

# 全コードの生成 (ent + api + openapi)
make gen

# 全サービスのビルド
make build

# テストの実行
make test
```

> 開発ワークフローの詳細は [バックエンドドキュメント](./backend/README.md) と [スクリプトガイド](./backend/scripts/WORKFLOWS_AND_BEST_PRACTICES.md) を参照してください。

## バックエンドスクリーンショット

<table>
    <tr>
        <td><img src="./docs/images/admin_login.png" alt="バックエンドログイン"/></td>
        <td><img src="./docs/images/admin_register.png" alt="バックエンド登録"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_post_list.png" alt="バックエンド投稿リスト"/></td>
        <td><img src="./docs/images/admin_post_edit.png" alt="バックエンド投稿編集"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_category_list.png" alt="バックエンドカテゴリリスト"/></td>
        <td><img src="./docs/images/admin_category_edit.png" alt="バックエンドカテゴリ編集"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_tag_list.png" alt="バックエンドタグリスト"/></td>
        <td><img src="./docs/images/admin_tag_edit.png" alt="バックエンドタグ編集"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_comment_list.png" alt="バックエンドコメントリスト"/></td>
        <td><img src="./docs/images/admin_site_list.png" alt="バックエンドサイトリスト"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_site_setting_list.png" alt="バックエンドサイト設定"/></td>
        <td><img src="./docs/images/admin_navigation_list.png" alt="バックエンドナビゲーションリスト"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/admin_analysis.png" alt="バックエンドデータ分析"/></td>
        <td><img src="./docs/images/admin_media_asset_list.png" alt="バックエンドメディア資産"/></td>
    </tr>
</table>

## フロントエンドスクリーンショット

<table>
    <tr>
        <td><img src="./docs/images/react_app_login.png" alt="React フロントエンドログイン"/></td>
        <td><img src="./docs/images/react_app_register.png" alt="React フロントエンド登録"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_homepage.png" alt="React フロントエンドホーム"/></td>
        <td><img src="./docs/images/react_app_about.png" alt="React フロントエンド概要"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_post_list.png" alt="React フロントエンド投稿リスト"/></td>
        <td><img src="./docs/images/react_app_post_detail.png" alt="React フロントエンド投稿詳細"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_category_list.png" alt="React フロントエンドカテゴリリスト"/></td>
        <td><img src="./docs/images/react_app_category_detail.png" alt="React フロントエンドカテゴリ詳細"/></td>
    </tr>
    <tr>
        <td><img src="./docs/images/react_app_tag_list.png" alt="React フロントエンドタグリスト"/></td>
        <td><img src="./docs/images/react_app_tag_detail.png" alt="React フロントエンドタグ詳細"/></td>
    </tr>
</table>

## お問い合わせ

- WeChat 個人アカウント：`yang_lin_bo`（備考：`go-wind-cms`）
- 掘金コラム：[go-wind-cms](https://juejin.cn/column/7541283508041826367)

## コントリビューション

以下のような貢献を歓迎します：

- [Issue](../../issues) の提出：バグ報告や機能提案
- [Pull Request](../../pulls) の提出：修正や新機能の追加
- ドキュメントや翻訳の改善
- 利用経験の共有

## ライセンス

このプロジェクトは [MIT License](./LICENSE) の下で公開されています。

## 謝辞

[![JetBrains](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://jb.gg/OpenSource)

[JetBrains](https://jb.gg/OpenSource) から無料の GoLand & WebStorm オープンソースライセンスを提供していただき、感謝いたします。

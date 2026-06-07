<div align="center">

# GoWind Content Hub

### 風行 · すぐに使える企業向けフロントエンド・バックエンド一体型コンテンツプラットフォーム

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org/)
[![React](https://img.shields.io/badge/React-19.x-61DAFB?logo=react&logoColor=white)](https://react.dev/)
[![Flutter](https://img.shields.io/badge/Flutter-3.x-02569B?logo=flutter&logoColor=white)](https://flutter.dev/)
[![Kratos](https://img.shields.io/badge/Kratos-2.9-00ADD8?logo=go&logoColor=white)](https://go-kratos.dev/)

**[English](./README.en-US.md)** · **[中文](./README.md)** · **日本語**

</div>

---

GoWind HCH（風行）は、すぐに使える企業向けの Golang フルスタック Headless コンテンツプラットフォーム（HCH = Headless Content Hub）であり、柔軟で拡張性の高いコンテンツ管理と配信ソリューションを企業に提供します。

**主な特徴：**

- **API ファースト** — 完全な RESTful / gRPC デュアルプロトコルインターフェース、OpenAPI ドキュメント自動生成
- **マルチプラットフォーム** — Vue、React、Taro（ミニプログラム）、Flutter の 4 つのフロントエンドをサポート
- **マルチテナント** — テナントデータの分離、部門・ロール・管理者の自動初期化
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
| DI | [Wire](https://github.com/google/wire) | コンパイル時依存性注入 |
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
| [Vue 3](https://vuejs.org/) | プログレッシブフロントエンドフレームワーク |
| [TypeScript](https://www.typescriptlang.org/) | 型安全開発 |
| [Ant Design Vue](https://antdv.com/) | エンタープライズ UI コンポーネントライブラリ |
| [Vben Admin](https://doc.vben.pro/) | 管理画面フレームワーク |

### アプリケーションフロントエンド

| バージョン | 技術スタック | 用途 |
|:---|:---|:---|
| Vue | [Nuxt.js](https://nuxt.com/) + [shadcn-vue](https://www.shadcn-vue.com/) | Web アプリ / SSR |
| React | [Next.js](https://nextjs.org/) + [shadcn/ui](https://ui.shadcn.com/) | Web アプリ / SSR |
| Taro | [Taro](https://docs.taro.zone/en/docs/) + React + [shadcn/ui](https://ui.shadcn.com/) | WeChat ミニプログラム / H5 |
| Flutter | [Flutter](https://flutter.dev/) + [BLoC](https://bloclibrary.dev/) | クロスプラットフォームネイティブアプリ |

## コア機能

### 組織 & 権限

| 機能 | 説明 |
|:---|:---|
| マルチテナント管理 | テナントの追加、有効化/無効化、パッケージ設定、データ分離。新規テナントの部門、ロール、管理者を自動初期化 |
| ユーザー管理 | ライフサイクル全体を管理、複数ロール/部署のバインディング、代理ログイン、高度な検索 |
| ロール管理 | ロール & ロールグループ管理、メニュー/API/データ権限のきめ細かい設定 |
| 権限管理 | 権限グループ、メニューノード、権限ポイント、ボタンレベルの制御に対応 |
| メニュー管理 | ディレクトリ/ページ/ボタンの 3 レベル設定、権限に基づく動的レンダリング |
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
| サイト設定 | 基本情報、SEO、アップロード制限、キャッシュポリシー、メール/SMS グローバル設定 |

### システム & 運用

| 機能 | 説明 |
|:---|:---|
| ファイル管理 | 画像/ドキュメント/動画の管理、ローカルまたは OSS ストレージ、プレビュー、ダウンロード、グループ管理 |
| 辞書管理 | データ辞書カテゴリ & アイテム、連携検索、並べ替え、インポート/エクスポート |
| API 管理 | バックエンド API の保守 & 自動同期、ツリー表示 |
| タスクスケジューリング | Cron ジョブ管理、開始/一時停止/即時実行、実行履歴 & ログ |
| メッセージ通知 | 複数レベルのメッセージカテゴリ、対象ユーザーへの送信、既読ステータス追跡 |
| 站内メッセージ | 個人受信トレイ、表示/削除、個別/一括既読 |
| キャッシュ管理 | リアルタイムキャッシュ照会、キーベースまたは一括クリーンアップ |
| ログインログ | ログイン成功/失敗ログ、アカウント、IP、所在地、デバイス、時刻 |
| 操作ログ | 全リンク監査証跡、操作者、IP、パラメータ & 結果 |
| パーソナルセンター | プロフィール編集、アバター、パスワードリセット、ログイン履歴 |
| Headless API | コンテンツ CRUD の完全な OpenAPI、マルチターミナル対応 |

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
│   ├── admin/                      # 管理画面フロントエンド (Vue3 + Ant Design Vue + Vben Admin)
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

# 全コードの生成 (ent + wire + api + openapi)
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

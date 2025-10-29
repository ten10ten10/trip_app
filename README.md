# trip_app

旅行計画管理WebアプリケーションのバックエンドAPI

## 📋 目次
- [概要](#概要)
- [技術スタック](#技術スタック)
- [実装済み機能](#実装済み機能)
- [プロジェクト構造](#プロジェクト構造)
- [テスト](#テスト)
- [セットアップ](#セットアップ)

## 概要

このプロジェクトは、旅行計画を管理するためのREST APIです。ユーザー認証、旅行情報の管理、スケジュール管理、共有リンク機能を提供します。

## 技術スタック

### バックエンド
- **言語**: Go 1.23+
- **フレームワーク**: Echo v4
- **ORM**: GORM v1.31
- **認証**: JWT (golang-jwt/jwt v5)
- **バリデーション**: go-playground/validator v10

### データベース
- **RDBMS**: PostgreSQL 15
- **ドライバ**: pgx v5

### API設計
- **仕様**: OpenAPI 3.0
- **コード生成**: oapi-codegen

### テスト
- **フレームワーク**: testify (mock, assert)
- **手法**: E2Eシナリオテスト

### インフラ・ツール
- **コンテナ**: Docker / Docker Compose
- **メール送信**: gomail v2
- **環境変数管理**: godotenv
- **セキュリティ**: bcrypt (golang.org/x/crypto)

## 実装済み機能

### ✅ 全23エンドポイント実装完了

#### ユーザー認証系 (6エンドポイント)
- `POST /signup` - ユーザー登録
- `POST /login` - ログイン
- `POST /logout` - ログアウト
- `POST /users/verify/{verificationToken}` - メール認証
- `GET /me` - 自分の情報取得
- `PUT /me/password` - パスワード変更

#### 旅行管理（要認証） (6エンドポイント)
- `GET /trips` - 旅行一覧取得
- `POST /trips` - 旅行作成
- `GET /trips/{tripId}` - 旅行詳細取得
- `PUT /trips/{tripId}` - 旅行更新
- `DELETE /trips/{tripId}` - 旅行削除
- `GET /trips/{tripId}/details` - 旅行詳細（スケジュール含む）取得

#### スケジュール管理（要認証） (5エンドポイント)
- `GET /trips/{tripId}/schedules` - スケジュール一覧取得
- `POST /trips/{tripId}/schedules` - スケジュール作成
- `GET /trips/{tripId}/schedules/{scheduleId}` - スケジュール詳細取得
- `PATCH /trips/{tripId}/schedules/{scheduleId}` - スケジュール更新
- `DELETE /trips/{tripId}/schedules/{scheduleId}` - スケジュール削除

#### 共有リンク (1エンドポイント)
- `POST /trips/{tripId}/share` - 共有リンク作成

#### 旅行情報（認証不要） (3エンドポイント)
- `GET /public/trips/{shareToken}` - 共有旅行情報取得
- `PUT /public/trips/{shareToken}` - 共有旅行情報更新
- `GET /public/trips/{shareToken}/details` - 共有旅行詳細取得

#### スケジュール管理（認証不要） (5エンドポイント)
- `GET /public/trips/{shareToken}/schedules` - 共有スケジュール一覧取得
- `POST /public/trips/{shareToken}/schedules` - 共有スケジュール作成
- `GET /public/trips/{shareToken}/schedules/{scheduleId}` - 共有スケジュール詳細取得
- `PATCH /public/trips/{shareToken}/schedules/{scheduleId}` - 共有スケジュール更新
- `DELETE /public/trips/{shareToken}/schedules/{scheduleId}` - 共有スケジュール削除

## プロジェクト構造

```
trip_app/
├── api/                      # OpenAPI仕様と自動生成コード
│   ├── openapi.yaml         # API仕様書
│   ├── server.gen.go        # 自動生成されたサーバーインターフェース
│   └── types.gen.go         # 自動生成された型定義
├── cmd/
│   └── app/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── domain/              # ドメインモデル
│   ├── handler/             # HTTPハンドラー層
│   ├── infrastructure/      # インフラ層（メール送信など）
│   ├── middleware/          # ミドルウェア
│   ├── repository/          # リポジトリ層（データアクセス）
│   ├── security/            # セキュリティ関連（JWT、パスワードハッシュなど）
│   └── usecase/             # ユースケース層（ビジネスロジック）
├── docker-compose.yml       # Docker構成
├── Dockerfile              # Dockerイメージ定義
└── go.mod                  # Go依存関係管理
```

## アーキテクチャ

このプロジェクトは**クリーンアーキテクチャ**を採用しています：

```
Handler層 → Usecase層 → Repository層 → Database
   ↓           ↓            ↓
Validator  Validator    GORM/PostgreSQL
```

### レイヤー責務

- **Handler層**: HTTPリクエスト/レスポンスの処理、入力バリデーション
- **Usecase層**: ビジネスロジック、トランザクション管理
- **Repository層**: データベースアクセス
- **Domain層**: ドメインモデル定義

### 特徴的な設計

1. **共有リンク機能**
   - `ShareTokenOwnershipMiddleware`でトークン検証
   - 認証不要エンドポイントで旅行情報を共有

2. **依存性注入**
   - コンストラクタインジェクションを使用
   - テストしやすい設計

3. **バリデーション二重化**
   - Handler層: 入力形式チェック
   - Usecase層: ビジネスルールチェック

## テスト

### ✅ E2Eシナリオテスト（全6シナリオ成功）

全23エンドポイントを網羅する統合テストを実装済み。

#### 実装済みシナリオ

1. **基本ユーザーフロー** - ユーザー登録 → メール認証 → ログイン
2. **旅行管理フロー** - 旅行のCRUD操作
3. **スケジュール管理フロー** - スケジュールのCRUD操作
4. **共有リンクフロー** - 共有リンク機能
5. **認可フロー** - アクセス制御
6. **パスワード変更フロー** - パスワード変更機能

#### テスト方針

- ✅ **E2Eシナリオテストのみ採用**
- ❌ **ユニットテストは実装しない**
- 理由: 開発効率重視、実際のユースケースに基づいたテスト

**詳細は [test/README.md](test/README.md) を参照**

### テスト設計方針

このプロジェクトでは**E2Eシナリオテストのみ**を採用しています：

- ✅ **実装済み**: 全機能の統合テスト
- ❌ **実装しない**: ユニットテスト（各層の単体テスト）

**理由**: 
- シナリオテストで全エンドポイントの動作を保証
- 開発効率重視（ポートフォリオプロジェクト）
- 実際のユースケースに基づいたテスト

## セットアップ

### 前提条件
- Go 1.21以上
- Docker & Docker Compose
- PostgreSQL 15以上（Dockerで起動可）

### 環境変数

`.env`ファイルを作成し、以下の環境変数を設定してください：

```env
DATABASE_URL=postgres://user:password@localhost:5432/trip_app?sslmode=disable
JWT_SECRET=your-secret-key
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=your-email@gmail.com
```

### 起動手順

```bash
# 依存関係のインストール
go mod download

# データベース起動（Docker）
docker compose up -d

# アプリケーション起動
go run cmd/app/main.go
```

サーバーは `http://localhost:8080` で起動します。

## ライセンス

MIT License

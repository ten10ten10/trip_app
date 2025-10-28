# テストガイド

trip_appのテストに関する全体的な情報をまとめています。

## 📁 ディレクトリ構成

```
test/
├── README.md             # このファイル（テスト全体のガイド）
├── e2e/                  # E2Eシナリオテスト
│   └── scenario_test.go # 6つの主要シナリオテスト
└── mock/                 # モック実装
    └── email_sender.go  # メール送信モック
```

## 🧪 テスト方針

このプロジェクトでは**E2Eシナリオテストのみ**を採用しています：

- ✅ **実装済み**: 全機能の統合テスト
- ❌ **実装しない**: ユニットテスト（各層の単体テスト）

**理由**:
- シナリオテストで全エンドポイントの動作を保証
- 開発効率重視（ポートフォリオプロジェクト）
- 実際のユースケースに基づいたテスト

## 📝 テストシナリオ一覧

### 1. TestScenario_BasicUserFlow
基本的なユーザー登録・認証フロー
- ユーザー登録 → メール認証 → ログイン → 自分の情報取得

### 2. TestScenario_TripManagementFlow
旅行管理の基本操作
- 旅行作成 → 一覧取得 → 詳細取得 → 更新 → 削除

### 3. TestScenario_ScheduleManagementFlow
スケジュール管理の基本操作
- スケジュール作成 → 一覧取得 → 詳細取得 → 更新 → 削除

### 4. TestScenario_ShareLinkFlow
共有リンク機能のテスト
- 共有リンク作成 → 認証なしで旅行取得 → スケジュール操作

### 5. TestScenario_AuthorizationFlow
認可機能のテスト
- 他人の旅行へのアクセス拒否確認

### 6. TestScenario_PasswordChangeFlow
パスワード変更機能のテスト
- パスワード変更 → 古いパスワード無効化確認

## 🚀 テスト実行方法

### 1. データベースの起動

```bash
# プロジェクトルートで実行
docker compose up -d db
```

### 2. 環境変数の設定

テストは以下のデータベース接続情報を使用します：

```bash
# .envファイルの設定に合わせて環境変数をエクスポート
export DATABASE_URL="postgres://<USER>:<PASSWORD>@localhost:<PORT>/<DATABASE>?sslmode=disable"
```

**注意**:
- `.env`ファイルに記載されている認証情報（`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`）と一致させてください
- `docker-compose.yml`で設定したポート（デフォルトは`15432`）を使用してください
- 環境変数が設定されていない場合、デフォルトで `postgres://user:password@localhost:5432/trip_app?sslmode=disable` が使用されます

### 3. テストの実行

```bash
# プロジェクトルートから全テスト実行
export DATABASE_URL="postgres://<USER>:<PASSWORD>@localhost:<PORT>/<DATABASE>?sslmode=disable"
go test ./test/e2e -v -count=1

# 特定のテストのみ実行
go test ./test/e2e -v -run TestScenario_BasicUserFlow -count=1

# カバレッジを確認
go test ./test/e2e -coverprofile=coverage.out -count=1
go tool cover -html=coverage.out
```

**オプション説明:**
- `-v`: 詳細な出力
- `-count=1`: キャッシュを無効化して常に実行
- `-run <pattern>`: 特定のテストパターンのみ実行

### 実行例

```bash
$ export DATABASE_URL="postgres://<USER>:<PASSWORD>@localhost:15432/<DATABASE>?sslmode=disable"
$ go test ./test/e2e -v -count=1

=== RUN   TestScenario_BasicUserFlow
--- PASS: TestScenario_BasicUserFlow (0.47s)
=== RUN   TestScenario_TripManagementFlow
--- PASS: TestScenario_TripManagementFlow (0.56s)
=== RUN   TestScenario_ScheduleManagementFlow
--- PASS: TestScenario_ScheduleManagementFlow (0.53s)
=== RUN   TestScenario_ShareLinkFlow
--- PASS: TestScenario_ShareLinkFlow (0.63s)
=== RUN   TestScenario_AuthorizationFlow
--- PASS: TestScenario_AuthorizationFlow (0.68s)
=== RUN   TestScenario_PasswordChangeFlow
--- PASS: TestScenario_PasswordChangeFlow (0.94s)
PASS
ok      trip_app/test/e2e    3.663s
```

## 🔧 テスト設計の特徴

- **実際のHTTPリクエスト**: `httptest.ResponseRecorder`を使用
- **DBを使った統合テスト**: 実際のPostgreSQLデータベースに接続
- **各テスト独立**: 各テスト前後でDBをクリーンアップ
- **モックメール送信**: 実際のメール送信をスキップしつつトークンを取得
- **完全なフロー**: ユーザー登録から各機能の操作まで実際のシナリオを再現

### シナリオテストとは

シナリオテストは、ユーザーの実際の操作フローを再現する統合テストです。

**ユニットテストとの違い:**
- **ユニットテスト**: 各関数・メソッドを個別にテスト
- **シナリオテスト**: 複数の機能を組み合わせた一連の流れをテスト

**メリット:**
- 実際のユースケースに基づいたテスト
- エンドポイント間の連携も確認可能
- リファクタリングに強い（内部実装が変わっても動作が同じならOK）

## 🐛 トラブルシューティング

### データベース接続エラー

```
dial tcp 127.0.0.1:5432: connect: connection refused
```

**解決方法**: PostgreSQLが起動していません。`docker compose up -d db` を実行してください。

### 認証エラー

```
FATAL: password authentication failed for user "user"
```

**解決方法**: `DATABASE_URL`の認証情報が`.env`ファイルと一致していません。環境変数を確認してください。

### テーブルが存在しないエラー

**解決方法**: `setupTestDB`関数で自動マイグレーションが実行されます。テーブルは自動作成されるため、手動での対応は不要です。

### ポート競合エラー

```
bind: address already in use
```

**解決方法**: 既に別のPostgreSQLが起動しています。`docker compose down`で停止するか、`.env`と`docker-compose.yml`のポート設定を変更してください。

## 🎯 テストデータ

テストで使用されるサンプルデータ：

- **旅行**:
  - 沖縄旅行 (2025-10-01 ~ 2025-10-05)
  - 京都旅行 (2025-10-10 ~ 2025-10-12)
  - 博多旅行 (2025-10-15 ~ 2025-10-20)

- **スケジュール例**:
  - 清水寺観光
  - 太宰府天満宮観光
  - 中洲屋台巡り

- **ユーザー**:
  - testuser / test@example.com
  - tripuser / trip@example.com
  - scheduleuser / schedule@example.com
  - など

## 📈 テストカバレッジ

現在のE2Eシナリオテストで以下をカバー：

- ✅ 全23エンドポイント
- ✅ ユーザー認証フロー（登録、認証、ログイン、パスワード変更）
- ✅ 旅行管理（CRUD操作）
- ✅ スケジュール管理（CRUD操作）
- ✅ 共有リンク機能
- ✅ 認可機能（アクセス制御）

## 🔄 CI/CDでの実行

GitHub Actionsでの自動テスト実行は、今後実装予定です。

```yaml
# .github/workflows/test.yml (予定)
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: test_user
          POSTGRES_PASSWORD: test_pass
          POSTGRES_DB: test_db
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go test ./test/e2e -v -count=1
```

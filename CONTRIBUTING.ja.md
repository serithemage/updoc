# updoc コントリビューションガイド

[English](CONTRIBUTING.md) | [한국어](CONTRIBUTING.ko.md)

updocへの貢献に興味を持っていただきありがとうございます！このドキュメントは、貢献のためのガイドラインと手順を提供します。

## 目次

- [行動規範](#行動規範)
- [はじめに](#はじめに)
- [開発環境のセットアップ](#開発環境のセットアップ)
- [貢献方法](#貢献方法)
- [プルリクエストプロセス](#プルリクエストプロセス)
- [コーディング規約](#コーディング規約)
- [コミットメッセージ規約](#コミットメッセージ規約)
- [テスト](#テスト)
- [ドキュメント](#ドキュメント)

## 行動規範

このプロジェクトは、すべての貢献者が守るべき行動規範に従います。相互作用において、敬意を払い、建設的な態度を維持してください。

## はじめに

### 前提条件

- Go 1.21以上
- Git
- Make（Makefileコマンド使用時、オプション）
- golangci-lint（リンティング用）

### ForkとClone

1. GitHubでリポジトリをForkします
2. Forkしたリポジトリをクローンします：
   ```bash
   git clone https://github.com/YOUR_USERNAME/updoc.git
   cd updoc
   ```
3. upstreamリモートを追加します：
   ```bash
   git remote add upstream https://github.com/serithemage/updoc.git
   ```

## 開発環境のセットアップ

### クイックセットアップ

```bash
# 開発依存関係のインストールとGitフックのセットアップ
make dev-setup
```

### 手動セットアップ

```bash
# golangci-lintのインストール
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# プロジェクトのビルド
go build -o updoc ./cmd/updoc

# テストの実行
go test ./...
```

### プロジェクト構造

```
updoc/
├── cmd/updoc/           # アプリケーションエントリーポイント
├── internal/
│   ├── api/             # Upstage APIクライアント
│   ├── cmd/             # CLIコマンド実装
│   ├── config/          # 設定管理
│   └── output/          # 出力フォーマッター
├── test/e2e/            # E2Eテスト
├── docs/                # ドキュメント
├── .github/             # GitHubテンプレートとワークフロー
└── Makefile             # ビルド自動化
```

## 貢献方法

### バグレポート

バグレポートを作成する前に：
1. 重複を避けるため、既存のイシューを確認してください
2. 関連情報を収集してください：
   - updocバージョン（`updoc version`）
   - OSとバージョン
   - 再現手順
   - 期待される動作 vs 実際の動作
   - エラーメッセージまたはログ

イシュー作成時は[Bug Reportテンプレート](.github/ISSUE_TEMPLATE/bug_report.md)を使用してください。

### 機能リクエスト

機能の提案を歓迎します！提出する前に：
1. その機能がすでにリクエストされていないか確認してください
2. プロジェクトの目標と一致するか検討してください
3. 明確なユースケースを提供してください

イシュー作成時は[Feature Requestテンプレート](.github/ISSUE_TEMPLATE/feature_request.md)を使用してください。

### コードの貢献

1. 作業する**イシューを見つける**か、議論のために新しいイシューを作成します
2. 他の人に作業中であることを知らせるためにイシューに**コメント**します
3. `main`から**ブランチを作成**します：
   ```bash
   git checkout -b feature/your-feature-name
   # または
   git checkout -b fix/bug-description
   ```
4. コーディング規約に従って**変更を行います**
5. 新機能の**テストを作成**します
6. **テストとリンティングを実行**します：
   ```bash
   make test
   make lint
   ```
7. コミット規約に従って**コミット**します
8. Forkに**Push**してプルリクエストを作成します

## プルリクエストプロセス

### 提出前の確認

- [ ] コードがエラーなくコンパイルされる
- [ ] すべてのテストが通過する（`make test`）
- [ ] リンティングが通過する（`make lint`）
- [ ] 必要に応じてドキュメントが更新されている
- [ ] コミットメッセージが規約に従っている

### PRガイドライン

1. **タイトル**: 明確で説明的なタイトルを使用
2. **説明**: 何をなぜ変更したか説明
3. **イシューをリンク**: 関連イシューを参照（例：「Fixes #123」）
4. **焦点を絞る**: 1つのPRは1つの問題を扱う
5. **迅速に対応**: レビューフィードバックに迅速に対応

### レビュープロセス

1. 自動化されたチェックの通過が必要（CI/CD）
2. 少なくとも1人のメンテナーの承認が必要
3. すべての議論が解決されている必要がある
4. ブランチが`main`と最新状態である必要がある

## コーディング規約

### Goコードスタイル

- [Effective Go](https://golang.org/doc/effective_go.html)に従う
- `gofmt`でフォーマット
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)に従う

### ガイドライン

- 関数を小さく焦点を絞って維持
- 説明的な変数名と関数名を使用
- エクスポートされる関数と複雑なロジックにコメントを追加
- エラーを明示的に処理し、無視しない
- グローバル状態を避ける

### 例

```go
// ParseDocumentはドキュメントファイルを解析し、構造化されたコンテンツを返します。
// ファイル形式がサポートされていないか、解析に失敗した場合はエラーを返します。
func ParseDocument(filePath string, opts ...Option) (*Result, error) {
    if filePath == "" {
        return nil, errors.New("file path cannot be empty")
    }

    // ... 実装
}
```

## コミットメッセージ規約

[Conventional Commits](https://www.conventionalcommits.org/)形式に従います：

### フォーマット

```
<type>(<scope>): <description>

[オプションの本文]

[オプションのフッター]
```

### タイプ

| タイプ | 説明 |
|--------|------|
| `feat` | 新機能 |
| `fix` | バグ修正 |
| `docs` | ドキュメントのみ |
| `style` | コードスタイル（フォーマットなど） |
| `refactor` | コードリファクタリング |
| `test` | テストの追加または更新 |
| `chore` | メンテナンスタスク |
| `perf` | パフォーマンス改善 |
| `ci` | CI/CDの変更 |

### 例

```bash
feat(parse): HWPXファイル形式のサポートを追加

fix(config): 環境変数からのAPIキー読み込みの問題を解決

docs(readme): インストール手順を更新

test(api): 非同期解析のユニットテストを追加
```

## テスト

### テストの実行

```bash
# ユニットテスト
make test

# E2Eテスト（UPSTAGE_API_KEYが必要）
export UPSTAGE_API_KEY="your-api-key"
make test-e2e

# カバレッジ付きの全テスト
go test -cover ./...
```

### テストの作成

- `*_test.go`ファイルにテストを配置
- 適切な場所でテーブル駆動テストを使用
- 外部依存関係をモック
- 数値だけでなく意味のあるカバレッジを目標に

### テスト例

```go
func TestParseRequest_Validate(t *testing.T) {
    tests := []struct {
        name    string
        req     *ParseRequest
        wantErr bool
    }{
        {
            name:    "有効なリクエスト",
            req:     &ParseRequest{FilePath: "test.pdf"},
            wantErr: false,
        },
        {
            name:    "空のファイルパス",
            req:     &ParseRequest{FilePath: ""},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.req.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## ドキュメント

### ドキュメントの更新タイミング

- 新機能やコマンドの追加時
- 既存の動作の変更時
- 新しい設定オプションの追加時
- 不明確または誤ったドキュメントの修正時

### ドキュメントファイル

| ファイル | 目的 |
|----------|------|
| `README.md` | プロジェクト概要とクイックスタート |
| `docs/CLI_MANUAL.md` | 詳細なCLIリファレンス |
| `CONTRIBUTING.md` | 貢献ガイドライン |

### 多言語サポート

このプロジェクトは英語、韓国語、日本語でドキュメントを管理しています。ドキュメントの更新時：

1. まず英語版を更新
2. `/translate-docs`を使用して翻訳を同期、または
3. 一貫性を保ちながら手動で翻訳を更新

## 質問がありますか？

- 質問は[Discussion](https://github.com/serithemage/updoc/discussions)を開いてください
- まず既存のイシューと議論を確認してください
- 明確に記述し、コンテキストを提供してください

updocへの貢献ありがとうございます！

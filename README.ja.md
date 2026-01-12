# updoc

[![Go Version](https://img.shields.io/github/go-mod/go-version/serithemage/updoc)](https://go.dev/)
[![CI](https://github.com/serithemage/updoc/actions/workflows/ci.yaml/badge.svg)](https://github.com/serithemage/updoc/actions/workflows/ci.yaml)
[![Release](https://img.shields.io/github/v/release/serithemage/updoc)](https://github.com/serithemage/updoc/releases)
[![License](https://img.shields.io/github/license/serithemage/updoc)](LICENSE)

[English](README.md) | [한국어](README.ko.md)

Upstage Document Parse APIのCLIツールです。

## 概要

`updoc`は、PDF、画像、Officeドキュメントなどを構造化されたテキスト（HTML、Markdown、Text）に変換するUpstage Document Parse APIのコマンドラインインターフェースです。Go言語で作成され、単一バイナリとして配布され、クロスプラットフォームをサポートします。

### 主な機能

- 多様なドキュメント形式をサポート：PDF、DOCX、PPTX、XLSX、HWP、HWPX
- 画像/スキャンドキュメントのOCR処理：JPEG、PNG、BMP、TIFF、HEIC
- 複数の出力形式：HTML、Markdown、Text、JSON
- 要素別（見出し、段落、表、図など）の構造化された結果
- バッチ処理およびディレクトリの再帰的探索をサポート
- 同期/非同期処理をサポート（最大1,000ページ）
- 単一バイナリ、外部依存関係なし

## インストール

Go 1.21以上が必要です。

```bash
go install github.com/serithemage/updoc/cmd/updoc@latest
```

またはソースからビルド：

```bash
git clone https://github.com/serithemage/updoc.git
cd updoc
make build
```

## クイックスタート

### 1. APIキーの設定

[Upstage Console](https://console.upstage.ai)でAPIキーを取得し、環境変数として設定します。

```bash
export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"
```

または設定コマンドを使用：

```bash
updoc config set api-key up_xxxxxxxxxxxxxxxxxxxx
```

### プライベートエンドポイント設定（オプション）

AWS Bedrockやプライベートホスティング環境を使用する場合、カスタムエンドポイントを設定できます。

```bash
# 環境変数で設定
export UPSTAGE_API_ENDPOINT="https://your-private-endpoint.com/v1"

# または設定コマンドを使用
updoc config set endpoint https://your-private-endpoint.com/v1

# またはコマンドオプションで指定
updoc parse document.pdf --endpoint https://your-private-endpoint.com/v1
```

### 2. ドキュメントの解析

```bash
# PDFをMarkdownに変換（デフォルト）
updoc parse document.pdf

# 結果をファイルに保存
updoc parse document.pdf -o result.md

# HTML形式に変換
updoc parse document.pdf -f html -o result.html
```

## 使い方

### 基本的な解析

```bash
# 標準出力に結果を出力
updoc parse report.pdf

# ファイルに保存
updoc parse report.pdf -o report.md

# 出力形式を指定：markdown（デフォルト）、html、text、json
updoc parse report.pdf -f html -o report.html
```

### 解析モード

| モード | 説明 | 用途 |
|--------|------|------|
| `standard` | 高速処理（デフォルト） | シンプルなレイアウトのドキュメント |
| `enhanced` | 精密分析 | 複雑な表、チャート、スキャンドキュメント |
| `auto` | 自動選択 | ドキュメントの特性に応じて自動決定 |

```bash
# 複雑な表やチャートがあるドキュメント
updoc parse financial-report.pdf --mode enhanced

# スキャンされたドキュメント（OCRを強制適用）
updoc parse scanned.pdf --ocr force --mode enhanced
```

### バッチ処理

```bash
# 複数のファイルを一度に処理
updoc parse *.pdf --output-dir ./results/

# ディレクトリ内のすべてのドキュメントを再帰的に処理
updoc parse ./documents/ --output-dir ./results/ --recursive

# 特定のパターンに一致するファイルのみ処理
updoc parse ./docs/**/*.pdf --output-dir ./converted/
```

### 詳細オプション

```bash
# チャートを表に変換
updoc parse report.pdf --chart-recognition

# 複数ページにまたがるテーブルを結合
updoc parse spreadsheet.pdf --merge-tables

# 座標情報を含める
updoc parse document.pdf --coordinates

# 要素のみ出力（全体の内容を除外）
updoc parse document.pdf --elements-only

# API応答全体をJSONで出力
updoc parse document.pdf --json -o result.json
```

### 非同期処理（大容量ドキュメント）

100ページを超える大容量ドキュメントは非同期APIを使用します。

```bash
# 非同期リクエストを開始
updoc parse large-document.pdf --async
# 出力: Request ID: req_abc123def456

# ステータスを確認
updoc status req_abc123def456

# リアルタイムでステータスを監視
updoc status req_abc123def456 --watch

# 結果を取得
updoc result req_abc123def456 -o output.md

# 完了まで待機してから結果を取得
updoc result req_abc123def456 --wait -o output.md
```

### 設定管理

```bash
# 現在の設定を確認
updoc config list

# デフォルトの出力形式を変更
updoc config set default-format html

# 利用可能なモデルを確認
updoc models
```

## コマンド概要

| コマンド | 説明 |
|----------|------|
| `updoc parse <file>` | ドキュメントを解析 |
| `updoc status <id>` | 非同期リクエストのステータスを確認 |
| `updoc result <id>` | 非同期リクエストの結果を取得 |
| `updoc config` | 設定を管理 |
| `updoc models` | 利用可能なモデル一覧 |
| `updoc version` | バージョン情報 |

詳細なオプションと使い方については、[CLIマニュアル](docs/CLI_MANUAL.md)を参照してください。

## サポートされるファイル形式

| カテゴリ | 形式 |
|----------|------|
| ドキュメント | PDF、DOCX、PPTX、XLSX、HWP、HWPX |
| 画像 | JPEG、PNG、BMP、TIFF、HEIC |

## API制限

| 項目 | 同期API | 非同期API |
|------|---------|-----------|
| 最大ページ数 | 100 | 1,000 |
| 推奨用途 | 小規模ドキュメント | 大容量ドキュメント、バッチ処理 |

## コントリビューション

プロジェクトへの貢献をありがとうございます！詳細は[コントリビューションガイド](CONTRIBUTING.ja.md)を参照してください。

### 開発環境のセットアップ

```bash
# リポジトリをクローン
git clone https://github.com/serithemage/updoc.git
cd updoc

# 開発環境をセットアップ（Gitフック、リンターのインストール）
make dev-setup

# ビルド
make build

# テストを実行
make test

# E2Eテスト（APIキーが必要）
export UPSTAGE_API_KEY="your-api-key"
make test-e2e

# リント
make lint
```

### 貢献方法

1. 既存のイシューを確認するか、新しいイシューを作成
2. リポジトリをフォーク
3. 機能ブランチを作成（`git checkout -b feature/amazing-feature`）
4. 変更をコミット（`git commit -m 'feat: Add amazing feature'`）
5. ブランチにプッシュ（`git push origin feature/amazing-feature`）
6. プルリクエストを作成

### コミットメッセージの規則

[Conventional Commits](https://www.conventionalcommits.org/)形式に従います：

- `feat:` 新機能
- `fix:` バグ修正
- `docs:` ドキュメントの変更
- `test:` テストの追加/修正
- `refactor:` リファクタリング
- `chore:` その他の変更

### プロジェクト構造

```
updoc/
├── cmd/updoc/           # エントリーポイント
├── internal/
│   ├── api/             # Upstage APIクライアント
│   ├── cmd/             # CLIコマンド実装
│   ├── config/          # 設定管理
│   └── output/          # 出力フォーマッター
├── test/e2e/            # E2Eテスト
├── docs/                # ドキュメント
└── Makefile
```

## ライセンス

MIT License

## 参考資料

- [CLIマニュアル](docs/CLI_MANUAL.md)
- [Upstage Document Parse公式ドキュメント](https://console.upstage.ai/docs/capabilities/document-parse)
- [Upstage APIリファレンス](https://console.upstage.ai/api-reference)
- [Upstage Console](https://console.upstage.ai)

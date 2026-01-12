# updoc CLIマニュアル

[English](CLI_MANUAL.md) | [한국어](CLI_MANUAL.ko.md)

## 目次

1. [はじめに](#はじめに)
2. [インストール](#インストール)
3. [設定](#設定)
4. [コマンド](#コマンド)
5. [使用例](#使用例)
6. [トラブルシューティング](#トラブルシューティング)
7. [APIリファレンス](#apiリファレンス)

---

## はじめに

`updoc`は、UpstageのDocument Parse APIをコマンドラインで使用できるCLIツールです。Go言語で作成され、単一バイナリとして配布され、様々な形式のドキュメントを構造化されたテキスト（HTML、Markdown、Text）に変換します。

### サポート機能

| 機能 | 説明 |
|------|------|
| ドキュメント変換 | PDF、Office、HWPをHTML/Markdown/Textに変換 |
| OCR | スキャンドキュメントや画像からテキストを抽出 |
| 構造分析 | 見出し、段落、表、図などの要素を分離 |
| レイアウト認識 | 複数カラムや複雑なレイアウトを処理 |
| 座標抽出 | 各要素のページ内位置情報を提供 |

---

## インストール

### 要件

- OS: macOS、Linux、Windows
- ビルド時: Go 1.21以上

### バイナリダウンロード

[Releases](https://github.com/serithemage/updoc/releases)ページからOSに合ったバイナリをダウンロードします。

#### macOS

```bash
# Apple Silicon (M1/M2/M3)
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-darwin-arm64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/

# Intel
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-darwin-amd64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/
```

#### Linux

```bash
# amd64
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-linux-amd64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/

# arm64
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-linux-arm64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/
```

#### Windows

```powershell
# PowerShell
Invoke-WebRequest -Uri https://github.com/serithemage/updoc/releases/latest/download/updoc-windows-amd64.exe -OutFile updoc.exe

# PATHに追加するか、任意の場所に移動
Move-Item updoc.exe C:\Users\$env:USERNAME\bin\
```

### Homebrew (macOS/Linux)

```bash
brew install serithemage/tap/updoc
```

### Goでインストール

```bash
go install github.com/serithemage/updoc@latest
```

### ソースからビルド

```bash
git clone https://github.com/serithemage/updoc.git
cd updoc

# ビルド
go build -o updoc ./cmd/updoc

# インストール（オプション）
sudo mv updoc /usr/local/bin/
```

### インストールの確認

```bash
updoc version
updoc --help
```

---

## 設定

### APIキーの設定

Document Parse APIを使用するには、Upstage APIキーが必要です。

#### 1. APIキーの取得

1. [Upstage Console](https://console.upstage.ai)にログイン
2. 新しいプロジェクトを作成するか、既存のプロジェクトを選択
3. API Keysメニューで新しいキーを生成
4. 生成されたキーをコピー

#### 2. キー設定方法

**方法A: 環境変数（推奨）**

```bash
# Linux/macOS
export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"

# Windows (PowerShell)
$env:UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"

# Windows (CMD)
set UPSTAGE_API_KEY=up_xxxxxxxxxxxxxxxxxxxx
```

永続的な設定のためにシェル設定ファイルに追加：

```bash
# ~/.bashrcまたは~/.zshrcに追加
echo 'export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"' >> ~/.zshrc
source ~/.zshrc
```

**方法B: 設定コマンド**

```bash
updoc config set api-key up_xxxxxxxxxxxxxxxxxxxx
```

**方法C: コマンドオプション**

```bash
updoc parse document.pdf --api-key up_xxxxxxxxxxxxxxxxxxxx
```

### プライベートエンドポイント設定

AWS Bedrock、プライベートホスティングなどのカスタムエンドポイントを使用する場合は、以下の方法で設定します。

**方法A: 環境変数**

```bash
export UPSTAGE_API_ENDPOINT="https://your-private-endpoint.com/v1"
```

**方法B: 設定コマンド**

```bash
updoc config set endpoint https://your-private-endpoint.com/v1
```

**方法C: コマンドオプション**

```bash
updoc parse document.pdf --endpoint https://your-private-endpoint.com/v1
```

優先順位: コマンドオプション > 環境変数 > 設定ファイル > デフォルト

### 設定ファイル

設定ファイルの場所:
- Linux/macOS: `~/.config/updoc/config.yaml`
- Windows: `%APPDATA%\updoc\config.yaml`

```yaml
api_key: "up_xxxxxxxxxxxxxxxxxxxx"
endpoint: ""  # デフォルトを使用する場合は空欄
default_format: markdown
default_mode: standard
default_ocr: auto
output_dir: "./output"
```

### 設定管理

```bash
# 現在の設定を表示
updoc config list

# 設定を変更
updoc config set default-format html
updoc config set default-mode enhanced

# 設定を照会
updoc config get default-format

# 設定をリセット
updoc config reset

# 設定ファイルのパスを表示
updoc config path
```

---

## コマンド

### updoc parse

ドキュメントを解析し、構造化されたテキストに変換します。

```
updoc parse <file> [options]
```

#### 引数

| 引数 | 説明 |
|------|------|
| `<file>` | 解析するドキュメントファイルのパス（必須） |

#### オプション

| オプション | 短縮形 | 説明 | デフォルト |
|------------|--------|------|------------|
| `--format <type>` | `-f` | 出力形式: html, markdown, text | markdown |
| `--output <path>` | `-o` | 出力ファイルパス | stdout |
| `--mode <mode>` | `-m` | 解析モード: standard, enhanced, auto | standard |
| `--model <name>` | | モデル名 | document-parse |
| `--ocr <type>` | | OCR設定: auto, force | auto |
| `--chart-recognition` | | チャートを表に変換 | true |
| `--no-chart-recognition` | | チャート変換を無効化 | |
| `--merge-tables` | | 複数ページにまたがるテーブルを結合 | false |
| `--coordinates` | | 座標情報を含める | true |
| `--no-coordinates` | | 座標情報を除外 | |
| `--elements-only` | `-e` | 要素のみを出力 | false |
| `--json` | `-j` | JSON形式で出力 | false |
| `--async` | `-a` | 非同期処理を使用 | false |
| `--output-dir` | `-d` | バッチ処理時の出力ディレクトリ | . |
| `--recursive` | `-r` | ディレクトリを再帰的に探索 | false |
| `--quiet` | `-q` | 進行メッセージを抑制 | false |
| `--verbose` | `-v` | 詳細出力 | false |
| `--api-key <key>` | | APIキーを指定 | 環境変数 |
| `--endpoint <url>` | | APIエンドポイントURL | デフォルトエンドポイント |

#### 例

```bash
# 基本的な使用方法
updoc parse report.pdf

# HTMLに変換してファイルに保存
updoc parse report.pdf -f html -o report.html

# enhancedモードで複雑なドキュメントを処理
updoc parse complex-form.pdf --mode enhanced

# スキャンドキュメントにOCRを強制適用
updoc parse scanned.pdf --ocr force

# JSON出力
updoc parse document.pdf --json -o result.json

# バッチ処理
updoc parse ./documents/*.pdf --output-dir ./results/
```

---

### updoc status

非同期リクエストのステータスを確認します。

```
updoc status <request-id> [options]
```

#### 引数

| 引数 | 説明 |
|------|------|
| `<request-id>` | 非同期リクエストID（必須） |

#### オプション

| オプション | 短縮形 | 説明 | デフォルト |
|------------|--------|------|------------|
| `--json` | `-j` | JSON形式で出力 | false |
| `--watch` | `-w` | リアルタイムでステータスを監視 | false |
| `--interval` | `-i` | 監視間隔（秒） | 5 |

#### 例

```bash
# ステータスを確認
updoc status abc123def456

# JSON形式で出力
updoc status abc123def456 --json

# リアルタイム監視
updoc status abc123def456 --watch
```

#### 出力例

```
Request ID: abc123def456
Status: processing
Progress: 45%
Pages processed: 45/100
Elapsed time: 1m 23s
```

---

### updoc result

非同期リクエストの結果を取得します。

```
updoc result <request-id> [options]
```

#### 引数

| 引数 | 説明 |
|------|------|
| `<request-id>` | 非同期リクエストID（必須） |

#### オプション

| オプション | 短縮形 | 説明 | デフォルト |
|------------|--------|------|------------|
| `--output <path>` | `-o` | 出力ファイルパス | stdout |
| `--format <type>` | `-f` | 出力形式 | markdown |
| `--wait` | `-w` | 完了まで待機 | false |
| `--timeout <sec>` | `-t` | 待機タイムアウト（秒） | 300 |
| `--json` | `-j` | JSON形式で出力 | false |

#### 例

```bash
# 結果を取得
updoc result abc123def456 -o output.md

# 完了まで待機して結果を取得
updoc result abc123def456 --wait -o output.md

# タイムアウト付き
updoc result abc123def456 --wait --timeout 600 -o output.md
```

---

### updoc models

利用可能なモデルを表示します。

```
updoc models [options]
```

#### オプション

| オプション | 短縮形 | 説明 | デフォルト |
|------------|--------|------|------------|
| `--json` | `-j` | JSON形式で出力 | false |

#### 出力例

```
Available Models:

  document-parse          デフォルトモデル（推奨、エイリアス）
  document-parse-250618   特定バージョン（2025-06-18）
  document-parse-nightly  最新テストバージョン

Tip: 'document-parse'エイリアスを使用すると、自動的に最新の安定版が適用されます。
```

---

### updoc config

設定を管理します。

```
updoc config <command> [key] [value]
```

#### サブコマンド

| コマンド | 説明 |
|----------|------|
| `list` | すべての設定を表示 |
| `get <key>` | 特定の設定を照会 |
| `set <key> <value>` | 設定を変更 |
| `reset` | 設定をリセット |
| `path` | 設定ファイルのパスを表示 |

#### 設定キー

| キー | 説明 | 値 |
|------|------|-----|
| `api-key` | APIキー | 文字列 |
| `endpoint` | APIエンドポイントURL | URL |
| `default-format` | デフォルト出力形式 | html, markdown, text |
| `default-mode` | デフォルト解析モード | standard, enhanced, auto |
| `default-ocr` | デフォルトOCR設定 | auto, force |
| `output-dir` | デフォルト出力ディレクトリ | パス |

#### 例

```bash
# すべての設定を表示
updoc config list

# 特定の設定を照会
updoc config get api-key

# 設定を変更
updoc config set default-format html
updoc config set default-mode enhanced

# 設定をリセット
updoc config reset
```

---

### updoc version

バージョン情報を表示します。

```
updoc version [options]
```

#### オプション

| オプション | 短縮形 | 説明 |
|------------|--------|------|
| `--short` | `-s` | バージョン番号のみ出力 |
| `--json` | `-j` | JSON形式で出力 |

#### 出力例

```
updoc version 1.0.0
  Commit: abc1234
  Built: 2025-01-11T10:00:00Z
  Go version: go1.21.5
  OS/Arch: darwin/arm64
```

---

## 使用例

### 基本的なワークフロー

```bash
# 1. APIキーを設定
export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"

# 2. PDFをMarkdownに変換
updoc parse report.pdf -o report.md

# 3. 結果を確認
cat report.md
```

### スキャンドキュメントの処理

```bash
# スキャンされたPDFをOCR処理
updoc parse scanned-document.pdf --ocr force --mode enhanced -o output.md
```

### 複雑なレイアウトのドキュメント

```bash
# 多くの表やチャートを含むドキュメントをenhancedモードで処理
updoc parse financial-report.pdf \
  --mode enhanced \
  --chart-recognition \
  --merge-tables \
  -o report.html \
  -f html
```

### 要素別分析

```bash
# ドキュメントを要素に分解してJSON出力
updoc parse document.pdf --elements-only --json -o elements.json

# jqで表のみを抽出
cat elements.json | jq '.elements[] | select(.category == "table")'

# jqで見出しのみを抽出
cat elements.json | jq '.elements[] | select(.category | startswith("heading"))'
```

### 大容量ドキュメントの処理

```bash
# 非同期リクエストを開始
updoc parse large-document.pdf --async
# 出力: Request ID: req_abc123

# リアルタイムでステータスを監視
updoc status req_abc123 --watch

# 完了後に結果を取得
updoc result req_abc123 -o result.md

# または完了まで待機
updoc result req_abc123 --wait --timeout 600 -o result.md
```

### バッチ処理

```bash
# 現在のディレクトリのすべてのPDFを処理
updoc parse *.pdf --output-dir ./results/

# ディレクトリ内のドキュメントを再帰的に処理
updoc parse ./documents/ --output-dir ./results/ --recursive

# シェルスクリプトで細かく制御
for file in *.pdf; do
  echo "Processing: $file"
  updoc parse "$file" -o "${file%.pdf}.md" --quiet
done
```

### パイプラインの活用

```bash
# 変換結果を他のツールにパイプ
updoc parse document.pdf | grep -i "important"

# 複数の形式に同時出力
updoc parse document.pdf --json | tee result.json | jq -r '.content.markdown' > result.md

# 特定の要素を抽出して処理
updoc parse document.pdf --json | jq -r '.elements[] | select(.category == "table") | .content.markdown'
```

### 自動化スクリプト例

```bash
#!/bin/bash
# batch_convert.sh - バッチドキュメント変換スクリプト

INPUT_DIR="${1:-.}"
OUTPUT_DIR="${2:-./output}"
FORMAT="${3:-markdown}"

mkdir -p "$OUTPUT_DIR"

find "$INPUT_DIR" -type f \( -name "*.pdf" -o -name "*.docx" -o -name "*.hwp" \) | while read -r file; do
  filename=$(basename "$file")
  output_file="$OUTPUT_DIR/${filename%.*}.${FORMAT}"

  echo "Converting: $filename"
  updoc parse "$file" -f "$FORMAT" -o "$output_file" --quiet

  if [ $? -eq 0 ]; then
    echo "  -> $output_file"
  else
    echo "  -> Failed"
  fi
done

echo "Done!"
```

---

## トラブルシューティング

### よくあるエラー

#### APIキーエラー

```
Error: Invalid API key
```

解決方法:
1. APIキーが正しいか確認
2. 環境変数が設定されているか確認: `echo $UPSTAGE_API_KEY`
3. キー前後の空白を削除
4. 設定を確認: `updoc config get api-key`

#### ファイル形式エラー

```
Error: Unsupported file format: .xyz
```

サポートされている形式:
- ドキュメント: PDF、DOCX、PPTX、XLSX、HWP
- 画像: JPEG、PNG、BMP、TIFF、HEIC

#### ページ数超過

```
Error: Document exceeds maximum page limit (100 pages for sync API)
```

同期APIは最大100ページ、非同期APIは最大1,000ページをサポート。
大容量ドキュメントには`--async`オプションを使用:

```bash
updoc parse large-document.pdf --async
```

#### タイムアウト

```
Error: Request timeout after 120s
```

解決方法:
- 非同期モードを使用: `--async`
- ネットワーク状態を確認
- ファイルサイズを確認

#### ファイルアクセスエラー

```
Error: Cannot read file: permission denied
```

解決方法:
```bash
# ファイル権限を確認
ls -la document.pdf

# 権限を変更（必要に応じて）
chmod 644 document.pdf
```

### デバッグ

```bash
# 詳細出力
updoc parse document.pdf --verbose

# リクエスト/レスポンスを確認
updoc parse document.pdf --verbose 2>&1 | tee debug.log

# 設定を確認
updoc config list
```

### ログレベル

`--verbose`フラグを使用すると、以下の情報が出力されます:
- APIリクエストURLとヘッダー
- リクエストパラメータ
- レスポンスステータスコード
- 処理時間

---

## APIリファレンス

### エンドポイント

| 用途 | URL |
|------|-----|
| 同期解析 | `POST https://api.upstage.ai/v1/document-digitization` |
| 非同期解析 | `POST https://api.upstage.ai/v1/document-digitization/async` |
| ステータス確認 | `GET https://api.upstage.ai/v1/document-digitization/async/{id}` |

### 認証

```
Authorization: Bearer <UPSTAGE_API_KEY>
```

### リクエスト形式

`Content-Type: multipart/form-data`

| フィールド | 型 | 必須 | 説明 |
|------------|------|------|------|
| `model` | string | はい | モデル名（document-parse） |
| `document` | file | はい | ドキュメントファイル |
| `mode` | string | | standard, enhanced, auto |
| `ocr` | string | | auto, force |
| `output_formats` | string | | 出力形式 |
| `chart_recognition` | boolean | | チャート変換 |
| `merge_multipage_tables` | boolean | | テーブル結合 |
| `coordinates` | boolean | | 座標を含める |

### レスポンス構造

```json
{
  "api": "document-parse",
  "model": "document-parse-250618",
  "content": {
    "html": "<h1>Title</h1>...",
    "markdown": "# Title\n...",
    "text": "Title\n..."
  },
  "elements": [
    {
      "id": 1,
      "category": "heading1",
      "page": 1,
      "content": {
        "html": "<h1>Title</h1>",
        "markdown": "# Title",
        "text": "Title"
      },
      "coordinates": [
        {"x": 0.1, "y": 0.05},
        {"x": 0.9, "y": 0.05},
        {"x": 0.9, "y": 0.08},
        {"x": 0.1, "y": 0.08}
      ]
    }
  ],
  "usage": {
    "pages": 10
  }
}
```

### 要素カテゴリ

| カテゴリ | 説明 |
|----------|------|
| `heading1` ~ `heading6` | 見出しレベル |
| `paragraph` | 段落 |
| `table` | 表 |
| `figure` | 図 |
| `chart` | チャート |
| `equation` | 数式 |
| `list_item` | リスト項目 |
| `header` | ヘッダー |
| `footer` | フッター |
| `caption` | キャプション |

---

## 付録

### 環境変数

| 変数 | 説明 |
|------|------|
| `UPSTAGE_API_KEY` | API認証キー |
| `UPSTAGE_API_ENDPOINT` | APIエンドポイントURL（プライベートホスティング用） |
| `UPDOC_CONFIG_PATH` | 設定ファイルパス（オプション） |
| `UPDOC_LOG_LEVEL` | ログレベル: debug, info, warn, error |

### 終了コード

| コード | 意味 |
|--------|------|
| 0 | 成功 |
| 1 | 一般エラー |
| 2 | 引数エラー |
| 3 | APIエラー |
| 4 | ファイルI/Oエラー |
| 5 | 認証エラー |

### 関連リンク

- [Upstage Console](https://console.upstage.ai)
- [Document Parse公式ドキュメント](https://console.upstage.ai/docs/capabilities/document-parse)
- [APIリファレンス](https://console.upstage.ai/api-reference)
- [Upstageブログ](https://upstage.ai/blog)
- [GitHubリポジトリ](https://github.com/serithemage/updoc)

### バージョン履歴

| バージョン | 日付 | 変更内容 |
|------------|------|----------|
| 1.0.0 | - | 初期リリース |

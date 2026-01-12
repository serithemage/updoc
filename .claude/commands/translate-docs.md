# /translate-docs

Synchronize documentation translations across English, Korean, and Japanese.

## Usage
```
/translate-docs [options]
```

## Options
- `--all`: Translate all documentation files (full sync)
- `--check`: Only show what would be translated without making changes
- `--file <path>`: Translate only the specified file

## What it does

1. **Detects changes** in English documentation files:
   - `README.md` → `README.ko.md`, `README.ja.md`
   - `docs/CLI_MANUAL.md` → `docs/CLI_MANUAL.ko.md`, `docs/CLI_MANUAL.ja.md`

2. **Analyzes diffs** to find modified sections

3. **Translates** only the changed portions to:
   - Korean (한국어)
   - Japanese (日本語)

4. **Preserves**:
   - Code blocks (unchanged)
   - Command examples (unchanged)
   - URLs and links (unchanged)
   - Markdown formatting
   - Technical terms

## Examples

```bash
# Translate recent changes
/translate-docs

# Full translation sync
/translate-docs --all

# Check what would be translated
/translate-docs --check

# Translate specific file
/translate-docs --file README.md
```

## Workflow

```
┌─────────────────┐
│  README.md      │ (English - Source)
└────────┬────────┘
         │ git diff
         ▼
┌─────────────────┐
│ Detect Changes  │
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌────────┐ ┌────────┐
│ Korean │ │Japanese│
│ .ko.md │ │ .ja.md │
└────────┘ └────────┘
```

## Translation Guidelines

### Technical Terms (Keep in English)
- `updoc`, `parse`, `config`, `status`, `result`
- `--format`, `--output`, `--async`, `--endpoint`
- `API key`, `endpoint`, `OAuth`, `JWT`
- File extensions: `.pdf`, `.md`, `.html`

### Section Titles

| English | Korean | Japanese |
|---------|--------|----------|
| Overview | 개요 | 概要 |
| Installation | 설치 | インストール |
| Quick Start | 빠른 시작 | クイックスタート |
| Usage | 사용법 | 使い方 |
| Commands | 명령어 | コマンド |
| Options | 옵션 | オプション |
| Examples | 예제 | 例 |
| Troubleshooting | 문제 해결 | トラブルシューティング |

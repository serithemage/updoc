# updoc Project Instructions

## Project Overview

`updoc` is a CLI tool for the Upstage Document Parse API, written in Go.

## Documentation Translation

This project maintains documentation in three languages:
- **English** (default): `README.md`, `docs/CLI_MANUAL.md`
- **Korean**: `README.ko.md`, `docs/CLI_MANUAL.ko.md`
- **Japanese**: `README.ja.md`, `docs/CLI_MANUAL.ja.md`

### Auto-Translation Workflow

When English documentation is modified:

1. **Detect changes** using git diff:
   ```bash
   git diff --name-only | grep -E "^(README\.md|docs/CLI_MANUAL\.md)$"
   ```

2. **Analyze specific changes**:
   ```bash
   git diff README.md
   ```

3. **Apply translations** to corresponding files:
   - `README.md` → `README.ko.md`, `README.ja.md`
   - `docs/CLI_MANUAL.md` → `docs/CLI_MANUAL.ko.md`, `docs/CLI_MANUAL.ja.md`

### Translation Guidelines

#### Do NOT translate:
- Code blocks and command examples
- CLI commands: `updoc`, `parse`, `config`, `status`, `result`
- Flags: `--format`, `--output`, `--async`, `--endpoint`, `-f`, `-o`
- Technical terms: `API key`, `endpoint`, `async`, `sync`
- File paths and extensions
- URLs and links
- Badge images

#### Section title mappings:

| English | Korean | Japanese |
|---------|--------|----------|
| Overview | 개요 | 概要 |
| Installation | 설치 | インストール |
| Quick Start | 빠른 시작 | クイックスタート |
| Usage | 사용법 | 使い方 |
| Configuration | 설정 | 設定 |
| Commands | 명령어 | コマンド |
| Examples | 예제 | 例 |
| Options | 옵션 | オプション |
| Arguments | 인자 | 引数 |
| Troubleshooting | 문제 해결 | トラブルシューティング |
| Contributing | 컨트리뷰션 | コントリビューション |
| License | 라이선스 | ライセンス |
| References | 참고 자료 | 参考資料 |
| API Reference | API 레퍼런스 | APIリファレンス |
| Appendix | 부록 | 付録 |
| Key Features | 주요 기능 | 主な機能 |
| Supported Formats | 지원 파일 형식 | サポートされるファイル形式 |
| Environment Variables | 환경 변수 | 環境変数 |

### Custom Command

Use `/translate-docs` to sync translations:
```
/translate-docs           # Translate recent changes
/translate-docs --all     # Full sync all docs
/translate-docs --check   # Preview without changes
```

## Development

### Build
```bash
make build
```

### Test
```bash
make test        # Unit tests
make test-e2e    # E2E tests (requires UPSTAGE_API_KEY)
```

### Lint
```bash
make lint
```

## Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Test changes
- `refactor:` Refactoring
- `chore:` Other changes

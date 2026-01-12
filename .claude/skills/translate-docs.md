# Translate Documentation Skill

## Description
Automatically translates documentation changes to Korean and Japanese when English documentation is updated.

## Trigger
Use this skill when:
- README.md is modified
- docs/CLI_MANUAL.md is modified
- User requests document translation sync
- User runs `/translate-docs`

## Workflow

### 1. Detect Changes
```bash
# Check for changes in English documentation
git diff --name-only HEAD~1 HEAD | grep -E "^(README\.md|docs/CLI_MANUAL\.md)$"

# Or check uncommitted changes
git diff --name-only | grep -E "^(README\.md|docs/CLI_MANUAL\.md)$"
```

### 2. Analyze Changed Sections
For each changed English document:
```bash
# Get diff for specific file
git diff HEAD~1 HEAD -- README.md
# Or for uncommitted changes
git diff -- README.md
```

### 3. Translation Mapping

| English File | Korean File | Japanese File |
|--------------|-------------|---------------|
| README.md | README.ko.md | README.ja.md |
| docs/CLI_MANUAL.md | docs/CLI_MANUAL.ko.md | docs/CLI_MANUAL.ja.md |

### 4. Translation Rules

#### Header Preservation
- Keep language switcher links unchanged
- Update the language links if file structure changes

#### Technical Terms (Do Not Translate)
- CLI commands: `updoc`, `parse`, `config`, etc.
- Flags: `--format`, `--output`, `-f`, `-o`, etc.
- File extensions: `.pdf`, `.md`, `.html`, etc.
- API terms: `API key`, `endpoint`, `async`, etc.
- Code blocks: Keep all code unchanged
- URLs: Keep unchanged

#### Section Title Translations

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
| Troubleshooting | 문제 해결 | トラブルシューティング |
| Contributing | 컨트리뷰션 | コントリビューション |
| License | 라이선스 | ライセンス |
| References | 참고 자료 | 参考資料 |
| API Reference | API 레퍼런스 | APIリファレンス |
| Appendix | 부록 | 付録 |

### 5. Translation Process

1. Read the changed sections from the English document
2. Identify corresponding sections in Korean/Japanese documents
3. Translate only the changed content
4. Preserve markdown formatting, code blocks, and links
5. Update the translated documents

### 6. Output Format

After translation:
```
Translation Summary:
- README.md changes applied to:
  - README.ko.md (Korean)
  - README.ja.md (Japanese)
- docs/CLI_MANUAL.md changes applied to:
  - docs/CLI_MANUAL.ko.md (Korean)
  - docs/CLI_MANUAL.ja.md (Japanese)
```

## Example Execution

```bash
# User command
/translate-docs

# Or automatically triggered after editing README.md
```

### Sample Translation

**English (README.md):**
```markdown
## Installation

Requires Go 1.21 or later.
```

**Korean (README.ko.md):**
```markdown
## 설치

Go 1.21 이상이 필요합니다.
```

**Japanese (README.ja.md):**
```markdown
## インストール

Go 1.21以上が必要です。
```

## Notes

- Always preserve the language switcher at the top of each document
- Keep code examples identical across all languages
- Maintain consistent formatting (tables, lists, code blocks)
- Do not translate proper nouns (Upstage, GitHub, etc.)
- Keep URLs unchanged
- Badge images should remain identical

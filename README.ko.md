# updoc

[![Go Version](https://img.shields.io/github/go-mod/go-version/serithemage/updoc)](https://go.dev/)
[![CI](https://github.com/serithemage/updoc/actions/workflows/ci.yaml/badge.svg)](https://github.com/serithemage/updoc/actions/workflows/ci.yaml)
[![Release](https://img.shields.io/github/v/release/serithemage/updoc)](https://github.com/serithemage/updoc/releases)
[![License](https://img.shields.io/github/license/serithemage/updoc)](LICENSE)

[English](README.md) | [日本語](README.ja.md)

업스테이지(Upstage) Document Parse API를 CLI로 사용할 수 있는 도구입니다.

## 개요

`updoc`은 PDF, 이미지, 오피스 문서 등을 구조화된 텍스트(HTML, Markdown, Text)로 변환하는 업스테이지 Document Parse API의 커맨드라인 인터페이스입니다. Go 언어로 작성되어 단일 바이너리로 배포되며, 크로스 플랫폼을 지원합니다.

### 주요 기능

- PDF, DOCX, PPTX, XLSX, HWP, HWPX 등 다양한 문서 형식 지원
- JPEG, PNG, BMP, TIFF, HEIC 등 이미지/스캔 문서 OCR 처리
- HTML, Markdown, Text, JSON 출력 형식 선택
- 요소별(제목, 단락, 표, 그림 등) 구조화된 결과 제공
- 배치 처리 및 디렉토리 재귀 탐색 지원
- 동기/비동기 처리 지원 (최대 1,000페이지)
- 단일 바이너리, 외부 의존성 없음

## 설치

Go 1.21 이상이 필요합니다.

```bash
go install github.com/serithemage/updoc/cmd/updoc@latest
```

또는 소스에서 빌드:

```bash
git clone https://github.com/serithemage/updoc.git
cd updoc
make build
```

## 빠른 시작

### 1. API 키 설정

[Upstage Console](https://console.upstage.ai)에서 API 키를 발급받은 후 환경 변수로 설정합니다.

```bash
export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"
```

또는 설정 명령어 사용:

```bash
updoc config set api-key up_xxxxxxxxxxxxxxxxxxxx
```

### 프라이빗 엔드포인트 설정 (선택)

AWS Bedrock이나 프라이빗 호스팅 환경을 사용하는 경우 커스텀 엔드포인트를 설정할 수 있습니다.

```bash
# 환경 변수로 설정
export UPSTAGE_API_ENDPOINT="https://your-private-endpoint.com/v1"

# 또는 설정 명령어 사용
updoc config set endpoint https://your-private-endpoint.com/v1

# 또는 명령어 옵션으로 지정
updoc parse document.pdf --endpoint https://your-private-endpoint.com/v1
```

### 2. 문서 파싱

```bash
# PDF를 Markdown으로 변환 (기본)
updoc parse document.pdf

# 결과를 파일로 저장
updoc parse document.pdf -o result.md

# HTML 형식으로 변환
updoc parse document.pdf -f html -o result.html
```

## 사용법

### 기본 파싱

```bash
# 표준 출력으로 결과 확인
updoc parse report.pdf

# 파일로 저장
updoc parse report.pdf -o report.md

# 출력 형식 지정: markdown (기본), html, text, json
updoc parse report.pdf -f html -o report.html
```

### 파싱 모드

| 모드 | 설명 | 용도 |
|------|------|------|
| `standard` | 빠른 처리 (기본) | 단순한 레이아웃 문서 |
| `enhanced` | 정밀 분석 | 복잡한 표, 차트, 스캔 문서 |
| `auto` | 자동 선택 | 문서 특성에 따라 자동 결정 |

```bash
# 복잡한 표와 차트가 있는 문서
updoc parse financial-report.pdf --mode enhanced

# 스캔된 문서 (OCR 강제 적용)
updoc parse scanned.pdf --ocr force --mode enhanced
```

### 배치 처리

```bash
# 여러 파일 한번에 처리
updoc parse *.pdf --output-dir ./results/

# 디렉토리 내 모든 문서 재귀 처리
updoc parse ./documents/ --output-dir ./results/ --recursive

# 특정 패턴의 파일만 처리
updoc parse ./docs/**/*.pdf --output-dir ./converted/
```

### 고급 옵션

```bash
# 차트를 표로 변환
updoc parse report.pdf --chart-recognition

# 다중 페이지에 걸친 테이블 병합
updoc parse spreadsheet.pdf --merge-tables

# 요소별 좌표 정보 포함
updoc parse document.pdf --coordinates

# 요소별 결과만 출력 (전체 내용 제외)
updoc parse document.pdf --elements-only

# JSON 형태로 전체 API 응답 출력
updoc parse document.pdf --json -o result.json
```

### 비동기 처리 (대용량 문서)

100페이지를 초과하는 대용량 문서는 비동기 API를 사용합니다.

```bash
# 비동기 요청 시작
updoc parse large-document.pdf --async
# 출력: Request ID: req_abc123def456

# 상태 확인
updoc status req_abc123def456

# 실시간 상태 모니터링
updoc status req_abc123def456 --watch

# 결과 가져오기
updoc result req_abc123def456 -o output.md

# 완료까지 대기 후 결과 가져오기
updoc result req_abc123def456 --wait -o output.md
```

### 설정 관리

```bash
# 현재 설정 확인
updoc config list

# 기본 출력 형식 변경
updoc config set default-format html

# 사용 가능한 모델 확인
updoc models
```

## 명령어 요약

| 명령어 | 설명 |
|--------|------|
| `updoc parse <file>` | 문서 파싱 |
| `updoc status <id>` | 비동기 요청 상태 확인 |
| `updoc result <id>` | 비동기 요청 결과 가져오기 |
| `updoc config` | 설정 관리 |
| `updoc models` | 사용 가능한 모델 목록 |
| `updoc version` | 버전 정보 |

상세한 옵션과 사용법은 [CLI 매뉴얼](docs/CLI_MANUAL.md)을 참조하세요.

## 지원 파일 형식

| 분류 | 형식 |
|------|------|
| 문서 | PDF, DOCX, PPTX, XLSX, HWP, HWPX |
| 이미지 | JPEG, PNG, BMP, TIFF, HEIC |

## API 제한

| 항목 | 동기 API | 비동기 API |
|------|----------|------------|
| 최대 페이지 수 | 100 | 1,000 |
| 권장 용도 | 소규모 문서 | 대용량 문서, 배치 처리 |

## 컨트리뷰션

프로젝트에 기여해 주셔서 감사합니다! 자세한 내용은 [기여 가이드](CONTRIBUTING.ko.md)를 참조하세요.

### 개발 환경 설정

```bash
# 저장소 클론
git clone https://github.com/serithemage/updoc.git
cd updoc

# 개발 환경 설정 (Git hooks, 린터 설치)
make dev-setup

# 빌드
make build

# 테스트 실행
make test

# E2E 테스트 (API 키 필요)
export UPSTAGE_API_KEY="your-api-key"
make test-e2e

# 린트
make lint
```

### 기여 방법

1. 이슈를 확인하거나 새 이슈를 생성합니다
2. 저장소를 Fork합니다
3. 기능 브랜치를 생성합니다 (`git checkout -b feature/amazing-feature`)
4. 변경사항을 커밋합니다 (`git commit -m 'feat: Add amazing feature'`)
5. 브랜치에 Push합니다 (`git push origin feature/amazing-feature`)
6. Pull Request를 생성합니다

### 커밋 메시지 규칙

[Conventional Commits](https://www.conventionalcommits.org/) 형식을 따릅니다:

- `feat:` 새로운 기능
- `fix:` 버그 수정
- `docs:` 문서 변경
- `test:` 테스트 추가/수정
- `refactor:` 리팩토링
- `chore:` 기타 변경

### 프로젝트 구조

```
updoc/
├── cmd/updoc/           # 엔트리포인트
├── internal/
│   ├── api/             # Upstage API 클라이언트
│   ├── cmd/             # CLI 명령어 구현
│   ├── config/          # 설정 관리
│   └── output/          # 출력 포매터
├── test/e2e/            # E2E 테스트
├── docs/                # 문서
└── Makefile
```

## 라이선스

MIT License

## 참고 자료

- [CLI 상세 매뉴얼](docs/CLI_MANUAL.md)
- [Upstage Document Parse 공식 문서](https://console.upstage.ai/docs/capabilities/document-parse)
- [Upstage API Reference](https://console.upstage.ai/api-reference)
- [Upstage Console](https://console.upstage.ai)

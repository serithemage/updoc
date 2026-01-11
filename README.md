# updoc

업스테이지(Upstage) Document Parse API를 CLI로 사용할 수 있는 도구입니다.

## 개요

`updoc`은 PDF, 이미지, 오피스 문서 등을 구조화된 텍스트(HTML, Markdown, Text)로 변환하는 업스테이지 Document Parse API의 커맨드라인 인터페이스입니다. Go 언어로 작성되어 단일 바이너리로 배포되며, 크로스 플랫폼을 지원합니다.

### 주요 기능

- PDF, DOCX, PPTX, XLSX, HWP 등 다양한 문서 형식 지원
- JPEG, PNG, BMP, TIFF, HEIC 등 이미지/스캔 문서 OCR 처리
- HTML, Markdown, Text 출력 형식 선택
- 요소별(제목, 단락, 표, 그림 등) 구조화된 결과 제공
- 동기/비동기 처리 지원 (최대 1,000페이지)
- 단일 바이너리, 외부 의존성 없음

## 설치

### 바이너리 다운로드

[Releases](https://github.com/serithemage/updoc/releases) 페이지에서 OS에 맞는 바이너리를 다운로드합니다.

```bash
# macOS (Apple Silicon)
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-darwin-arm64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-darwin-amd64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/

# Linux (amd64)
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-linux-amd64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/serithemage/updoc/releases/latest/download/updoc-windows-amd64.exe -OutFile updoc.exe
```

### Homebrew (macOS/Linux)

```bash
brew install serithemage/tap/updoc
```

### Go로 설치

```bash
go install github.com/serithemage/updoc@latest
```

### 소스에서 빌드

```bash
git clone https://github.com/serithemage/updoc.git
cd updoc
go build -o updoc ./cmd/updoc
```

## 설정

### API 키 발급

1. [Upstage Console](https://console.upstage.ai)에서 계정 생성
2. 프로젝트 생성 후 API 키 발급
3. 환경 변수 또는 설정 파일에 API 키 저장

### 환경 변수 설정

```bash
export UPSTAGE_API_KEY="your-api-key-here"
```

또는 설정 명령어 사용:

```bash
updoc config set api-key your-api-key-here
```

## 사용법

### 기본 사용

```bash
# 문서 파싱 (기본: Markdown 출력)
updoc parse document.pdf

# 출력 형식 지정
updoc parse document.pdf --format html
updoc parse document.pdf --format text
updoc parse document.pdf --format markdown

# 결과를 파일로 저장
updoc parse document.pdf -o result.md
updoc parse document.pdf --output result.html --format html
```

### 파싱 모드 선택

```bash
# Standard 모드 (기본, 단순한 레이아웃 문서용)
updoc parse document.pdf --mode standard

# Enhanced 모드 (복잡한 표, 차트, 스캔 문서용)
updoc parse document.pdf --mode enhanced

# Auto 모드 (자동 선택)
updoc parse document.pdf --mode auto
```

### OCR 설정

```bash
# 자동 OCR (기본, 필요시에만 OCR 적용)
updoc parse scanned.pdf --ocr auto

# 강제 OCR (모든 페이지에 OCR 적용)
updoc parse scanned.pdf --ocr force
```

### 고급 옵션

```bash
# 차트를 표로 변환
updoc parse report.pdf --chart-recognition

# 다중 페이지 테이블 병합
updoc parse spreadsheet.pdf --merge-tables

# 좌표 정보 포함
updoc parse document.pdf --coordinates

# 요소별 결과만 출력
updoc parse document.pdf --elements-only

# JSON 형태로 전체 응답 출력
updoc parse document.pdf --json
```

### 비동기 처리 (대용량 문서)

```bash
# 비동기 요청 시작
updoc parse large-document.pdf --async
# 출력: Request ID: abc123

# 상태 확인
updoc status abc123

# 결과 가져오기
updoc result abc123 -o output.md
```

### 배치 처리

```bash
# 여러 파일 처리
updoc parse *.pdf --output-dir ./results

# 디렉토리 내 모든 문서 처리
updoc parse ./documents/ --output-dir ./results --recursive
```

## 명령어 레퍼런스

| 명령어 | 설명 |
|--------|------|
| `updoc parse <file>` | 문서 파싱 실행 |
| `updoc status <id>` | 비동기 요청 상태 확인 |
| `updoc result <id>` | 비동기 요청 결과 가져오기 |
| `updoc config` | 설정 관리 |
| `updoc models` | 사용 가능한 모델 목록 |
| `updoc version` | 버전 정보 출력 |

## 옵션 레퍼런스

| 옵션 | 단축 | 설명 | 기본값 |
|------|------|------|--------|
| `--format` | `-f` | 출력 형식 (html, markdown, text) | markdown |
| `--output` | `-o` | 출력 파일 경로 | stdout |
| `--mode` | `-m` | 파싱 모드 (standard, enhanced, auto) | standard |
| `--ocr` | | OCR 설정 (auto, force) | auto |
| `--model` | | 사용할 모델 | document-parse |
| `--chart-recognition` | | 차트를 표로 변환 | true |
| `--merge-tables` | | 다중 페이지 테이블 병합 | false |
| `--coordinates` | | 좌표 정보 포함 | true |
| `--elements-only` | `-e` | 요소별 결과만 출력 | false |
| `--json` | `-j` | JSON 형태로 전체 응답 출력 | false |
| `--async` | `-a` | 비동기 처리 사용 | false |
| `--quiet` | `-q` | 진행 메시지 숨김 | false |
| `--verbose` | `-v` | 상세 로그 출력 | false |

## 지원 파일 형식

### 문서
- PDF
- Microsoft Office: DOCX, PPTX, XLSX
- 한글: HWP

### 이미지
- JPEG, PNG, BMP, TIFF, HEIC

## API 제한

| 항목 | 동기 API | 비동기 API |
|------|----------|------------|
| 최대 페이지 수 | 100 | 1,000 |
| 권장 용도 | 소규모 문서 | 대용량 문서, 배치 처리 |

## 출력 예시

### Markdown 출력

```markdown
# 문서 제목

## 1. 서론

문서의 첫 번째 단락입니다.

| 항목 | 값 |
|------|-----|
| A    | 100 |
| B    | 200 |
```

### 요소별 JSON 출력

```json
{
  "elements": [
    {
      "id": 1,
      "category": "heading1",
      "page": 1,
      "content": {
        "text": "문서 제목",
        "markdown": "# 문서 제목"
      }
    },
    {
      "id": 2,
      "category": "paragraph",
      "page": 1,
      "content": {
        "text": "문서의 첫 번째 단락입니다."
      }
    }
  ]
}
```

## 요소 카테고리

Document Parse API가 인식하는 요소 유형:

- `heading1` ~ `heading6`: 제목 레벨
- `paragraph`: 단락
- `table`: 표
- `figure`: 그림
- `chart`: 차트
- `equation`: 수식
- `list_item`: 리스트 항목
- `header`: 머리말
- `footer`: 꼬리말
- `caption`: 캡션

## 개발

### 요구 사항

- Go 1.21 이상

### 빌드

```bash
# 현재 플랫폼용 빌드
go build -o updoc ./cmd/updoc

# 모든 플랫폼용 빌드
make build-all

# 테스트 실행
go test ./...

# 린트
golangci-lint run
```

### 프로젝트 구조

```
updoc/
├── cmd/
│   └── updoc/
│       └── main.go          # 엔트리포인트
├── internal/
│   ├── api/
│   │   └── client.go        # Upstage API 클라이언트
│   ├── cmd/
│   │   ├── parse.go         # parse 명령어
│   │   ├── status.go        # status 명령어
│   │   ├── result.go        # result 명령어
│   │   ├── config.go        # config 명령어
│   │   └── models.go        # models 명령어
│   ├── config/
│   │   └── config.go        # 설정 관리
│   └── output/
│       └── formatter.go     # 출력 포매터
├── docs/
│   └── CLI_MANUAL.md        # 상세 매뉴얼
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 라이선스

MIT License

## 참고 자료

- [Upstage Document Parse 공식 문서](https://console.upstage.ai/docs/capabilities/document-parse)
- [Upstage API Reference](https://console.upstage.ai/api-reference)
- [Upstage Console](https://console.upstage.ai)

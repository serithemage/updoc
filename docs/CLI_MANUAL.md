# updoc CLI 매뉴얼

## 목차

1. [소개](#소개)
2. [설치](#설치)
3. [설정](#설정)
4. [명령어](#명령어)
5. [사용 예제](#사용-예제)
6. [문제 해결](#문제-해결)
7. [API 레퍼런스](#api-레퍼런스)

---

## 소개

`updoc`은 업스테이지(Upstage)의 Document Parse API를 커맨드라인에서 사용할 수 있게 해주는 CLI 도구입니다. Go 언어로 작성되어 단일 바이너리로 배포되며, 다양한 형식의 문서를 구조화된 텍스트(HTML, Markdown, Text)로 변환합니다.

### 지원 기능

| 기능 | 설명 |
|------|------|
| 문서 변환 | PDF, Office, HWP 등을 HTML/Markdown/Text로 변환 |
| OCR | 스캔 문서 및 이미지에서 텍스트 추출 |
| 구조 분석 | 제목, 단락, 표, 그림 등 요소별 분리 |
| 레이아웃 인식 | 다단 구성, 복잡한 레이아웃 처리 |
| 좌표 추출 | 각 요소의 페이지 내 위치 정보 제공 |

---

## 설치

### 요구 사항

- 운영체제: macOS, Linux, Windows
- 빌드 시: Go 1.21 이상

### 바이너리 다운로드

[Releases](https://github.com/serithemage/updoc/releases) 페이지에서 OS에 맞는 바이너리를 다운로드합니다.

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

# PATH에 추가하거나 원하는 위치로 이동
Move-Item updoc.exe C:\Users\$env:USERNAME\bin\
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

# 빌드
go build -o updoc ./cmd/updoc

# 설치 (선택)
sudo mv updoc /usr/local/bin/
```

### 설치 확인

```bash
updoc version
updoc --help
```

---

## 설정

### API 키 설정

Document Parse API를 사용하려면 Upstage API 키가 필요합니다.

#### 1. API 키 발급

1. [Upstage Console](https://console.upstage.ai)에 로그인
2. 새 프로젝트 생성 또는 기존 프로젝트 선택
3. API Keys 메뉴에서 새 키 생성
4. 생성된 키 복사

#### 2. 키 설정 방법

**방법 A: 환경 변수 (권장)**

```bash
# Linux/macOS
export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"

# Windows (PowerShell)
$env:UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"

# Windows (CMD)
set UPSTAGE_API_KEY=up_xxxxxxxxxxxxxxxxxxxx
```

쉘 설정 파일에 추가하여 영구 설정:

```bash
# ~/.bashrc 또는 ~/.zshrc에 추가
echo 'export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"' >> ~/.zshrc
source ~/.zshrc
```

**방법 B: 설정 명령어**

```bash
updoc config set api-key up_xxxxxxxxxxxxxxxxxxxx
```

**방법 C: 명령어 옵션**

```bash
updoc parse document.pdf --api-key up_xxxxxxxxxxxxxxxxxxxx
```

### 프라이빗 엔드포인트 설정

AWS Bedrock, 프라이빗 호스팅 등 커스텀 엔드포인트를 사용하는 경우 다음 방법으로 설정합니다.

**방법 A: 환경 변수**

```bash
export UPSTAGE_API_ENDPOINT="https://your-private-endpoint.com/v1"
```

**방법 B: 설정 명령어**

```bash
updoc config set endpoint https://your-private-endpoint.com/v1
```

**방법 C: 명령어 옵션**

```bash
updoc parse document.pdf --endpoint https://your-private-endpoint.com/v1
```

우선순위: 명령어 옵션 > 환경 변수 > 설정 파일 > 기본값

### 설정 파일

설정 파일 위치:
- Linux/macOS: `~/.config/updoc/config.yaml`
- Windows: `%APPDATA%\updoc\config.yaml`

```yaml
api_key: "up_xxxxxxxxxxxxxxxxxxxx"
endpoint: ""  # 기본값 사용 시 비워둠
default_format: markdown
default_mode: standard
default_ocr: auto
output_dir: "./output"
```

### 설정 관리

```bash
# 현재 설정 보기
updoc config list

# 설정 값 변경
updoc config set default-format html
updoc config set default-mode enhanced

# 설정 값 조회
updoc config get default-format

# 설정 초기화
updoc config reset

# 설정 파일 위치 확인
updoc config path
```

---

## 명령어

### updoc parse

문서를 파싱하여 구조화된 텍스트로 변환합니다.

```
updoc parse <file> [options]
```

#### 인자

| 인자 | 설명 |
|------|------|
| `<file>` | 파싱할 문서 파일 경로 (필수) |

#### 옵션

| 옵션 | 단축 | 설명 | 기본값 |
|------|------|------|--------|
| `--format <type>` | `-f` | 출력 형식: html, markdown, text | markdown |
| `--output <path>` | `-o` | 출력 파일 경로 | stdout |
| `--mode <mode>` | `-m` | 파싱 모드: standard, enhanced, auto | standard |
| `--model <name>` | | 모델 지정 | document-parse |
| `--ocr <type>` | | OCR 설정: auto, force | auto |
| `--chart-recognition` | | 차트를 표로 변환 | true |
| `--no-chart-recognition` | | 차트 변환 비활성화 | |
| `--merge-tables` | | 다중 페이지 테이블 병합 | false |
| `--coordinates` | | 좌표 정보 포함 | true |
| `--no-coordinates` | | 좌표 정보 제외 | |
| `--elements-only` | `-e` | 요소별 결과만 출력 | false |
| `--json` | `-j` | JSON 형식으로 출력 | false |
| `--async` | `-a` | 비동기 처리 사용 | false |
| `--output-dir` | `-d` | 배치 처리 시 출력 디렉토리 | . |
| `--recursive` | `-r` | 디렉토리 재귀 탐색 | false |
| `--quiet` | `-q` | 진행 메시지 숨김 | false |
| `--verbose` | `-v` | 상세 로그 출력 | false |
| `--api-key <key>` | | API 키 직접 지정 | 환경변수 |
| `--endpoint <url>` | | API 엔드포인트 URL | 기본 엔드포인트 |

#### 예제

```bash
# 기본 사용
updoc parse report.pdf

# HTML로 변환하여 파일 저장
updoc parse report.pdf -f html -o report.html

# Enhanced 모드로 복잡한 문서 처리
updoc parse complex-form.pdf --mode enhanced

# 스캔 문서 강제 OCR
updoc parse scanned.pdf --ocr force

# JSON 형식 출력
updoc parse document.pdf --json -o result.json

# 배치 처리
updoc parse ./documents/*.pdf --output-dir ./results/
```

---

### updoc status

비동기 요청의 상태를 확인합니다.

```
updoc status <request-id> [options]
```

#### 인자

| 인자 | 설명 |
|------|------|
| `<request-id>` | 비동기 요청 ID (필수) |

#### 옵션

| 옵션 | 단축 | 설명 | 기본값 |
|------|------|------|--------|
| `--json` | `-j` | JSON 형식으로 출력 | false |
| `--watch` | `-w` | 실시간 상태 모니터링 | false |
| `--interval` | `-i` | 모니터링 간격 (초) | 5 |

#### 예제

```bash
# 상태 확인
updoc status abc123def456

# JSON 형식으로 출력
updoc status abc123def456 --json

# 실시간 모니터링
updoc status abc123def456 --watch
```

#### 출력 예시

```
Request ID: abc123def456
Status: processing
Progress: 45%
Pages processed: 45/100
Elapsed time: 1m 23s
```

---

### updoc result

비동기 요청의 결과를 가져옵니다.

```
updoc result <request-id> [options]
```

#### 인자

| 인자 | 설명 |
|------|------|
| `<request-id>` | 비동기 요청 ID (필수) |

#### 옵션

| 옵션 | 단축 | 설명 | 기본값 |
|------|------|------|--------|
| `--output <path>` | `-o` | 출력 파일 경로 | stdout |
| `--format <type>` | `-f` | 출력 형식 | markdown |
| `--wait` | `-w` | 완료까지 대기 | false |
| `--timeout <sec>` | `-t` | 대기 타임아웃(초) | 300 |
| `--json` | `-j` | JSON 형식으로 출력 | false |

#### 예제

```bash
# 결과 가져오기
updoc result abc123def456 -o output.md

# 완료까지 대기 후 결과 가져오기
updoc result abc123def456 --wait -o output.md

# 타임아웃 설정
updoc result abc123def456 --wait --timeout 600 -o output.md
```

---

### updoc models

사용 가능한 모델 목록을 표시합니다.

```
updoc models [options]
```

#### 옵션

| 옵션 | 단축 | 설명 | 기본값 |
|------|------|------|--------|
| `--json` | `-j` | JSON 형식으로 출력 | false |

#### 출력 예시

```
Available Models:

  document-parse          기본 모델 (권장, alias)
  document-parse-250618   특정 버전 (2025-06-18)
  document-parse-nightly  최신 테스트 버전

Tip: 'document-parse' alias를 사용하면 자동으로 최신 안정 버전이 적용됩니다.
```

---

### updoc config

설정을 관리합니다.

```
updoc config <command> [key] [value]
```

#### 하위 명령어

| 명령어 | 설명 |
|--------|------|
| `list` | 모든 설정 표시 |
| `get <key>` | 특정 설정 조회 |
| `set <key> <value>` | 설정 변경 |
| `reset` | 설정 초기화 |
| `path` | 설정 파일 경로 표시 |

#### 설정 키

| 키 | 설명 | 값 |
|----|------|-----|
| `api-key` | API 키 | 문자열 |
| `endpoint` | API 엔드포인트 URL | URL |
| `default-format` | 기본 출력 형식 | html, markdown, text |
| `default-mode` | 기본 파싱 모드 | standard, enhanced, auto |
| `default-ocr` | 기본 OCR 설정 | auto, force |
| `output-dir` | 기본 출력 디렉토리 | 경로 |

#### 예제

```bash
# 모든 설정 보기
updoc config list

# 특정 설정 조회
updoc config get api-key

# 설정 변경
updoc config set default-format html
updoc config set default-mode enhanced

# 설정 초기화
updoc config reset
```

---

### updoc version

버전 정보를 표시합니다.

```
updoc version [options]
```

#### 옵션

| 옵션 | 단축 | 설명 |
|------|------|------|
| `--short` | `-s` | 버전 번호만 출력 |
| `--json` | `-j` | JSON 형식으로 출력 |

#### 출력 예시

```
updoc version 1.0.0
  Commit: abc1234
  Built: 2025-01-11T10:00:00Z
  Go version: go1.21.5
  OS/Arch: darwin/arm64
```

---

## 사용 예제

### 기본 워크플로우

```bash
# 1. API 키 설정
export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"

# 2. PDF를 Markdown으로 변환
updoc parse report.pdf -o report.md

# 3. 변환 결과 확인
cat report.md
```

### 스캔 문서 처리

```bash
# 스캔된 PDF를 OCR 처리
updoc parse scanned-document.pdf --ocr force --mode enhanced -o output.md
```

### 복잡한 레이아웃 문서

```bash
# 표, 차트가 많은 문서를 enhanced 모드로 처리
updoc parse financial-report.pdf \
  --mode enhanced \
  --chart-recognition \
  --merge-tables \
  -o report.html \
  -f html
```

### 요소별 분석

```bash
# 문서를 요소별로 분리하여 JSON 출력
updoc parse document.pdf --elements-only --json -o elements.json

# jq로 표만 추출
cat elements.json | jq '.elements[] | select(.category == "table")'

# jq로 제목만 추출
cat elements.json | jq '.elements[] | select(.category | startswith("heading"))'
```

### 대용량 문서 처리

```bash
# 비동기 요청 시작
updoc parse large-document.pdf --async
# 출력: Request ID: req_abc123

# 상태 실시간 모니터링
updoc status req_abc123 --watch

# 완료 후 결과 가져오기
updoc result req_abc123 -o result.md

# 또는 완료까지 대기
updoc result req_abc123 --wait --timeout 600 -o result.md
```

### 배치 처리

```bash
# 현재 디렉토리의 모든 PDF 처리
updoc parse *.pdf --output-dir ./results/

# 특정 디렉토리의 문서 재귀 처리
updoc parse ./documents/ --output-dir ./results/ --recursive

# 셸 스크립트로 세부 제어
for file in *.pdf; do
  echo "Processing: $file"
  updoc parse "$file" -o "${file%.pdf}.md" --quiet
done
```

### 파이프라인 활용

```bash
# 변환 결과를 다른 도구로 전달
updoc parse document.pdf | grep -i "important"

# 여러 형식으로 동시 출력
updoc parse document.pdf --json | tee result.json | jq -r '.content.markdown' > result.md

# 특정 요소만 추출하여 처리
updoc parse document.pdf --json | jq -r '.elements[] | select(.category == "table") | .content.markdown'
```

### 자동화 스크립트 예제

```bash
#!/bin/bash
# batch_convert.sh - 문서 일괄 변환 스크립트

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

## 문제 해결

### 자주 발생하는 오류

#### API 키 오류

```
Error: Invalid API key
```

해결:
1. API 키가 올바른지 확인
2. 환경 변수가 설정되어 있는지 확인: `echo $UPSTAGE_API_KEY`
3. 키 앞뒤 공백 제거
4. 설정 확인: `updoc config get api-key`

#### 파일 형식 오류

```
Error: Unsupported file format: .xyz
```

지원 형식:
- 문서: PDF, DOCX, PPTX, XLSX, HWP
- 이미지: JPEG, PNG, BMP, TIFF, HEIC

#### 페이지 수 초과

```
Error: Document exceeds maximum page limit (100 pages for sync API)
```

동기 API는 최대 100페이지, 비동기 API는 최대 1,000페이지 지원.
대용량 문서는 `--async` 옵션 사용:

```bash
updoc parse large-document.pdf --async
```

#### 타임아웃

```
Error: Request timeout after 120s
```

해결:
- 비동기 모드 사용: `--async`
- 네트워크 상태 확인
- 파일 크기 확인

#### 파일 접근 오류

```
Error: Cannot read file: permission denied
```

해결:
```bash
# 파일 권한 확인
ls -la document.pdf

# 권한 수정 (필요시)
chmod 644 document.pdf
```

### 디버깅

```bash
# 상세 로그 출력
updoc parse document.pdf --verbose

# 요청/응답 확인
updoc parse document.pdf --verbose 2>&1 | tee debug.log

# 설정 상태 확인
updoc config list
```

### 로그 레벨

`--verbose` 플래그를 사용하면 다음 정보가 출력됩니다:
- API 요청 URL 및 헤더
- 요청 파라미터
- 응답 상태 코드
- 처리 시간

---

## API 레퍼런스

### 엔드포인트

| 용도 | URL |
|------|-----|
| 동기 파싱 | `POST https://api.upstage.ai/v1/document-digitization` |
| 비동기 파싱 | `POST https://api.upstage.ai/v1/document-digitization/async` |
| 상태 확인 | `GET https://api.upstage.ai/v1/document-digitization/async/{id}` |

### 인증

```
Authorization: Bearer <UPSTAGE_API_KEY>
```

### 요청 형식

`Content-Type: multipart/form-data`

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| `model` | string | O | 모델명 (document-parse) |
| `document` | file | O | 문서 파일 |
| `mode` | string | | standard, enhanced, auto |
| `ocr` | string | | auto, force |
| `output_formats` | string | | 출력 형식 |
| `chart_recognition` | boolean | | 차트 변환 |
| `merge_multipage_tables` | boolean | | 테이블 병합 |
| `coordinates` | boolean | | 좌표 포함 |

### 응답 구조

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

### 요소 카테고리

| 카테고리 | 설명 |
|----------|------|
| `heading1` ~ `heading6` | 제목 레벨 |
| `paragraph` | 단락 |
| `table` | 표 |
| `figure` | 그림 |
| `chart` | 차트 |
| `equation` | 수식 |
| `list_item` | 리스트 항목 |
| `header` | 머리말 |
| `footer` | 꼬리말 |
| `caption` | 캡션 |

---

## 부록

### 환경 변수

| 변수 | 설명 |
|------|------|
| `UPSTAGE_API_KEY` | API 인증 키 |
| `UPSTAGE_API_ENDPOINT` | API 엔드포인트 URL (프라이빗 호스팅용) |
| `UPDOC_CONFIG_PATH` | 설정 파일 경로 (선택) |
| `UPDOC_LOG_LEVEL` | 로그 레벨: debug, info, warn, error |

### 종료 코드

| 코드 | 의미 |
|------|------|
| 0 | 성공 |
| 1 | 일반 오류 |
| 2 | 인자 오류 |
| 3 | API 오류 |
| 4 | 파일 I/O 오류 |
| 5 | 인증 오류 |

### 관련 링크

- [Upstage Console](https://console.upstage.ai)
- [Document Parse 공식 문서](https://console.upstage.ai/docs/capabilities/document-parse)
- [API Reference](https://console.upstage.ai/api-reference)
- [Upstage 블로그](https://upstage.ai/blog)
- [GitHub 저장소](https://github.com/serithemage/updoc)

### 버전 히스토리

| 버전 | 날짜 | 변경 사항 |
|------|------|-----------|
| 1.0.0 | - | 최초 릴리스 |

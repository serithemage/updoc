# updoc - Product Requirements Document

## 1. 개요

### 1.1 프로젝트 요약
`updoc`은 업스테이지(Upstage) Document Parse API를 CLI로 래핑하여 커맨드라인에서 문서 파싱을 수행할 수 있게 해주는 도구입니다.

### 1.2 목표
- 다양한 문서 형식(PDF, Office, HWP, 이미지)을 HTML/Markdown/Text로 변환
- 단일 바이너리로 배포하여 설치 간소화
- 동기/비동기 API 모두 지원
- 배치 처리 및 자동화 스크립트 연동 지원

### 1.3 대상 사용자
- 문서를 프로그래밍 방식으로 처리해야 하는 개발자
- RAG(Retrieval-Augmented Generation) 파이프라인 구축자
- 문서 자동화 워크플로우 담당자

---

## 2. 기능 요구사항

### 2.1 핵심 명령어

| 우선순위 | 명령어 | 설명 |
|----------|--------|------|
| P0 | `updoc parse <file>` | 문서 파싱 (동기) |
| P0 | `updoc config` | 설정 관리 |
| P1 | `updoc parse --async` | 비동기 파싱 |
| P1 | `updoc status <id>` | 비동기 요청 상태 확인 |
| P1 | `updoc result <id>` | 비동기 결과 가져오기 |
| P2 | `updoc models` | 모델 목록 조회 |
| P2 | `updoc version` | 버전 정보 |

### 2.2 parse 명령어 옵션

| 우선순위 | 옵션 | 설명 |
|----------|------|------|
| P0 | `--format, -f` | 출력 형식 (html, markdown, text) |
| P0 | `--output, -o` | 출력 파일 경로 |
| P0 | `--mode, -m` | 파싱 모드 (standard, enhanced, auto) |
| P1 | `--ocr` | OCR 설정 (auto, force) |
| P1 | `--json, -j` | JSON 형식 출력 |
| P1 | `--async, -a` | 비동기 처리 |
| P2 | `--chart-recognition` | 차트를 표로 변환 |
| P2 | `--merge-tables` | 다중 페이지 테이블 병합 |
| P2 | `--coordinates` | 좌표 정보 포함 |
| P2 | `--elements-only, -e` | 요소별 결과만 출력 |
| P2 | `--output-dir, -d` | 배치 출력 디렉토리 |
| P2 | `--recursive, -r` | 디렉토리 재귀 탐색 |
| P3 | `--verbose, -v` | 상세 로그 |
| P3 | `--quiet, -q` | 진행 메시지 숨김 |

### 2.3 config 명령어

| 우선순위 | 하위 명령어 | 설명 |
|----------|-------------|------|
| P0 | `config set <key> <value>` | 설정 값 저장 |
| P0 | `config get <key>` | 설정 값 조회 |
| P1 | `config list` | 전체 설정 표시 |
| P2 | `config reset` | 설정 초기화 |
| P2 | `config path` | 설정 파일 경로 |

---

## 3. 기술 요구사항

### 3.1 개발 환경
- 언어: Go 1.21+
- CLI 프레임워크: cobra
- HTTP 클라이언트: net/http (표준 라이브러리)
- 설정 관리: viper
- 테스트: testing (표준 라이브러리)

### 3.2 프로젝트 구조

```
updoc/
├── cmd/
│   └── updoc/
│       └── main.go              # 엔트리포인트
├── internal/
│   ├── api/
│   │   ├── client.go            # API 클라이언트
│   │   ├── types.go             # 요청/응답 타입
│   │   └── client_test.go
│   ├── cmd/
│   │   ├── root.go              # 루트 명령어
│   │   ├── parse.go             # parse 명령어
│   │   ├── config.go            # config 명령어
│   │   ├── status.go            # status 명령어
│   │   ├── result.go            # result 명령어
│   │   ├── models.go            # models 명령어
│   │   └── version.go           # version 명령어
│   ├── config/
│   │   ├── config.go            # 설정 로드/저장
│   │   └── config_test.go
│   └── output/
│       ├── formatter.go         # 출력 포매터
│       └── formatter_test.go
├── docs/
│   ├── PRD.md
│   └── CLI_MANUAL.md
├── .gitignore
├── .goreleaser.yaml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 3.3 API 연동

**엔드포인트:**
- 동기: `POST https://api.upstage.ai/v1/document-digitization`
- 비동기: `POST https://api.upstage.ai/v1/document-digitization/async`
- 상태: `GET https://api.upstage.ai/v1/document-digitization/async/{id}`

**인증:**
- `Authorization: Bearer <UPSTAGE_API_KEY>`

**요청 형식:**
- `Content-Type: multipart/form-data`

### 3.4 설정 파일

**위치:**
- Linux/macOS: `~/.config/updoc/config.yaml`
- Windows: `%APPDATA%\updoc\config.yaml`

**스키마:**
```yaml
api_key: string
default_format: string  # html, markdown, text
default_mode: string    # standard, enhanced, auto
default_ocr: string     # auto, force
output_dir: string
```

### 3.5 환경 변수

| 변수 | 설명 | 우선순위 |
|------|------|----------|
| `UPSTAGE_API_KEY` | API 키 | 설정 파일보다 높음 |
| `UPDOC_CONFIG_PATH` | 설정 파일 경로 | |
| `UPDOC_LOG_LEVEL` | 로그 레벨 | |

---

## 4. 비기능 요구사항

### 4.1 성능
- 파일 업로드 시 스트리밍 처리
- 대용량 응답 처리 시 메모리 효율적 파싱

### 4.2 에러 처리
- API 오류 시 명확한 에러 메시지
- 재시도 로직 (네트워크 오류 시)
- 종료 코드 정의 (0: 성공, 1: 일반 오류, 2: 인자 오류, 3: API 오류, 4: 파일 I/O 오류, 5: 인증 오류)

### 4.3 보안
- API 키는 설정 파일에 저장 시 파일 권한 제한 (0600)
- 상세 로그에서 API 키 마스킹

### 4.4 배포
- goreleaser를 통한 크로스 컴파일
- 지원 플랫폼: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64

---

## 5. 마일스톤

### M1: MVP (P0 기능)
- [x] 프로젝트 구조 설정
- [ ] API 클라이언트 구현
- [ ] parse 명령어 (동기, 기본 옵션)
- [ ] config 명령어 (set, get)
- [ ] 기본 에러 처리

### M2: 비동기 지원 (P1 기능)
- [ ] 비동기 파싱 구현
- [ ] status 명령어
- [ ] result 명령어
- [ ] OCR 옵션
- [ ] JSON 출력

### M3: 고급 기능 (P2 기능)
- [ ] 배치 처리 (output-dir, recursive)
- [ ] 추가 파싱 옵션 (chart-recognition, merge-tables, coordinates)
- [ ] elements-only 옵션
- [ ] models 명령어
- [ ] config list, reset, path

### M4: 완성도 (P3 기능)
- [ ] verbose/quiet 모드
- [ ] 완전한 테스트 커버리지
- [ ] goreleaser 설정
- [ ] GitHub Actions CI/CD

---

## 6. 구현 태스크

### Phase 1: 프로젝트 초기화
1. Go 모듈 초기화 및 의존성 설정
2. 프로젝트 디렉토리 구조 생성
3. Makefile 작성
4. .gitignore 설정

### Phase 2: 핵심 인프라
5. 설정 관리 모듈 구현 (config 패키지)
6. API 클라이언트 기본 구조 구현
7. API 타입 정의 (요청/응답 구조체)
8. 루트 명령어 및 CLI 프레임워크 설정

### Phase 3: MVP 명령어
9. config set/get 명령어 구현
10. parse 명령어 구현 (동기, 기본 옵션)
11. 출력 포매터 구현 (html, markdown, text)
12. 기본 에러 처리 및 종료 코드

### Phase 4: 비동기 지원
13. 비동기 API 클라이언트 구현
14. parse --async 옵션 구현
15. status 명령어 구현
16. result 명령어 구현

### Phase 5: 고급 기능
17. 배치 처리 구현 (output-dir, recursive)
18. 추가 파싱 옵션 구현
19. elements-only 옵션 구현
20. models 명령어 구현
21. config list/reset/path 구현

### Phase 6: 완성도
22. version 명령어 구현
23. verbose/quiet 모드 구현
24. 테스트 작성
25. goreleaser 설정
26. CI/CD 파이프라인 구성

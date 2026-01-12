# updoc 기여 가이드

[English](CONTRIBUTING.md) | [日本語](CONTRIBUTING.ja.md)

updoc에 기여해 주셔서 감사합니다! 이 문서는 기여 방법에 대한 가이드라인을 제공합니다.

## 목차

- [행동 강령](#행동-강령)
- [시작하기](#시작하기)
- [개발 환경 설정](#개발-환경-설정)
- [기여 방법](#기여-방법)
- [Pull Request 프로세스](#pull-request-프로세스)
- [코딩 표준](#코딩-표준)
- [커밋 메시지 규칙](#커밋-메시지-규칙)
- [테스트](#테스트)
- [문서화](#문서화)

## 행동 강령

이 프로젝트는 모든 기여자가 지켜야 할 행동 강령을 따릅니다. 상호작용에서 존중하고 건설적인 태도를 유지해 주세요.

## 시작하기

### 사전 요구 사항

- Go 1.21 이상
- Git
- Make (Makefile 명령어 사용 시, 선택)
- golangci-lint (린팅용)

### Fork 및 Clone

1. GitHub에서 저장소를 Fork합니다
2. Fork한 저장소를 Clone합니다:
   ```bash
   git clone https://github.com/YOUR_USERNAME/updoc.git
   cd updoc
   ```
3. upstream 원격 저장소를 추가합니다:
   ```bash
   git remote add upstream https://github.com/serithemage/updoc.git
   ```

## 개발 환경 설정

### 빠른 설정

```bash
# 개발 의존성 설치 및 Git 훅 설정
make dev-setup
```

### 수동 설정

```bash
# golangci-lint 설치
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 프로젝트 빌드
go build -o updoc ./cmd/updoc

# 테스트 실행
go test ./...
```

### 프로젝트 구조

```
updoc/
├── cmd/updoc/           # 애플리케이션 진입점
├── internal/
│   ├── api/             # Upstage API 클라이언트
│   ├── cmd/             # CLI 명령어 구현
│   ├── config/          # 설정 관리
│   └── output/          # 출력 포매터
├── test/e2e/            # E2E 테스트
├── docs/                # 문서
├── .github/             # GitHub 템플릿 및 워크플로우
└── Makefile             # 빌드 자동화
```

## 기여 방법

### 버그 리포트

버그 리포트를 작성하기 전에:
1. 중복을 피하기 위해 기존 이슈를 확인하세요
2. 관련 정보를 수집하세요:
   - updoc 버전 (`updoc version`)
   - 운영 체제 및 버전
   - 재현 단계
   - 예상 동작 vs 실제 동작
   - 에러 메시지 또는 로그

이슈 생성 시 [Bug Report 템플릿](.github/ISSUE_TEMPLATE/bug_report.md)을 사용하세요.

### 기능 제안

기능 제안을 환영합니다! 제출하기 전에:
1. 해당 기능이 이미 요청되었는지 확인하세요
2. 프로젝트 목표와 일치하는지 고려하세요
3. 명확한 사용 사례를 제공하세요

이슈 생성 시 [Feature Request 템플릿](.github/ISSUE_TEMPLATE/feature_request.md)을 사용하세요.

### 코드 기여

1. 작업할 **이슈를 찾거나** 논의를 위해 새 이슈를 생성하세요
2. 다른 사람들에게 작업 중임을 알리기 위해 이슈에 **댓글**을 남기세요
3. `main`에서 **브랜치를 생성**하세요:
   ```bash
   git checkout -b feature/your-feature-name
   # 또는
   git checkout -b fix/bug-description
   ```
4. 코딩 표준을 따라 **변경 사항을 작성**하세요
5. 새 기능에 대한 **테스트를 작성**하세요
6. **테스트 및 린팅을 실행**하세요:
   ```bash
   make test
   make lint
   ```
7. 커밋 규칙에 따라 **커밋**하세요
8. Fork에 **Push**하고 Pull Request를 생성하세요

## Pull Request 프로세스

### 제출 전 확인

- [ ] 코드가 에러 없이 컴파일됨
- [ ] 모든 테스트 통과 (`make test`)
- [ ] 린팅 통과 (`make lint`)
- [ ] 필요한 경우 문서가 업데이트됨
- [ ] 커밋 메시지가 규칙을 따름

### PR 가이드라인

1. **제목**: 명확하고 설명적인 제목 사용
2. **설명**: 무엇을 왜 변경했는지 설명
3. **이슈 연결**: 관련 이슈 참조 (예: "Fixes #123")
4. **집중**: 하나의 PR은 하나의 문제를 다룸
5. **신속한 응답**: 리뷰 피드백에 신속하게 대응

### 리뷰 프로세스

1. 자동화된 검사 통과 필요 (CI/CD)
2. 최소 한 명의 메인테이너 승인 필요
3. 모든 논의가 해결되어야 함
4. 브랜치가 `main`과 최신 상태여야 함

## 코딩 표준

### Go 코드 스타일

- [Effective Go](https://golang.org/doc/effective_go.html) 준수
- `gofmt`을 사용한 포매팅
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) 준수

### 가이드라인

- 함수를 작고 집중적으로 유지
- 설명적인 변수 및 함수 이름 사용
- 내보내는 함수와 복잡한 로직에 주석 추가
- 에러를 명시적으로 처리, 무시하지 않음
- 전역 상태 피하기

### 예시

```go
// ParseDocument는 문서 파일을 파싱하고 구조화된 컨텐츠를 반환합니다.
// 파일 형식이 지원되지 않거나 파싱이 실패하면 에러를 반환합니다.
func ParseDocument(filePath string, opts ...Option) (*Result, error) {
    if filePath == "" {
        return nil, errors.New("file path cannot be empty")
    }

    // ... 구현
}
```

## 커밋 메시지 규칙

[Conventional Commits](https://www.conventionalcommits.org/) 형식을 따릅니다:

### 형식

```
<type>(<scope>): <description>

[선택적 본문]

[선택적 푸터]
```

### 타입

| 타입 | 설명 |
|------|------|
| `feat` | 새로운 기능 |
| `fix` | 버그 수정 |
| `docs` | 문서만 변경 |
| `style` | 코드 스타일 (포매팅 등) |
| `refactor` | 코드 리팩토링 |
| `test` | 테스트 추가 또는 업데이트 |
| `chore` | 유지보수 작업 |
| `perf` | 성능 개선 |
| `ci` | CI/CD 변경 |

### 예시

```bash
feat(parse): HWPX 파일 형식 지원 추가

fix(config): 환경 변수에서 API 키 로딩 문제 해결

docs(readme): 설치 방법 업데이트

test(api): 비동기 파싱 유닛 테스트 추가
```

## 테스트

### 테스트 실행

```bash
# 유닛 테스트
make test

# E2E 테스트 (UPSTAGE_API_KEY 필요)
export UPSTAGE_API_KEY="your-api-key"
make test-e2e

# 커버리지와 함께 전체 테스트
go test -cover ./...
```

### 테스트 작성

- `*_test.go` 파일에 테스트 배치
- 적절한 곳에 테이블 기반 테스트 사용
- 외부 의존성 모킹
- 숫자보다 의미 있는 커버리지 목표

### 테스트 예시

```go
func TestParseRequest_Validate(t *testing.T) {
    tests := []struct {
        name    string
        req     *ParseRequest
        wantErr bool
    }{
        {
            name:    "유효한 요청",
            req:     &ParseRequest{FilePath: "test.pdf"},
            wantErr: false,
        },
        {
            name:    "빈 파일 경로",
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

## 문서화

### 문서 업데이트 시점

- 새로운 기능이나 명령어 추가 시
- 기존 동작 변경 시
- 새로운 설정 옵션 추가 시
- 불명확하거나 잘못된 문서 수정 시

### 문서 파일

| 파일 | 목적 |
|------|------|
| `README.md` | 프로젝트 개요 및 빠른 시작 |
| `docs/CLI_MANUAL.md` | 상세 CLI 레퍼런스 |
| `CONTRIBUTING.md` | 기여 가이드라인 |

### 다국어 지원

이 프로젝트는 영어, 한국어, 일본어로 문서를 관리합니다. 문서 업데이트 시:

1. 영어 버전을 먼저 업데이트
2. `/translate-docs`를 사용하여 번역 동기화, 또는
3. 일관성을 유지하며 수동으로 번역 업데이트

## 질문이 있으신가요?

- 질문은 [Discussion](https://github.com/serithemage/updoc/discussions)을 열어주세요
- 먼저 기존 이슈와 토론을 확인해 주세요
- 명확하게 작성하고 맥락을 제공해 주세요

updoc에 기여해 주셔서 감사합니다!

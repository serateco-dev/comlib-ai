# WIZ Common Library (Go)

WIZ Platform 마이크로서비스를 위한 공통 유틸리티 및 컴포넌트 라이브러리입니다.

## 기능

- **JWT 유틸리티**: JWT 토큰 생성, 검증, 파싱
- **HTTP 응답**: 표준화된 API 응답 구조체
- **에러 핸들링**: 공통 예외 처리 및 에러 응답
- **데이터베이스**: GORM 기반 데이터베이스 유틸리티
- **암호화**: AES256 암호화/복호화 유틸리티
- **미들웨어**: Gin 기반 공통 미들웨어

## 설치

```bash
go get github.com/serateco-dev/comlib-ai
```

## 사용법

```go
import (
    "github.com/serateco-dev/comlib-ai/dto"
    "github.com/serateco-dev/comlib-ai/jwt"
    "github.com/serateco-dev/comlib-ai/crypto"
)

// API 응답
response := dto.NewSuccessResponse("Success", data)

// JWT 토큰 생성
token, err := jwt.GenerateToken(userID, "secret")

// 암호화
encrypted, err := crypto.Encrypt("plaintext", "key")
```

## 라이센스

MIT License
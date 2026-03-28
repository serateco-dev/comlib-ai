// jwt/jwt.go - JwtUtil.java 대응
// JWT 토큰 생성 및 검증 유틸리티
package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTUtil JWT 유틸리티 구조체
type JWTUtil struct {
	SecretKey        string
	Expiration       time.Duration
	RefreshExpiration time.Duration
	Issuer           string
}

// Claims JWT 클레임 구조체
type Claims struct {
	UserNo        int64    `json:"userNo"`
	ServiceID     string   `json:"serviceId,omitempty"`
	Role          string   `json:"role,omitempty"`
	Auth          []string `json:"auth,omitempty"`
	NickName      string   `json:"nickName,omitempty"`
	Provider      string   `json:"provider,omitempty"`
	AccessibleAPI []string `json:"accessibleApi,omitempty"`
	DeviceCD      string   `json:"deviceCd,omitempty"`
	DeviceStr     string   `json:"deviceStr,omitempty"`
	TokenType     string   `json:"tokenType,omitempty"`
	jwt.RegisteredClaims
}

// NewJWTUtil JWT 유틸리티 생성
func NewJWTUtil(secretKey string, expiration, refreshExpiration time.Duration, issuer string) *JWTUtil {
	return &JWTUtil{
		SecretKey:         secretKey,
		Expiration:        expiration,
		RefreshExpiration: refreshExpiration,
		Issuer:            issuer,
	}
}

// GenerateAccessToken Access Token 생성
func (j *JWTUtil) GenerateAccessToken(
	userNo int64,
	userID string,
	serviceID string,
	role string,
	auth []interface{},
	nickName string,
	provider string,
	accessibleAPI []string,
	additionalClaims map[string]interface{},
) (string, error) {
	now := time.Now()
	expiryTime := now.Add(j.Expiration)

	// auth를 문자열 배열로 변환
	var authStrings []string
	if auth != nil {
		for _, a := range auth {
			switch v := a.(type) {
			case int:
				authStrings = append(authStrings, strconv.Itoa(v))
			case int64:
				authStrings = append(authStrings, strconv.FormatInt(v, 10))
			case float64:
				authStrings = append(authStrings, strconv.FormatInt(int64(v), 10))
			case string:
				authStrings = append(authStrings, v)
			default:
				authStrings = append(authStrings, fmt.Sprintf("%v", v))
			}
		}
	}

	// accessibleAPI가 비어있으면 "all"로 설정
	if len(accessibleAPI) == 0 {
		accessibleAPI = []string{"all"}
	}

	claims := Claims{
		UserNo:        userNo,
		ServiceID:     serviceID,
		Role:          role,
		Auth:          authStrings,
		NickName:      nickName,
		Provider:      provider,
		AccessibleAPI: accessibleAPI,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    j.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	}

	// 추가 클레임 설정
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 추가 클레임이 있으면 설정
	if additionalClaims != nil {
		for key, value := range additionalClaims {
			token.Claims.(jwt.MapClaims)[key] = value
		}
	}

	return token.SignedString([]byte(j.SecretKey))
}

// GenerateAccessTokenSimple 간소화된 Access Token 생성
func (j *JWTUtil) GenerateAccessTokenSimple(
	userNo int64,
	userID string,
	serviceID string,
	role string,
	accessibleAPI []string,
	additionalClaims map[string]interface{},
) (string, error) {
	return j.GenerateAccessToken(userNo, userID, serviceID, role, nil, "", "", accessibleAPI, additionalClaims)
}

// GenerateRefreshToken Refresh Token 생성 (기본 유효기간)
func (j *JWTUtil) GenerateRefreshToken(userNo int64, userID string) (string, error) {
	return j.GenerateRefreshTokenWithExpiration(userNo, userID, j.RefreshExpiration)
}

// GenerateRefreshTokenWithExpiration Refresh Token 생성 (유효기간 지정)
func (j *JWTUtil) GenerateRefreshTokenWithExpiration(userNo int64, userID string, expiration time.Duration) (string, error) {
	now := time.Now()
	expiryTime := now.Add(expiration)

	claims := Claims{
		UserNo:    userNo,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userNo, 10), // Refresh Token은 userNo를 subject로
			Issuer:    j.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	}

	// userID가 있으면 추가
	if userID != "" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token.Claims.(jwt.MapClaims)["userId"] = userID
		return token.SignedString([]byte(j.SecretKey))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// ParseToken 토큰 파싱
func (j *JWTUtil) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ExtractUserNo userNo 추출
func (j *JWTUtil) ExtractUserNo(tokenString string) (int64, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserNo, nil
}

// ExtractUserID userID 추출
func (j *JWTUtil) ExtractUserID(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}

// ExtractNickName 닉네임 추출
func (j *JWTUtil) ExtractNickName(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.NickName, nil
}

// ExtractAuth 권한 추출
func (j *JWTUtil) ExtractAuth(tokenString string) ([]string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return claims.Auth, nil
}

// ExtractRole 역할 추출
func (j *JWTUtil) ExtractRole(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}

// ExtractAccessibleAPI 접근 가능 API 목록 추출
func (j *JWTUtil) ExtractAccessibleAPI(tokenString string) ([]string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return []string{"all"}, err
	}
	if len(claims.AccessibleAPI) == 0 {
		return []string{"all"}, nil
	}
	return claims.AccessibleAPI, nil
}

// ExtractServiceID 서비스 ID 추출
func (j *JWTUtil) ExtractServiceID(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.ServiceID, nil
}

// ExtractProvider 제공자 추출
func (j *JWTUtil) ExtractProvider(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Provider, nil
}

// ExtractDeviceCD 디바이스 코드 추출
func (j *JWTUtil) ExtractDeviceCD(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.DeviceCD, nil
}

// ExtractDeviceStr 디바이스 문자열 추출
func (j *JWTUtil) ExtractDeviceStr(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.DeviceStr, nil
}

// ValidateToken 토큰 검증
func (j *JWTUtil) ValidateToken(tokenString string) error {
	_, err := j.ParseToken(tokenString)
	return err
}

// IsTokenValid 토큰 유효성 검사
func (j *JWTUtil) IsTokenValid(tokenString string) bool {
	if tokenString == "" {
		return false
	}
	return j.ValidateToken(tokenString) == nil
}

// IsTokenExpired 토큰 만료 확인
func (j *JWTUtil) IsTokenExpired(tokenString string) bool {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return true
	}
	return claims.ExpiresAt.Before(time.Now())
}

// IsRefreshToken Refresh Token 여부 확인
func (j *JWTUtil) IsRefreshToken(tokenString string) bool {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return false
	}
	return claims.TokenType == "refresh"
}
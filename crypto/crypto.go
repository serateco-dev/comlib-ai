// crypto/crypto.go - CryptoUtil.java 대응
// AES256 암호화/복호화 유틸리티
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

const (
	KeySize = 32 // AES-256: 32 bytes
	IVSize  = 16 // AES 블록 크기: 16 bytes
)

// CryptoUtil AES256 암호화/복호화 유틸리티
type CryptoUtil struct {
	key []byte
	iv  []byte
}

// NewCryptoUtil CryptoUtil 생성
func NewCryptoUtil(secretKey, iv string) (*CryptoUtil, error) {
	keyBytes := []byte(secretKey)
	ivBytes := []byte(iv)

	// 키 길이 검증 (AES-256은 32바이트 필요)
	if len(keyBytes) != KeySize {
		return nil, fmt.Errorf("secret key must be %d bytes for AES-256", KeySize)
	}

	// IV는 16바이트 (AES 블록 크기)
	if len(ivBytes) != IVSize {
		return nil, fmt.Errorf("IV must be %d bytes", IVSize)
	}

	return &CryptoUtil{
		key: keyBytes,
		iv:  ivBytes,
	}, nil
}

// Encrypt AES256 암호화
func (c *CryptoUtil) Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	// AES 블록 암호화 생성
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 평문을 바이트로 변환
	plainBytes := []byte(plainText)

	// PKCS5 패딩 적용
	paddedBytes := pkcs5Padding(plainBytes, aes.BlockSize)

	// CBC 모드 암호화
	mode := cipher.NewCBCEncrypter(block, c.iv)
	encrypted := make([]byte, len(paddedBytes))
	mode.CryptBlocks(encrypted, paddedBytes)

	// Base64 인코딩
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt AES256 복호화
func (c *CryptoUtil) Decrypt(encrypted string) (string, error) {
	if encrypted == "" {
		return "", nil
	}

	// Base64 디코딩
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// AES 블록 암호화 생성
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// CBC 모드 복호화
	mode := cipher.NewCBCDecrypter(block, c.iv)
	decrypted := make([]byte, len(encryptedBytes))
	mode.CryptBlocks(decrypted, encryptedBytes)

	// PKCS5 패딩 제거
	unpaddedBytes, err := pkcs5Unpadding(decrypted)
	if err != nil {
		return "", fmt.Errorf("failed to remove padding: %w", err)
	}

	return string(unpaddedBytes), nil
}

// EncryptUserID userId 암호화 (Gateway -> 헤더 전달 -> MSA에서 활용)
func (c *CryptoUtil) EncryptUserID(userID string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("userID cannot be empty")
	}

	encrypted, err := c.Encrypt(userID)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt userID: %w", err)
	}

	if encrypted == "" {
		return "", fmt.Errorf("encryption returned empty string for userID")
	}

	return encrypted, nil
}

// DecryptUserID userId 복호화 (MSA에서 사용)
func (c *CryptoUtil) DecryptUserID(encryptedUserID string) (string, error) {
	if encryptedUserID == "" {
		return "", fmt.Errorf("encryptedUserID cannot be empty")
	}

	decrypted, err := c.Decrypt(encryptedUserID)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt userID: %w", err)
	}

	if decrypted == "" {
		return "", fmt.Errorf("decryption returned empty string for userID")
	}

	return decrypted, nil
}

// pkcs5Padding PKCS5 패딩 적용
func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// pkcs5Unpadding PKCS5 패딩 제거
func pkcs5Unpadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	// 패딩 검증
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}
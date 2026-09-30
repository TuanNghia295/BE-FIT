package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewJWTService(privateKeyPEM, publicKeyPEM []byte) (*JWTService, error) {
	privateBlock, _ := pem.Decode(privateKeyPEM)
	if privateBlock == nil {
		return nil, fmt.Errorf("JWT private key is not valid PEM")
	}
	privateKey, err := x509.ParseECPrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse JWT private key: %w", err)
	}
	if privateKey.Curve != elliptic.P256() {
		return nil, fmt.Errorf("JWT private key must use the P-256 curve for ES256")
	}

	publicBlock, _ := pem.Decode(publicKeyPEM)
	if publicBlock == nil {
		return nil, fmt.Errorf("JWT public key is not valid PEM")
	}
	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse JWT public key: %w", err)
	}
	publicKey, ok := parsedPublicKey.(*ecdsa.PublicKey)
	if !ok || publicKey.Curve != elliptic.P256() {
		return nil, fmt.Errorf("JWT public key must be ECDSA P-256 for ES256")
	}
	if privateKey.X.Cmp(publicKey.X) != 0 || privateKey.Y.Cmp(publicKey.Y) != 0 {
		return nil, fmt.Errorf("JWT private and public keys do not match")
	}

	return &JWTService{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func LoadJWTService(privateKeyPath, publicKeyPath string) (*JWTService, error) {
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read JWT private key: %w", err)
	}
	publicKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read JWT public key: %w", err)
	}
	return NewJWTService(privateKeyPEM, publicKeyPEM)
}

/*
JWT bao gồm thêm một bộ claim, dùng để truyền thông tin giữa client và server.
Các claim này có thể cho biết người phát hành token, thời gian có hiệu lực, các quyền mà token được cấp,…
nói chung sẽ phụ thuộc chủ yếu vào mục đích sử dụng.
*/
type AccessTokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *JWTService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	now := time.Now()

	// Percision: the quality of being exact -> Sự chính xác, Độ chính xác
	// Create claims
	claims := AccessTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}

	// creating a new token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	// SignedString creates and returns a complete, signed JWT.
	return token.SignedString(s.privateKey)
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}))

	if err != nil {
		return nil, err
	}

	//kiểm tra claims có đúng kiểu mong đợi không.
	claims, ok := token.Claims.(*AccessTokenClaims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func GenerateRefreshToken() (string, error) {
	// create a slice have 32 bytes
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *JWTService) GenerateRefreshToken() (string, error) {
	return GenerateRefreshToken()
}

func HashRefreshToken(token string) string {
	// Chuyển chuỗi token thành byte rồi tính SHA-256.
	// Kết quả là một mảng cố định [32]byte
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}

func (s *JWTService) HashRefreshToken(token string) string {
	return HashRefreshToken(token)
}

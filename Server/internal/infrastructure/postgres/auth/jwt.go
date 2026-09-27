package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: []byte(secret),
	}
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// SignedString creates and returns a complete, signed JWT.
	return token.SignedString(s.secret)
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

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

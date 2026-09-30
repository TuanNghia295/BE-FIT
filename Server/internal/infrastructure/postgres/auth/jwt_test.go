package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func testKeyPair(t *testing.T) ([]byte, []byte) {
	t.Helper()
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	privateDER, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	privatePEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privateDER})
	publicPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})
	return privatePEM, publicPEM
}

func TestJWTServiceSignsAndValidatesES256(t *testing.T) {
	privatePEM, publicPEM := testKeyPair(t)
	service, err := NewJWTService(privatePEM, publicPEM)
	if err != nil {
		t.Fatalf("NewJWTService returned an error: %v", err)
	}

	userID := uuid.New()
	tokenString, err := service.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken returned an error: %v", err)
	}

	claims, err := service.ValidateAccessToken(tokenString)
	if err != nil {
		t.Fatalf("ValidateAccessToken returned an error: %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("expected user ID %s, got %s", userID, claims.UserID)
	}
}

func TestNewJWTServiceRejectsMismatchedKeys(t *testing.T) {
	privatePEM, _ := testKeyPair(t)
	_, otherPublicPEM := testKeyPair(t)

	if _, err := NewJWTService(privatePEM, otherPublicPEM); err == nil {
		t.Fatal("expected mismatched private/public keys to be rejected")
	}
}

func TestLoadJWTServiceReadsPEMFiles(t *testing.T) {
	privatePEM, publicPEM := testKeyPair(t)
	privatePath := filepath.Join(t.TempDir(), "jwt_private.pem")
	publicPath := filepath.Join(t.TempDir(), "jwt_public.pem")
	if err := os.WriteFile(privatePath, privatePEM, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(publicPath, publicPEM, 0600); err != nil {
		t.Fatal(err)
	}

	service, err := LoadJWTService(privatePath, publicPath)
	if err != nil {
		t.Fatalf("LoadJWTService returned an error: %v", err)
	}
	token, err := service.GenerateAccessToken(uuid.New())
	if err != nil {
		t.Fatalf("GenerateAccessToken returned an error: %v", err)
	}
	if _, err := service.ValidateAccessToken(token); err != nil {
		t.Fatalf("ValidateAccessToken returned an error: %v", err)
	}
}

func TestJWTServiceRejectsHS256Token(t *testing.T) {
	privatePEM, publicPEM := testKeyPair(t)
	service, err := NewJWTService(privatePEM, publicPEM)
	if err != nil {
		t.Fatal(err)
	}

	claims := AccessTokenClaims{
		UserID: uuid.New(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	hsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hsTokenString, err := hsToken.SignedString([]byte("test-only-secret"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.ValidateAccessToken(hsTokenString); err == nil {
		t.Fatal("expected HS256 token to be rejected")
	}
}

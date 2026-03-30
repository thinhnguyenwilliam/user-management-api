// user-management-api/pkg/auth/jwt.go
package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/thinhnguyenwilliam/user-management-api/internal/utils"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/rediscache"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

type JWTService struct {
	signingKey string
	encryptKey []byte
	cache      rediscache.Cache
}

func NewJWTService(signingKey string, encryptKey []byte, cache rediscache.Cache) ITokenService {
	return &JWTService{
		signingKey: signingKey,
		encryptKey: encryptKey,
		cache:      cache,
	}
}

type TokenPayload struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
}

type Claims struct {
	Data string `json:"data"` // encrypted payload
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (js *JWTService) GenerateRefreshToken(userID string) (string, string, error) {
	jti := uuid.NewString()

	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-management-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(js.signingKey))
	if err != nil {
		return "", "", err
	}

	// 🔥 LƯU REDIS NGAY TẠI ĐÂY
	key := "refresh_token:" + jti
	err = js.cache.Set(context.Background(), key, userID, RefreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return tokenStr, jti, nil
}

func (js *JWTService) ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	var claims RefreshClaims

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(js.signingKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return &claims, nil
}

// 🔐 Generate Access Token
func (js *JWTService) GenerateAccessToken(payload TokenPayload) (string, error) {
	// 1. Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// 2. Encrypt payload
	encData, err := utils.Encrypt(string(payloadBytes), js.encryptKey)
	if err != nil {
		return "", err
	}

	// 3. Create claims
	claims := Claims{
		Data: encData,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(), // jti
			Subject:   payload.UserID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-management-api",
		},
	}

	// 4. Sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(js.signingKey))
}

// 🔓 Parse & Decrypt Token
func (js *JWTService) ParseAccessToken(tokenStr string) (*TokenPayload, error) {
	var claims Claims

	// 1. Parse token
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(js.signingKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	// 2. Decrypt data
	decrypted, err := utils.Decrypt(claims.Data, js.encryptKey)
	if err != nil {
		return nil, err
	}

	// 3. Unmarshal payload
	var payload TokenPayload
	if err := json.Unmarshal([]byte(decrypted), &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

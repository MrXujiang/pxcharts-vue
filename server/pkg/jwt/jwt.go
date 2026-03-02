package jwt

import (
	"errors"
	"mvtable/internal/pkg/config"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret      []byte
	accessTime  time.Duration
	refreshTime time.Duration
	issuer      string
}

var (
	instance *Manager
	once     sync.Once
)

func Init(cfg config.JWTConfig) *Manager {
	once.Do(func() {
		instance = &Manager{
			secret:      []byte(cfg.Secret),
			accessTime:  cfg.AccessTime,
			refreshTime: cfg.RefreshTime,
			issuer:      cfg.Issuer,
		}
	})
	return instance
}

func GetInstance() *Manager {
	return instance
}

// GenerateAccessToken 生成访问令牌
func (m *Manager) GenerateAccessToken(userID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTime)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// GenerateRefreshToken 生成刷新JWT令牌
func (m *Manager) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTime)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken 解析JWT令牌
func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证JWT令牌
func (m *Manager) ValidateToken(tokenString string) error {
	_, err := m.ParseToken(tokenString)
	return err
}

// GetUserID 从令牌中获取用户ID
func (m *Manager) GetUserID(tokenString string) (string, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

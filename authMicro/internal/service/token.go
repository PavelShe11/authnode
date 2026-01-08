package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PavelShe11/studbridge/authMicro/internal/config"
	"github.com/PavelShe11/studbridge/authMicro/internal/entity"
	"github.com/PavelShe11/studbridge/authMicro/internal/port"
	"github.com/PavelShe11/studbridge/authMicro/internal/repository"
	commonEntity "github.com/PavelShe11/studbridge/common/entity"
	"github.com/PavelShe11/studbridge/common/logger"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	refreshTokenSessionRepo *repository.RefreshTokenSessionRepository
	accountProvider         port.AccountProvider
	jwtConfig               config.JWTConfig
	logger                  logger.Logger
}

func NewTokenService(
	refreshTokenSessionRepo *repository.RefreshTokenSessionRepository,
	accountProvider port.AccountProvider,
	logger logger.Logger,
	jwtConfig config.JWTConfig,
) *TokenService {
	return &TokenService{
		jwtConfig:               jwtConfig,
		refreshTokenSessionRepo: refreshTokenSessionRepo,
		accountProvider:         accountProvider,
		logger:                  logger,
	}
}

func (s *TokenService) CreateTokens(ctx context.Context, accountId string) (*entity.Tokens, error) {
	s.cleanupExpiredSessions(ctx)

	claimsResult, err := s.accountProvider.GetAccessTokenPayload(ctx, accountId)
	if err != nil {
		s.logger.Error(fmt.Errorf("failed to get token payload: %w", err))
		return nil, err
	}

	if claimsResult != nil {
		if sub, ok := claimsResult["sub"].(string); ok && sub != "" {
			accountId = sub
		}
	}

	now := time.Now()
	refreshExpiry := now.Add(s.jwtConfig.RefreshTokenExpiration)
	accessExpiry := now.Add(s.jwtConfig.AccessTokenExpiration)

	refreshTokenString, accessTokenString, err := s.generateTokenPair(
		accountId,
		claimsResult,
		now,
		refreshExpiry,
		accessExpiry,
	)
	if err != nil {
		s.logger.Error(err)
		return nil, commonEntity.NewInternalError()
	}

	session := &entity.RefreshTokenSession{
		AccountID:    accountId,
		RefreshToken: refreshTokenString,
		ExpiresAt:    refreshExpiry,
	}

	if err := s.refreshTokenSessionRepo.Save(ctx, session); err != nil {
		s.logger.Error(fmt.Errorf("failed to save refresh token session: %w", err))
		return nil, commonEntity.NewInternalError()
	}

	return &entity.Tokens{
		AccessToken:         accessTokenString,
		AccessTokenExpires:  accessExpiry.Unix(),
		RefreshToken:        refreshTokenString,
		RefreshTokenExpires: accessExpiry.Unix(),
	}, nil
}

func (s *TokenService) RefreshTokens(ctx context.Context, refreshTokenString string) (*entity.Tokens, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(
		refreshTokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.jwtConfig.Secret), nil
		},
	)

	if err != nil || !token.Valid {
		s.logger.Debug(err)
		return nil, entity.NewInvalidRefreshTokenError()
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		s.logger.Error("sub not found in refresh token")
		return nil, entity.NewInvalidRefreshTokenError()
	}

	session, err := s.refreshTokenSessionRepo.FindByToken(ctx, refreshTokenString)
	if err != nil {
		s.logger.Error(err)
		return nil, commonEntity.NewInternalError()
	}
	if session == nil {
		return nil, entity.NewInvalidRefreshTokenError()
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.refreshTokenSessionRepo.DeleteByToken(ctx, refreshTokenString)
		return nil, entity.NewRefreshTokenExpiredError()
	}

	if err := s.refreshTokenSessionRepo.DeleteByToken(ctx, refreshTokenString); err != nil {
		s.logger.Error(err)
	}

	return s.CreateTokens(ctx, sub)
}

func (s *TokenService) generateTokenPair(
	accountId string,
	claimsResult map[string]interface{},
	now time.Time,
	refreshExpiry time.Time,
	accessExpiry time.Time,
) (refreshToken string, accessToken string, err error) {
	baseClaims := jwt.MapClaims{
		"sub": accountId,
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}

	if claimsResult != nil {
		for key, value := range claimsResult {
			baseClaims[key] = value
		}
	}

	refreshClaims := jwt.MapClaims{
		"sub": accountId,
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": refreshExpiry.Unix(),
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJWT.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	accessClaims := make(jwt.MapClaims, len(baseClaims)+1)
	for key, value := range baseClaims {
		accessClaims[key] = value
	}
	accessClaims["exp"] = accessExpiry.Unix()

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJWT.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return refreshToken, accessToken, nil
}

func (s *TokenService) cleanupExpiredSessions(ctx context.Context) {
	if err := s.refreshTokenSessionRepo.CleanExpired(ctx); err != nil {
		s.logger.Error(fmt.Errorf("error cleaning expired refresh token sessions: %w", err))
	}
}

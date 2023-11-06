package jwt

import (
	"context"
	"time"

	"kredit-plus/app/constants"
	customerDBModels "kredit-plus/app/db/dto/customer"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type IJWTService interface {
	GenerateToken(ctx context.Context, user customerDBModels.Customer) (TokenDetails, error)
	ParseToken(ctx context.Context, tokenString string) (*JWTToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (TokenDetails, error)
}

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}

type JWTToken struct {
	UserUUID string `json:"user_uuid"`
	jwt.StandardClaims
}

func (s *JWTService) GenerateToken(ctx context.Context, user customerDBModels.Customer) (TokenDetails, error) {
	tokenDetails := TokenDetails{}

	// Generate access token
	accessToken, atExpires, err := generateAccessToken(user)
	if err != nil {
		return tokenDetails, err
	}

	// Generate refresh token
	refreshToken, rtExpires, err := generateRefreshToken(user)
	if err != nil {
		return tokenDetails, err
	}

	tokenDetails.AccessToken = accessToken
	tokenDetails.AtExpires = atExpires
	tokenDetails.RefreshToken = refreshToken
	tokenDetails.RtExpires = rtExpires

	return tokenDetails, nil
}

func (s *JWTService) ParseToken(ctx context.Context, tokenString string) (*JWTToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.Config.JwtConfig.JWT_ACCESS_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTToken); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func (s *JWTService) RefreshToken(ctx context.Context, refreshToken string) (TokenDetails, error) {
	// Parse the refresh token to extract user information
	claims, err := parseRefreshToken(refreshToken)
	if err != nil {
		return TokenDetails{}, err
	}

	// Check if the refresh token is still valid
	if time.Now().Unix() > claims.ExpiresAt {
		return TokenDetails{}, jwt.ErrInvalidKey
	}

	// Generate a new access token
	userUUIDString := claims.UserUUID

	userUUID, err := uuid.Parse(userUUIDString)
	if err != nil {
		return TokenDetails{}, err
	}

	user := customerDBModels.Customer{UUID: userUUID} // Create a user object with the user UUID
	newAccessToken, atExpires, err := generateAccessToken(user)
	if err != nil {
		return TokenDetails{}, err
	}

	// Generate a new refresh token (optional, depending on your use case)
	newRefreshToken, rtExpires, err := generateRefreshToken(user)
	if err != nil {
		return TokenDetails{}, err
	}

	tokenDetails := TokenDetails{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		AtExpires:    atExpires,
		RtExpires:    rtExpires,
	}

	return tokenDetails, nil
}

func generateAccessToken(user customerDBModels.Customer) (string, int64, error) {
	claims := jwt.MapClaims{
		"user_uuid": user.UUID.String(),
		"exp":       time.Now().Add(time.Minute * time.Duration(constants.Config.JwtConfig.JWT_ACCESS_EXP)).Unix(),
		"iat":       time.Now().Unix(),
	}

	return generateToken(claims, constants.Config.JwtConfig.JWT_ACCESS_SECRET)
}

func generateRefreshToken(user customerDBModels.Customer) (string, int64, error) {
	claims := jwt.MapClaims{
		"user_uuid": user.UUID.String(),
		"exp":       time.Now().Add(time.Minute * time.Duration(constants.Config.JwtConfig.JWT_REFRESH_EXP)).Unix(),
		"iat":       time.Now().Unix(),
	}

	return generateToken(claims, constants.Config.JwtConfig.JWT_REFRESH_SECRET)
}

func generateToken(claims jwt.Claims, secretKey string) (string, int64, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}
	return tokenString, claims.(jwt.MapClaims)["exp"].(int64), nil
}

func parseRefreshToken(refreshToken string) (*JWTToken, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &JWTToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.Config.JwtConfig.JWT_REFRESH_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTToken); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

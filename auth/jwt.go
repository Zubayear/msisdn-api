package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"huspass/model"
	"os"
	"strconv"
	"time"
)

type JWTClaim struct {
	Username string `json:"usr"`
	Roles    model.Role
	jwt.RegisteredClaims
}

func GenerateJWT(user *model.User) (string, error) {
	exp, err := strconv.Atoi(os.Getenv("EXPTIME"))
	if err != nil {
		return "", err
	}
	claims := &JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "https://jwt.token.com",
			Subject:  "token",
			Audience: nil,
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Duration(exp) * time.Hour),
			},
			IssuedAt: &jwt.NumericDate{Time: time.Now()},
			ID:       uuid.NewString(),
		},
	}
	claims.Roles = user.Roles
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func ValidateToken(signedToken string) (model.Role, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return "", errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return "", errors.New("token expired")
	}
	return claims.Roles, nil
}

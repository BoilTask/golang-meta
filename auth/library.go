package auth

import (
	"github.com/golang-jwt/jwt/v5"
	metaerror "meta/meta-error"
)

func GetToken(claims jwt.Claims, secret []byte) (
	*string,
	error,
) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func ValidateJWT(token string, claims jwt.Claims, keyFunc jwt.Keyfunc) error {
	result, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		return err
	}
	if !result.Valid {
		return metaerror.New("invalid token")
	}
	return nil
}

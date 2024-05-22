package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JWTCustomClaims struct {
	UID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func NewToken(id int64, secret string) (string, error) {
	claims := &JWTCustomClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func NewJWTConfig(secret string) echojwt.Config {
	// Configure middleware with the custom claims type.
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTCustomClaims)
		},
		SigningKey:     []byte(secret),
		SuccessHandler: successHandler,
	}
}

func successHandler(c echo.Context) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return
	}

	claims := token.Claims.(*JWTCustomClaims)
	if !ok {
		return
	}

	c.Set("uid", claims.UID)
}

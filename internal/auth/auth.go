package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
)

type Claims struct {
	jwt.RegisteredClaims
	User string
}

func NewToken(user, password string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(model.Expire)),
		},
		User: user,
	})

	auth, err := token.SignedString([]byte(password))
	if err != nil {
		return "", err
	}

	return auth, nil
}

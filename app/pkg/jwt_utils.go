package pkg

import (
	"fmt"
	"os"

	//"pm_go_version/app/domain/dto"
	"pm_go_version/app/domain/entity"

	//"pm_go_version/app/domain/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret []byte = []byte(os.Getenv("SECRET"))

type CustomClaims struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(request *entity.User) (string, error) {
	claims := CustomClaims{
		ID:       request.ID,
		UserName: request.UserName,
		Email:    request.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(JwtSecret)
}

func Verfiy(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

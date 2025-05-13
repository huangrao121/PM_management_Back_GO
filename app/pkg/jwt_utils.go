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
var TempJwtSecret []byte = []byte(os.Getenv("TEMP_SECRET"))

type CustomClaims struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type TempClaims struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Sub      string `json:"sub"`
	Provider string `json:"provider"`
	jwt.RegisteredClaims
}

func GenerateTempJWT(email, name, sub, provider string) (string, error) {
	claims := TempClaims{
		Email:    email,
		Name:     name,
		Sub:      sub,
		Provider: provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(TempJwtSecret)
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

func Verify(tokenString string) (*CustomClaims, error) {
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

func VerifyTemp(tokenString string) (*TempClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TempClaims{}, func(t *jwt.Token) (interface{}, error) {
		return TempJwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*TempClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

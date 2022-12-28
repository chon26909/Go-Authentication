package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	JwtSecretKey     = []byte(os.Getenv("JWT_SECRET_KEY"))
	JwtSigningMethod = jwt.SigningMethodHS256
)

func NewToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(JwtSigningMethod, claims)
	return token.SignedString(JwtSecretKey)
}

// func validateSignMethod(token *jwt.Token) ([]byte, error) {
// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 	}

// 	return JwtSecretKey, nil
// }

// func ParseToken(tokenString string) (*jwt.StandardClaims, error) {
// 	claims := new(jwt.StandardClaims)
// 	token, err := jwt.ParseWithClaims(tokenString, claims, validateSignMethod)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var ok bool
// 	claims, ok = token.Claims.(*jwt.StandardClaims)
// 	if !ok || !token.Valid {
// 		return nil, util.ErrInvalidAuthToken
// 	}

// 	return claims, nil
// }

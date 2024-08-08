package Infrastructure

import (
	"main/Domain"
	"main/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func SignJwt(existingUser Domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  existingUser.ID,
		"email":    existingUser.Email,
		"is_admin": existingUser.Is_Admin,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
		// for the purpose of this task, we will set the expiration time to 1 minute
	})

	jwtToken, err := token.SignedString(config.JwtSecret)
	return jwtToken, err
}

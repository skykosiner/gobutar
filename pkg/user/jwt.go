package user

import "github.com/golang-jwt/jwt/v5"

var secretKey = []byte("kv'qwjmjvzqwhumeotnpuecrl&yf[{)}]*&(=[{]*}=)ga")

func CreateJWT(email string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	}).SignedString(secretKey)
}

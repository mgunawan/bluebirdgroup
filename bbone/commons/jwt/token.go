package jwt

import "github.com/dgrijalva/jwt-go"

//ExtractClaims ...
func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	token, _ := jwt.Parse(tokenStr, nil)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, true
	}
	return nil, false
}

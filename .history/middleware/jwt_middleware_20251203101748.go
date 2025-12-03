package middleware

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/*
	getSecretKey
	Mengambil JWT secret dari environment (.env).
	Jika tidak ada, gunakan fallback default.
*/
func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "SECRET_KEY_UAS_PBE" // fallback jika .env belum dibuat
	}
	return []byte(secret)
}

/*
	GenerateJWT
	Membuat token JWT yang berisi:
	- user_id
	- role
	- expired 24 jam

	Function ini dipakai oleh AuthService saat login.
*/
func GenerateJWT(userID string, role string) (string, error) {

	// Claims = data yang akan dikirim di token
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	// Membuat token baru
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key
	return token.SignedString(getSecretKey())
}

/*
	ValidateToken
	Memvalidasi token JWT dan mengembalikan data claims.
	Dipakai oleh JWT middleware saat membaca header Authorization.
*/
func ValidateToken(tokenString string) (jwt.MapClaims, error) {

	// Parse token dan cek apakah valid
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	// Ambil claims (data di dalam token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot parse token claims")
	}

	return claims, nil
}
ctx.Set("user_id", claims["user_id"])
ctx.Set("role", claims["role"])

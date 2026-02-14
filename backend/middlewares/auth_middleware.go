package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"backend-commerce/config"
)

// Ambil secret key dari environment variable
// Jika tidak ada, gunakan default "secret_key"
var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret_key"))

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Ambil header Authorization dari request
		tokenString := c.GetHeader("Authorization")

		// Jika token kosong, kembalikan respons 401 Unauthorized
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}

		// Hapus prefix "Bearer " dari token
		// Header biasanya berbentuk: "Bearer <token>"
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Buat struct untuk menampung klaim token
		claims := &jwt.RegisteredClaims{}

		// Parse token dan verifikasi tanda tangan dengan jwtKey
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Kembalikan kunci rahasia untuk memverifikasi token
			return jwtKey, nil
		})

		// Jika token tidak valid atau terjadi error saat parsing
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthenticated",
			})
			c.Abort()
			return
		}

		// Simpan klaim "sub" (username/id) ke dalam context
		// Nanti bisa diambil di controller dengan c.Get("username")
		c.Set("username", claims.Subject)

		// Lanjut ke handler berikutnya
		c.Next()
	}
}

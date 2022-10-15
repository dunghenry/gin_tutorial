package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rs = c.Request.Header["Token"][0]
		if rs == "Bearer" || rs == "" || rs == "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failure",
				"message": "You're not authenticated",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			tokenString := strings.Split(c.Request.Header["Token"][0], " ")[1]
			if len(tokenString) == 195 {
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
				})
				if token.Valid {
					c.Next()
				} else if ve, ok := err.(*jwt.ValidationError); ok {
					if ve.Errors&jwt.ValidationErrorMalformed != 0 {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  "failure",
							"message": "Token is not valid",
						})
						c.AbortWithStatus(http.StatusUnauthorized)
					} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
						// Token is either expired or not active yet
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  "failure",
							"message": "Token is expired",
						})
						c.AbortWithStatus(http.StatusUnauthorized)
					} else {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  "failure",
							"message": "Couldn't handle this token",
						})
						c.AbortWithStatus(http.StatusUnauthorized)
					}
				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  "failure",
						"message": "Couldn't handle this token",
					})
					c.AbortWithStatus(http.StatusUnauthorized)
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failure",
					"message": "Token is not valid",
				})
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		}

	}
}

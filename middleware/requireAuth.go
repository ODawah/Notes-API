package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Notes-App/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("In middleware")
	// GET Cookie  off req
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("there's an error with the signing method")
			return nil, fmt.Errorf("there's an error with the signing method")
		}
		return []byte(os.Getenv("SECRET")), nil

	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("Token expired")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		user, err := models.FindUserByUUID(claims["sub"].(string))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user", user)
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}

package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"pkg/db/go/internal/models"
	"pkg/db/go/pkg/db"
	"time"
)

func refreshAccesToken(c *gin.Context) {

	c.SetSameSite(http.SameSiteLaxMode)

	refreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		c.SetCookie("logged", "false", 0, "/", "localhost", false, true)
	}

	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.SetCookie("logged", "false", 0, "/", "localhost", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "refresh token expired"})
		}

		var user models.UserModel
		db.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.SetCookie("logged", "false", 0, "/", "localhost", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "id invalid refreshtoken"})
		}
		ctx := context.Background()
		db.RDB.Del(ctx, refreshTokenString)

		accesTokenString, refreshTokenString, err := generateTokenPair(user.ID)
		if err != nil {
			c.SetCookie("logged", "false", 0, "/", "localhost", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "cannot regenerate tokens"})
		}

		db.RDB.Set(ctx, refreshTokenString, user.ID, 60*15)

		c.SetCookie("acces_token", accesTokenString, 60*5, "", "", false, true)
		c.SetCookie("refresh_token", refreshTokenString, 60*15, "", "", false, true)

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "refresh token invalid or no claims"})

	}
}

func ValidateToken(c *gin.Context) {

	accesTokenString, err := c.Cookie("acces_token")

	if err != nil {
		refreshAccesToken(c)
		accesTokenString, err = c.Cookie("acces_token")
	}

	token, err := jwt.Parse(accesTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		var user models.UserModel
		db.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.SetCookie("logged", "false", 0, "/", "localhost", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "id invalid accestoken"})
		}

	} else {
		c.SetCookie("logged", "false", 0, "/", "localhost", false, true)
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid accestoken or no claims"})

	}

}

func generateTokenPair(ID uint) (string, string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": ID,
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	})

	secretKey := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": ID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err

	}
	return tokenString, refreshTokenString, nil
}

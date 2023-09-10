package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"pkg/db/go/internal/models"
	"pkg/db/go/pkg/db"
	"time"
)

type ITokenRepo interface {
	GetToken(token string) (uint, error)
	CreateToken(token string, id uint, expiration time.Duration) error
}

type RedisTokenService struct {
	Repo ITokenRepo
}

func NewTokenService(repo ITokenRepo) *RedisTokenService {
	return &RedisTokenService{
		Repo: repo,
	}
}

func (svc *RedisTokenService) GenerateTokenPair(ID uint) (string, string, error) {

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

func (svc *RedisTokenService) InsertToken(token string, id uint, expiration time.Duration) error {
	err := svc.Repo.CreateToken(token, id, expiration)
	return err
}

func (svc *RedisTokenService) GetID(token string) (uint, error) {
	id, err := svc.Repo.GetToken(token)
	return id, err

}

func (svc *RedisTokenService) Validate(accesToken string, refreshToken string) (string, string, uint) {

	if accesToken == "" {
		newAccesToken, newRefreshToken, id := refreshAccesToken(refreshToken)
		return newAccesToken, newRefreshToken, id
	}

	token, _ := jwt.Parse(accesToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		var user models.UserModel
		db.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			return "", "", 0
		}

		accesTokenString, refreshTokenString, err := generateTokenPair(user.ID)

		if err != nil {
			return "", "", 0
		}
		ctx := context.Background()
		db.RDB.Set(ctx, refreshTokenString, user.ID, 60*15)

		return accesTokenString, refreshTokenString, user.ID

	} else {
		return "", "", 0

	}

}

func refreshAccesToken(refreshTokenString string) (string, string, uint) {
	ctx := context.Background()
	exists, _ := db.RDB.Exists(ctx, refreshTokenString).Result()
	if exists != 1 {
		return "", "", 0
	}

	token, _ := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return "", "", 0
		}

		var user models.UserModel
		db.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			return "", "", 0
		}
		db.RDB.Del(ctx, refreshTokenString)

		accesTokenString, refreshTokenString, err := generateTokenPair(user.ID)
		if err != nil {
			return "", "", 0
		}

		db.RDB.Set(ctx, refreshTokenString, user.ID, 60*15)

		return accesTokenString, refreshTokenString, user.ID

	} else {
		return "", "", 0

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

package services

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"intership/internal/models"
	"intership/pkg/db"
	"os"
	"time"
)

type ITokenRepo interface {
	CreateToken(token string, id primitive.ObjectID, expiration time.Duration) error
}

type RedisTokenService struct {
	Repo ITokenRepo
}

func NewTokenService(repo ITokenRepo) *RedisTokenService {
	return &RedisTokenService{
		Repo: repo,
	}
}

func (svc *RedisTokenService) GenerateTokenPair(ID primitive.ObjectID) (string, string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": ID,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	})

	secretKey := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": ID,
		"exp": time.Now().Add(time.Minute * 60).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", err

	}
	return tokenString, refreshTokenString, nil
}

func (svc *RedisTokenService) InsertToken(token string, id primitive.ObjectID, expiration time.Duration) error {
	err := svc.Repo.CreateToken(token, id, expiration)
	return err
}

func (svc *RedisTokenService) Validate(refreshToken string) (string, string, error) {
	ctx := context.Background()
	exists, err := db.RDB.Exists(ctx, refreshToken).Result()
	if exists != 1 {
		fmt.Println("0002")
		return "", "", err
	}

	token, _ := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {

			return "", "", fmt.Errorf("expired token")
		}

		var user models.UserModel
		db.DB.First(&user, claims["sub"])
		db.RDB.Del(ctx, refreshToken)

		accesTokenString, refreshTokenString, err := svc.GenerateTokenPair(user.ID)
		if err != nil {
			fmt.Println(0001)
			return "", "", err
		}

		db.RDB.Set(ctx, refreshTokenString, user.ID, 60*15)

		//err = svc.InsertToken(refreshToken, user.ID, 60*60)
		//if err != nil {
		//	fmt.Println(0000)
		//	return "", "", err
		//}

		return accesTokenString, refreshTokenString, nil

	} else {
		return "", "", err

	}

}

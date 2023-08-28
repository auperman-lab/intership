package service

import (
	"github.com/golang-jwt/jwt"
	"os"
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

func (rsvc *RedisTokenService) GenerateTokenPair(ID uint) (string, string, error) {

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

func (rsvc *RedisTokenService) InsertToken(token string, id uint, expiration time.Duration) error {
	err := rsvc.Repo.CreateToken(token, id, expiration)
	return err
}

func (rsvc *RedisTokenService) GetID(token string) (uint, error) {
	id, err := rsvc.Repo.GetToken(token)
	return id, err

}

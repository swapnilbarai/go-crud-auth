package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"swapnilbarai/go-crud-auth/models"
	"swapnilbarai/go-crud-auth/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func IsEmailValid(email string) bool {
	return emailRegex.MatchString(email)
}

func IsMobileNumberValid(mobileNumber string) bool {
	return mobileNumberRegex.MatchString(mobileNumber)
}

func FormatInvalidMeesage(field, value string) string {
	return fmt.Sprintf("Provided %s : %s is not in correct format\n", field, value)
}

func HashPassword(password string) string {

	hasher := sha256.New()
	hasher.Write([]byte(password))

	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateJWTToken(userName string, admin bool, tokenType string, tokenDuration int, issueAt int64) (string, error) {

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userName,
			"iss": "go-todo-app",
			"exp": time.Now().Add(time.Duration(tokenDuration)).Unix(),
			"iat": issueAt,
			"aud": admin,
		})
	secret := AccessSecret
	if tokenType == utils.RefrehTokenType {
		secret = RefreshSecret
	}
	accessTokenString, err := accessToken.SignedString(secret)
	if err == nil {
		models.InsertToken(issueAt, userName, tokenType)
	}
	return accessTokenString, err
}

func VerifyJWTToken(token, tokenType string) (*jwt.Token, error) {

	secret := AccessSecret
	if tokenType == utils.RefrehTokenType {
		secret = RefreshSecret
	}
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err == nil {
		claims := jwtToken.Claims
		issueAt, issueError := claims.GetIssuedAt()
		if issueError != nil {
			return nil, issueError
		}
		subject, subjectError := claims.GetSubject()
		if subjectError != nil {
			return nil, subjectError
		}
		if okay := models.IsTokenValid(issueAt.Unix(), subject, tokenType); !okay {
			return nil, errors.New("invalid token")
		}
		if !jwtToken.Valid {
			models.InvalidateToken(issueAt.Unix(), subject, tokenType)
		}
	}
	return jwtToken, err
}

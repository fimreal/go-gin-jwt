package rest

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(userID int, expire int64, secret string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["expire"] = expire
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return at.SignedString([]byte(secret))
}

func CreateTokens(userID int) (string, string, error) {

	expire := time.Now().Add(ExpTime).Unix()
	accessToken, err := CreateToken(userID, expire, "access"+Secret)
	if err != nil {
		return "", "", err
	}
	refresh := time.Now().Add(RefreshTime).Unix()
	refreshToken, err := CreateToken(userID, refresh, "refresh"+Secret)

	return accessToken, refreshToken, err
}

func NullStr2Str(str string) (nullStr sql.NullString) {
	if str == "" {
		nullStr.String = ""
		nullStr.Valid = false
	} else {
		nullStr.String = str
		nullStr.Valid = true
	}
	return
}

// sha256 to string
func EncodePassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	md := hash.Sum(nil)
	encodedPassword := hex.EncodeToString(md)
	return encodedPassword
}

func DecodeToken(tokenString string, secret string) (bool, jwt.MapClaims, error) {
	Claims := jwt.MapClaims{}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	}
	token, err := jwt.ParseWithClaims(tokenString, Claims, keyFunc)
	return token.Valid, Claims, err
}

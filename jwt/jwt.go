package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"context"
	"log"
	"mywon/students_reports/graph/model"
	"net/http"
	"os"
	"strings"
	"time"
)
var contextKey = "contextMap"
var Secret string
var SecretKey []byte
func setSecret() {
	Secret = os.Getenv("SECRET")
	SecretKey = []byte(Secret)
}

type Payload struct {
	Id int `json:"id"`
	UserName string `json:"user_name"`
	AuthToken string `json:"auth_token"`
	Jwt string `json:"jwt"`
}

var contextMap Payload

func GenerateToken(jwtDetails model.GenerateJWTDetails) (string, error) {
	setSecret()
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["id"] = jwtDetails.Id
	claims["username"] = jwtDetails.Username
	claims["role"] = jwtDetails.Role
	//claims["isAdmin"] = jwtDetails.IsAdmin
	claims["rights"] = jwtDetails.Rights
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

func Auth()func(http.Handler)http.Handler{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					contextMap.AuthToken = "INVALID"
					ctx := context.WithValue(r.Context(), contextKey,contextMap)
					r = r.WithContext(ctx)
					next.ServeHTTP(w, r)
					return
				}
			}()
			setSecret()
			authToken := r.Header.Get("authorization")
			split := strings.Split(authToken, "")
			if len(split) > 1 {
				authToken = split[1]
			} else {
				authToken = split[0]
			}
			contextMap.Jwt = authToken
			if authToken == ""{
				contextMap.AuthToken = "INVALID"
				ctx :=context.WithValue(r.Context(),contextKey,contextMap)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
			_, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return SecretKey, nil
			})

			if err != nil {
				log.Println("err", err)
				contextMap.AuthToken = "INVALID"
				ctx := context.WithValue(r.Context(), contextKey, contextMap)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		})
	}
}

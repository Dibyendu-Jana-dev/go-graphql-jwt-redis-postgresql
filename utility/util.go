package utility

import (
	b64 "encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func RecoverFromPanic(panic string) {
	if r := recover(); r != nil {
		log.Printf("panic: %+v, recovered from %s\n", r, panic)
	}
}

func CheckAscOrDesc(orderBy *string) string {
	if orderBy == nil {
		return "asc"
	}
	order := strings.TrimSpace(*orderBy)
	if strings.EqualFold(order, "desc") {
		return "desc"
	} else if strings.EqualFold(order, "asc") {
		return "asc"
	} else {
		return "asc"
	}
}

func Base64Encoder(input string) string {
	encodedString := b64.URLEncoding.EncodeToString([]byte(input))
	return encodedString
}

func Base64Decoder(input string) string {
	decodedString, err := b64.URLEncoding.DecodeString(input)
	if err != nil {
		log.Println("error while decoding the string : ", input)
	}

	return string(decodedString)
}

func PasswordHashAndSalt(Password string) (string, error) {
	pwd := []byte(Password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

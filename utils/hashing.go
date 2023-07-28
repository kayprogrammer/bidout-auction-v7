package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	log.Println(err)
    return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
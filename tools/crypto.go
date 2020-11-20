/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */


package tools

import (
    "math/rand"
    "strings"
    "errors"

    "github.com/GehirnInc/crypt"
    _ "github.com/GehirnInc/crypt/sha256_crypt"

)

func RandString(n int) string {
    const letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    arr := make([]byte, n)
    for i := range arr {
        arr[i] = letters[rand.Intn(len(letters))]
    }
    return string(arr)
}

func CreateHash(key string) (string, error) {
    crypt := crypt.SHA256.New()
    return crypt.Generate([]byte(key), []byte("$5$" + RandString(12)))
}

func CheckHash(hash, password string) error {
    arr := strings.Split(hash, "$")
    if len(arr) < 3 {
        return errors.New("incorrect hash structure")
    }
    hashType := arr[1]
    hashSalt := arr[2]

    if hashType != "5" {
        return errors.New("incorrect hash type")
    }
    crypt := crypt.SHA256.New()
    newHash, _ := crypt.Generate([]byte(password), []byte("$" + hashType + "$" + hashSalt))

    if hash != newHash {
        return errors.New("password is incorrect")
    }
    return nil
}

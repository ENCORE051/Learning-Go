package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func TestCustomClaims(t *testing.T) string {
	secret := []byte("mysecretkey")

	// Create a new token object, specifying signing method and the claims
	claims := MyClaims{
		User: "jack051",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Token expires in 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "myapp",
		},
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	fmt.Println("Token:", token)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secret)
	if err != nil {
		t.Fatal("Error signing token:", err)
	}
	fmt.Println("Signed Token:", signedToken)
	return signedToken
}

func TestParseCustomClaims(t *testing.T) {
	secret := []byte("mysecretkey")

	// 假设通过HS256算法创建并签名生成了一个token  这个token是通过上面的TestCustomClaims生成的signedToken
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiamFjayIsImlzcyI6Im15YXBwIiwiZXhwIjoxNzk0Mjk5NzM4LCJuYmYiOjE3OTQyMTMzMzgsImlhdCI6MTc5NDIxMzMzOH0.d0aOCZZSfzwL3aDDVe1aKItzaQ3sR7ZigOoHLGoC8B4"
	// tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzNDU2LCJuYW1lIjoiamFjayJ9.QxLw9NkFgZW3BluyXIofe4efp1IAy61s8b2fe3Eo86M"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法是否为HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		t.Fatal("Error parsing token:", err)
	}
	if !token.Valid {
		t.Fatal("Token is not valid")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	} else {
		fmt.Println(err)
	}
}

func TestProcess(t *testing.T) {
	secret := []byte("my secret")
	// 假设通过HS256算法创建并签名生成了一个token
	// 这个token是验证原始的claims结构体的token
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzNDU2LCJuYW1lIjoiamFjayJ9.QxLw9NkFgZW3BluyXIofe4efp1IAy61s8b2fe3Eo86M"

	// 传入token字符串和验证钩子函数，返回值就是一个Token结构体
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否匹配
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不匹配的签名算法: %s", token.Header["alg"])
		}
		// 返回验证密钥
		return secret, nil
	})

	if token.Valid {
		fmt.Println("token合法")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		fmt.Println("传入的字符串甚至连一个token都不是...")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		fmt.Println("token已经过期或者还没有生效")
	} else {
		fmt.Println("token处理异常...")
	}
}

func TestRsa(t *testing.T) {

	// 创建密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey

	if err != nil {
		fmt.Println(err)
		return
	}

	// claims
	claims := MyClaims{
		User: "114514",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// 私钥加密
	signedString, err := token.SignedString(privateKey)

	fmt.Println(signedString, err)

	// 公钥解密
	token, err = jwt.ParseWithClaims(signedString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		fmt.Println(err)
	} else if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		fmt.Println(claims)
	}
}

func TestCustomClaimsParse(t *testing.T, s string) {
	secret := []byte("mysecretkey")
	tokenstring := s
	token, err := jwt.ParseWithClaims(tokenstring, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{"HS256"})) // 使用option进行验证

	// 类型断言
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		fmt.Println(claims)
	} else {
		fmt.Println(err)
	}
}

// func main() {
// 	signedToken := TestCustomClaims(nil)
// 	fmt.Println("--------------------------------------------------")
// 	TestParseCustomClaims(nil)
// 	fmt.Println("--------------------------------------------------")
// 	TestProcess(nil)
// 	fmt.Println("--------------------------------------------------")
// 	TestCustomClaimsParse(nil, signedToken)
// 	fmt.Println("--------------------------------------------------")
// 	TestRsa(nil)
// }

func sign(input, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func main() {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"user1","exp":1735689600}`))

	input := header + "." + payload
	token := input + "." + sign(input, "secret")

	parts := strings.Split(token, ".")
	fmt.Println("segments:", len(parts))
	fmt.Println("signature matches:", sign(parts[0]+"."+parts[1], "secret") == parts[2])
	fmt.Println("wrong key rejected:", sign(parts[0]+"."+parts[1], "wrong") != parts[2])
}

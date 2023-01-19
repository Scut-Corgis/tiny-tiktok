package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var mySigningKey = []byte("tiny-tiktok") // 密钥

type MyCustomClaims struct {
	Name string `json:"name"` //token中只放入了用户名
	jwt.RegisteredClaims
}

func GenerateToken(name string) string {
	fmt.Printf("generate token: %v\n", name)
	// Create the claims
	claims := MyCustomClaims{
		"name",
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if ss, err := token.SignedString(mySigningKey); err == nil {
		fmt.Println("generate token success!")
		fmt.Println("token : ", ss)
		return ss
	} else {
		println("generate token fail\n")
		return "fail"
	}
}

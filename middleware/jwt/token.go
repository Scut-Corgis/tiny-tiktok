package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var mySigningKey = []byte("tiny-tiktok") // 密钥

type MyCustomClaims struct {
	Name string `json:"name"` // 用户名
	Id   int64  `json:"id"`   // id
	jwt.RegisteredClaims
}

func GenerateToken(name string, id int64) string {
	log.Printf("generate token —— name:%v id:%v\n", name, id)
	// Create the claims
	claims := MyCustomClaims{
		name,
		id,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if ss, err := token.SignedString(mySigningKey); err == nil {
		log.Println("generate token success!")
		log.Println("token:", ss)
		return ss
	} else {
		println("generate token fail\n")
		return "fail"
	}
}

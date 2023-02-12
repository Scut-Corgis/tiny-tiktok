package jwt

import (
	"testing"
)

func TestAuth(*testing.T) {
	username := "corgis"
	token := GenerateToken(username, 1000)
	println(token)
}

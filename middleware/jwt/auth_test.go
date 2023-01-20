package jwt

import (
	"testing"
)

func TestAuth(*testing.T) {
	username := "corgis"
	token := GenerateToken(username)
	println(token)
}

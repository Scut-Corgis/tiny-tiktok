package dao

import (
	"fmt"
	"testing"
)

func TestQueryUserByName(t *testing.T) {
	Init()
	user, err := QueryUserByName("test")
	fmt.Println(user)
	fmt.Println(err)
}

func TestQueryUserById(t *testing.T) {
	Init()
	user, err := QueryUserById(1000)
	fmt.Println(user)
	fmt.Println(err)
}

func TestInsertUser(t *testing.T) {
	Init()
	user := User{
		Name:     "test",
		Password: "123",
	}
	newUser, err := InsertUser(user)
	fmt.Println(newUser)
	fmt.Println(err)
}

func TestQueryAllNames(t *testing.T) {
	Init()
	names := QueryAllNames()
	fmt.Println(names)
}

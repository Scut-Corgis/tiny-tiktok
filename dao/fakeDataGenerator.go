package dao

import (
	"github.com/brianvoe/gofakeit/v6"
)

func fakeUsers(num int) {
	gofakeit.Seed(0)
	for i := 0; i < num; i++ {
		user := User{}
		user.Name = gofakeit.Username()
		user.Password = gofakeit.Password(false, false, true, false, false, 8)
		// fmt.Println(user)
		InsertUser(&user)
	}
}

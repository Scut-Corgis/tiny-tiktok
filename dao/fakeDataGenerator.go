package dao

import (
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func FakeUsers(num int) {
	gofakeit.Seed(0)
	Db.Migrator().DropTable(&User{})
	Db.Migrator().CreateTable(&User{})
	for i := 0; i < num; i++ {
		user := User{}
		user.Name = gofakeit.Username()
		user.Password = gofakeit.Password(false, false, true, false, false, 8)
		// fmt.Println(user)
		InsertUser(&user)
	}
}

func FakeFollows(num int) {
	rand.Seed(time.Now().Unix())
	Db.Migrator().DropTable(&Follow{})
	err := Db.Migrator().CreateTable(&Follow{})
	if err != nil {
		log.Println(err.Error())
	}
	for i := 0; i < num; i++ {
		//follow := FollowTable{}
		var count int64
		Db.Model(&User{}).Count(&count)
		// fmt.Println(count)
		var a, b int64
		a = rand.Int63n(count)
		b = rand.Int63n(count)
		for a == b {
			b = rand.Int63n(count)
		}
		InsertFollow(a+1, b+1)
	}
}

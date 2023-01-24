package dao

import (
	"bytes"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func RebuildTable() bool {
	cmd := exec.Command("sh", "/Users/zaizai/Projects/GolandProjects/tiny-tiktok/config/rebuildTable.sh")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err, stderr.String())
		return false
	}
	return true
}

func FakeUsers(num int) {
	gofakeit.Seed(0)
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
	for i := 0; i < num; i++ {
		var count int64
		Db.Model(&User{}).Count(&count)
		var a, b int64
		a = rand.Int63n(count)
		b = rand.Int63n(count)
		for a == b {
			b = rand.Int63n(count)
		}
		InsertFollow(a+1001, b+1001)
	}
}

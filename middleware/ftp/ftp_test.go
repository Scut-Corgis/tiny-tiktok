package ftp

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestFfmpeg(t *testing.T) {
	Init()
	for i := 1; i <= 1; i++ {
		name := "bear.mp4"
		file, err := os.Open("/Users/zaizai/Projects/GolandProjects/tiny-tiktok/data/" + name)
		if err != nil {
			t.Fatal("Can't open the file")
		}
		err = SendVideoFile(strconv.Itoa(i), file)
		fmt.Println(err)
	}
	time.Sleep(2 * time.Second)
}

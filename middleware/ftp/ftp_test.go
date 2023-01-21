package ftp

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestFfmpeg(t *testing.T) {
	Init()
	for i := 1; i <= 20; i++ {
		file, err := os.Open("/home/hjg/go/src/tiny-tiktok/middleware/ftp/ftp_test.go")
		if err != nil {
			t.Fatal("打开不了文件")
		}
		go SendVideoFile(strconv.Itoa(i), file)
	}
	time.Sleep(2 * time.Second)
}

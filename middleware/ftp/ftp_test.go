package ftp

import (
	"os"
	"testing"
)

func TestFfmpeg(t *testing.T) {
	Init()
	file, err := os.Open("/home/hjg/go/src/tiny-tiktok/middleware/ftp/ftp_test.go")
	if err != nil {
		t.Fatal("打开不了文件")
	}
	err = SendVideoFile("ftp测试", file)
	if err != nil {
		t.Fatal("文件发送失败")
	}
}

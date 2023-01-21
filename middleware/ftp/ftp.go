package ftp

import (
	"io"
	"log"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"github.com/jlaffaye/ftp"
)

var ClientFtp *ftp.ServerConn

func Init() {
	var err error
	ClientFtp, err = ftp.Dial(config.Ftp_addr_port, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = ClientFtp.Login(config.Ftp_username, config.Ftp_password)
	if err != nil {
		log.Fatal(err)
	}

}

func SendVideoFile(filename string, file io.Reader) error {
	err := ClientFtp.Stor("./videos/"+filename+".mp4", file)
	if err != nil {
		log.Fatalln("视频发送失败 : ", filename)
	} else {
		log.Println("视频发送成功！")
	}
	return err
}

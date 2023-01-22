package ftp

import (
	"io"
	"log"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"github.com/jlaffaye/ftp"
)

// 原ftp开源库无法单连接实现并发安全，因此本项目实现了并发ftp传输
var FtpChan chan *ftp.ServerConn
var ftpConnList [config.Ftp_max_concurrent_cnt]*ftp.ServerConn

func Init() {
	FtpChan = make(chan *ftp.ServerConn, 20)
	var err error
	for _, conn := range ftpConnList {
		conn, err = ftp.Dial(config.Ftp_addr_port, ftp.DialWithTimeout(5*time.Second))
		if err != nil {
			log.Fatal(err)
		}

		err = conn.Login(config.Ftp_username, config.Ftp_password)
		if err != nil {
			log.Fatal(err)
		}
		FtpChan <- conn
	}
	log.Println("ftp初始化成功！")
	go keepalive()
}

func SendVideoFile(filename string, file io.Reader) error {
	conn := <-FtpChan
	err := conn.Stor("./videos/"+filename+".mp4", file)
	if err != nil {
		log.Fatalln("视频发送失败 : ", filename, "Error : ", err)
	} else {
		log.Println("视频发送成功！")
	}
	FtpChan <- conn
	return err
}

// vsftpd 配置为永不断开空闲连接， 因此keepalive未实现
func keepalive() {

}

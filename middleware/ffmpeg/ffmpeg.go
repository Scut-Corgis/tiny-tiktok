package ffmpeg

import (
	"bytes"
	"log"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"golang.org/x/crypto/ssh"
)

type Ffmsg struct {
	VideoName string
	ImageName string
}

var ClientSSH *ssh.Client
var Ffchan chan Ffmsg

func Init() {
	var err error
	clientConfig := &ssh.ClientConfig{
		Timeout: time.Second * 5, //最多5s建立连接
		User:    config.Ssh_username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Ssh_password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //不检查公钥
	}
	ClientSSH, err = ssh.Dial("tcp", config.Ssh_addr_port, clientConfig)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	Ffchan = make(chan Ffmsg, config.Ssh_max_taskCnt)
	go dispatcher()

	log.Println("ssh初始化成功")
}

// 将ffmpeg命令以及其参数通过ssh运行，若运行失败，则重新放入channel运行
func dispatcher() {
	for ffmsg := range Ffchan {
		go func(f Ffmsg) {
			err := Ffmpeg(f.VideoName, f.ImageName)
			if err != nil {
				Ffchan <- f
				log.Println("ffmpeg调用失败，等待重新执行")
			}
		}(ffmsg)
	}
}

func Ffmpeg(videoName, imageName string) error {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := ClientSSH.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	// 调试代码
	if err := session.Run("ls"); err != nil {
		log.Println("SSH failed to run: " + err.Error())
		log.Println("远端输出错误信息 : ", b.String())
		return err
	}
	log.Println("远端输出list信息 : ", b.String())

	// if err := session.Run("/usr/local/ffmpeg/bin/ffmpeg -ss 00:00:01 -i " + config.Ftp_video_path + videoName + ".mp4 -vframes 1 " + config.Ftp_image_path + imageName + ".jpg"); err != nil {
	// 	log.Println("SSH failed to run: " + err.Error())
	// 	log.Println("远端输出错误信息 : ", b.String())
	// 	return err
	// }
	log.Println("ffmpeg 执行成功")
	return nil
}

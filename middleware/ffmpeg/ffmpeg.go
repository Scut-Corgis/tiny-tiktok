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
	go keepalive()
	log.Println("ssh init successfully!")
}

// 将ffmpeg命令以及其参数通过ssh运行，若运行失败，则重新放入channel运行
func dispatcher() {
	for ffmsg := range Ffchan {
		go func(f Ffmsg) {
			err := Ffmpeg(f.VideoName, f.ImageName)
			if err != nil {
				Ffchan <- f
				log.Println("ffmpeg call failed, wait for re-execute")
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

	if err := session.Run("ffmpeg -ss 00:00:01 -i " + config.Ftp_video_path + videoName + ".mp4 -vframes 1 " + config.Ftp_image_path + imageName + ".jpg"); err != nil {
		log.Println("SSH failed to run: ", err.Error())
		log.Println("remote fail message: ", b.String())
		return err
	}
	log.Println("ffmpeg start successfully")
	return nil
}

// openssl 默认为60s断开连接， 因此设置10s一次心跳
func keepalive() {
	for {
		s, err := ClientSSH.NewSession()
		if err != nil {
			log.Println("ssh is disconnected err: ", err)
		}
		time.Sleep(10 * time.Second)
		s.Close()
	}
}

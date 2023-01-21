package ffmpeg

import (
	"testing"
	"time"
)

func TestFfmpeg(t *testing.T) {
	Init()
	Ffchan <- Ffmsg{
		VideoName: "bear",
		ImageName: "bear",
	}
	time.Sleep(20 * time.Second)
	Ffchan <- Ffmsg{
		VideoName: "bear",
		ImageName: "bear2",
	}
	time.Sleep(2 * time.Second)
}

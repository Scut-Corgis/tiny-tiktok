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
	time.Sleep(2 * time.Second)
}

package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/config"
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/ffmpeg"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/ftp"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}

// Publish save upload file to ftp server
func Publish(c *gin.Context) {
	username := c.GetString("username")
	title := c.PostForm("title")

	user, err := dao.QueryUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	timeToDB := time.Now()
	pubilishTImeIntNano := timeToDB.UnixNano()
	publishTimeStr := strconv.FormatInt(pubilishTImeIntNano, 10)
	//	timeStr := time.Now().Format("2006-01-02 15:04:05")
	videoName := username + "_" + publishTimeStr
	imageName := username + "_" + publishTimeStr

	videoFile, err := data.Open()
	if err != nil {
		log.Fatalln("pulish意料之外的错误， data.Open()失败")
	}
	//ftp发送视频文件
	err = ftp.SendVideoFile(videoName, videoFile)
	if err != nil {
		log.Println(username, "的视频ftp失败")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "视频ftp失败！"})
		return
	}
	//插入数据库
	video := &dao.Video{
		AuthorId:    user.Id,
		PlayUrl:     config.Url_addr_port + config.Url_Play_prefix + videoName + ".mp4",
		CoverUrl:    config.Url_addr_port + config.Url_Image_prefix + imageName + ".jpg",
		PublishTime: timeToDB,
		Title:       title,
	}
	err = dao.InsertVideosTable(video)
	if err != nil {
		log.Println(username, "的视频，数据库插入失败")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "视频数据库插入失败！"})
		return
	} else {
		log.Println("视频数据库插入成功")
	}
	// ssh调用ffmpeg
	ffmpeg.Ffchan <- ffmpeg.Ffmsg{
		VideoName: videoName,
		ImageName: imageName,
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}

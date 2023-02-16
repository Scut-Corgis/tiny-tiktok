package controller

import (
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/service"
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

// Feed GET/douyin/feed/ 视频流接口
func Feed(c *gin.Context) {
	vsi := service.VideoServiceImpl{}
	queryUserId := c.GetInt64("id")
	latestTimeStr := c.Query("latest_time")

	if latestTimeStr == "" {
		latestTimeStr = strconv.FormatInt(time.Now().Unix(), 10)
	}
	latestTimeInt, _ := strconv.ParseInt(latestTimeStr, 10, 64)
	// 时间戳校准
	if latestTimeInt > time.Now().Unix() {
		latestTimeInt = time.Now().Unix()
	}
	latestTime := time.Unix(latestTimeInt, 0)
	videoIdList := vsi.GetMost30videosIdList(latestTime)

	var videoList = make([]Video, 0, len(videoIdList))
	var nextTimeInt = time.Now().Unix()
	for _, videoId := range videoIdList {
		videoDetail, publishTime := vsi.QueryVideoDetailByVideoId(videoId, queryUserId)
		publishTimeInt := publishTime.Unix()
		if publishTimeInt < nextTimeInt {
			nextTimeInt = publishTimeInt
		}
		video := Video{
			Id:            videoDetail.Id,
			Author:        User(videoDetail.Author),
			PlayUrl:       videoDetail.PlayUrl,
			CoverUrl:      videoDetail.CoverUrl,
			FavoriteCount: videoDetail.FavoriteCount,
			CommentCount:  videoDetail.CommentCount,
			IsFavorite:    videoDetail.IsFavorite,
			Title:         videoDetail.Title,
		}
		videoList = append(videoList, video)
	}
	log.Println(videoList)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTimeInt,
	})
}

// Publish POST/douyin/publish/action/ 投稿接口
func Publish(c *gin.Context) {
	vsi := service.VideoServiceImpl{}
	username := c.GetString("username")
	id := c.GetInt64("id")
	title := c.PostForm("title")

	if !redis.CuckooFilterUserName.Contain([]byte(username)) {
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
	videoName := username + "_" + publishTimeStr
	imageName := username + "_" + publishTimeStr

	videoFile, err := data.Open()
	if err != nil {
		log.Fatalln("data open failed")
	}
	// ftp 发送视频文件
	err = ftp.SendVideoFile(videoName, videoFile)
	if err != nil {
		log.Println(username, "video ftp failed")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video ftp failed!"})
		return
	}
	//插入数据库
	video := &dao.Video{
		AuthorId:    id,
		PlayUrl:     "http://" + config.Url_addr + config.Url_Play_prefix + videoName + ".mp4",
		CoverUrl:    "http://" + config.Url_addr + config.Url_Image_prefix + imageName + ".jpg",
		PublishTime: timeToDB,
		Title:       title,
	}
	flag := vsi.InsertVideosTable(video)
	if !flag {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Insert video failed!"})
	} else {
		log.Println("Insert video successfully!")
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

// PublishList GET/douyin/publish/list/ 发布列表
func PublishList(c *gin.Context) {
	vsi := service.VideoServiceImpl{}
	queryUserId := c.GetInt64("id")
	userIdStr := c.Query("user_id")
	authorId, _ := strconv.ParseInt(userIdStr, 10, 64)
	var videoList = make([]Video, 0)
	videoIdList := dao.GetVideoIdListByUserId(authorId)
	for _, videoId := range videoIdList {
		videoDetail, _ := vsi.QueryVideoDetailByVideoId(videoId, queryUserId)
		video := Video{
			Id:            videoDetail.Id,
			Author:        User(videoDetail.Author),
			PlayUrl:       videoDetail.PlayUrl,
			CoverUrl:      videoDetail.CoverUrl,
			FavoriteCount: videoDetail.FavoriteCount,
			CommentCount:  videoDetail.CommentCount,
			IsFavorite:    videoDetail.IsFavorite,
			Title:         videoDetail.Title,
		}
		videoList = append(videoList, video)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

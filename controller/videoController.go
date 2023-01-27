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
	username := c.GetString("username")
	var queryUserId int64
	if username == "" {
		queryUserId = -1
	} else {
		queryUser, _ := dao.QueryUserByName(username)
		queryUserId = queryUser.Id
	}

	latestTimeStr := c.Query("latest_time")
	if latestTimeStr == "" {
		latestTimeStr = strconv.FormatInt(time.Now().Unix(), 10)
	}
	latestTimeInt, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		log.Fatalln("timeStr 转 timeInt 出现了意料之外的错误")
	}
	latestTime := time.Unix(latestTimeInt, 0)
	videoIdList := dao.GetMost30videosIdList(latestTime)
	log.Println(videoIdList)
	var videoList []Video = make([]Video, 0, len(videoIdList))
	var nextTimeInt int64 = time.Now().Unix()
	for _, videoId := range videoIdList {
		videoDetail, publishTime := dao.QueryVideoDetailByVideoId(videoId, queryUserId)
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
		log.Println(video)
		videoList = append(videoList, video)
	}
	log.Println(videoList)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTimeInt,
	})
}

// Publish save upload file to ftp server
func Publish(c *gin.Context) {
	username := c.GetString("username")
	title := c.PostForm("title")

	user, err := dao.QueryUserByName(username)
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
	username := c.GetString("username")
	queryUser, _ := dao.QueryUserByName(username)
	queryUserId := queryUser.Id
	userIdStr := c.Query("user_id")
	authorId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		log.Fatalln("strconv.ParseInt(userIdStr, 10, 64) 失败")
	}
	var videoList []Video = make([]Video, 0)
	videoIdList := dao.GetVideoIdListByUserId(authorId, queryUserId)

	for _, videoId := range videoIdList {
		videoDetail, _ := dao.QueryVideoDetailByVideoId(videoId, queryUserId)
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

package service

import (
	"fmt"
	"log"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

type Video struct {
	Author        dao.UserTable `json:"author"`         // 视频作者信息
	CommentCount  int64         `json:"comment_count"`  // 视频的评论总数
	CoverURL      string        `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64         `json:"favorite_count"` // 视频的点赞总数
	ID            int64         `json:"id"`             // 视频唯一标识
	IsFavorite    bool          `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string        `json:"play_url"`       // 视频播放地址
	Title         string        `json:"title"`          // 视频标题
}

// 本模块使用：
// 点赞--在likes表中插入数据
func Like(userId int64, videoId int64) error {
	likeData := dao.Like{UserId: userId, VideoId: videoId}
	return dao.InsertLike(likeData)
}

// 取消点赞--在likes表中删除数据
func Unlike(userId int64, videoId int64) error {
	return dao.DeleteLike(userId, videoId)
}

// 获取视频
func GetVideo(videoId int64, userId int64) (Video, error) {
	//视频信息
	var video Video
	VideoData, err := dao.GetVideoByVideoId(videoId)
	if err != nil {
		log.Printf("方法dao.GetVideoByVideoId(videoId) 失败：%v", err)
		return video, err
	} else {
		log.Printf("方法dao.GetVideoByVideoId(videoId) 成功")
	}
	video.ID = VideoData.Id             // 视频唯一标识
	video.PlayURL = VideoData.PlayUrl   // 视频播放地址
	video.CoverURL = VideoData.CoverUrl // 视频封面地址
	video.Title = VideoData.Title       // 视频标题

	// 视频作者信息
	AuthorInfo, err := dao.QueryUserTableById(VideoData.AuthorId)
	if err != nil {
		log.Printf("dao.QueryUserTableById(VideoData.AuthorId)失败：%v", err)
		return video, err
	} else {
		log.Printf("方法dao.QueryUserTableById(VideoData.AuthorId)成功")
	}
	video.Author = AuthorInfo

	// 视频的点赞总数
	LikeUserIdList, err := dao.GetLikeUserIdList(videoId)
	if err != nil {
		log.Printf("方法dao.GetLikeUserIdList(videoId) 失败：%v", err)
		return video, err
	} else {
		log.Printf("方法dao.GetLikeUserIdList(videoId)成功")
	}
	video.FavoriteCount = int64(len(LikeUserIdList))

	// 视频的评论总数
	var temp int64
	temp = 0
	video.CommentCount = temp

	// true-已点赞，false-未点赞
	video.IsFavorite = true

	return video, nil
}

// 获取点赞列表
func GetLikeLists(userId int64) ([]Video, error) {
	var result []Video

	videoIdList, err := dao.GetLikeVideoIdList(userId)
	if err != nil {
		fmt.Println("获取用户点赞视频Id出错！")
	}
	for _, videoId := range videoIdList {
		video, err := GetVideo(videoId, userId)
		if err != nil {
			fmt.Println("获取用户视频数据出错！")
			return result, err
		}
		result = append(result, video)
	}
	return result, nil
}

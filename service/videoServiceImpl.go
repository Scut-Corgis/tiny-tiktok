package service

import (
	"log"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

type VideoServiceImpl struct {
	CommentService
}

func (VideoServiceImpl) QueryVideoById(id int64) dao.Video {
	video, err := dao.QueryVideoById(id)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("Video not found!")
		return video
	}
	log.Println("Query video successfully!")
	return video
}

func (VideoServiceImpl) QueryVideoDetailByVideoId(videoId int64, queryUserId int64) (dao.VideoDetail, time.Time) {
	usi := UserServiceImpl{}
	csi := CommentServiceImpl{}
	videoDetail := dao.VideoDetail{}
	video, err := dao.QueryVideoById(videoId)
	if err != nil {
		log.Println("Video not found!")
		return videoDetail, time.Time{}
	}
	userResp, err := usi.QueryUserRespById(video.AuthorId)
	if err != nil {
		log.Println("QueryUser not found!")
		return videoDetail, time.Time{}
	}
	userResp.IsFollow = usi.JudgeIsFollowById(queryUserId, userResp.Id)
	videoDetail.Id = video.Id
	videoDetail.Author = userResp
	videoDetail.PlayUrl = video.PlayUrl
	videoDetail.CoverUrl = video.CoverUrl
	videoDetail.Title = video.Title
	videoDetail.FavoriteCount, _ = usi.LikeCount(videoId)
	videoDetail.CommentCount, _ = csi.CountComments(videoId)
	videoDetail.IsFavorite, _ = usi.IsLike(videoId, queryUserId)
	return videoDetail, video.PublishTime
}

func (VideoServiceImpl) GetVideoIdListByUserId(id int64) []int64 {
	return dao.GetVideoIdListByUserId(id)
}

func (VideoServiceImpl) GetMost30videosIdList(latestTime time.Time) []int64 {
	return dao.GetMost30videosIdList(latestTime)
}

func (VideoServiceImpl) InsertVideosTable(video *dao.Video) bool {
	err := dao.InsertVideosTable(video)
	return err == nil
}

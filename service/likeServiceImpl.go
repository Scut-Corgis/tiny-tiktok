package service

import (
	"log"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

type LikeServiceImpl struct {
	UserService //出错
}

func (LikeServiceImpl) Like(userId int64, videoId int64) error {
	IsFavorite := dao.JudgeIsFavorite(userId, videoId)
	//如果点赞关系已经存在，则直接返回
	if IsFavorite {
		return nil
	}
	likeData := dao.Like{UserId: userId, VideoId: videoId}
	return dao.InsertLike(&likeData)
}

func (LikeServiceImpl) Unlike(userId int64, videoId int64) error {
	return dao.DeleteLike(userId, videoId)
}

func (LikeServiceImpl) GetVideo(videoId int64, userId int64) dao.VideoDetail {
	VideoDetail, _ := dao.QueryVideoDetailByVideoId(videoId, userId)
	return VideoDetail
}

func (LikeServiceImpl) GetLikeLists(userId int64) ([]dao.VideoDetail, error) {
	var result []dao.VideoDetail

	videoIdList, err := dao.GetLikeVideoIdList(userId)
	if err != nil {
		log.Println("获取用户点赞视频Id出错！")
	}
	//fmt.Println(videoIdList)
	for _, videoId := range videoIdList {
		video := GetVideo(videoId, userId)
		result = append(result, video)
	}
	return result, nil
}

func (LikeServiceImpl) IsFavourite(videoId int64, userId int64) (bool, error) {
	likeinfon, err := dao.GetLikInfo(videoId, userId)
	if err != nil {
		log.Println("get likeInfo failed")
		return false, nil
	}
	if likeinfon.UserId != userId || likeinfon.VideoId != videoId {
		return false, nil
	}
	return true, nil
}

func (LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	return dao.GetLikeCountByVideoId(videoId)
}

package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

/*
本模块使用：
*/

type LikeService interface {
	/*其他模块使用*/
	//判断用户userId是否点赞视频videoId
	IsLike(videoId int64, userId int64) (bool, error)
	//获取视频videoId的点赞数
	LikeCount(videoId int64) (int64, error)
	//获取用户userId发布视频的总被赞数
	TotalLiked(userId int64) int64
	//获取用户userId喜欢的视频数量
	LikeVideoCount(userId int64) (int64, error)

	/*点赞模块使用*/
	// 功能：点赞--在likes表中插入数据
	Like(userId int64, videoId int64) error
	// 功能：取消点赞--在likes表中删除数据
	Unlike(userId int64, videoId int64) error
	// 功能：获取视频详细信息
	//GetVideo(videoId int64, userId int64) dao.VideoDetail
	// 功能：获取点赞列表
	GetLikeLists(userId int64) ([]dao.VideoDetail, error)
}

// // 功能：点赞--在likes表中插入数据
// // 实现：如果点赞关系已经在likes表中存在，则直接返回；否则将点赞关系插入到likes表中
// func Like(userId int64, videoId int64) error {
// 	IsFavorite := dao.JudgeIsFavorite(userId, videoId)
// 	//如果点赞关系已经存在，则直接返回
// 	if IsFavorite {
// 		return nil
// 	}
// 	likeData := dao.Like{UserId: userId, VideoId: videoId}
// 	return dao.InsertLike(&likeData)
// }

// // 功能：取消点赞--在likes表中删除数据
// func Unlike(userId int64, videoId int64) error {
// 	return dao.DeleteLike(userId, videoId)
// }

// // 功能：获取视频详细信息
// func GetVideo(videoId int64, userId int64) dao.VideoDetail {
// 	vsi := VideoServiceImpl{}
// 	VideoDetail, _ := vsi.QueryVideoDetailByVideoId(videoId, userId)
// 	return VideoDetail
// }

// // 功能：获取点赞列表
// func GetLikeLists(userId int64) ([]dao.VideoDetail, error) {
// 	var result []dao.VideoDetail

// 	videoIdList, err := dao.GetLikeVideoIdList(userId)
// 	if err != nil {
// 		fmt.Println("获取用户点赞视频Id出错！")
// 	}
// 	//fmt.Println(videoIdList)
// 	for _, videoId := range videoIdList {
// 		video := GetVideo(videoId, userId)
// 		result = append(result, video)
// 	}
// 	return result, nil
// }

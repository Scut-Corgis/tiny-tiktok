package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

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

	// 功能：点赞
	Like(userId int64, videoId int64) error
	// 功能：取消点赞
	Unlike(userId int64, videoId int64) error
	// 功能：获取视频详细信息
	//GetVideo(videoId int64, userId int64) dao.VideoDetail
	// 功能：获取点赞列表
	GetLikeLists(userId int64) ([]dao.VideoDetail, error)
}

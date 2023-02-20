package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
	"log"
	"strconv"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

type VideoServiceImpl struct {
	CommentService
}

func (VideoServiceImpl) QueryVideoById(id int64) (dao.Video, error) {
	video, err := dao.QueryVideoById(id)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("Video not found!")
		return video, err
	}
	log.Println("Query video successfully!")
	return video, nil
}

func (VideoServiceImpl) QueryVideoDetailByVideoId(videoId int64, queryUserId int64) (dao.VideoDetail, time.Time) {
	usi := UserServiceImpl{}
	csi := CommentServiceImpl{}
	videoDetail := dao.VideoDetail{}
	video, err1 := dao.QueryVideoById(videoId)
	if err1 != nil {
		log.Println("Video not found!")
		return videoDetail, time.Time{}
	}
	userResp, err2 := usi.QueryUserRespById(video.AuthorId)
	if err2 != nil {
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
	// 添加到过滤器
	redis.CuckooFilterVideoId.Add([]byte(strconv.FormatInt(video.Id, 10)))
	err := dao.InsertVideosTable(video)
	return err == nil
}

func (VideoServiceImpl) CountWorks(id int64) int64 {
	redisUserKey := util.Video_User_key + strconv.FormatInt(id, 10)
	if cnt, err := redis.RedisDb.SCard(redis.Ctx, redisUserKey).Result(); cnt > 0 {
		if err != nil {
			log.Println("redis query error!")
			return -1
		}
		redis.RedisDb.Expire(redis.Ctx, redisUserKey, redis.RandomTime())
		return cnt
	}
	cnt := int64(len(dao.GetVideoIdListByUserId(id)))
	redis.RedisDb.Set(redis.Ctx, redisUserKey, cnt, util.Relation_FollowingCnt_TTL)
	return cnt
}

func (VideoServiceImpl) ExpireWorks(id int64) {
	redisUserKey := util.Video_User_key + strconv.FormatInt(id, 10)
	redis.RedisDb.Del(redis.Ctx, redisUserKey)
}

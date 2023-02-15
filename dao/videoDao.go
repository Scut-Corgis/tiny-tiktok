package dao

import (
	"log"
	"time"
)

type Video struct {
	Id          int64
	AuthorId    int64
	PlayUrl     string
	CoverUrl    string
	PublishTime time.Time
	Title       string
}

// VideoDetail 和 controller/common.go的Video结构保持一致
type VideoDetail struct {
	Id            int64
	Author        UserResp
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
	Title         string
}

type Like struct {
	Id      int64 `gorm:"column:id"`
	UserId  int64 `gorm:"column:user_id"`
	VideoId int64 `gorm:"column:video_id"`
}

type Comment struct {
	Id          int64     `gorm:"column:id"`
	UserId      int64     `gorm:"column:user_id"`
	VideoId     int64     `gorm:"column:video_id"`
	CommentText string    `gorm:"column:comment_text"`
	CreateDate  time.Time `gorm:"column:create_date"`
}

// QueryVideoById 根据视频id查询视频
func QueryVideoById(id int64) (Video, error) {
	video := Video{}
	if err := Db.Where("id = ?", id).First(&video).Error; err != nil {
		log.Println(err.Error())
		return video, err
	}
	return video, nil
}

// InsertVideosTable 将video插入videos表内
func InsertVideosTable(video *Video) error {
	if err := Db.Create(&video).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// QueryVideoDetailByVideoId 根据视频id和查询用户id查询视频的详细信息
func QueryVideoDetailByVideoId(videoId int64, queryUserId int64) (VideoDetail, time.Time) {
	var err error
	var detailVideo VideoDetail
	videoShort, err := QueryVideoById(videoId)
	publishTIme := videoShort.PublishTime
	if err != nil {
		log.Fatalln("QueryVideoDetailByVideoId : 参数id可能有误")
	}
	userTable, err := QueryUserRespById(videoShort.AuthorId)
	if err != nil {
		log.Fatalln("QueryUserRespById : 视频作者id可能有误")
	}
	// 若queryUserId不为-1，则多一步查询是否关注了视频发布的作者
	if queryUserId != -1 && JudgeIsFollowById(queryUserId, userTable.Id) {
		userTable.IsFollow = true
	}

	detailVideo.Id = videoShort.Id
	detailVideo.Author = userTable
	detailVideo.PlayUrl = videoShort.PlayUrl
	detailVideo.CoverUrl = videoShort.CoverUrl
	detailVideo.Title = videoShort.Title

	Db.Model(&Like{}).Where("video_id = ?", detailVideo.Id).Count(&detailVideo.FavoriteCount)   // 统计点赞数量
	Db.Model(&Comment{}).Where("video_id = ?", detailVideo.Id).Count(&detailVideo.CommentCount) // 统计评论数量
	//查询是否点赞了视频
	if queryUserId != -1 && JudgeIsFavorite(queryUserId, detailVideo.Id) {
		detailVideo.IsFavorite = true
	}
	return detailVideo, publishTIme
}

func GetMost30videosIdList(latestTime time.Time) []int64 {
	var videoIdList = make([]int64, 0, 30)
	//Db.Raw("SELECT id FROM videos WHERE publish_time < ? ORDER BY publish_time desc LIMIT 30", latestTime).Scan(&videoIdList)
	if err := Db.Select("id").Where("publish_time < ?", latestTime).Order("publish_time desc").Limit(30).Find(&videoIdList).Error; err != nil {
		log.Println(err.Error())
		return videoIdList
	}
	return videoIdList
}

func GetVideoIdListByUserId(authorId int64) []int64 {
	var videoIdList = make([]int64, 0)
	//Db.Raw("SELECT id FROM videos WHERE author_id = ?", authorId).Scan(&videoIdList)
	if err := Db.Select("id").Where("author_id = ?", authorId).Find(&videoIdList).Error; err != nil {
		log.Println(err.Error())
		return videoIdList
	}
	return videoIdList
}

func JudgeIsFavorite(userId int64, videoId int64) bool { // 判断userid是否点赞了VideoId
	var count int64
	Db.Model(&Like{}).Where("user_id = ? and video_id = ?", userId, videoId).Count(&count)
	return count > 0
}

//// JudgeVideoIsExist 判断videoId的视频是否存在
//func JudgeVideoIsExist(videoId int64) bool {
//	// 布谷鸟过滤器
//	return redis.CuckooFilterVideoId.Contain([]byte(strconv.FormatInt(videoId, 10)))
//}

// QueryAllVideoIds 查询所有的用户名
func QueryAllVideoIds() []int64 {
	videos := make([]int64, 0)
	if err := Db.Table("videos").Select("id").Find(&videos).Error; err != nil {
		log.Println(err.Error())
		return videos
	}
	return videos
}

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
func InsertVideosTable(video Video) (Video, error) {
	if err := Db.Create(&video).Error; err != nil {
		log.Println(err.Error())
		return Video{}, err
	}
	return video, nil
}

func GetMost30videosIdList(latestTime time.Time) []int64 {
	var videoIdList = make([]int64, 0, 30)
	if err := Db.Table("videos").Select("id").Where("publish_time < ?", latestTime).Order("publish_time desc").Limit(30).Find(&videoIdList).Error; err != nil {
		log.Println(err.Error())
		return videoIdList
	}
	return videoIdList
}

func GetVideoIdListByUserId(authorId int64) []int64 {
	var videoIdList = make([]int64, 0)
	if err := Db.Table("videos").Select("id").Where("author_id = ?", authorId).Find(&videoIdList).Error; err != nil {
		log.Println(err.Error())
		return videoIdList
	}
	return videoIdList
}

// QueryAllVideoIds 查询所有的用户名
func QueryAllVideoIds() []int64 {
	videos := make([]int64, 0)
	if err := Db.Table("videos").Select("id").Find(&videos).Error; err != nil {
		log.Println(err.Error())
		return videos
	}
	return videos
}

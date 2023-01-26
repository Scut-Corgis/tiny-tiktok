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

// 和 controller/common.go的Video结构保持一致
type VideoDetail struct {
	Id            int64
	Author        UserTable
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

// QueryVideoById 根据视频id查询视频， 请确保id存在！
func QueryVideoById(id int64) (Video, error) {
	video := Video{}
	if err := Db.Where("id = ?", id).First(&video).Error; err != nil {
		log.Println(err.Error())
		if video.Id == 0 {
			log.Fatalln("查询了不存在的视频id")
		}
		return video, err
	}
	return video, nil
}

func InsertVideosTable(video *Video) error {
	if err := Db.Create(&video).Error; err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// 参数为VideoId， 以及调用该函数的查询作者的ID，若无查询作者，请定为-1。
// 返回videoDetail,视频发布时间;
// 该函数不会返回err，因为参数确保是有效的，若无效会直接os.exit()
// 多了一个发布时间，主要是方便处理feed流回复返回的next_time,不需要可以丢弃
func QueryVideoDetailByVideoId(videoId int64, queryUserId int64) (VideoDetail, time.Time) {
	var err error
	var detailVideo VideoDetail
	videoshort, err := QueryVideoById(videoId)
	publishTIme := videoshort.PublishTime
	if err != nil {
		log.Fatalln("QueryVideoDetailByVideoId : 参数id可能有误")
	}
	userTable, err := QueryUserTableById(videoshort.AuthorId)
	if err != nil {
		log.Fatalln("QueryUserTableById : 视频作者id可能有误")
	}
	// 若queryUserId不为-1，则多一步查询是否关注了视频发布的作者
	if queryUserId != -1 && JudgeIsFollowById(queryUserId, userTable.Id) {
		userTable.IsFollow = true
	}

	detailVideo.Id = videoshort.Id
	detailVideo.Author = userTable
	detailVideo.PlayUrl = videoshort.PlayUrl
	detailVideo.CoverUrl = videoshort.CoverUrl
	detailVideo.Title = videoshort.Title

	Db.Model(&Like{}).Where("video_id = ?", detailVideo.Id).Count(&detailVideo.FavoriteCount)   // 统计点赞数量
	Db.Model(&Comment{}).Where("video_id = ?", detailVideo.Id).Count(&detailVideo.CommentCount) // 统计评论数量
	//查询是否点赞了视频
	if queryUserId != -1 && JudgeIsFavorite(queryUserId, detailVideo.Id) {
		detailVideo.IsFavorite = true
	}
	return detailVideo, publishTIme
}

func GetMost30videosIdList(latestTime time.Time) []int64 {
	var videoIdList []int64 = make([]int64, 0, 30)
	Db.Raw("SELECT id FROM videos WHERE publish_time < ? ORDER BY publish_time desc LIMIT 30", latestTime).Scan(&videoIdList)
	return videoIdList
}

func GetVideoIdListByUserId(authorId, queryUserId int64) []int64 {
	var videoIdList []int64 = make([]int64, 0)
	Db.Raw("SELECT id FROM videos WHERE author_id = ?", authorId).Scan(&videoIdList)
	return videoIdList
}

func JudgeIsFavorite(userid int64, videoId int64) bool { // 判断userid是否点赞了VideoId
	var count int64
	Db.Model(&Like{}).Where("user_id = ? and video_id = ?", userid, videoId).Count(&count)
	return count > 0
}

// 判断videoId的视频是否存在
func JudgeVideoIsExist(videoId int64) bool {
	var count int64
	Db.Model(&Video{}).Where(map[string]interface{}{"id": videoId}).Count(&count)
	return count > 0
}

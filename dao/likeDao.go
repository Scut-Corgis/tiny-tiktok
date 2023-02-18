package dao

import (
	"errors"
	"log"
)

// TableName 获取点赞表名
func (Like) TableName() string {
	return "likes"
}

// InsertLike 插入点赞数据
func InsertLike(likeData *Like) error {
	err := Db.Model(&Like{}).Create(&likeData).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert like data failed")
	}
	return nil
}

// DeleteLike 删除点赞数据
func DeleteLike(userId int64, videoId int64) error {
	err := Db.Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Delete(&Like{}).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("delete like data failed")
	}
	return nil
}

// GetLikeVideoIdList 根据userId查询其点赞全部videoId ==>> 得到每个userId用户的所有点赞视频videoId
func GetLikeVideoIdList(userId int64) ([]int64, error) {
	var likeVideoIdList []int64
	err := Db.Model(&Like{}).Where(map[string]interface{}{"user_id": userId}).Pluck("video_id", &likeVideoIdList).Error
	if err != nil {
		//查询数据为0，返回空likeVideoIdList切片，以及返回无错误
		if "record not found" == err.Error() {
			log.Println("there are no likeVideoIds")
			return likeVideoIdList, nil
		} else {
			//如果查询数据库失败，返回获取likeVideoIdList失败
			log.Println(err.Error())
			return likeVideoIdList, errors.New("get likeVideoIdList failed")
		}
	}
	return likeVideoIdList, nil
}

// GetLikeUserIdList 根据videoId查询点赞该视频的全部user_id
func GetLikeUserIdList(videoId int64) ([]int64, error) {
	var likeUserIdList []int64 //存所有该视频点赞用户id；
	//查询likes表对应视频id点赞用户，返回查询结果
	err := Db.Model(Like{}).Where(map[string]interface{}{"video_id": videoId}).Pluck("user_id", &likeUserIdList).Error
	//查询过程出现错误，返回默认值0，并输出错误信息
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("get likeUserIdList failed")
	} else {
		//没查询到或者查询到结果，返回数量以及无报错
		return likeUserIdList, nil
	}
}

// GetLikeCountByVideoId 根据video_id获取该视频的点赞数
func GetLikeCountByVideoId(videoId int64) (int64, error) {
	var result int64
	err := Db.Model(Like{}).Where(map[string]interface{}{"video_id": videoId}).Count(&result).Error
	if err != nil {
		log.Println(err.Error())
		return result, errors.New("get likeCount failed")
	}
	return result, nil
}

//// GetLikInfo 获取点赞信息
//func GetLikInfo(userId int64, videoId int64) (Like, error) {
//	var likeInfo Like
//	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).First(&likeInfo).Error
//	if err != nil {
//		if "record not found" == err.Error() {
//			log.Println("can't find data")
//			return Like{}, nil
//		} else {
//			log.Println(err.Error())
//			return likeInfo, errors.New("get likeInfo failed")
//		}
//	}
//	return likeInfo, nil
//}

package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"log"
)

type CommentServiceImpl struct {
	UserService
}

func (CommentServiceImpl) QueryCommentsByVideoId(id int64) []dao.Comment {
	commentsList, err := dao.QueryCommentsByVideoId(id)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("Video not found!")
		return commentsList
	}
	log.Println("Query comments successfully!")
	return commentsList
}

func (CommentServiceImpl) InsertComment(comment *dao.Comment) bool {
	flag := dao.InsertComment(comment)
	if flag == false {
		log.Println("Insert comment failed!")
		return flag
	}
	log.Println("Insert comment successfully!")
	return flag
}

func (CommentServiceImpl) DeleteComment(id int64) bool {
	flag := dao.DeleteComment(id)
	if flag == false {
		log.Println("Delete comment failed!")
		return flag
	}
	log.Println("Delete comment successfully!")
	return flag
}

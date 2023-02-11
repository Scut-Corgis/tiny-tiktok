package service

import "github.com/Scut-Corgis/tiny-tiktok/dao"

type CommentService interface {
	// QueryCommentsByVideoId 获取评论列表
	QueryCommentsByVideoId(id int64) []dao.Comment

	// PostComment 发布评论
	PostComment(comment *dao.Comment) (int32, string)

	// InsertComment 插入评论
	InsertComment(comment *dao.Comment) bool

	// DeleteComment 删除评论
	DeleteComment(id int64) bool
}

package service

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
)

type CommentService interface {
	// CommentCount 根据视频id统计评论数量
	CommentCount(id int64) (int64, error)

	// QueryCommentsByVideoId 获取评论列表
	QueryCommentsByVideoId(id int64) []dao.Comment

	// PostComment 发布评论
	PostComment(comment dao.Comment) (int64, int32, string)

	// DeleteComment 删除评论
	DeleteComment(id int64) bool
}

package util

import (
	"time"
)

const Day = time.Hour * 24
const Month = Day * 30

// user模块
const User_Id_Key = "user:id:" // key:user_id value:username relation 1:1

// relation 模块
const Relation_Follow_Key = "relation:follow:"
const Relation_Follow_TTL = Day

const Relation_FollowerList_Key = "relation:followerList:"
const Relation_FollowerList_TTL = Day

// comment 模块
const Relation_Comment_Key = "relation:comment"
const Relation_Comment_TTL = Day

// video 模块
const Relation_Video_Key = "relation:video"

// 点赞模块
const Like_User_Key = "like:user"
const Like_Video_key = "like:video"

// Comment 模块
const Comment_Comment_Key = "comment:comment:" // key:comment_id value:video_id relation 1:1
const Comment_Video_Key = "comment:video:"     // key:video_id value:comment_id ralation 1:n

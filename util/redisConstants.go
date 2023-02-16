package util

import (
	"math/rand"
	"time"
)

const Day = time.Hour * 24
const Month = Day * 30

var HourRandnum int64 = rand.Int63n(24)
var DayRandnum int64 = rand.Int63n(30)

// user模块
const User_Id_Key = "user:id:" // key:user_id value:username relation 1:1

// relation 模块
const Relation_Follow_Key = "relation:follow:"

var Relation_Follow_TTL = Day + Day*time.Duration(DayRandnum)

const Relation_FollowerCnt_Key = "relation:followersCount:"

var Relation_FollowerCnt_TTL = Day + Day*time.Duration(DayRandnum)

const Relation_FollowingCnt_Key = "relation:followingsCount:"

var Relation_FollowingCnt_TTL = Day + Day*time.Duration(DayRandnum)

// 点赞模块
const Like_User_Key = "like:user"
const Like_Video_key = "like:video"

// Comment 模块
const Comment_Comment_Key = "comment:comment:" // key:comment_id value:video_id relation 1:1
const Comment_Video_Key = "comment:video:"     // key:video_id value:comment_id ralation 1:n

// message 模块
const Message_LatestMsg_Key = "message:latestMessage:"

var Message_LatestMsg_TTL = Month + time.Hour*time.Duration(HourRandnum)

const Message_MsgList_Key = "message:messageList:"

var Message_MsgList_TTL = Month + time.Hour*time.Duration(HourRandnum)

package util

import (
	"fmt"
	"math/rand"
	"time"
)

const Day = time.Hour * 24
const Month = Day * 30

var HourRandnum int64 = rand.Int63n(24)
var DayRandnum int64 = rand.Int63n(30)

// user 模块
const User_Id_Key = "user:id:" // key:user_id value:username relation 1:1
// video 模块
const Video_User_key = "video:user:"

// 点赞模块
const Like_User_Key = "like:user:"
const Like_Video_key = "like:video:"
const MyDefault = -1

// Comment 模块
const Comment_Comment_Key = "comment:comment:" // key:comment_id value:video_id relation 1:1
const Comment_Video_Key = "comment:video:"     // key:video_id value:comment_id ralation 1:n

// relation 模块
const Relation_Follow_Key = "relation:follow:"

var Relation_Follow_TTL = Day + Day*time.Duration(DayRandnum)

const Relation_FollowerCnt_Key = "relation:followersCount:"

var Relation_FollowerCnt_TTL = Day + Day*time.Duration(DayRandnum)

const Relation_FollowingCnt_Key = "relation:followingsCount:"

var Relation_FollowingCnt_TTL = Day + Day*time.Duration(DayRandnum)

// message 模块
const Message_LatestMsg_Key = "message:latestMessage:"

var Message_LatestMsg_TTL = Month + time.Hour*time.Duration(HourRandnum)

const Message_MsgList_Key = "message:messageList:"

var Message_MsgList_TTL = Month + time.Hour*time.Duration(HourRandnum)

const Message_MessageId_Key = "message:messageId:"

var Message_MessageId_TTL = Day + time.Hour*time.Duration(HourRandnum)

/*
生成消息key
参数：sendMsgUserId int64 发送消息用户id，recvMsgUserId int64 接收消息用户id
返回：msgKey string 消息key
*/
func GenMsgKey(sendMsgUserId int64, recvMsgUserId int64) string {
	return fmt.Sprintf("%d_%d", sendMsgUserId, recvMsgUserId)
}

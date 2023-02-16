package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
)

type MessageServiceImpl struct{}

type RedisMessage struct {
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
}

/*
发送消息
*/
func (MessageServiceImpl) SendMessage(userId int64, toUserId int64, content string) (bool, error) {
	msgKey := genMsgKey(userId, toUserId)
	createTime := time.Now().Format("2006-01-02 15:04:05")

	//redis缓存

	redisMsg := RedisMessage{
		Content:    content,
		CreateTime: createTime,
	}
	// 消息以hash类型保存
	data, err := json.Marshal(redisMsg)
	if err != nil {
		log.Println(err)
	}
	// 最新消息
	redisLatestMsgFrom := util.Message_LatestMsg_Key + strconv.FormatInt(userId, 10)
	if redis.RedisDb.Set(redis.Ctx, redisLatestMsgFrom, data, util.Message_LatestMsg_TTL).Err() != nil {
		log.Println(err)
	}

	redisLatestMsgTo := util.Message_LatestMsg_Key + strconv.FormatInt(toUserId, 10)
	if redis.RedisDb.Set(redis.Ctx, redisLatestMsgTo, data, util.Message_LatestMsg_TTL).Err() != nil {
		log.Println(err)
	}
	// 全部消息
	redisMsgListKey := util.Message_MsgList_Key + msgKey
	if redis.RedisDb.SAdd(redis.Ctx, redisMsgListKey, data).Err() != nil {
		log.Println(err)
	}
	redis.RedisDb.Expire(redis.Ctx, redisMsgListKey, util.Message_MsgList_TTL)
	//
	// redis.RedisDb.Do("hmset", redis.RedisDb.Args{redisMsgKey}.AddFlat(redisMsg)...)
	// //获取缓存
	// value, _ := redis.RedisDb.Values(redis.RedisDb.Do("hgetall", redisMsgKey))
	// //将values转成结构体
	// object := &RedisMessage{}
	// redis.RedisDb.ScanStruct(value, object)

	// redis.RedisDb.SAdd(redis.Ctx, redisMsgKey, redisMsg)

	//#优化 mysql的增加放到消息队列里
	return dao.InsertMessage(msgKey, content, createTime)
}

/*
读取聊天记录
*/
func (MessageServiceImpl) GetChatRecord(userId int64, toUserId int64) ([]dao.MessageResp, error) {
	msgKey := genMsgKey(userId, toUserId)
	msgList := make([]dao.MessageResp, 0)

	messages, err := dao.QueryMessagesByMsgKey(msgKey)
	if err != nil {
		return msgList, err
	}
	for _, message := range messages {
		msgResp := dao.MessageResp{}
		msgResp.Id = message.Id
		msgResp.ToUserId = toUserId
		msgResp.FromUserId = userId
		msgResp.Content = message.Content
		msgResp.CreateTime = message.CreateTime
		msgList = append(msgList, msgResp)
	}
	return msgList, nil
}

/*
生成消息key
参数：sendMsgUserId int64 发送消息用户id，recvMsgUserId int64 接收消息用户id
返回：msgKey string 消息key
*/
func genMsgKey(sendMsgUserId int64, recvMsgUserId int64) string {
	return fmt.Sprintf("%d_%d", sendMsgUserId, recvMsgUserId)
}

// /*
// 解析消息key，得到to_user_id, from_user_id
// 参数：msgKey string 消息key
// 返回：from_user_id int64 发送消息用户id, to_user_id int64 接收消息用户id
// */
// func parseMsgKey(msgKey string) (int64, int64) {
// 	ids := strings.Split(msgKey, "_")
// 	from_user_id, err := strconv.ParseInt(ids[0], 10, 64)
// 	if err != nil {
// 		return -1, -1
// 	}
// 	to_user_id, err := strconv.ParseInt(ids[1], 10, 64)
// 	if err != nil {
// 		return -1, -1
// 	}
// 	return from_user_id, to_user_id
// }

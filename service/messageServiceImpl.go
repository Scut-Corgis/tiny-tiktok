package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
)

type MessageServiceImpl struct{}

type LatestMessage struct {
	Id         int64  `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"createTime"`
	MsgType    int64  `json:"msgType"` // 1为发送方，0为接收方
}

/*
发送消息
*/
func (MessageServiceImpl) SendMessage(userId int64, toUserId int64, content string) (bool, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")

	// redis缓存 最新消息
	redisLatestMessage := LatestMessage{}
	redisLatestMessage.Content = content
	redisLatestMessage.CreateTime = createTime
	redisLatestMessage.MsgType = 1
	dataFrom, err := json.Marshal(redisLatestMessage)
	if err != nil {
		log.Println(err)
	}
	msgKey := genMsgKey(userId, toUserId)
	redisLatestMsgKey := util.Message_LatestMsg_Key + msgKey
	if redis.RedisDb.Set(redis.Ctx, redisLatestMsgKey, dataFrom, util.Message_LatestMsg_TTL).Err() != nil {
		log.Println(err)
	}
	// #优化：启用消息队列后，聊天记录只能访问到最新的记录，无法返回所有聊天记录
	// sb := strings.Builder{}
	// sb.WriteString(strconv.FormatInt(userId, 10))
	// sb.WriteString("#%#")
	// sb.WriteString(strconv.FormatInt(toUserId, 10))
	// sb.WriteString("#%#")
	// sb.WriteString(content)
	// sb.WriteString("#%#")
	// sb.WriteString(createTime)
	// rabbitmq.RabbitMQMessageAdd.Producer(sb.String())

	msgId, err := dao.InsertMessage(userId, toUserId, content, createTime)
	if err != nil || msgId < 0 {
		log.Println(err)
		return false, err
	}
	redisMessageIdKey := util.Message_MessageId_Key + genMsgKey(userId, toUserId)
	redis.RedisDb.SAdd(redis.Ctx, redisMessageIdKey, msgId)
	redis.RedisDb.Expire(redis.Ctx, redisMessageIdKey, util.Message_MessageId_TTL)

	return true, nil
}

/*
读取聊天记录
*/
func (MessageServiceImpl) GetChatRecord(userId int64, toUserId int64) ([]dao.MessageResp, error) {
	msgList := make([]dao.MessageResp, 0)
	messages, err := dao.QueryMessagesByMsgKey(userId, toUserId)
	if err != nil {
		return msgList, err
	}
	// redis 去重
	redisMessageIdKey := util.Message_MessageId_Key + genMsgKey(userId, toUserId)
	for _, message := range messages {
		if flag, err := redis.RedisDb.SIsMember(redis.Ctx, redisMessageIdKey, message.Id).Result(); flag {
			if err != nil {
				log.Println(err)
			}
			continue
		}
		msgResp := dao.MessageResp{}
		msgResp.Id = message.Id
		msgResp.ToUserId = message.ToUserId
		msgResp.FromUserId = message.FromUserId
		msgResp.Content = message.Content
		msgResp.CreateTime = timeStringToUnix(message.CreateTime) //接口文档有误
		redis.RedisDb.SAdd(redis.Ctx, redisMessageIdKey, message.Id)
		redis.RedisDb.Expire(redis.Ctx, redisMessageIdKey, util.Message_MessageId_TTL)
		msgList = append(msgList, msgResp)
	}
	// 再去重
	msi := MessageServiceImpl{}
	latestMsg, err := msi.GetLatestMessage(userId, toUserId)
	if err == nil && len(msgList) > 0 && latestMsg.Id == msgList[len(msgList)-1].Id {
		return make([]dao.MessageResp, 0), nil
	}
	return msgList, nil
}

/*
获取最新聊天记录
*/
func (MessageServiceImpl) GetLatestMessage(userId int64, toUserId int64) (LatestMessage, error) {

	// 当前用户为发送方
	sendLatestMsg := LatestMessage{}
	redisLatestMsgKey := util.Message_LatestMsg_Key + genMsgKey(userId, toUserId)
	sendStringCmd, err1 := redis.RedisDb.Get(redis.Ctx, redisLatestMsgKey).Result()
	if err1 == nil {
		err := json.Unmarshal([]byte(sendStringCmd), &sendLatestMsg)
		if err != nil {
			log.Println(err)
		}
	}
	// 当前用户为接收方
	recvLatestMsg := LatestMessage{}
	redisLatestMsgKey = util.Message_LatestMsg_Key + genMsgKey(toUserId, userId)
	recvStringCmd, err2 := redis.RedisDb.Get(redis.Ctx, redisLatestMsgKey).Result()
	if err2 == nil {
		err := json.Unmarshal([]byte(recvStringCmd), &recvLatestMsg)
		if err == nil {
			//当前用户为接收方进行的查询，所以改成0
			recvLatestMsg.MsgType = 0
		}
	}
	if err1 == nil && err2 == nil {
		sendTime := timeStringToUnix(sendLatestMsg.CreateTime)
		recvTime := timeStringToUnix(recvLatestMsg.CreateTime)
		if sendTime < recvTime {
			return recvLatestMsg, nil
		} else {
			return sendLatestMsg, nil
		}
	}
	if err1 == nil {
		return sendLatestMsg, nil
	}
	if err2 == nil {
		return recvLatestMsg, nil
	}
	// 缓存没数据，则去数据库查
	latestMsg := LatestMessage{}
	message, err := dao.QueryLatestMessageByUserId(userId, toUserId)
	if err != nil {
		log.Println(err)
		return latestMsg, err
	}
	latestMsg.Id = message.Id
	latestMsg.Content = message.Content
	latestMsg.CreateTime = message.CreateTime
	if message.FromUserId == userId {
		latestMsg.MsgType = 1
	} else if message.ToUserId == userId {
		latestMsg.MsgType = 0
	}
	return latestMsg, nil
}

/*
时间字符串转时间戳
*/
func timeStringToUnix(timeString string) int64 {
	loc, _ := time.LoadLocation("Local")
	timeDate, _ := time.ParseInLocation("2006-01-02 15:04:05", timeString, loc)
	return timeDate.Unix()
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

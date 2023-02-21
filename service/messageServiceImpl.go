package service

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/rabbitmq"
	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
)

type MessageServiceImpl struct{}

/*
发送消息
*/
func (MessageServiceImpl) SendMessage(userId int64, toUserId int64, content string) (bool, error) {
	createTime := time.Now()
	createTimeStr := createTime.Format("2006-01-02 15:04:05")
	// 启用消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.FormatInt(userId, 10))
	sb.WriteString("#%#")
	sb.WriteString(strconv.FormatInt(toUserId, 10))
	sb.WriteString("#%#")
	sb.WriteString(content)
	sb.WriteString("#%#")
	sb.WriteString(createTimeStr)
	rabbitmq.RabbitMQMessageAdd.Producer(sb.String())
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
	redisMessageIdKey := util.Message_MessageId_Key + util.GenMsgKey(userId, toUserId)
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
		msgResp.CreateTime = util.TimeToUnix(message.CreateTime) //接口文档有误,返回类型为时间戳
		redis.RedisDb.SAdd(redis.Ctx, redisMessageIdKey, message.Id)
		redis.RedisDb.Expire(redis.Ctx, redisMessageIdKey, util.Message_MessageId_TTL)
		msgList = append(msgList, msgResp)
	}
	// 再根据最新消息时间进行去重
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
func (MessageServiceImpl) GetLatestMessage(userId int64, toUserId int64) (dao.LatestMessage, error) {
	// 若当前用户为发送方
	sendLatestMsg := dao.LatestMessage{}
	redisLatestMsgKey := util.Message_LatestMsg_Key + util.GenMsgKey(userId, toUserId)
	sendStringCmd, err1 := redis.RedisDb.Get(redis.Ctx, redisLatestMsgKey).Result()
	if err1 == nil {
		err := json.Unmarshal([]byte(sendStringCmd), &sendLatestMsg)
		if err != nil {
			log.Println(err)
		}
	}
	// 若当前用户为接收方
	recvLatestMsg := dao.LatestMessage{}
	redisLatestMsgKey = util.Message_LatestMsg_Key + util.GenMsgKey(toUserId, userId)
	recvStringCmd, err2 := redis.RedisDb.Get(redis.Ctx, redisLatestMsgKey).Result()
	if err2 == nil {
		err := json.Unmarshal([]byte(recvStringCmd), &recvLatestMsg)
		if err == nil {
			//当前用户为接收方进行的查询，所以改成0
			recvLatestMsg.MsgType = 0
		}
	}
	if err1 == nil && err2 == nil {
		sendTime := util.TimeStrToUnix(sendLatestMsg.CreateTime)
		recvTime := util.TimeStrToUnix(recvLatestMsg.CreateTime)
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
	latestMsg := dao.LatestMessage{}
	message, err := dao.QueryLatestMessageByUserId(userId, toUserId)
	if err != nil {
		log.Println(err)
		return latestMsg, err
	}
	latestMsg.Id = message.Id
	latestMsg.Content = message.Content
	latestMsg.CreateTime = util.TimeToTimeStr(message.CreateTime)
	if message.FromUserId == userId {
		latestMsg.MsgType = 1
	} else if message.ToUserId == userId {
		latestMsg.MsgType = 0
	}
	return latestMsg, nil
}

package rabbitmq

import (
	"log"
	"strconv"
	"strings"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"
	"github.com/streadway/amqp"
)

type MessageMQ struct {
	RabbitMQ
	Channel    *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
}

var RabbitMQMessageAdd *MessageMQ

// NewMessageRabbitMQ 获取messageMQ的对应队列。
func NewMessageRabbitMQ(queueName string) *MessageMQ {
	messageMQ := &MessageMQ{
		RabbitMQ:  *MyRabbitMQ,
		QueueName: queueName,
	}
	var err error
	messageMQ.Channel, err = messageMQ.Conn.Channel()
	MyRabbitMQ.failOnErr(err, "Failed to get channel!")
	return messageMQ
}
func InitMessageRabbitMQ() {
	RabbitMQMessageAdd = NewMessageRabbitMQ("Message Add")
	go RabbitMQMessageAdd.Consumer()
	log.Println("RabbitMQMessageAdd init successfully!")
}

// Producer 生产
func (c *MessageMQ) Producer(message string) {
	_, err := c.Channel.QueueDeclare(
		c.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		log.Println(err.Error())
	}

	c.Channel.Publish(
		c.Exchange,
		c.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

// Consumer 消费
func (c *MessageMQ) Consumer() {

	_, err := c.Channel.QueueDeclare(
		c.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err.Error())
	}

	//2、接收消息
	msgs, err := c.Channel.Consume(
		c.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go c.consumerMessageAdd(msgs)
	<-forever

}

// 添加聊天记录 消息队列
func (c *MessageMQ) consumerMessageAdd(messages <-chan amqp.Delivery) {
	for message := range messages {
		// 参数解析
		params := strings.Split(string(message.Body), "#%#")
		msgId, _ := strconv.ParseInt(params[2], 10, 64)

		// 更新聊天记录缓存
		redisMessageIdKey := util.Message_MessageId_Key + params[0] + "_" + params[1]
		redis.RedisDb.SAdd(redis.Ctx, redisMessageIdKey, msgId)
		redis.RedisDb.Expire(redis.Ctx, redisMessageIdKey, util.Message_MessageId_TTL)

	}
}

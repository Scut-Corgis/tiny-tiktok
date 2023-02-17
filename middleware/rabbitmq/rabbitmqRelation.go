package rabbitmq

import (
	"log"
	"strconv"
	"strings"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/streadway/amqp"
)

type RelationMQ struct {
	RabbitMQ
	Channel    *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
}

var RabbitMQRelationAdd *RelationMQ
var RabbitMQRelationDel *RelationMQ

// NewFollowRabbitMQ 获取relationMQ的对应队列。
func NewRelationRabbitMQ(queueName string) *RelationMQ {
	relationMQ := &RelationMQ{
		RabbitMQ:  *MyRabbitMQ,
		QueueName: queueName,
	}
	var err error
	relationMQ.Channel, err = relationMQ.Conn.Channel()
	MyRabbitMQ.failOnErr(err, "获取通道失败")
	return relationMQ
}
func InitRelationRabbitMQ() {

	RabbitMQRelationAdd = NewRelationRabbitMQ("Relation Add")
	go RabbitMQRelationAdd.Consumer()
	log.Println("RabbitMQRelationAdd init successfully!")

	RabbitMQRelationDel = NewRelationRabbitMQ("Relation Del")
	go RabbitMQRelationDel.Consumer()
	log.Println("RabbitMQRelationDel init successfully!")
}

// Producer 生产
func (c *RelationMQ) Producer(message string) {
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
func (c *RelationMQ) Consumer() {

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
	switch c.QueueName {
	case "Relation Add":
		go c.consumerFollowAdd(msgs)
	case "Relation Del":
		go c.consumerFollowDel(msgs)

	}
	<-forever

}

// 关系添加的消费方式。
func (c *RelationMQ) consumerFollowAdd(messages <-chan amqp.Delivery) {
	for message := range messages {
		// 参数解析
		params := strings.Split(string(message.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		followId, _ := strconv.ParseInt(params[1], 10, 64)

		if err := dao.InsertFollow(userId, followId); nil != err {
			log.Println(err.Error())
		}
	}
}

// 关系删除的消费方式。
func (c *RelationMQ) consumerFollowDel(messages <-chan amqp.Delivery) {
	for message := range messages {
		// 参数解析。
		params := strings.Split(string(message.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		followId, _ := strconv.ParseInt(params[1], 10, 64)

		if err := dao.DeleteFollow(userId, followId); nil != err {
			log.Println(err.Error())
		}
		// 再删Redis里的信息，防止脏数据，保证最终一致性。
		//updateRedisWithDel(userId, targetId)
	}
}

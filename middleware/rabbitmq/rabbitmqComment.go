package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type CommentMQ struct {
	RabbitMQ
	Channel    *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
}

var RabbitMQCommentDel *CommentMQ

// NewCommentRabbitMQ 获取Comment的消息队列
func NewCommentRabbitMQ(queueName string) *CommentMQ {
	commentMQ := &CommentMQ{
		RabbitMQ:  *MyRabbitMQ,
		QueueName: queueName,
	}
	var err error
	commentMQ.Channel, err = commentMQ.Conn.Channel()
	MyRabbitMQ.failOnErr(err, "Failed to get channel!")
	return commentMQ
}

func InitCommentRabbitMQ() {
	RabbitMQCommentDel = NewCommentRabbitMQ("Comment Del")
	go RabbitMQCommentDel.Consumer()
	log.Println("RabbitMQCommentDel init successfully！")
}

// Producer 生产
func (c *CommentMQ) Producer(message string) {
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
	err = c.Channel.Publish(
		c.Exchange,
		c.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Println(err.Error())
	}
}

// Consumer 消费
func (c *CommentMQ) Consumer() {
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
	messages, err := c.Channel.Consume(
		c.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	forever := make(chan bool)
	go consumerCommentDel(messages)
	<-forever
}

func consumerCommentDel(messages <-chan amqp.Delivery) {
	for message := range messages {
		log.Println(message.Body)
	}
}

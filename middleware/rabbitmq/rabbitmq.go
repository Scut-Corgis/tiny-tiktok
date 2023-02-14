package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const MQURL = "amqp://tiktok:123456@127.0.0.1:5672/"

type RabbitMQ struct {
	Conn *amqp.Connection
	//MQ链接字符串
	Mqurl string
}

var MyRabbitMQ *RabbitMQ

// 创建结构体实例
func Init() {
	MyRabbitMQ = &RabbitMQ{
		Mqurl: MQURL,
	}
	var err error
	MyRabbitMQ.Conn, err = amqp.Dial(MyRabbitMQ.Mqurl)
	MyRabbitMQ.failOnErr(err, "创建连接失败")
	log.Println("rabbitmq TCP连接成功")
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// 断开channel和connection
func (r *RabbitMQ) Destory() {
	r.Conn.Close()
}

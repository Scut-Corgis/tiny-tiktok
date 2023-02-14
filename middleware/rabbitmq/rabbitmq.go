package rabbitmq

import (
	"fmt"
	"github.com/Scut-Corgis/tiny-tiktok/config"
	"log"

	"github.com/streadway/amqp"
)

var MQURL = "amqp://" + config.RabbitMQ_username + ":" + config.RabbitMQ_passsword + "@" + config.RabbitMQ_IP + ":" + config.RabbitMQ_host + "/"

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

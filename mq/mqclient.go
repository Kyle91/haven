// @Author Eric
// @Date 2024/6/10 0:02:00
// @Desc
package mq

import (
	"github.com/Kyle91/haven/log"
	"github.com/Kyle91/haven/routine"
	"github.com/streadway/amqp"
)

// MQClient 是Client接口的一个简单实现
type MQClient struct {
	url  string
	conn *amqp.Connection
}

// NewMQClient
//
//	@Description: 初始化一个mq client
//	@param url
//	@return *MQClient
//	@return error
func NewMQClient(url string) (*MQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MQClient{url: url, conn: conn}, nil
}

// ensureQueueAndExchange
//
//	@Description: 确保队列和交换器存在
//	@receiver c
//	@param exchange
//	@param queueName
//	@param routingKey
//	@return error
func (c *MQClient) ensureQueueAndExchange(exchange, queueName, routingKey string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 如果交换器不存在，那就创建一个
	err = ch.ExchangeDeclare(
		exchange, // 交换器名称
		"direct", // 交换器类型
		true,     // 持久性，true保证RabbitMQ重启后交换器依然存在
		false,    //自动删除，交换器会在所有绑定的队列都不再绑定时自动删除
		false,    // 内部。如果设置为true，这个交换器不能被生产者直接发送消息
		false,    // 不等待。如果设置为true，则不等待服务器确认交换器是否成功声明
		nil,      // 键值对的额外参数
	)
	if err != nil {
		return err
	}

	// 如果队列不存在，那就创建一个
	_, err = ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 是否持久化队列，true为持久化，false为临时的，持久化则意味着重启Rabbitmq服务器消息也不会丢失
		false,     // 当设置为true时，队列在没有任何消费者（Consumer）时会自动删除。这通常用于临时队列
		false,     // 当设置为true时，队列将被标记为独占队列。独占队列只能被声明它的连接访问，并且队列在连接关闭时会自动删除。
		false,     // 当设置为true时，服务器不会对QueueDeclare命令发送应答
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// 绑定队列到交换器
	err = ch.QueueBind(
		queueName,  // 队列名称
		routingKey, // 路由键
		exchange,   // 交换器名称
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

// Publish
//
//	@Description: 发布数据
//	@receiver c
//	@param exchange 交换器名称
//	@param routingKey 路由键名称
//	@param queueName 队列名称
//	@param message  消息内容
//	@return error
func (c *MQClient) Publish(exchange, routingKey, queueName string, message []byte) error {
	// 确保交换器、队列存在
	err := c.ensureQueueAndExchange(exchange, queueName, routingKey)
	if err != nil {
		return err
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		exchange,   // 交换器
		routingKey, // 路由键
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		return err
	}

	log.Infof("Published message: %s", message)
	return nil
}

// Subscribe
//
//	@Description: 订阅消息
//	@receiver c
//	@param exchange  交换器
//	@param queueName 队列名
//	@param routingKey 路由键
//	@param handler
//	@return error
func (c *MQClient) Subscribe(exchange, queueName, routingKey string, handler func(amqp.Delivery)) error {
	// 确保交换器和路由键存在
	err := c.ensureQueueAndExchange(exchange, queueName, routingKey)
	if err != nil {
		return err
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	routine.Go(func() {
		for d := range msgs {
			handler(d)
		}
	})

	log.Infof("Subscribed to topic: %s", routingKey)
	return nil
}

func (c *MQClient) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

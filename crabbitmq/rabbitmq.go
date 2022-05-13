package crabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type RabbitmqClient struct {
	*amqp.Connection
	Chan *amqp.Channel
	//下面需要直接初始化
	uri               string
	queueChans        map[string]<-chan amqp.Delivery //队列对应的消费者chan
	queueAutoAckChans map[string]<-chan amqp.Delivery //队列对应的消费者chan
	l                 *sync.RWMutex
}
type RabbitmqConf struct {
	User string
	Pwd  string
	Host string
	Port string
}

func NewRabbitmqClient(c *RabbitmqConf) (rc *RabbitmqClient, err error) {
	rc = &RabbitmqClient{
		uri:               "amqp://" + c.User + ":" + c.Pwd + "@" + c.Host + ":" + c.Port + "/",
		queueChans:        map[string]<-chan amqp.Delivery{},
		queueAutoAckChans: map[string]<-chan amqp.Delivery{},
		l:                 &sync.RWMutex{},
	}
	// 新建一个连接
	rc.Connection, err = amqp.Dial(rc.uri)
	if err != nil {
		return
	}
	rc.Chan, err = rc.Channel() //获取一个通道
	if err != nil {
		return
	}
	go rc.checkConn() //连通性检测
	return
}

//检查连通性，并自动重连
func (rc *RabbitmqClient) checkConn1() {
	for {
		time.Sleep(3 * time.Second)
		if rc.IsClosed() { //判断连接是否关闭，关闭则重连
			rc.reset() //重置
			fmt.Println("Rabbitmq Conn重连中。。。")
			conn, err := amqp.Dial(rc.uri) // 新建一个连接
			if err != nil {
				continue
			}
			rc.Connection = conn
		}

		//检查chan
		chanErr := make(chan *amqp.Error)
		rc.Chan.NotifyClose(chanErr)
		<-chanErr
		rc.reset() //重置
		c, err := rc.Channel()
		if err != nil {
			continue
		}
		rc.Chan = c
	}
}
func (rc *RabbitmqClient) checkConn() {
	for {
		time.Sleep(3 * time.Second)
		if rc.IsClosed() || rc.Chan.IsClosed() { //判断连接是否关闭，关闭则重连
			rc.reset() //重置
			fmt.Println("Rabbitmq Conn重连中。。。")
			conn, err := amqp.Dial(rc.uri) // 新建一个连接
			if err != nil {
				continue
			}
			rc.Connection = conn

			//重连chan
			c, err := rc.Connection.Channel()
			if err != nil {
				continue
			}
			rc.Chan = c
		}

	}
}
func (rc *RabbitmqClient) reset() {
	rc.l.Lock()
	defer rc.l.Unlock()
	if rc.Chan != nil {
		rc.Chan.Close()
	}
	rc.queueChans = map[string]<-chan amqp.Delivery{} //重置
	rc.queueAutoAckChans = map[string]<-chan amqp.Delivery{}
}
func (rc *RabbitmqClient) CreateQueue(queueName string) (err error) {
	rc.l.Lock()
	defer rc.l.Unlock()
	_, err = rc.Chan.QueueDeclare(
		queueName, // queueName
		true,      // 队列持久化
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return
	}
	return
}

func (rc *RabbitmqClient) Send(queueName, messageId string, body []byte) (err error) {
	rc.l.Lock()
	defer rc.l.Unlock()

	return rc.Chan.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			Body:         body,
			MessageId:    messageId,
			DeliveryMode: 2, //消息持久化
		})
}

type mqData struct {
	DeliveryTag uint64 //此次消费者读取此条数据的唯一标识，消费者重置后，此id会置零
	MessageId   string
	Body        []byte
}

func (rc *RabbitmqClient) Get(queueName string) (md *mqData, err error) {
	rc.l.Lock()
	defer rc.l.Unlock()
	queueChan, ok := rc.queueChans[queueName]
	if !ok {
		queueChan, err = rc.Chan.Consume(
			queueName, // queue
			queueName, // 此消费者名字，方便销毁用
			false,     // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			return
		}
		rc.queueChans[queueName] = queueChan
		time.Sleep(time.Second)
	}

	select {
	case m := <-queueChan: //查看是否有消息
		md = &mqData{
			DeliveryTag: m.DeliveryTag,
			MessageId:   m.MessageId,
			Body:        m.Body,
		}
	default:

	}
	return
}

func (rc *RabbitmqClient) GetAndAutoAck(queueName string) (md *mqData, err error) {
	rc.l.Lock()
	defer rc.l.Unlock()
	queueChan, ok := rc.queueAutoAckChans[queueName]
	if !ok {
		queueChan, err = rc.Chan.Consume(
			queueName, // queue
			queueName, // 此消费者名字，方便销毁用
			true,      // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)
		if err != nil {
			return
		}
		rc.queueAutoAckChans[queueName] = queueChan
		time.Sleep(time.Second)
	}

	select {
	case m := <-queueChan: //查看是否有消息
		md = &mqData{
			DeliveryTag: m.DeliveryTag,
			MessageId:   m.MessageId,
			Body:        m.Body,
		}
	default:

	}
	return
}

package crabbitmq_test

import (
	"github.com/streadway/amqp"
	"go_common/clogs"
	"go_common/crabbitmq"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	log      = clogs.NewLog("7", false)
	mqClient *crabbitmq.RabbitmqClient
)

func init() {
	var err error
	mqClient, err = crabbitmq.NewRabbitmqClient(&crabbitmq.RabbitmqConf{
		User: "rabbitmq",
		Pwd:  "rabbitmq",
		Host: "172.16.5.137",
		Port: "5672",
	})
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
}
func TestName(t *testing.T) {
agen:
	mqCh, err := mqClient.Channel()
	if err != nil {
		log.Error(err)
		time.Sleep(time.Second)
		goto agen
	}
	defer mqCh.Close()
	closeChan := make(chan *amqp.Error)
	notifyClose := mqCh.NotifyClose(closeChan)
	msgs, err := mqCh.Consume(
		"qq",    // queue
		"asdfg", // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Error(err)
		time.Sleep(time.Second)
		goto agen
	}
	closeFlag := false
	for {
		select {
		case e := <-notifyClose:
			log.Error(e.Error())
			closeFlag = true
		case d := <-msgs:
			//ch.Cancel("asdfg", true)
			log.Info("count", d.MessageCount, "id", d.MessageId, d.DeliveryTag, string(d.Body))
			err = mqCh.Ack(d.DeliveryTag, false)
			//err=ch.Reject(d.DeliveryTag, true)
			if err != nil {
				log.Error(err)
			}
		default:

		}
		if closeFlag {
			time.Sleep(time.Second)
			break
		}
	}
	goto agen
}
func TestSend(t *testing.T) {
	err := mqClient.Send("qq", "a", []byte("aaaa"))
	err = mqClient.Send("qq", "b", []byte("bbbb"))
	if err != nil {
		log.Error(err)
	}

}
func TestCreateQueue(t *testing.T) {
	log.Info(mqClient.CreateQueue("errUrls"))
}

func TestGetData(t *testing.T) {
agen:
	time.Sleep(1 * time.Second)
	md, err := mqClient.GetAndAutoAck("qq")
	if err != nil {
		log.Error(err)
		goto agen
	}
	if md == nil { //暂时没有消息
		log.Warn(md)
		goto agen
	}
	log.Info(md.DeliveryTag, md.MessageId, string(md.Body))
	//mqClient.Chan.Nack(md.DeliveryTag, false, true)
	goto agen
}

func TestHttp(t *testing.T) {
	http.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {
		md, err := mqClient.Get("qq")
		if err != nil {
			log.Error(err)
			w.Write([]byte("错误"))
			return
		}
		if md == nil { //暂时没有消息
			log.Warn(md)
			w.Write([]byte("没数据"))
			return
		}
		log.Info(md.DeliveryTag, md.MessageId, string(md.Body))
	})
	http.HandleFunc("/2", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.Form.Get("id")
		intNum, _ := strconv.Atoi(id)
		int64Num := uint64(intNum)
		log.Info(int64Num)
		mqClient.Chan.Nack(int64Num, false, true)
	})
	http.ListenAndServe(":9999", nil)
}

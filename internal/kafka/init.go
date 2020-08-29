package kafka

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"time"
	//cluster "github.com/bsm/sarama-cluster"
)

type Kafka struct {
	addrs []string
	c     sarama.Client
	sp    sarama.SyncProducer
	ap    sarama.AsyncProducer
	sc    *cluster.Client
	ss    sarama.Consumer
}

func NewKafka(adds []string) *Kafka {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	config.Producer.Return.Errors = true
	c, err := sarama.NewClient(adds, config)
	if err != nil {
		logx.Fatal(err.Error())
	}
	cconfig := cluster.NewConfig()
	cconfig.Consumer.Return.Errors = true
	cconfig.Group.Return.Notifications = true
	sc, err := cluster.NewClient(adds, cconfig)
	if err != nil {
		logx.Fatal(err.Error())
	}
	sp, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		logx.Fatalf("sarama.NewSyncProducer err, message=%s \n", err)
	}
	ap, err := sarama.NewAsyncProducerFromClient(c)
	if err != nil {
		logx.Fatalf("sarama.NewAsyncProducer err, message=%s \n", err)
	}
	ss, err := sarama.NewConsumerFromClient(c)
	if err != nil {
		logx.Fatalf("sarama.NewConsumer err, message=%s \n", err)
	}
	return &Kafka{addrs: adds, c: c, sc: sc, sp: sp, ap: ap, ss: ss}
}

func (k *Kafka) Close() {
	err := k.c.Close()
	if err != nil {
		logx.Error("kafka close:", err.Error())
	}
}

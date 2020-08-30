package cinit

import (
	"github.com/janiokq/Useless-blog/internal/kafka"
	"strings"
)

var Kf *kafka.Kafka

func KafkaInit() {
	addrs := strings.Split(Config.Kafka.Addr, ",")
	Kf = kafka.NewKafka(addrs)
}

func KafkaClose() {
	Kf.Close()
}

package user

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
)

func msgNotify(ctx context.Context, msgStr string) {
	topic := cinit.TopicUserSrvChange
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgStr),
	}
	part, offset, err := cinit.Kf.SyncProducer(msg)
	if err != nil {
		logx.Error("msg notify error:"+err.Error(), ctx)
	}
	logx.Infof("topic:%s,part:%d,offset:%d", topic, part, offset, ctx)
}

package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addr, "demo", cfg)
	assert.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	start := time.Now()
	err = consumer.Consume(ctx, []string{"test_topic"}, ConsumerHandler{})
	assert.NoError(t, err)
	t.Log(time.Since(start).String())
}

type ConsumerHandler struct {
}

func (c ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("这是Setup")
	//partitions := session.Claims()["test_topic"]
	//var offset int64 = 0
	//for _, part := range partitions {
	//	session.ResetOffset("test_topic",
	//		part, offset, "")
	//}
	return nil
}

func (c ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("这是Cleanup")
	return nil
}

/*
*
该处实现了消息按批次处理，一秒钟内一次处理10条（异步消费），不足10条直接处理。然后统一提交
*/
func (c ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10
	for {
		log.Println("一个批次开始")
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		var done = false
		var eg errgroup.Group
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 超时了
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				batch = append(batch, msg)
				eg.Go(func() error {
					// 并发处理（异步消费）
					log.Println("接收到消息：", string(msg.Value))
					return nil
				})
			}
		}
		cancel()
		err := eg.Wait()
		if err != nil {
			log.Println(err)
			continue
		}
		// 批量消费，凑够了一批，然后处理
		//log.Println(batch)

		// 批量提交
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}

	}
}

func (c ConsumerHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		log.Println("接收到消息：", string(msg.Value))
		session.MarkMessage(msg, "")
	}
	return nil
}

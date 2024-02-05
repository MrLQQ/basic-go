package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var addr = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	/**
	 * Sarama使用入门：指定分区
	 * 默认消息都发送到了Partition0
	 * 正常来说，在Sarama里面，可以通过指定config中的Partitioner来指定最终的目标分区。
	 * *Random：随机挑选一个。
	 * *RoundRobin：轮询
	 * *Hash：根据key的哈希值来筛选一个
	 * *Manual：根据Message中的partition字段来选择
	 * *ConsistentHash：一致性哈希，用的CRC16算法。
	 * *Custom：实际上不Custom，而是自定义一部分Hash参数，本质上是一个Hash的实现
	 */
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(addr, cfg)
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner // 轮询
	//cfg.Producer.Partitioner = sarama.NewRandomPartitioner            // 随机跳转
	//cfg.Producer.Partitioner = sarama.NewHashPartitioner              // 使用hash算法
	//cfg.Producer.Partitioner = sarama.NewManualPartitioner            // 使用指定的partition（SendMessage中指定）
	//cfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner // 使用一致性HASH算法
	//cfg.Producer.Partitioner = sarama.NewCustomPartitioner()          // 自定义HASH

	// 这个是为了兼容JAVA，不要使用
	//cfg.Producer.Partitioner = sarama.NewReferenceHashPartitioner

	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "test_topic",
			Value: sarama.StringEncoder("这是一条消息"),
			// 会在生产者和消费者之间传递的
			Headers: []sarama.RecordHeader{
				{
					Key:   []byte("key1"),
					Value: []byte("value1"),
				},
			},
			Metadata: "这是 metadata",
		})
	}
}

func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addr, cfg)
	assert.NoError(t, err)
	msgs := producer.Input()
	msgs <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("TestAsyncProducer这是一条消息"),
		// 会在生产者和消费者之间传递的
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
		Metadata: "这是 metadata",
	}

	select {
	case msg := <-producer.Successes():
		t.Log("发送成功", string(msg.Value.(sarama.StringEncoder)))
	case err := <-producer.Errors():
		t.Log("发送失败", err.Err, err.Msg)

	}
}

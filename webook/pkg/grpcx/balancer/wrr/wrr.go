package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

const Name = "custom_weighted_round_robin"

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &PickBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type PickBuilder struct {
}

func (p *PickBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for sc, scInfo := range info.ReadySCs {
		metadata, _ := scInfo.Address.Metadata.(map[string]any)
		weightVal, _ := metadata["weight"]
		weight, _ := weightVal.(float64)

		conns = append(conns, &weightConn{
			SubConn:       sc,
			weight:        int(weight),
			currentWeight: int(weight),
		})
	}
	return &Picker{}
}

type Picker struct {
	conns []*weightConn
	// 为了避免并发同时操作了权重，所以加锁
	lock sync.Mutex
}

// Pick 实现加权平滑轮询
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	// 判断是否为空
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var total int
	var maxCC *weightConn
	for _, c := range p.conns {
		// 1.计算所有的权重之和
		total += c.weight
		// 2.每个节点当前权重和初始权重相加得到最新的当前权重
		c.currentWeight += c.weight
		// 3.选择当前权重最大的节点
		if maxCC == nil || maxCC.currentWeight < c.currentWeight {
			maxCC = c
		}
	}
	// 4.成功后将选中节点的当前权重减去总权重
	maxCC.currentWeight -= total

	return balancer.PickResult{
		SubConn: maxCC.SubConn,
		// 这个是调用回调接口
		Done: func(info balancer.DoneInfo) {
			// 要在这里进一步调整weight/currentWeight
		},
	}, nil
}

type weightConn struct {
	balancer.SubConn
	weight        int
	currentWeight int
}

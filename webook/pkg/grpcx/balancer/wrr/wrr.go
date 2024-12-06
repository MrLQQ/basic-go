package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
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
	return &Picker{
		conns: conns,
	}
}

type Picker struct {
	conns []*weightConn
	// 为了避免并发同时操作了权重，所以加锁
	lock sync.Mutex
}

// Pick 实现加权平滑轮询
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	//p.lock.Lock()
	//defer p.lock.Unlock()
	// 判断是否为空
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var total int
	var maxCC *weightConn
	for _, c := range p.conns {
		c.mutex.Lock()
		// 1.计算所有的权重之和
		total += c.weight
		// 2.每个节点当前权重和初始权重相加得到最新的当前权重
		c.currentWeight += c.weight
		// 3.选择当前权重最大的节点
		if maxCC == nil || maxCC.currentWeight < c.currentWeight {
			maxCC = c
		}
		c.mutex.Unlock()
	}
	maxCC.mutex.Lock()
	// 4.成功后将选中节点的当前权重减去总权重
	maxCC.currentWeight -= total
	maxCC.mutex.Unlock()

	return balancer.PickResult{
		SubConn: maxCC.SubConn,
		// 这个是调用回调接口
		Done: func(info balancer.DoneInfo) {
			// 要在这里进一步调整weight/currentWeight
			// failover 在这里做文章
			// 根据调用结果的具体错误信息进行容错
			// 1.如果要是触发了限流
			// 1.1 你可以考虑直接挪走这个节点，后面再挪回来
			// 1.2 你可考虑直接将 weight/currentWight 调整到极低
			// 2.如果触发了熔断 可以优先考虑1.1的方法
			// 3.降级呢 可以优先考虑1.1的方法

			//------------简单实现动态调整权重---------------
			maxCC.mutex.Lock()
			defer maxCC.mutex.Unlock()
			// 权重最小下限校验
			if info.Err != nil && maxCC.weight == 1 {
				return
			}
			// 权重最大上限校验
			if info.Err == nil && maxCC.weight == math.MaxUint32 {
				return
			}

			// 动态调整权重，调用成功提升权重，调用失败降低权重
			if info.Err != nil {
				maxCC.weight--
			} else {
				maxCC.weight++
			}
		},
	}, nil
}

type weightConn struct {
	mutex sync.Mutex

	balancer.SubConn
	weight        int
	currentWeight int
}

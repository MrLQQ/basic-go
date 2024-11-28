package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var (
	//go:embed lua/incr_cnt.lua
	luaIncrCnt string
	//go:embed lua/interactive_ranking_incr.lua
	luaRankingIncr string
	//go:embed lua/interactive_ranking_set.lua
	luaRankingSet string
)

var RankingUpdateErr = errors.New("指定的元素不存在")

const fieldReadCnt = "read_cnt"
const fieldLikeCnt = "like_cnt"
const fieldCollectCnt = "collect_cnt"

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, id int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, id int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, id int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, res domain.Interactive) error
	// IncrRankingIfPresent 如果排名数据存在就+1
	IncrRankingIfPresent(ctx context.Context, biz string, bizId int64) error
	// SetRankingScore 如果排名数据不存在就把数据库中读取到的更新到缓存，如果更新过就+1
	SetRankingScore(ctx context.Context, biz string, bizId int64, score int64) error
	// LikeTop 基本实现，是借助 zset
	// 1.前100名是一个高频数据，你可以结合本地缓存。
	//	你可以定时刷新本地缓存，比如说每5s调用LikeTop,放进去本地缓存
	// 2.如果你有一亿数据，怎样维护？zset放一亿个元素， redis撑不住
	//   2.1 不是真的维护一亿，而是维护近期的数据的点赞数，比如三天内的数据
	//   2.2 你要分 key。这是Redis解决大数据结构常见的方案
	// 3.借助定时任务，每分钟计算一次，如果有很多数据，一分钟不够便利一遍
	// 4.我每次计算，算1000名，然后借助zset来实时维护这1000名的分数
	LikeTop(ctx context.Context, biz string) ([]domain.Interactive, error)
}

type InteractiveRedisCache struct {
	client redis.Cmdable
}

func NewInteractiveRedisCache(client redis.Cmdable) InteractiveCache {
	return &InteractiveRedisCache{
		client: client,
	}
}

func (r *InteractiveRedisCache) IncrRankingIfPresent(ctx context.Context, biz string, bizId int64) error {
	res, err := r.client.Eval(ctx, luaRankingIncr, []string{r.rankingKey(biz)}, bizId).Result()
	if err != nil {
		return err
	}
	if res.(int64) == 0 {
		return RankingUpdateErr
	}
	return nil
}

func (r *InteractiveRedisCache) SetRankingScore(ctx context.Context, biz string, bizId int64, count int64) error {
	return r.client.Eval(ctx, luaRankingSet, []string{r.rankingKey(biz)}, bizId, count).Err()
}

// BatchSetRankingScore 设置整个数据
func (r *InteractiveRedisCache) BatchSetRankingScore(ctx context.Context, biz string, interactives []domain.Interactive) error {
	members := make([]redis.Z, 0, len(interactives))
	for _, interactive := range interactives {
		members = append(members, redis.Z{
			Score:  float64(interactive.LikeCnt),
			Member: interactive.BizId,
		})
	}
	return r.client.ZAdd(ctx, r.rankingKey(biz), members...).Err()
}

func (r *InteractiveRedisCache) LikeTop(ctx context.Context, biz string) ([]domain.Interactive, error) {
	var start int64 = 0
	var end int64 = 99
	key := fmt.Sprintf("too_100_%s", biz)
	res, err := r.client.ZRevRangeWithScores(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	interactives := make([]domain.Interactive, 0, 100)
	for i := 0; i < len(res); i++ {
		id, _ := strconv.ParseInt(res[i].Member.(string), 10, 64)
		interactives = append(interactives, domain.Interactive{
			Biz:     biz,
			BizId:   id,
			LikeCnt: int64(res[i].Score),
		})
	}
	return interactives, nil
}

func (i *InteractiveRedisCache) Set(ctx context.Context,
	biz string, bizId int64,
	res domain.Interactive) error {
	key := i.key(biz, bizId)
	err := i.client.HSet(ctx, key, fieldCollectCnt, res.CollectCnt,
		fieldReadCnt, res.ReadCnt,
		fieldLikeCnt, res.LikeCnt,
	).Err()
	if err != nil {
		return err
	}
	return i.client.Expire(ctx, key, time.Minute*15).Err()
}

func (i *InteractiveRedisCache) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	key := i.key(biz, id)
	res, err := i.client.HGetAll(ctx, key).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, ErrKeyNotExist
	}
	var intr domain.Interactive
	intr.BizId = id
	// 这边是可以忽略错误的
	intr.CollectCnt, _ = strconv.ParseInt(res[fieldCollectCnt], 10, 64)
	intr.LikeCnt, _ = strconv.ParseInt(res[fieldLikeCnt], 10, 64)
	intr.ReadCnt, _ = strconv.ParseInt(res[fieldReadCnt], 10, 64)
	return intr, nil
}

func (i *InteractiveRedisCache) IncrCollectCntIfPresent(ctx context.Context,
	biz string, id int64) error {
	key := i.key(biz, id)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldCollectCnt, 1).Err()
}

func (i *InteractiveRedisCache) IncrLikeCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, 1).Err()
}

func (i *InteractiveRedisCache) DecrLikeCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := i.key(biz, bizId)
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldLikeCnt, -1).Err()
}

func (i *InteractiveRedisCache) IncrReadCntIfPresent(ctx context.Context,
	biz string, bizId int64) error {
	key := i.key(biz, bizId)
	// 不是特别需要处理 res
	//res, err := i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Int()
	return i.client.Eval(ctx, luaIncrCnt, []string{key}, fieldReadCnt, 1).Err()
}

func (i *InteractiveRedisCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}

func (r *InteractiveRedisCache) rankingKey(biz string) string {
	return fmt.Sprintf("top_100_%s", biz)
}

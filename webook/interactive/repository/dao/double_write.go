package dao

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

type DoubleWriteDAO struct {
	src     InteractiveDAO
	dst     InteractiveDAO
	pattern *atomicx.Value[string]
	l       logger.LoggerV1
}

func NewDoubleWriteDAO(src *gorm.DB, dst *gorm.DB, l logger.LoggerV1) *DoubleWriteDAO {
	return &DoubleWriteDAO{
		src:     NewGORMInteractiveDAO(src),
		dst:     NewGORMInteractiveDAO(dst),
		l:       l,
		pattern: atomicx.NewValueOf(PatternSrcFirst),
	}
}

func (d *DoubleWriteDAO) UpdatePattern(pattern string) {
	d.pattern.Store(pattern)
}

var errUnknownPatten = errors.New("未知的写模式")

func (d *DoubleWriteDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcOnly:
		return d.src.IncrReadCnt(ctx, biz, bizId)
	case PatternSrcFirst:
		err := d.src.IncrReadCnt(ctx, biz, bizId)
		if err != nil {
			return err
		}
		err = d.dst.IncrReadCnt(ctx, biz, bizId)
		if err != nil {
			// 要不要return？
			// 正常来说，我们认为双写阶段，src成功了就算业务上成功了
			d.l.Error("双写写入dst，失败", logger.Error(err),
				logger.Int64("biz_id", bizId),
				logger.String("biz", biz))
		}
		return nil
	case PatternDstFirst:
		err := d.dst.IncrReadCnt(ctx, biz, bizId)
		if err == nil {
			err1 := d.src.IncrReadCnt(ctx, biz, bizId)
			if err1 != nil {
				d.l.Error("双写写入src，失败", logger.Error(err1),
					logger.Int64("biz_id", bizId),
					logger.String("biz", biz))
			}
		}
		return err
	case PatternDstOnly:
		return d.dst.IncrReadCnt(ctx, biz, bizId)
	default:
		return errUnknownPatten
	}
}

func (d *DoubleWriteDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	pattern := d.pattern.Load()
	switch pattern {
	case PatternSrcFirst, PatternSrcOnly:
		return d.src.Get(ctx, biz, id)
	case PatternDstOnly, PatternDstFirst:
		return d.dst.Get(ctx, biz, id)
	default:
		return Interactive{}, errUnknownPatten
	}
}

func (d *DoubleWriteDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	//TODO implement me
	panic("implement me")
}

const (
	PatternSrcOnly  = "src_only"
	PatternSrcFirst = "src_first"
	PatternDstFirst = "dst_first"
	PatternDstOnly  = "dst_only"
)

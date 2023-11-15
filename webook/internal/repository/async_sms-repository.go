package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/dao"
	"context"
	"github.com/ecodeclub/ekit/sqlx"
)

var ErrWaitingSMSNotFound = dao.ErrWaitingSMSNotFound

type AsyncSmsRepository interface {
	PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error)
	ReportScheduleResult(ctx context.Context, id int64, res bool) error
	Add(ctx context.Context, as domain.AsyncSms) error
}

type asyncSmsRepository struct {
	dao dao.AsyncSmsDAO
}

func NewAsyncSMSRepository(dao dao.AsyncSmsDAO) AsyncSmsRepository {
	return &asyncSmsRepository{
		dao: dao,
	}
}

func (a *asyncSmsRepository) PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error) {
	as, err := a.dao.GetWaitingSMS(ctx)
	if err != nil {
		return domain.AsyncSms{}, err
	}
	return domain.AsyncSms{
		Id:       as.Id,
		TplId:    as.Config.Val.TplId,
		Numbers:  as.Config.Val.Numbers,
		Args:     as.Config.Val.Args,
		RetryMax: as.RetryMax,
	}, nil
}

func (a *asyncSmsRepository) ReportScheduleResult(ctx context.Context, id int64, success bool) error {
	if success {
		return a.dao.MarkSuccess(ctx, id)
	}
	return a.dao.MarkFailed(ctx, id)
}

func (a *asyncSmsRepository) Add(ctx context.Context, as domain.AsyncSms) error {
	return a.dao.Insert(ctx, dao.AsyncSms{
		Config: sqlx.JsonColumn[dao.SmsConfig]{
			Val: dao.SmsConfig{
				TplId:   as.TplId,
				Args:    as.Args,
				Numbers: as.Numbers,
			},
			Valid: true,
		},
		RetryMax: as.RetryMax,
	})

}

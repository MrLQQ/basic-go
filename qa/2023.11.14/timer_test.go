//go:build demo

package demo

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	// 每秒钟的定时器
	//tm := time.NewTimer(time.Second)
	// 这里就是触发了
	//for now := range tm.C {
	// 在这里异步发送
	//}
}

type smsHomework struct {
}

func (s *smsHomework) StartAsyncCycle() {
	for {
		req, err := s.preempt()
		if err == ReqNotFound {
			// 数据没找到，我睡一秒
			time.Sleep(time.Second)
			continue
		} else if err != nil {
			// 记录日志
		}
		// 记得给一个 ctx
		err := s.Send(ctx, req.TplId, req.Args, req.Numbers...)
		// 这里记得处理 error
		// 如果你设计了重试，这里要考虑重试的问题
		if err == nil {
			// 我这个请求处理完毕了
			s.repo.MarkSuccess(req)
			continue
		}
		// 你在这边要增加已重试次数
	}
}

func (s *smsHomework) preempt() (AsyncReq, error) {
	// 你这边可以考虑多个实例的问题
	// 例如，有一个请求发送给 152xxxxxx 一条短信了，
	panic("从数据库里面捞一个异步发送请求过来")
}

func (s *smsHomework) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	if s.needAsync() {
		// 需要转异步了
		// 转储到数据库里面去
		s.repo.Store(AsyncReq{})
		return nil
	}
	return s.svc.Send(ctx, tplId, args, numbers...)
}

func (s *smsHomework) needAsync() bool {
	// 你要核心考虑的判定方式
}

type AsyncReq struct {
	TplId   string
	Args    []string
	Numbers []string
	// 你还可以有其他的
}

var l sync.Mutex

func Lock() {
	l.Lock()
	DoSomething()
	l.Unlock()
}

func LockV1() {
	l.Lock()
	defer l.Unlock()
	DoSomething()

}

func DoSomething() {
	panic("abc")
}

package inventory

import (
	"fmt"
	"github.com/xiaobudongzhang/seata-golang/client/context"
	"sync"
)

var (
	s *service
	m sync.RWMutex
)

type service struct {

}

type Service interface {
	Sell(bookId, userId int64, ctx *context.RootContext) (id int64, err error)

	Confirm(id int64, state int) (err error)
}


func GetService()(Service, error)  {
	if s == nil {
		return nil, fmt.Errorf("getservice 为初始化")
	}
	return s, nil
}

func Init()  {
	m.Lock()
	defer m.Unlock()

	if s != nil {
		return
	}
	s = &service{}
}
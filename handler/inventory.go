package handler

import (
	"context"

	"fmt"
	"github.com/micro/go-micro/v2/util/log"
	inv "github.com/xiaobudongzhang/micro-inventory-srv/model/inventory"
	proto "github.com/xiaobudongzhang/micro-inventory-srv/proto/inventory"
	"github.com/micro/go-micro/v2/metadata"
	context2 "github.com/xiaobudongzhang/seata-golang/client/context"
)

var (
	invService inv.Service
)

type Service struct {
}

// Init 初始化handler
func Init() {
	invService, _ = inv.GetService()
}

// Sell 库存销存
func (e *Service) Sell(ctx context.Context, req *proto.Request, rsp *proto.Response) (err error) {

	rex, erx:= metadata.FromContext(ctx)
	fmt.Printf("meta:%v-%v", rex, erx)


	rootContext := &context2.RootContext{Context:ctx}
	rootContext.Bind(rex["Xid"])


	id, err := invService.Sell(req.BookId, req.UserId, rootContext)
	if err != nil {
		log.Logf("[Sell] 销存失败，bookId：%d，userId: %d，%s", req.BookId, req.UserId, err)
		rsp.Success = false
		return
	}

	rsp.InvH = &proto.InvHistory{
		Id: id,
	}

	rsp.Success = true
	return nil
}

// Confirm 库存销存 确认
func (e *Service) Confirm(ctx context.Context, req *proto.Request, rsp *proto.Response) (err error) {
	err = invService.Confirm(req.HistoryId, int(req.HistoryState))
	if err != nil {
		log.Logf("[Confirm] 确认销存失败，%s", err)
		rsp.Success = false
		return
	}

	rsp.Success = true
	return nil
}

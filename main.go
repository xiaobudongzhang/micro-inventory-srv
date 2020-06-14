package main

import (
	"fmt"

	"github.com/xiaobudongzhang/micro-basic/common"
	"github.com/xiaobudongzhang/micro-inventory-srv/handler"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-plugins/config/source/grpc/v2"

	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	log "github.com/micro/go-micro/v2/util/log"
	"github.com/xiaobudongzhang/micro-basic/basic"
	"github.com/xiaobudongzhang/micro-basic/config"
	"github.com/xiaobudongzhang/micro-inventory-srv/model"

	proto "github.com/xiaobudongzhang/micro-inventory-srv/proto/inventory"

	openTrace "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	tracer "github.com/xiaobudongzhang/micro-plugins/tracer/myjaeger"
)

var (
	appName = "inventory_service"
	cfg     = &appCfg{}
)

type appCfg struct {
	common.AppCfg
}

func main() {
	initCfg()
	micReg := etcd.NewRegistry(registryOptions)

	t, io, err1 := tracer.NewTracer(cfg.Name, "")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// New Service
	service := micro.NewService(
		micro.Name("mu.micro.book.service.inventory"),
		micro.Registry(micReg),
		micro.Version("latest"),
		micro.WrapHandler(openTrace.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	// Initialise service
	service.Init(micro.Action(func(c *cli.Context) error {
		model.Init()
		handler.Init()
		return nil
	}),
	)

	// Register Handler
	proto.RegisterInventoryHandler(service.Server(), new(handler.Service))

	// Register Struct as Subscriber
	//micro.RegisterSubscriber("mu.micro.book.service.inventory", service.Server(), new(subscriber.Inventory))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	etcdCfg := &common.Etcd{}
	err := config.C().App("etcd", etcdCfg)
	if err != nil {

		log.Log(err)
		panic(err)
	}
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.Host, etcdCfg.Port)}
}

func initCfg() {
	source := grpc.NewSource(
		grpc.WithAddress("127.0.0.1:9600"),
		grpc.WithPath("micro"),
	)

	basic.Init(config.WithSource(source))

	err := config.C().App(appName, cfg)
	if err != nil {
		panic(err)
	}

	log.Logf("配置 cfg:%v", cfg)

	return
}

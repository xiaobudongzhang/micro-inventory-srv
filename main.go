package main

import (
	"fmt"

	"github.com/xiaobudongzhang/micro-inventory-srv/handler"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"

	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	log "github.com/micro/go-micro/v2/util/log"
	"github.com/xiaobudongzhang/micro-basic/basic"
	"github.com/xiaobudongzhang/micro-basic/config"
	"github.com/xiaobudongzhang/micro-inventory-srv/model"

	proto "github.com/xiaobudongzhang/micro-inventory-srv/proto/inventory"
)

func main() {
	basic.Init()

	micReg := etcd.NewRegistry(registryOptions)

	// New Service
	service := micro.NewService(
		micro.Name("mu.micro.book.service.inventory"),
		micro.Registry(micReg),
		micro.Version("latest"),
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
	etcdCfg := config.GetEtcdConfig()
	ops.Addrs = []string{fmt.Sprintf("%s:%d", etcdCfg.GetHost(), etcdCfg.GetPort())}
}

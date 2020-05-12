module github.com/xiaobudongzhang/micro-inventory-srv

go 1.14

replace github.com/xiaobudongzhang/micro-basic => /data/ndemo/micro-basic

replace github.com/xiaobudongzhang/micro-inventory-srv => /data/ndemo/micro-inventory-srv

replace github.com/xiaobudongzhang/micro-payment-srv => /data/ndemo/micro-payment-srv

replace github.com/xiaobudongzhang/micro-order-srv => /data/ndemo/micro-order-srv

replace github.com/xiaobudongzhang/micro-plugins => /data/ndemo/micro-plugins

require (
	github.com/golang/protobuf v1.4.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.5.0
	github.com/micro/go-plugins v1.5.1 // indirect
	github.com/micro/go-plugins/config/source/grpc/v2 v2.5.0 // indirect
	github.com/xiaobudongzhang/micro-basic v1.1.5
	github.com/xiaobudongzhang/micro-plugins v0.0.0-00010101000000-000000000000
)

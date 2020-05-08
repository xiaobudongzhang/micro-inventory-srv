module github.com/xiaobudongzhang/micro-inventory-srv

go 1.14

replace github.com/xiaobudongzhang/micro-basic => /data/ndemo/micro-basic

replace github.com/xiaobudongzhang/micro-inventory-srv => /data/ndemo/micro-inventory-srv

require (
	github.com/golang/protobuf v1.4.0
	github.com/micro-in-cn/tutorials/microservice-in-micro v0.0.0-20200415151649-6b5af13cdcea
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.5.0
	github.com/xiaobudongzhang/micro-basic v1.1.5
)

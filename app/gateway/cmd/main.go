package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
	"micro-todoList-k8s/app/common"
	"micro-todoList-k8s/app/gateway/router"
	"micro-todoList-k8s/app/gateway/rpc"
	"micro-todoList-k8s/app/user/repository/cache"
	"micro-todoList-k8s/config"
	log "micro-todoList-k8s/pkg/logger"
	"time"
)

func main() {
	config.Init()
	rpc.InitRPC()
	cache.InitCache()
	log.InitLog()

	fmt.Println(config.C.Service.Env)
	// consul
	consulReg := consul.NewRegistry(registry.Addrs(fmt.Sprintf("%s:%s", config.C.Consul.ConsulHost, config.C.Consul.ConsulPort)))

	// 初始化 Tracer
	tracer := common.GetTracer(config.C.Server.GateWayServiceName, config.C.Server.GateWayServiceAddress)
	// 初始化 Prometheus
	common.PrometheusBoot(config.C.Prometheus.PrometheusGateWayPath, config.C.Prometheus.PrometheusGateWayAddress)

	// 创建微服务实例，使用gin暴露http接口并注册到consul
	server := web.NewService(
		web.Name(config.C.Server.GateWayServiceName),
		web.Address(config.C.Server.GateWayServiceAddress),
		// 将服务调用实例使用gin处理
		web.Handler(router.NewRouter(tracer)),
		web.Registry(consulReg),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*15),
		web.Metadata(map[string]string{"protocol": "http"}),
	)
	// 接收命令行参数
	_ = server.Init()
	_ = server.Run()
}

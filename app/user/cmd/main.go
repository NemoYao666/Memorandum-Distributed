package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"micro-todoList-k8s/app/common"
	"micro-todoList-k8s/app/user/repository/db/dao"
	"micro-todoList-k8s/app/user/service"
	"micro-todoList-k8s/config"
	"micro-todoList-k8s/idl/pb"
)

func main() {
	config.Init()
	dao.InitDB()

	// consul
	consulReg := consul.NewRegistry(registry.Addrs(fmt.Sprintf("%s:%s", config.C.Consul.ConsulHost, config.C.Consul.ConsulPort)))

	// 初始化 Tracer
	tracer := common.GetTracer(config.C.Server.UserServiceName, config.C.Server.UserServiceAddress)
	tracerHandler := opentracing.NewHandlerWrapper(tracer)
	// 初始化 Prometheus
	common.PrometheusBoot(config.C.Prometheus.PrometheusUserServicePath, config.C.Prometheus.PrometheusUserServiceAddress)

	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name(config.C.Server.UserServiceName), // 微服务名字
		micro.Address(config.C.Server.UserServiceAddress),
		micro.Registry(consulReg), // consul注册件
		micro.WrapHandler(tracerHandler),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)
	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterUserServiceHandler(microService.Server(), service.GetUserSrv())
	// 启动微服务
	_ = microService.Run()
}

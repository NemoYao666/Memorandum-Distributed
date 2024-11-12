package main

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/consul"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"micro-todoList-k8s/app/common"

	"micro-todoList-k8s/app/task/repository/db/dao"
	"micro-todoList-k8s/app/task/repository/mq"
	"micro-todoList-k8s/app/task/script"
	"micro-todoList-k8s/app/task/service"
	"micro-todoList-k8s/config"
	"micro-todoList-k8s/idl/pb"
	log "micro-todoList-k8s/pkg/logger"
)

func main() {
	config.Init()
	dao.InitDB()
	mq.InitRabbitMQ()
	log.InitLog()

	// 启动一些脚本
	loadingScript()

	// consul
	consulReg := consul.NewRegistry(registry.Addrs(fmt.Sprintf("%s:%s", config.ConsulHost, config.ConsulPort)))

	// 初始化 Tracer
	tracer := common.GetTracer(config.TaskServiceName, config.TaskServiceAddress)
	tracerHandler := opentracing.NewHandlerWrapper(tracer)
	// 初始化 Prometheus
	common.PrometheusBoot(config.PrometheusTaskServicePath, config.PrometheusTaskServiceAddress)

	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name(config.TaskServiceName), // 微服务名字
		micro.Address(config.TaskServiceAddress),
		micro.Registry(consulReg), // consul注册件
		micro.WrapHandler(tracerHandler),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterTaskServiceHandler(microService.Server(), service.GetTaskSrv())
	// 启动微服务
	_ = microService.Run()
}

// 异步启动监听消息队列
func loadingScript() {
	ctx := context.Background()
	go script.TaskCreateSync(ctx)
}

package rpc

import (
	"github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"go-micro.dev/v4"
	"micro-todoList-k8s/app/common"
	"micro-todoList-k8s/config"

	"micro-todoList-k8s/app/gateway/wrappers"
	"micro-todoList-k8s/idl/pb"
)

var (
	UserService pb.UserService
	TaskService pb.TaskService
)

// InitRPC 初始化RPC客户端
func InitRPC() {
	userTracer := common.GetTracer(config.C.Server.UserClientName, config.C.Server.UserServiceAddress)
	userTracerClient := opentracing.NewClientWrapper(userTracer)
	// 用户
	userMicroService := micro.NewService(
		micro.Name(config.C.Server.UserClientName),
		micro.WrapClient(wrappers.NewUserWrapper),
		micro.WrapClient(userTracerClient),
	)
	// 用户服务调用实例
	userService := pb.NewUserService(config.C.Server.UserServiceName, userMicroService.Client())

	taskTracer := common.GetTracer(config.C.Server.TaskClientName, config.C.Server.TaskServiceAddress)
	taskTracerClient := opentracing.NewClientWrapper(taskTracer)
	// task
	taskMicroService := micro.NewService(
		micro.Name(config.C.Server.TaskClientName),
		micro.WrapClient(wrappers.NewTaskWrapper),
		micro.WrapClient(taskTracerClient),
	)
	taskService := pb.NewTaskService(config.C.Server.TaskServiceName, taskMicroService.Client())

	UserService = userService
	TaskService = taskService
}

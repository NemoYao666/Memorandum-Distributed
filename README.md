# Memorandum
# 简单备忘录

基于Golang V1.21，Go-Micro v4，Gin，Gorm  
基于Consul实现服务注册中心及配置中心  
支持gateway和各模块之间的Protobuf RPC通信  
支持限流熔断机制，JWT token验证，swagger API文档生成  
基于CentOS7 Docker部署Zipkin链路追踪，Prometheus Grafana监控，Redis登陆缓存，RabbitMQ任务创建消息队列  

****

# Running Environment

## Win + Linux
  服务列表 与 WEB UI
```shell
# win
mysql 
swagger      http://127.0.0.1:4000/swagger/index.html
# linux docker
redis   
rabbitMQ     http://127.0.0.1:15672
zipkin       http://127.0.0.1:9411/zipkin
prometheus   http://127.0.0.1:9090/targets
# linux server
consul       http://127.0.0.1:8500
grafana      http://127.0.0.1:3000/login
```

  启动流程
```shell
# docker linux
systemctl start docker # redis rabbitMQ auto start
docker start # zipkin_container_id
```  

```shell
# consul linux
cd /opt/consul
./start_consul.sh
```  

```shell
# prometheus linux
# 配置抓取路径 
cd /opt/prometheus
gedit prometheus.yml  # web服务所在的IP
docker start # prometheus_container_id
# grafana
systemctl start grafana-server.service
```
  
****

# Project Architecture
## 1.micro_todolist 项目总体
```
micro-todolist/
├── app                   // 各个微服务
│   ├── common            // 链路追踪、监控
│   ├── docs              // swagger文档
│   ├── gateway           // 网关
│   ├── task              // 任务模块微服务
│   └── user              // 用户模块服务
├── bin                   // 编译后的二进制文件模块
├── config                // 配置文件
├── consts                // 定义的常量
├── doc                   // 接口文档
├── idl                   // protoc文件
│   └── pb                // 放置生成的pb文件
├── logs                  // 放置打印日志模块
├── pkg                   // 各种包
│   ├── ctl               // 用户操作
│   ├── e                 // 统一错误状态码
│   ├── logger            // 日志
│   └── util              // 各种工具、JWT等等..
└── types                 // 定义各种结构体
```

## 2.gateway 网关部分
```
gateway/
├── cmd                   // 启动入口
├── http                  // HTTP请求头
├── metrics               // 监控指标
├── handler               // 视图层
├── logs                  // 放置打印日志模块
├── middleware            // 中间件：跨域，鉴权，监控，链路追踪
├── router                // http 路由模块
├── rpc                   // rpc 调用
└── wrappers              // 熔断
```

## 3.user && task 用户与任务模块
```
task/
├── cmd                   // 启动入口
├── metrics               // 监控指标
├── service               // 业务服务
├── repository            // 持久层
│    ├── db               // 视图层
│    │    ├── dao         // 对数据库进行操作
│    │    └── model       // 定义数据库的模型
│    └── mq               // 放置 mq
├── script                // 监听 mq 的脚本
└── service               // 服务
```


根据`config.yaml.example`编写`config/config.yaml`



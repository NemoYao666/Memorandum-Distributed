# Memorandum-Distributed
# 备忘录-分布式微服务架构

基于Golang V1.21，Go-Micro v4，Gin，Gorm  
基于CentOS7 K8S集群 实现微服务分布式部署  
基于Consul实现服务注册中心及配置中心  
支持Gateway和各微服务模块之间的Protobuf RPC通信  
支持限流熔断机制，JWT token验证，swagger API文档生成  
基于Docker部署Zipkin链路追踪，Prometheus Grafana监控，Redis登陆缓存，RabbitMQ任务创建消息队列  

****

# Running Environment

## Dev: Win + Linux
  服务列表 与 WEB UI
```shell
# win
mysql 
# linux docker IP:105
redis   
rabbitMQ       http://192.168.88.105:15672
zipkin         http://192.168.88.105:9411/zipkin
prometheus     http://192.168.88.105:9090/targets
# linux server IP:105
consul         http://192.168.88.105:8500
grafana        http://192.168.88.105:3000/login
# linux k8s IP:106,107,108
swagger        http://192.168.88.106:4000/swagger/index.html
k8s dashboard  https://192.168.88.106:<port>
```

  启动流程
```shell
# docker linux IP:105
systemctl start docker # redis rabbitMQ auto start
docker start # zipkin_container_id
```  

```shell
# consul linux IP:105
cd /opt/consul
./start_consul.sh
```

```shell
# prometheus linux IP:105
# 配置抓取路径 
cd /opt/prometheus
gedit prometheus.yml  # web服务所在的IP
docker start # prometheus_container_id
# grafana linux IP:105
systemctl start grafana-server.service
```

windows主机启动相关服务即可

## Pre: Win + Linux Docker
IP: 192.168.88.105 机器与windows开启所有服务
镜像已经打包上传到 IP:192.168.88.106 的虚拟机  
docker start # 三个微服务  
读取的consul的环境由启动容器的时候挂载上去的config.yaml里的Env字段为主  
Env=pre

备份：
```shell
# 构建镜像
docker build -t micro-todoList-k8s/gateway:latest -f ./DockerfileGateway .
# 容器启动
cd /root/options/micro-todoList-k8s
docker run -d \
  --name gateway \
  -p 4000:4000 \
  -v $(pwd)/config/config.yaml:/usr/local/bin/config/config.yaml \
  micro-todoList-k8s/gateway:latest
```  

## Prod: Win + Linux K8S
根据打包的镜像部署到K8S集群  
IP:192.168.88.106  Type:Master  
IP:192.168.88.107  Type:Worker  
IP:192.168.88.108  Type:Worker  

IP: 192.168.88.105 机器与windows开启所有服务  

集群初始化 IP:106
```shell
# 所有节点
kubeadm reset
# 只在主节点执行
kubeadm init \
--apiserver-advertise-address=192.168.88.106 \
--control-plane-endpoint=cluster-endpoint \
--image-repository  registry.cn-hangzhou.aliyuncs.com/google_containers \
--kubernetes-version v1.20.9 \
--service-cidr=10.96.0.0/12 \
--pod-network-cidr=172.20.0.0/16
# 执行：
export KUBECONFIG=/etc/kubernetes/admin.conf
```

从节点上执行加入集群 IP:107 IP:108
```shell
kubeadm join cluster-endpoint:6443 --token
```

主节点上修改从节点角色 IP:106
```shell
kubectl label node k8snode2 node-role.kubernetes.io/worker=worker
kubectl label node k8snode3 node-role.kubernetes.io/worker=worker
# 查看集群：
kubectl get nodes
```

应用网络插件 IP:106
```shell
cd /root/options/calico
kubectl apply -f calico.yaml
kubectl get pod -A | grep calico
```

k8s dashboard IP:106
```shell
cd /root/options/k8sDashboard
kubectl apply -f recommended.yaml
# 对外暴露访问端口，修改配置，将--type=ClusterIP修改为--type=NodePort
kubectl edit svc kubernetes-dashboard -n kubernetes-dashboard
# 创建 k8s dashboard 访问账号
kubectl apply -f dashboard-token.yaml
# 查看用户列表
kubectl get secret -n kubernetes-dashboard
# 查看名为admin-user-token-????? 的secret 作为 token
kubectl describe secret admin-user-token-????? -n kubernetes-dashboard
```
应用 k8s 配置 IP:106
```shell
cd /root/options/k8s-server
kubectl apply -f configmap.yaml
```

k8s deploy IP:106
```shell
kubectl apply -f gatewayDeployment.yaml
kubectl apply -f userDeployment.yaml
kubectl apply -f taskDeployment.yaml
```

k8s service IP:106
```shell
# app
kubectl apply -f gatewayService.yaml
kubectl apply -f userService.yaml
kubectl apply -f taskService.yaml
# prometheus
kubectl apply -f gatewayPromethusService.yaml
kubectl apply -f userPromethusService.yaml
kubectl apply -f taskPromethusService.yaml

kubectl get services -A
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
├── bak                   // k8s集群运行相关文件备份
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



package config

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var C *Config

// Config 是用来解析 YAML 文件的结构体
type Config struct {
	Service struct {
		Env string `yaml:"Env"`
	} `yaml:"service"`

	Mysql struct {
		Db         string `yaml:"Db"`
		DbHost     string `yaml:"DbHost"`
		DbPort     string `yaml:"DbPort"`
		DbUser     string `yaml:"DbUser"`
		DbPassWord string `yaml:"DbPassWord"`
		DbName     string `yaml:"DbName"`
		Charset    string `yaml:"Charset"`
	} `yaml:"mysql"`

	RabbitMQ struct {
		RabbitMQ         string `yaml:"RabbitMQ"`
		RabbitMQUser     string `yaml:"RabbitMQUser"`
		RabbitMQPassWord string `yaml:"RabbitMQPassWord"`
		RabbitMQHost     string `yaml:"RabbitMQHost"`
		RabbitMQPort     string `yaml:"RabbitMQPort"`
	} `yaml:"rabbitmq"`

	Consul struct {
		ConsulHost string `yaml:"ConsulHost"`
		ConsulPort string `yaml:"ConsulPort"`
		ConsulKey  string `yaml:"ConsulKey"`
	} `yaml:"consul"`

	Server struct {
		GateWayServiceName    string `yaml:"GateWayServiceName"`
		GateWayServiceAddress string `yaml:"GateWayServiceAddress"`
		UserServiceName       string `yaml:"UserServiceName"`
		UserClientName        string `yaml:"UserClientName"`
		UserServiceAddress    string `yaml:"UserServiceAddress"`
		TaskServiceName       string `yaml:"TaskServiceName"`
		TaskClientName        string `yaml:"TaskClientName"`
		TaskServiceAddress    string `yaml:"TaskServiceAddress"`
	} `yaml:"server"`

	Zipkin struct {
		ZipkinUrl string `yaml:"ZipkinUrl"`
	} `yaml:"zipkin"`

	Prometheus struct {
		PrometheusGateWayPath        string `yaml:"PrometheusGateWayPath"`
		PrometheusGateWayAddress     string `yaml:"PrometheusGateWayAddress"`
		PrometheusUserServicePath    string `yaml:"PrometheusUserServicePath"`
		PrometheusUserServiceAddress string `yaml:"PrometheusUserServiceAddress"`
		PrometheusTaskServicePath    string `yaml:"PrometheusTaskServicePath"`
		PrometheusTaskServiceAddress string `yaml:"PrometheusTaskServiceAddress"`
	} `yaml:"prometheus"`

	Redis struct {
		RedisHost     string `yaml:"RedisHost"`
		RedisPort     string `yaml:"RedisPort"`
		RedisPassword string `yaml:"RedisPassword"`
	} `yaml:"redis"`
}

func Init() {
	// 读取 YAML 文件内容
	fileContent, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件错误: %v", err)
	}

	// 解析 YAML 文件内容
	var config Config
	err = yaml.Unmarshal(fileContent, &config)
	if err != nil {
		log.Fatalf("解析配置文件错误: %v", err)
	}

	env := config.Service.Env
	if env == "dev" {
		// 设置全局配置
		SetGlobalConfig(&config)
	} else if env == "prod" {
		var configNew = &Config{}
		configNew.Consul.ConsulHost = config.Consul.ConsulHost
		configNew.Consul.ConsulPort = config.Consul.ConsulPort
		configNew.Consul.ConsulKey = config.Consul.ConsulKey
		LoadConsulConfig(configNew)
		SetGlobalConfig(configNew)
	}

}

func SetGlobalConfig(config *Config) {
	C = config
	fmt.Println("配置文件加载成功")
}

func LoadConsulConfig(config *Config) {
	// 初始化 Consul 客户端
	consulConfig := api.DefaultConfig()
	consulConfig.Address = fmt.Sprintf("%s:%s", config.Consul.ConsulHost, config.Consul.ConsulPort)
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("初始化 Consul 客户端失败: %v", err)
	}
	// 从 Consul 中获取配置
	kvPair, _, err := consulClient.KV().Get(config.Consul.ConsulKey, nil)
	if err != nil {
		log.Fatalf("从 Consul 获取配置失败: %v", err)
	}
	// 反序列化 YAML 配置
	err = yaml.Unmarshal(kvPair.Value, &config)
	if err != nil {
		log.Fatalf("解析 YAML 配置失败: %v", err)
	}
}

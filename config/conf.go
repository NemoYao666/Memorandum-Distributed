package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
	Charset    string

	RabbitMQ         string
	RabbitMQUser     string
	RabbitMQPassWord string
	RabbitMQHost     string
	RabbitMQPort     string

	ConsulHost string
	ConsulPort string

	GateWayServiceName    string
	GateWayServiceAddress string
	UserServiceName       string
	UserClientName        string
	UserServiceAddress    string
	TaskServiceName       string
	TaskClientName        string
	TaskServiceAddress    string

	ZipkinUrl string

	PrometheusGateWayPath        string
	PrometheusGateWayAddress     string
	PrometheusUserServicePath    string
	PrometheusUserServiceAddress string
	PrometheusTaskServicePath    string
	PrometheusTaskServiceAddress string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDbName   int
)

func Init() {
	file, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadMysqlData(file)
	LoadConsul(file)
	LoadRabbitMQ(file)
	LoadServer(file)
	LoadZipkin(file)
	LoadPrometheus(file)
	LoadRedisData(file)
}

func LoadMysqlData(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
	Charset = file.Section("mysql").Key("Charset").String()
}

func LoadRabbitMQ(file *ini.File) {
	RabbitMQ = file.Section("rabbitmq").Key("RabbitMQ").String()
	RabbitMQUser = file.Section("rabbitmq").Key("RabbitMQUser").String()
	RabbitMQPassWord = file.Section("rabbitmq").Key("RabbitMQPassWord").String()
	RabbitMQHost = file.Section("rabbitmq").Key("RabbitMQHost").String()
	RabbitMQPort = file.Section("rabbitmq").Key("RabbitMQPort").String()
}

func LoadConsul(file *ini.File) {
	ConsulHost = file.Section("consul").Key("ConsulHost").String()
	ConsulPort = file.Section("consul").Key("ConsulPort").String()
}

func LoadServer(file *ini.File) {
	GateWayServiceName = file.Section("server").Key("GateWayServiceName").String()
	GateWayServiceAddress = file.Section("server").Key("GateWayServiceAddress").String()
	UserServiceName = file.Section("server").Key("UserServiceName").String()
	UserClientName = file.Section("server").Key("UserClientName").String()
	UserServiceAddress = file.Section("server").Key("UserServiceAddress").String()
	TaskServiceName = file.Section("server").Key("TaskServiceName").String()
	TaskClientName = file.Section("server").Key("TaskClientName").String()
	TaskServiceAddress = file.Section("server").Key("TaskServiceAddress").String()
}

func LoadZipkin(file *ini.File) {
	ZipkinUrl = file.Section("zipkin").Key("ZipkinUrl").String()
}

func LoadPrometheus(file *ini.File) {
	PrometheusGateWayPath = file.Section("prometheus").Key("PrometheusGateWayPath").String()
	PrometheusGateWayAddress = file.Section("prometheus").Key("PrometheusGateWayAddress").String()
	PrometheusUserServicePath = file.Section("prometheus").Key("PrometheusUserServicePath").String()
	PrometheusUserServiceAddress = file.Section("prometheus").Key("PrometheusUserServiceAddress").String()
	PrometheusTaskServicePath = file.Section("prometheus").Key("PrometheusTaskServicePath").String()
	PrometheusTaskServiceAddress = file.Section("prometheus").Key("PrometheusTaskServiceAddress").String()

}

func LoadRedisData(file *ini.File) {
	RedisHost = file.Section("redis").Key("RedisHost").String()
	RedisPort = file.Section("redis").Key("RedisPort").String()
	RedisPassword = file.Section("redis").Key("RedisPassword").String()
}

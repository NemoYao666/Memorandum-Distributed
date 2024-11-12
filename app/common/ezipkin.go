package common

import (
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
	"micro-todoList-k8s/config"
)

// GetTracer 初始化追踪器
func GetTracer(serviceName string, host string) opentracing.Tracer {

	url := config.C.Zipkin.ZipkinUrl
	zipkinReporter := zipkinhttp.NewReporter(url)

	endpoint, err := zipkin.NewEndpoint(serviceName, host)
	log.Println("GetTracer:", serviceName, host, url)

	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	nativeTracer, err := zipkin.NewTracer(zipkinReporter, zipkin.WithLocalEndpoint(endpoint))

	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.InitGlobalTracer(tracer)

	return tracer

}

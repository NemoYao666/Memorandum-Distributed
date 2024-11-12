package common

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func PrometheusBoot(path string, address string) {
	http.Handle(fmt.Sprintf(path), promhttp.Handler())
	//启动web 服务
	go func() {
		address := address
		err := http.ListenAndServe(address, nil)
		if err != nil {
			log.Fatal("启动失败")
		}
	}()

}

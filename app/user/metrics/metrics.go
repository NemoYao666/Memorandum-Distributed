package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	QueryUserLoginCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "query_user_login_counter",
			Help: "Total number of user login queries",
		},
		[]string{"counts"},
	)
)

func init() {
	prometheus.MustRegister(QueryUserLoginCounter)
}

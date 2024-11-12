package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	QueryGetTaskListCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "query_get_task_list_counter",
			Help: "Total number of get task list queries",
		},
		[]string{"counts"},
	)
)

func init() {
	prometheus.MustRegister(QueryGetTaskListCounter)
}

package script

import (
	"context"

	"micro-todoList-k8s/app/task/repository/mq/task"
	log "micro-todoList-k8s/pkg/logger"
)

func TaskCreateSync(ctx context.Context) {
	tSync := new(task.SyncTask)
	err := tSync.RunTaskCreate(ctx)
	if err != nil {
		log.LogrusObj.Infof("RunTaskCreate:%s", err)
	}
}

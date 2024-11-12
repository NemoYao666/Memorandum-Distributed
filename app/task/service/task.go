package service

import (
	"context"
	"encoding/json"
	"micro-todoList-k8s/app/task/metrics"
	"sync"

	"micro-todoList-k8s/app/task/repository/db/dao"
	"micro-todoList-k8s/app/task/repository/db/model"
	"micro-todoList-k8s/app/task/repository/mq"
	"micro-todoList-k8s/idl/pb"
	"micro-todoList-k8s/pkg/e"
	log "micro-todoList-k8s/pkg/logger"
)

var TaskSrvIns *TaskSrv
var TaskSrvOnce sync.Once

type TaskSrv struct {
}

func GetTaskSrv() *TaskSrv {
	TaskSrvOnce.Do(func() {
		TaskSrvIns = &TaskSrv{}
	})
	return TaskSrvIns
}

// CreateTask 创建备忘录，将备忘录信息生产，放到rabbitMQ消息队列中
//
//	@Summary		CreateTask
//	@Description	CreateTaskDescription
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"token"
//	@Param			req				body		pb.TaskRequest	true	"task"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		500				{string}	string	"bad request"
//	@Router			/api/v1/task [POST]
func (t *TaskSrv) CreateTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) (err error) {
	body, _ := json.Marshal(req) // title，content
	resp.Code = e.SUCCESS
	err = mq.SendMessage2MQ(body)
	if err != nil {
		resp.Code = e.ERROR
		return
	}
	return
}

func TaskMQ2MySQL(ctx context.Context, req *pb.TaskRequest) error {
	m := &model.Task{
		Uid:       uint(req.Uid),
		Title:     req.Title,
		Status:    int(req.Status),
		Content:   req.Content,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	return dao.NewTaskDao(ctx).CreateTask(m)
}

// GetTasksList 实现备忘录服务接口 获取备忘录列表
//
//	@Summary		GetTasksList
//	@Description	GetTasksListDescription
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"token"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		500				{string}	string	"bad request"
//	@Router			/api/v1/tasks [get]
func (t *TaskSrv) GetTasksList(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskListResponse) (err error) {
	resp.Code = e.SUCCESS
	if req.Limit == 0 {
		req.Limit = 10
	}
	// 查找备忘录
	r, count, err := dao.NewTaskDao(ctx).ListTaskByUserId(req.Uid, int(req.Start), int(req.Limit))
	if err != nil {
		resp.Code = e.ERROR
		log.LogrusObj.Error("ListTaskByUserId err:%v", err)
		return
	}
	// 返回proto里面定义的类型
	var taskRes []*pb.TaskModel
	for _, item := range r {
		taskRes = append(taskRes, BuildTask(item))
	}
	resp.TaskList = taskRes
	resp.Count = uint32(count)

	metrics.QueryGetTaskListCounter.WithLabelValues("counts").Inc()

	return
}

// GetTask 获取详细的备忘录
//
//	@Summary		GetTask
//	@Description	GetTaskDescription
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"token"
//	@Param			id				path		string	true	"id"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		500				{string}	string	"bad request"
//	@Router			/api/v1/task/{id} [get]
func (t *TaskSrv) GetTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) (err error) {
	resp.Code = e.SUCCESS
	r, err := dao.NewTaskDao(ctx).GetTaskByTaskIdAndUserId(req.Id, req.Uid)
	if err != nil {
		resp.Code = e.ERROR
		log.LogrusObj.Error("GetTask err:%v", err)
		return
	}
	taskRes := BuildTask(r)
	resp.TaskDetail = taskRes
	return
}

// UpdateTask 修改备忘录
//
//	@Summary		UpdateTask
//	@Description	UpdateTaskDescription
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"token"
//	@Param			id				path		string	true	"id"
//	@Param			req				body		pb.TaskRequest	true	"task"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		500				{string}	string	"bad request"
//	@Router			/api/v1/task/{id} [PUT]
func (t *TaskSrv) UpdateTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) (err error) {
	// 查找该用户的这条信息
	resp.Code = e.SUCCESS
	taskData, err := dao.NewTaskDao(ctx).UpdateTask(req)
	if err != nil {
		resp.Code = e.ERROR
		log.LogrusObj.Error("UpdateTask err:%v", err)
		return
	}
	resp.TaskDetail = BuildTask(taskData)
	return
}

// DeleteTask 删除备忘录
//
//	@Summary		DeleteTask
//	@Description	DeleteTaskDescription
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"token"
//	@Param			id				path		string	true	"id"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		500				{string}	string	"bad request"
//	@Router			/api/v1/task/{id} [DELETE]
func (t *TaskSrv) DeleteTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) (err error) {
	resp.Code = e.SUCCESS
	err = dao.NewTaskDao(ctx).DeleteTaskByIdAndUserId(req.Id, req.Uid)
	if err != nil {
		resp.Code = e.ERROR
		log.LogrusObj.Error("DeleteTask err:%v", err)
		return
	}
	return
}

func BuildTask(item *model.Task) *pb.TaskModel {
	taskModel := pb.TaskModel{
		Id:         uint64(item.ID),
		Uid:        uint64(item.Uid),
		Title:      item.Title,
		Content:    item.Content,
		StartTime:  item.StartTime,
		EndTime:    item.EndTime,
		Status:     int64(item.Status),
		CreateTime: item.CreatedAt.Unix(),
		UpdateTime: item.UpdatedAt.Unix(),
	}
	return &taskModel
}

package dao

import (
	"micro-todoList-k8s/app/user/repository/db/model"
)

// 将结构体映射到数据库，就不用自己注册表了
func migration() {
	_db.Set(`gorm:table_options`, "charset=utf8mb4").
		AutoMigrate(&model.User{})
}

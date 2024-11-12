package service

import (
	"context"
	"errors"
	"sync"

	"gorm.io/gorm"
	"micro-todoList-k8s/app/user/metrics"
	"micro-todoList-k8s/app/user/repository/db/dao"
	"micro-todoList-k8s/app/user/repository/db/model"
	"micro-todoList-k8s/idl/pb"
	"micro-todoList-k8s/pkg/e"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
}

// GetUserSrv 懒汉式单例模式 lazy-loading --> 懒汉式:携程进入，只执行一次
func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

// GetUserSrvHungry 饿汉式式单例模式 --> 饿汉式:携程进入，没有直接生成，资源浪费
func GetUserSrvHungry() *UserSrv {
	if UserSrvIns == nil {
		UserSrvIns = new(UserSrv)
	}
	return UserSrvIns
}

// UserLogin
//
//	@Summary		UserLogin
//	@Description	UserLoginDescription
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			req	body		pb.UserRequest	true	"pb.UserRequest"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{string}	string	"bad request"
//	@Router			/api/v1/user/login [post]
func (u *UserSrv) UserLogin(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	resp.Code = e.SUCCESS
	user, err := dao.NewUserDao(ctx).FindUserByUserName(req.UserName)
	if err != nil {
		resp.Code = e.ERROR
		return
	}

	if user.ID == 0 {
		err = errors.New("用户不存在")
		resp.Code = e.ERROR
		return
	}

	if !user.CheckPassword(req.Password) {
		err = errors.New("用户密码错误")
		resp.Code = e.InvalidParams
		return
	}

	resp.UserDetail = BuildUser(user)
	metrics.QueryUserLoginCounter.WithLabelValues("counts").Inc()

	return
}

// UserRegister
//
//	@Summary		UserRegister
//	@Description	UserRegisterDescription
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			req	body		pb.UserRequest	true	"pb.UserRequest"
//	@Success		200	{object}	map[string]interface{}
//	@Failure		500	{string}	string	"bad request"
//	@Router			/api/v1/user/register [post]
func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRequest, resp *pb.UserDetailResponse) (err error) {
	if req.Password != req.PasswordConfirm {
		err = errors.New("两次密码输入不一致")
		resp.Code = e.ERROR
		return
	}
	resp.Code = e.SUCCESS
	_, err = dao.NewUserDao(ctx).FindUserByUserName(req.UserName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果不存在就继续下去
			// ...continue
		} else {
			resp.Code = e.ERROR
			return
		}
	}
	user := &model.User{
		UserName: req.UserName,
	}
	// 加密密码
	if err = user.SetPassword(req.Password); err != nil {
		resp.Code = e.ERROR
		return
	}
	if err = dao.NewUserDao(ctx).CreateUser(user); err != nil {
		resp.Code = e.ERROR
		return
	}

	resp.UserDetail = BuildUser(user)
	return
}

func BuildUser(item *model.User) *pb.UserModel {
	userModel := pb.UserModel{
		Id:        uint32(item.ID),
		UserName:  item.UserName,
		CreatedAt: item.CreatedAt.Unix(),
		UpdatedAt: item.UpdatedAt.Unix(),
	}
	return &userModel
}

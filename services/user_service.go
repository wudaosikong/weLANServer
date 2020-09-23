package services

import (
	"fmt"
	"weLANServer/models"
	"xorm.io/xorm"
)

/**
 * 用户服务
 * 标准的开发模式将每个实体的提供的功能以接口标准的形式定义,供控制层进行调用。
 *
 */
type UserService interface {
	//通过用户用户名+密码 获取用户实体 如果查询到，返回用户实体，并返回true
	//否则 返回 nil ，false
	GetByUserNameAndPassword(username, password string) (models.User, bool)
	GetByUserName(username string) bool
	//获取用户总数
	GetUserCount() (int64, error)
	AddUser(userRegister *models.User)bool
}

func NewUserService(db *xorm.Engine) UserService {
	return &userSevice{
		engine: db,
	}
}

/**
 * 用户的服务实现结构体
 */
type userSevice struct {
	engine *xorm.Engine
}

/**
 * 查询用户总数
 */
func (ac *userSevice) GetUserCount() (int64, error) {
	count, err := ac.engine.Count(new(models.User))

	if err != nil {
		panic(err.Error())
		return 0, err
	}
	return count, nil
}

/**
 * 通过用户名和密码查询用户
 */
func (ac *userSevice) GetByUserNameAndPassword(username, password string) (models.User, bool) {
	var User models.User

	ac.engine.Where(" User_name = ? and pwd = ? ", username, password).Get(&User)

	return User, User.UserId != 0
}

func (ac *userSevice)GetByUserName(username string) bool{
	var User models.User

	ac.engine.Where(" User_name = ? ", username).Get(&User)

	return User.UserId != 0
}

func (ac *userSevice)AddUser(userRegister *models.User)bool{
	rowNum, err := ac.engine.Insert(&userRegister)
	if err!=nil{
		panic(err.Error())
		return false
	}
	fmt.Println(rowNum) //rowNum 受影响的记录条数
	fmt.Println()
	return true
}
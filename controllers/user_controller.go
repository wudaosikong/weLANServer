package controllers

import (
	"encoding/json"
	"fmt"
	"weLANServer/models"
	"weLANServer/services"
	"weLANServer/utils"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

/**
 * 用户控制器
 */
type UserController struct {
	//iris/v12框架自动为每个请求都绑定上下文对象
	Ctx iris.Context

	//User功能实体
	Service services.UserService

	//session对象
	Session *sessions.Session
}

const (
	UserTABLENAME = "user"
	USER          = "user"
)

type UserLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

/**
 * 用户退出功能
 * 请求类型：Get
 * 请求url：User/singout
 */
func (ac *UserController) GetSingout() mvc.Result {

	//删除session，下次需要从新登录
	ac.Session.Delete(USER)
	return mvc.Response{
		Object: map[string]interface{}{
			"status":  utils.RECODE_OK,
			"success": utils.Recode2Text(utils.RESPMSG_SIGNOUT),
		},
	}
}

/**
 * 处理获取用户总数的路由请求
 * 请求类型：Get
 * 请求Url：User/count
 */
func (ac *UserController) GetCount() mvc.Result {

	count, err := ac.Service.GetUserCount()
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_FAIL,
				"message": utils.Recode2Text(utils.RESPMSG_ERRORUSERCOUNT),
				"count":   0,
			},
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RECODE_OK,
			"count":  count,
		},
	}
}

/**
 * 获取用户信息接口
 * 请求类型：Get
 * 请求url：/User/info
 */
func (ac *UserController) GetInfo() mvc.Result {

	//从session中获取信息
	userByte := ac.Session.Get(USER)

	//session为空
	if userByte == nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_UNLOGIN,
				"type":    utils.EEROR_UNLOGIN,
				"message": utils.Recode2Text(utils.EEROR_UNLOGIN),
			},
		}
	}

	//解析数据到User数据结构
	var user models.User
	err := json.Unmarshal(userByte.([]byte), &user)

	//解析失败
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"status":  utils.RECODE_UNLOGIN,
				"type":    utils.EEROR_UNLOGIN,
				"message": utils.Recode2Text(utils.EEROR_UNLOGIN),
			},
		}
	}

	//解析成功
	return mvc.Response{
		Object: map[string]interface{}{
			"status": utils.RECODE_OK,
			"data":   user.UserToRespDesc(),
		},
	}
}

/**
 * 用户登录功能
 * 接口：/User/login
 */
func (ac *UserController) PostLogin(context iris.Context) mvc.Result {

	iris.New().Logger().Info(" user login ")

	//var userLogin UserLogin
	//ac.Ctx.ReadJSON(&userLogin)
	var userLogin = &UserLogin{context.FormValue("user_name"), context.FormValue("password")}

	//数据参数检验
	if userLogin.UserName == "" || userLogin.Password == "" {
		return mvc.Response{
			//Object: map[string]interface{}{
			//	"status":  "0",
			//	"success": "登录失败",
			//	"message": "用户名或密码为空,请重新填写后尝试登录",
			//},
			Text: "用户名或密码为空,请重新填写后尝试登录",
		}
	}

	//根据用户名、密码到数据库中查询对应的管理信息
	user, exist := ac.Service.GetByUserNameAndPassword(userLogin.UserName, userLogin.Password)

	//用户不存在
	if !exist {
		return mvc.Response{
			//Object: map[string]interface{}{
			//	"status":  "0",
			//	"success": "登录失败",
			//	"message": "用户名或者密码错误,请重新登录",
			//},
			Text: "用户名或者密码错误,请重新登录",
		}
	}

	//用户存在 设置session
	userByte, _ := json.Marshal(user)
	ac.Session.Set(USER, userByte)

	return mvc.Response{
		//Object: map[string]interface{}{
		//	"status":  "1",
		//	"success": "登录成功",
		//	"message": "用户登录成功",
		//},
		Text: "登录成功",
	}

}

//用户注册页面配置、渲染
func (ac *UserController) GetRegister() mvc.View {
	//用户注册模板配置
	var registerView = mvc.View{
		//文件名,视图文件必须放在views文件夹下,因为这是app := iris.Default()默认的
		//当然你也可以自己设置存放位置
		Name: "user/register.html",
		//传入的数据
		Data: map[string]interface{}{},
	}
	return registerView
}

func (ac *UserController) PostRegisterSend(context iris.Context) mvc.View {
	var registerResult string

	userRegister := &models.User{UserName: context.FormValue("user_name"), Pwd: context.FormValue("password"), MyName: context.FormValue("name")}
	fmt.Println(userRegister.UserName)
	exist := ac.Service.GetByUserName(userRegister.UserName)
	if exist {
		registerResult = "此用户名已注册，请注册其他用户名！"
	} else if ac.Service.AddUser(userRegister) {
		registerResult = "用户注册成功！"
	}
	//用户注册结果模板配置
	var registerResultView = mvc.View{
		//文件名,视图文件必须放在views文件夹下,因为这是app := iris.Default()默认的
		//当然你也可以自己设置存放位置
		Name: "user/registerResult.html",
		//传入的数据
		Data: map[string]interface{}{
			"MyMessage": registerResult,
		},
	}
	return registerResultView
}

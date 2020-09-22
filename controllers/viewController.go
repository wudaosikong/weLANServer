package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type ViewController struct {
	Ctx iris.Context
}

//注意 这里的方法名不能随意取名
//名字的开头是url的method,后面以大写字母为分隔路径
//如果路径感觉过长可以继续再给路由组v1分组,用多个controller接收
//此处url：http://localhost:8080/v1/views/hello,method：GET
func (v *ViewController) Get() mvc.View {

	helloView := mvc.View{
		//文件名,视图文件必须放在views文件夹下,因为这是app := iris.Default()默认的
		//当然你也可以自己设置存放位置
		Name: "hello.html",
		//传入的数据
		Data: iris.Map{"content": "街角魔族是最好看的动漫"},
	}

	return helloView
}
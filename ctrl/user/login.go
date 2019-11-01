package user

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-admin/conf"
	"go-admin/models"
	"go-admin/modules/response"
	"go-admin/public/common"
	"strconv"
	"time"
)
type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password   string `form:"password" json:"password" binding:"required"`
}
func Login(c *gin.Context) {
	var u User
	err :=c.BindJSON(&u)
	if err!=nil	{
		response.ShowError(c, "fail")
		return
	}
	if u.Username == "" || u.Password == "" {
		response.ShowError(c, "fail")
		return
	}
	user := models.SystemUser{Nickname: u.Username}
	has := user.GetRowByNickname()
	if !has {
		response.ShowError(c, "fail")
		return
	}
	if common.Sha1En(u.Password+user.Salt) != user.Password {
		response.ShowError(c, "fail")
		return
	}
	session := sessions.Default(c)
	var data = make(map[string]interface{}, 0)

	v := session.Get(conf.Cfg.Token)
	if v == nil {
		cur := time.Now()
		//纳秒
		timestamps := cur.UnixNano()
		times := strconv.FormatInt(timestamps, 10)
		v = common.Md5En(common.GetRandomString(16) + times)
		session.Set(conf.Cfg.Token, v)
		session.Set(v, user.Id)
		err=session.Save()
		fmt.Println("设置成功")
	}
	data[conf.Cfg.Token] = v
	response.ShowData(c, data)
	return
}
func Loginout (c *gin.Context){
	session := sessions.Default(c)
	session.Clear()
	//清除session
	_=session.Save()
	response.ShowSuccess(c, "success")
	return

}
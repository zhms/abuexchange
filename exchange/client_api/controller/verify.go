package controller

import (
	"fmt"
	"math/rand"
	"xserver/abugo"
	"xserver/server"

	"github.com/spf13/viper"
)

type VerifyController struct {
	debug bool
}

func (c *VerifyController) Init() {
	c.debug = viper.GetBool("server.debug")
	gropu := server.Http().NewGroup("/verify")
	{
		gropu.PostNoAuth("/send",c.send)
	}
}
////////////////////////////////////////////////////////////////////////
//发送验证码
///////////////////////////////////////////////////////////////////////
type send_request struct{
	Account string `binding:"required"` //账号
	SellerId int `binding:"required"` //运营商
	UseType int `binding:"required"` //用途
}
func (c *VerifyController) send(ctx *abugo.AbuHttpContent) {
	reqdata := send_request{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(1,err.Error())
		return
	}
	VerifyCode := fmt.Sprint(rand.Intn(999999-100000)+100000)
	sql := "replace into ex_verify(Account,SellerId,UseType,VerifyCode)values(?,?,?,?)"
	server.Db().Conn().Query(sql,reqdata.Account,reqdata.SellerId,reqdata.UseType,VerifyCode)
	if c.debug {
		ctx.Put("VerifyCode",VerifyCode)
	}
	ctx.RespOK()
}
////////////////////////////////////////////////////////////////////////
package controller

import (
	"fmt"
	"math/rand"
	"xserver/abugo"
	"xserver/server"
)

type VerifyController struct {
}

func (c *VerifyController) Init() {
	gropu := server.Http().NewGroup("/verify")
	{
		gropu.PostNoAuth("/send", c.send)
	}
}

////////////////////////////////////////////////////////////////////////
//发送验证码
///////////////////////////////////////////////////////////////////////

func (c *VerifyController) send(ctx *abugo.AbuHttpContent) {
	type RequestData struct {
		Account  string `binding:"required"` //账号
		SellerId int    `binding:"required"` //运营商
		UseType  int    `binding:"required"` //用途
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	VerifyCode := fmt.Sprint(rand.Intn(999999-100000) + 100000)
	sql := "replace into ex_verify(Account,SellerId,UseType,VerifyCode)values(?,?,?,?)"
	server.Db().Conn().Query(sql, reqdata.Account, reqdata.SellerId, reqdata.UseType, VerifyCode)
	if server.Debug() {
		ctx.Put("VerifyCode", VerifyCode)
	}
	ctx.RespOK()
}

/////////////////////////////////////////////////////////////////////////

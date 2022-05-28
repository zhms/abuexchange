package controller

import (
	"xserver/abugo"
	"xserver/server"
)

type UserController struct {
}
func (c *UserController) Init() {
	gropu := server.Http().NewGroup("/user")
	{
		gropu.PostNoAuth("/register",c.register)
		gropu.PostNoAuth("/login_password",c.login_password)
		gropu.PostNoAuth("/login_verifycode",c.login_verifycode)
	}
}
////////////////////////////////////////////////////////////////////////
//注册
///////////////////////////////////////////////////////////////////////
type user_register_request struct{
	SellerId int `binding:"required"` //运营商
	Account string `binding:"required"` //账号
	Password string `binding:"required"` //密码
	VerifyCode string `binding:"required"` //验证码
}
type user_register_db_result struct{
	UserId int
}
func (c *UserController) register(ctx *abugo.AbuHttpContent) {
	reqdata := user_register_request{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1,err.Error())
		return
	}
	dbconn := server.Db().Conn()
	queryresult,err := dbconn.Query("call ex_api_user_register(?,?,?,?)",reqdata.Account,reqdata.SellerId,reqdata.Password,reqdata.VerifyCode)
	if err != nil {
		ctx.RespErr(-2,err.Error())
		return
	}
	queryresult.Next()
	dbresult := user_register_db_result{}
	dberr := abugo.GetDbResult(queryresult,&dbresult)
	if dberr != nil{
		ctx.RespErr(dberr.ErrCode,dberr.ErrMsg)
		return
	}
	queryresult.Close()
	ctx.RespOK(dbresult)
}
////////////////////////////////////////////////////////////////////////
//玩家登录(密码)
///////////////////////////////////////////////////////////////////////
type user_login_password_request struct{
	SellerId int `binding:"required"` //运营商
	Account string `binding:"required"` //账号
	Password string `binding:"required"` //密码
}
type user_login_password_db_result struct{
}
func (c *UserController) login_password(ctx *abugo.AbuHttpContent) {
	reqdata := user_login_password_request{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(1,err.Error())
		return
	}
	dbconn := server.Db().Conn()
	queryresult,err := dbconn.Query("call ex_api_user_login_password(?,?,?)",reqdata.Account,reqdata.SellerId,reqdata.Password)
	if err != nil {
		ctx.RespErr(-2,err.Error())
		return
	}
	queryresult.Next()
	dbresult := user_login_password_db_result{}
	dberr := abugo.GetDbResult(queryresult,&dbresult)
	if dberr != nil{
		ctx.RespErr(dberr.ErrCode,dberr.ErrMsg)
		return
	}
	queryresult.Close()
	ctx.RespOK()
}
////////////////////////////////////////////////////////////////////////
//玩家登录(密码+验证码)
///////////////////////////////////////////////////////////////////////
type user_login_verifycode_request struct{
	SellerId int `binding:"required"` //运营商
	Account string `binding:"required"` //账号
	Password string `binding:"required"` //密码
	VerifyCode string `binding:"required"` //机器码
}
type user_login_verifycode_result struct{
	UserId int
	SellerId int
	OldToken string
	NewToken string
}
func (c *UserController) login_verifycode(ctx *abugo.AbuHttpContent) {
	reqdata := user_login_verifycode_request{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(1,err.Error())
		return
	}
	dbconn := server.Db().Conn()
	queryresult,err := dbconn.Query("call ex_api_user_login_verifycode(?,?,?,?)",reqdata.Account,reqdata.SellerId,reqdata.Password,reqdata.VerifyCode)
	if err != nil {
		ctx.RespErr(2,err.Error())
		return
	}
	queryresult.Next()
	dbresult := user_login_verifycode_result{}
	dberr := abugo.GetDbResult(queryresult,&dbresult)
	if dberr != nil{
		ctx.RespErr(dberr.ErrCode,dberr.ErrMsg)
		return
	}
	queryresult.Close()
	tokendata := server.TokenData{}
	tokendata.UserId = dbresult.UserId
	tokendata.SellerId = dbresult.SellerId
	if len(dbresult.OldToken) > 0{
		server.Http().DelToken(dbresult.OldToken)
	}
	server.Http().SetToken(dbresult.NewToken,tokendata)
	ctx.Put("UserId",dbresult.UserId)
	ctx.Put("SellerId",dbresult.SellerId)
	ctx.Put("Token",dbresult.NewToken)
	ctx.RespOK()
}
///////////////////////////////////////////////////////////////////////
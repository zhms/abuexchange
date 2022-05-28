package server

import (
	"encoding/json"
	"fmt"
	"xserver/abugo"

	"github.com/spf13/viper"
)

var http* abugo.AbuHttp
var redis* abugo.AbuRedis
var db *abugo.AbuDb
var admindb *abugo.AbuDb
var websocket * abugo.AbuWebsocket
var debug bool = false

func Init() {
	abugo.Init()
	debug = viper.GetBool("server.debug")
	http =  new(abugo.AbuHttp)
	http.Init("server.http.http.port")
	redis = new(abugo.AbuRedis)
	redis.Init("server.redis")
	db = new(abugo.AbuDb)
	db.Init("server.db")
	admindb = new(abugo.AbuDb)
	admindb.Init("server.admindb")
	{
		gropu := http.NewGroup("/user")
		{
			gropu.PostNoAuth("/login",user_login)
		}
	}
	{
		gropu := http.NewGroup("/role")
		{
			gropu.PostNoAuth("/login",user_login)
		}
	}
}

func Http() *abugo.AbuHttp {
	return http
}

func Redis() *abugo.AbuRedis{
	return redis
}

func Db() *abugo.AbuDb{
	return db
}

func Debug() bool{
	return debug
}

func Run(){
	abugo.Run()
}

type TokenData struct{
	Account string
	SellerId int
	RoleData string
}

func GetToken(ctx *abugo.AbuHttpContent) *TokenData{
	td := TokenData{}
	err:= json.Unmarshal([]byte(ctx.TokenData),&td)
	if err != nil {
		return nil
	}
	return &td
}

type user_login_request struct{
	Account string `binding:"required"`
	Password string `binding:"required"`
	VerifyCode string `binding:"required"`
}

func user_login(ctx *abugo.AbuHttpContent) {
	reqdata := user_login_request{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1,err.Error())
		return
	}
	sqlstr := "select id,SellerId,`Password`,RoleName,Token,State,GoogleSecret from x_user where Account = ?"
	sqlparam := []interface{}{reqdata.Account}
	var password string
	var token string
	var rolename string
	var sellerid int
	var userstate int
	var userid int
	var googlesecret string
	qserr,qsresult := admindb.QueryScan(sqlstr,sqlparam,&userid,&sellerid,&password,&rolename,&token,&userstate,&googlesecret)
	if qserr != nil {
		ctx.RespErr(-2,qserr.Error())
		return
	}
	if !qsresult {
		ctx.RespErr(-3,"账号不存在")
		return
	}
	if userstate != 1 {
		ctx.RespErr(-4,"账号已被禁用")
		return
	}
	if password != reqdata.Password {
		ctx.RespErr(-5,"密码不正确")
		return
	}
	var roledata string
	var rolestate int
	sqlstr = "select RoleData,State from x_role where RoleName = ? and SellerId = ?"
	sqlparam = []interface{}{rolename,sellerid}
	qserr,qsresult = admindb.QueryScan(sqlstr,sqlparam,&roledata,&rolestate)
	if qserr != nil {
		ctx.RespErr(-2,qserr.Error())
		return
	}
	if !qsresult {
		ctx.RespErr(-6,"角色不存在")
		return
	}
	if rolestate != 1 {
		ctx.RespErr(-7,"角色已被禁用")
		return
	}
	fmt.Println(debug)
	if !debug && len(googlesecret) > 0 && !abugo.VerifyGoogleCode(googlesecret,reqdata.VerifyCode) {
		ctx.RespErr(-10,"验证码不正确")
		return
	}
	if len(token) > 0 { http.DelToken(token) }
	tokendata := TokenData{}
	tokendata.Account = reqdata.Account
	tokendata.SellerId = sellerid
	tokendata.RoleData = roledata
	token = abugo.GetUuid()
	http.SetToken(token,tokendata)
	sqlstr = "update x_user set Token = ?,LoginCount = LoginCount + 1,LoginTime = now(),LoginIp = ? where id = ?"
	err = admindb.QueryNoResult(sqlstr,token,ctx.GetIp(),userid)
	if err != nil {
		ctx.RespErr(-2,err.Error())
		return
	}
	sqlstr = "insert into x_login_log(UserId,Account,Token,LoginIp)values(?,?,?,?)"
	err = admindb.QueryNoResult(sqlstr,userid,reqdata.Account,token,ctx.GetIp())
	if err != nil {
		ctx.RespErr(-2,err.Error())
		return
	}
	ctx.Put("UserId",userid)
	ctx.Put("Token",token)
	ctx.RespOK()
}
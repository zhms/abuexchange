package controller

import (
	"fmt"
	"xserver/abugo"
	"xserver/server"

	"github.com/beego/beego/logs"
)

type UserController struct {
}

func (c *UserController) Init() {
	group := server.Http().NewGroup("/user")
	{
		group.Post("/list", user_list)
	}
}

func user_list(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		SellerId int `validate:"required"`
		UserId   int
		Account  string
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := server.GetToken(ctx)
	if !server.Auth2(token, "玩家管理", "账号管理", "查") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-1, "运营商不正确")
		return
	}
	if reqdata.SellerId == -1 {
		reqdata.SellerId = 0
	}
	where := abugo.AbuWhere{}
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	where.AddInt("and", "UserId", reqdata.UserId, 0)
	where.AddString("and", "Account", reqdata.Account, "")
	var total int
	server.Db().QueryScan(where.CountSql(fmt.Sprintf("%suser", server.DbPrefix)), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := server.Db().Conn().Query(where.Sql(fmt.Sprintf("%suser", server.DbPrefix), reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if err != nil {
		logs.Error(err)
		ctx.RespErr(-2, err.Error())
		return
	}
	type ReturnData struct {
		Id           int
		Account      string
		SellerId     int
		UserId       int
		NickName     string
		RegisterTime string
	}
	data := []ReturnData{}
	for dbresult.Next() {
		d := ReturnData{}
		abugo.GetDbResult(dbresult, &d)
		data = append(data, d)
	}
	dbresult.Close()
	ctx.Put("data", data)
	ctx.Put("page", reqdata.Page)
	ctx.Put("pagesize", reqdata.PageSize)
	ctx.Put("total", total)
	ctx.RespOK()
}

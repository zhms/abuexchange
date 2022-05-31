package controller

import (
	"xserver/abugo"
	"xserver/server"
)

type SellerController struct {
}

func (c *SellerController) Init() {
	group := server.Http().NewGroup("/seller")
	{
		group.Post("/list", c.list)
	}
}

type seller_data struct {
	SellerId   int
	SellerName string
}

func (c *SellerController) list(ctx *abugo.AbuHttpContent) {
	token := server.GetToken(ctx)
	if token.SellerId != -1 {
		ctx.Put("data", []interface{}{})
		ctx.RespOK()
		return
	}
	sql := "select * from ex_seller where state = 1"
	dbresult, err := server.Db().Conn().Query(sql)
	if err != nil {
		ctx.RespErr(1, err.Error())
		return
	}
	data := []seller_data{}
	data = append(data, seller_data{0, "全部"})
	data = append(data, seller_data{-1, "总后台"})
	for dbresult.Next() {
		d := seller_data{}
		abugo.GetDbResult(dbresult, &d)
		data = append(data, d)
	}
	dbresult.Close()
	ctx.RespOK(data)
}

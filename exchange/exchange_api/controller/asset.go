package controller

import (
	"xserver/abugo"
	"xserver/server"
)

type AssetController struct {
}

func (c *AssetController) Init() {
	gropu := server.Http().NewGroup("/asset")
	{
		gropu.Post("/list",c.list)
	}
}
////////////////////////////////////////////////////////////////////////
//资产列表
///////////////////////////////////////////////////////////////////////
type asset_data struct{
	Symbol string
	AssetType int
	AssetAmt int64
	FrozenAmt int64
}
func (c *AssetController) list(ctx *abugo.AbuHttpContent) {
	token := server.GetToken(ctx)
	queryresult,err := server.Db().Conn().Query("select Symbol,AssetType,AssetAmt,FrozenAmt from ex_asset where userid = ?",token.UserId)
	if err != nil {
		ctx.RespErr(-2,err.Error())
		return
	}
	assets := []asset_data{}
	for queryresult.Next() {
		dbresult := asset_data{}
		abugo.GetDbResult(queryresult,&dbresult)
		assets = append(assets,dbresult)
	}
	queryresult.Close()
	ctx.RespOK(assets)
}
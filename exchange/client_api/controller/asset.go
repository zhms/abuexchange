package controller

import (
	"xserver/abugo"
	"xserver/server"

	"github.com/spf13/viper"
)

type AssetController struct {
	debug bool
}

func (c *AssetController) Init() {
	c.debug = viper.GetBool("server.debug")
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
	dbconn := server.Db().Conn()
	queryresult,err := dbconn.Query("select Symbol,AssetType,AssetAmt,FrozenAmt from ex_asset where userid = ?",token.UserId)
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
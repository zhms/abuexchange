package controller

import (
	"strings"
	"xserver/abugo"
	"xserver/server"
)

var testpublickey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxkq6/DabE+yGHzdWZnH4
Tmy1iic5BiRmoTY3/qUjXC/keUy9mpHuYALzh5VkjGsVAEq8aFTzLDD6hQZ3nrEZ
I1ClcpAPSnN/W5dfwyfSzHTZOGd91m5WI2u5XRXJ0i5WnPY6LTAOzqdkI+T7iIST
5jcgpgsYecsCjEKc+j4AM4kLPBA7rQFl4Ld6FbSZxjsUQHqsJ7BAWW+utsgRHUsk
ZwiiKhPB4YN76S+YX1wcNbpdyTQz1p8xI4RjxSnLh7U93b8JVGZgHIuV6sa267Hm
o927y024KZQ+2nz2EcOpcHl5kq25BZP8kfnDPOqAiN1HC7haQfdqY9vradiUMGp9
dQIDAQAB
-----END PUBLIC KEY-----`

var testprivatekey = `-----BEGIN PRIVATE KEY-----
MIIEuwIBADALBgkqhkiG9w0BAQEEggSnMIIEowIBAAKCAQEAxkq6/DabE+yGHzdW
ZnH4Tmy1iic5BiRmoTY3/qUjXC/keUy9mpHuYALzh5VkjGsVAEq8aFTzLDD6hQZ3
nrEZI1ClcpAPSnN/W5dfwyfSzHTZOGd91m5WI2u5XRXJ0i5WnPY6LTAOzqdkI+T7
iIST5jcgpgsYecsCjEKc+j4AM4kLPBA7rQFl4Ld6FbSZxjsUQHqsJ7BAWW+utsgR
HUskZwiiKhPB4YN76S+YX1wcNbpdyTQz1p8xI4RjxSnLh7U93b8JVGZgHIuV6sa2
67Hmo927y024KZQ+2nz2EcOpcHl5kq25BZP8kfnDPOqAiN1HC7haQfdqY9vradiU
MGp9dQIDAQABAoIBAQCB0xOxZZ8K4GS/s19iG2RTFGGXXi67ksGO8wyeMEiCZujT
cicBOGI1gjknn+wA82sAke3g8R105WxgA+f1gkXPFVH9O8yHFi9gLC+KelxNeINv
pnKBHiRPdIwcIdPNRIF4qs1UII7RRk8OITF+JC8hHx7FF6aitwOoCUT1ofpHwl6+
rIgXDOfS5XjZ/c/NZ0k8zajO9PU/GfWLEUWffHIPvE+CKo5Eel9sltHG2mbWoZwP
Z7vynMWDKzXRDseWHxQ9KdBGZX24J2suXtsJMN1kL4TpMxJaP6jOGmj7A4YUzdNr
74WOfSaYCzKqa7uog+56VRpMNItVah/P9BhSMk4BAoGBAOvdhtos4twe+dqckpkX
ISXR4D4f4IwAP3hbpZLt4bZ7glq+a1w5ukxnPlvS91LHN4YXmRfrDBaKiATDs6Ck
885EDIPA5oHodF6Cckss2i4DjikYPJY0Gt05AlttIvwzDJfv0Nv8NoLpDXbc/6f4
Qq4AzK8EmSb3idws1vFnZjVVAoGBANc4GCwwQtHUj/6aMZ+hMYSjTcQJTebGWgUd
yOxjqfN9zYtfkd2tpyF7LNZTF72QijmCqiXboOWqVNU1IDOvw5f0ZWGod18QkqBa
7emJpexiFtz8qq7NZI2IKa7retgZfDON5EbLWQGKRQlBsKB07EUYjM09U9Aw74mU
8xOa9yehAoGANQ1l6KsR9/sLrg1rt/2I2i3j53VaF4Nyw0qx6mORUm9jea+9DEh7
y/WpBBRmxvrCVWn3aHbZCDTutUujmq2fnTSXMAdykaPMkPl8ZcZX6OOp1Tp7Xjmo
FTxeeFwK83k0CvTJIMIRLM7o1WSOKKThsqvBXliFbktdeeWoDzJ3veUCgYAXoL3u
pAHB/Z8taYpHJzBDipYwDoFx85bQdvunEC4JYGdfhMeD2du+bkJ26TpzAlpahADg
FSOt8yxLGAEIC6Timt9CgHdShYRDfJggCo+fywkNi+PFOyFt9GP3Iz0iWnHrSwmj
yBFmfOKlaIPWqd+CmnZi1ffMpIK3zRV+soWt4QKBgGLTUjM8cOtwS07s/T+rPDLg
3vLMtTwOJCJWlFYEjo1FfUme27GPiw/5u3BiDl4MshkHoZoh3C7hcXVu4LeV+3sC
QKLi49b8UghPw2VsqVDHC1cAxjmlWK1XS0wlYq9GaR/nvinzkfdHhB/+h9BgriLW
sRYCiEzODrdWCOU4m7VM
-----END PRIVATE KEY-----`

type ThirdController struct {
}

func verify_sign() {

}

func (c *ThirdController) Init() {
	group := server.Http().NewGroup("/v1/third")
	{
		group.PostNoAuth("/user_register", third_user_register)
		group.PostNoAuth("/get_balance", third_get_balance)
		group.PostNoAuth("/transfer_in", third_transfer_in)
		group.PostNoAuth("/transfer_out", third_transfer_out)
	}
}

func third_user_register(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UniqueId  string `validate:"required"`
		Password  string `validate:"required"`
		TimeStamp int64  `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	seller := server.GetSeller(reqdata.SellerId)
	if ctx.RespErrString(seller == nil, &errcode, "商户不存在") {
		return
	}
	if ctx.RespErrString(!server.Debug() && !abugo.RsaVerify(reqdata, seller.ApiThirdPublicKey), &errcode, "签名不正确") {
		return
	}
	sql := "call ex_api_third_register(?,?,?,?)"
	dbresult, err := server.Db().Conn().Query(sql, reqdata.UniqueId, reqdata.SellerId, reqdata.Password, "{}")
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		UserId int
	}
	returndata := ReturnData{}
	if dbresult.Next() {
		dberr := abugo.GetDbResult(dbresult, &returndata)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	ctx.RespOK(returndata)
}

func third_get_balance(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UserId    string `validate:"required"`
		Symbol    string `validate:"required"`
		AssetType int    `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	reqdata.Symbol = strings.ToLower(reqdata.Symbol)
	seller := server.GetSeller(reqdata.SellerId)
	if ctx.RespErrString(seller == nil, &errcode, "商户不存在") {
		return
	}
	if ctx.RespErrString(!server.Debug() && !abugo.RsaVerify(reqdata, seller.ApiThirdPublicKey), &errcode, "签名不正确") {
		return
	}
	sql := "select AssetAmt from ex_asset where UserId = ? and AssetType = ? and Symbol = ?"
	dbresult, err := server.Db().Conn().Query(sql, reqdata.UserId, reqdata.AssetType, reqdata.Symbol)
	if ctx.RespErr(err, &errcode) {
		return
	}
	var balance int64
	if dbresult.Next() {
		dbresult.Scan(&balance)
	}
	dbresult.Close()
	ctx.Put("balance", balance)
	ctx.RespOK()
}

func third_transfer_in(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UserId    string `validate:"required"`
		Symbol    string `validate:"required"`
		AssetType int    `validate:"required"`
		OrderId   int64  `validate:"required"`
		Amount    int64  `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	reqdata.Symbol = strings.ToLower(reqdata.Symbol)
	seller := server.GetSeller(reqdata.SellerId)
	if ctx.RespErrString(seller == nil, &errcode, "商户不存在") {
		return
	}
	if ctx.RespErrString(!server.Debug() && !abugo.RsaVerify(reqdata, seller.ApiThirdPublicKey), &errcode, "签名不正确") {
		return
	}
	sql := "call ex_api_transfer_in_in(?,?,?,?,?,?,?,?,?)"
	dbresult, err := server.Db().Conn().Query(sql, reqdata.OrderId, reqdata.UserId, reqdata.SellerId, reqdata.AssetType, reqdata.Symbol, reqdata.Amount, "{}", 1, "钱包转入")
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Balance int64
	}
	d := ReturnData{}
	if dbresult.Next() {
		dberr := abugo.GetDbResult(dbresult, &d)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	ctx.Put("balance", d.Balance)
	ctx.RespOK()
}

func third_transfer_out(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UserId    string `validate:"required"`
		Symbol    string `validate:"required"`
		AssetType int    `validate:"required"`
		OrderId   int64  `validate:"required"`
		Amount    int64  `validate:"required"`
	}
	errcode := 0
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if ctx.RespErr(err, &errcode) {
		return
	}
	reqdata.Symbol = strings.ToLower(reqdata.Symbol)
	seller := server.GetSeller(reqdata.SellerId)
	if ctx.RespErrString(seller == nil, &errcode, "商户不存在") {
		return
	}
	if ctx.RespErrString(!server.Debug() && !abugo.RsaVerify(reqdata, seller.ApiThirdPublicKey), &errcode, "签名不正确") {
		return
	}
	sql := "call ex_api_transfer_in_out(?,?,?,?,?,?,?,?,?)"
	dbresult, err := server.Db().Conn().Query(sql, reqdata.OrderId, reqdata.UserId, reqdata.SellerId, reqdata.AssetType, reqdata.Symbol, reqdata.Amount, "{}", 1, "钱包转入")
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Balance int64
	}
	d := ReturnData{}
	if dbresult.Next() {
		dberr := abugo.GetDbResult(dbresult, &d)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	ctx.Put("balance", d.Balance)
	ctx.RespOK()
}

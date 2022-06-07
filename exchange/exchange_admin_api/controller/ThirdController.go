package controller

import (
	"fmt"
	"strings"
	"time"
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

func (c *ThirdController) Init() {
	group := server.Http().NewGroup("/v1/third")
	{
		group.PostNoAuth("/user_register", third_user_register)
		group.PostNoAuth("/get_balance", third_get_balance)
		group.PostNoAuth("/transfer_in", third_transfer_in)
		group.PostNoAuth("/transfer_out", third_transfer_out)
		group.PostNoAuth("/server_login", third_server_login)
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
	sql := fmt.Sprintf("call %sapi_third_register(?,?,?,?)", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.UniqueId, reqdata.SellerId, reqdata.Password, "{}")
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		UserId    int
		Timestamp int64
		Sign      string
	}
	retdata := ReturnData{}
	if dbresult.Next() {
		dberr := abugo.GetDbResult(dbresult, &retdata)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	retdata.Timestamp = time.Now().Unix()
	retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	ctx.RespOK(retdata)
}

func third_get_balance(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UserId    string `validate:"required"`
		Symbol    string `validate:"required"`
		AssetType int    `validate:"required"`
		TimeStamp int64  `validate:"required"`
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
	sql := fmt.Sprintf("select AssetAmt from %sasset where UserId = ? and AssetType = ? and Symbol = ?", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.UserId, reqdata.AssetType, reqdata.Symbol)
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Balance   int64
		Timestamp int64
		Sign      string
	}
	retdata := ReturnData{}
	if dbresult.Next() {
		dbresult.Scan(&retdata.Balance)
	}
	dbresult.Close()
	retdata.Timestamp = time.Now().Unix()
	retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	ctx.RespOK(retdata)
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
		TimeStamp int64  `validate:"required"`
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
	sql := fmt.Sprintf("call %sapi_transfer_in_in(?,?,?,?,?,?,?,?,?)", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.OrderId, reqdata.UserId, reqdata.SellerId, reqdata.AssetType, reqdata.Symbol, reqdata.Amount, "{}", 1, "钱包转入")
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Balance   int64
		Timestamp int64
		Sign      string
	}
	retdata := ReturnData{}
	if dbresult.Next() {
		dberr := abugo.GetDbResult(dbresult, &retdata)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	retdata.Timestamp = time.Now().Unix()
	retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	ctx.Put("balance", retdata.Balance)
	ctx.RespOK(retdata)
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
		TimeStamp int64  `validate:"required"`
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
	sql := fmt.Sprintf("call %sapi_transfer_in_out(?,?,?,?,?,?,?,?,?)", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.OrderId, reqdata.UserId, reqdata.SellerId, reqdata.AssetType, reqdata.Symbol, reqdata.Amount, "{}", 1, "钱包转入")
	if ctx.RespErr(err, &errcode) {
		return
	}
	type ReturnData struct {
		Balance   int64
		Sign      string
		TimeStamp int64
	}
	retdata := ReturnData{}
	if dbresult.Next() {
		dberr := abugo.GetDbResult(dbresult, &retdata)
		if ctx.RespDbErr(dberr) {
			return
		}
	}
	retdata.TimeStamp = time.Now().Unix()
	retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	ctx.RespOK(retdata)
}

type third_server_token struct {
	UserId   int
	SellerId int
}

func third_server_login(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Sign      string `validate:"required"`
		SellerId  int    `validate:"required"`
		UserId    int    `validate:"required"`
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
	sql := fmt.Sprintf("select Password from %suser where UserId = ? and SellerId = ?", server.DbPrefix)
	dbresult, err := server.Db().Conn().Query(sql, reqdata.UserId, reqdata.SellerId)
	if ctx.RespErr(err, &errcode) {
		return
	}
	if dbresult.Next() {
		var password string
		dbresult.Scan(&password)
		if reqdata.Password != password {
			ctx.RespErrString(true, &errcode, "密码不正确")
			return
		}
	} else {
		ctx.RespErrString(true, &errcode, "账号不存在")
		return
	}
	type ReturnData struct {
		Token     string
		Sign      string
		TimeStamp int64
	}
	retdata := ReturnData{}
	retdata.Token = abugo.GetUuid()
	tokendata := third_server_token{}
	tokendata.UserId = reqdata.UserId
	tokendata.SellerId = reqdata.SellerId
	server.Redis().Set(fmt.Sprint("exchange:third:server:login:", retdata.Token), tokendata)
	retdata.TimeStamp = time.Now().Unix()
	retdata.Sign = abugo.RsaSign(retdata, seller.ApiPrivateKey)
	ctx.RespOK(retdata)
}

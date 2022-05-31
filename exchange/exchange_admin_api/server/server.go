package server

import (
	"encoding/json"
	"reflect"
	"strings"
	"xserver/abugo"

	"github.com/beego/beego/logs"
	"github.com/spf13/viper"
	"github.com/valyala/fastjson"
)

var http *abugo.AbuHttp
var redis *abugo.AbuRedis
var db *abugo.AbuDb
var admindb *abugo.AbuDb
var websocket *abugo.AbuWebsocket
var debug bool = false

func Init() {
	abugo.Init()
	debug = viper.GetBool("server.debug")
	http = new(abugo.AbuHttp)
	http.Init("server.http.http.port")
	redis = new(abugo.AbuRedis)
	redis.Init("server.redis")
	db = new(abugo.AbuDb)
	db.Init("server.db")
	admindb = new(abugo.AbuDb)
	admindb.Init("server.admindb")
	http.PostNoAuth("/user/login", user_login)
	http.Post("/system/login_log", login_log)
	http.Post("/role/list", role_list)
	http.Post("/role/listall", role_listall)
	http.Post("/role/roledata", role_data)
	sql := "select RoleData from x_role where SellerId = -1 and RoleName = '超级管理员'"
	var dbauthdata string
	admindb.QueryScan(sql, []interface{}{}, &dbauthdata)
	if dbauthdata != AuthData {
		sql = "select Id,SellerId,RoleName,RoleData from x_role"
		dbresult, _ := admindb.Conn().Query(sql)
		for dbresult.Next() {
			var roleid int
			var sellerid int
			var rolename string
			var roledata string
			dbresult.Scan(&roleid, &sellerid, &rolename, &roledata)
			if sellerid == -1 && rolename == "超级管理员" {
				continue
			}
			jnewdata := make(map[string]interface{})
			json.Unmarshal([]byte(AuthData), &jnewdata)
			clean_auth(jnewdata)
			jrdata := make(map[string]interface{})
			json.Unmarshal([]byte(roledata), &jrdata)
			for k, v := range jrdata {
				set_auth(k, jnewdata, v.(map[string]interface{}))
			}
			newauthbyte, _ := json.Marshal(&jnewdata)
			sql = "update x_role set RoleData = ? where id = ?"
			admindb.QueryNoResult(sql, string(newauthbyte), roleid)
		}
		dbresult.Close()
		sql = "update x_role set RoleData = ? where RoleName = '超级管理员'"
		admindb.QueryNoResult(sql, AuthData)
	}
}

func clean_auth(node map[string]interface{}) {
	for k, v := range node {
		if strings.Index(reflect.TypeOf(v).Name(), "float") >= 0 {
			node[k] = 0
		} else {
			clean_auth(v.(map[string]interface{}))
		}
	}
}

func set_auth(parent string, newdata map[string]interface{}, node map[string]interface{}) {
	for k, v := range node {
		if strings.Index(reflect.TypeOf(v).Name(), "float") >= 0 {
			if v.(float64) != 1 {
				continue
			}
			path := strings.Split(parent, ".")
			if len(path) == 0 {
				continue
			}
			fk, fok := newdata[path[0]]
			if !fok {
				continue
			}
			var pn *interface{} = &fk
			var finded bool = true
			for i := 1; i < len(path); i++ {
				tk := path[i]
				tv, ok := (*pn).(map[string]interface{})[tk]
				if !ok {
					finded = false
					break
				}
				pn = &tv
			}
			if finded {
				(*pn).(map[string]interface{})[k] = 1
			}

		} else {
			set_auth(parent+"."+k, newdata, v.(map[string]interface{}))
		}
	}
}

func Http() *abugo.AbuHttp {
	return http
}

func Redis() *abugo.AbuRedis {
	return redis
}

func Db() *abugo.AbuDb {
	return db
}

func Debug() bool {
	return debug
}

func Run() {
	abugo.Run()
}

type TokenData struct {
	Account  string
	SellerId int
	AuthData string
}

func GetToken(ctx *abugo.AbuHttpContent) *TokenData {
	td := TokenData{}
	err := json.Unmarshal([]byte(ctx.TokenData), &td)
	if err != nil {
		return nil
	}
	return &td
}

type user_login_request_struct struct {
	Account    string `binding:"required"`
	Password   string `binding:"required"`
	VerifyCode string `binding:"required"`
}

type menu_data_struct struct {
	Icon  string             `json:"icon"`
	Index string             `json:"index"`
	Title string             `json:"title"`
	Subs  []menu_data_struct `json:"subs"`
}

func user_login(ctx *abugo.AbuHttpContent) {
	reqdata := user_login_request_struct{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	sqlstr := "select id,SellerId,`Password`,RoleName,Token,State,GoogleSecret,LoginTime,LoginCount,NickName from x_user where Account = ?"
	sqlparam := []interface{}{reqdata.Account}
	var password string
	var token string
	var rolename string
	var sellerid int
	var userstate int
	var userid int
	var googlesecret string
	var logintime string
	var logincount int
	var nickname string
	qserr, qsresult := admindb.QueryScan(sqlstr, sqlparam, &userid, &sellerid, &password, &rolename, &token, &userstate, &googlesecret, &logintime, &logincount, &nickname)
	if qserr != nil {
		ctx.RespErr(-2, qserr.Error())
		return
	}
	if !qsresult {
		ctx.RespErr(-3, "账号不存在")
		return
	}
	if userstate != 1 {
		ctx.RespErr(-4, "账号已被禁用")
		return
	}
	if password != reqdata.Password {
		ctx.RespErr(-5, "密码不正确")
		return
	}
	var authdata string
	var rolestate int
	sqlstr = "select RoleData,State from x_role where RoleName = ? and SellerId = ?"
	sqlparam = []interface{}{rolename, sellerid}
	qserr, qsresult = admindb.QueryScan(sqlstr, sqlparam, &authdata, &rolestate)
	if qserr != nil {
		ctx.RespErr(-2, qserr.Error())
		return
	}
	if !qsresult {
		ctx.RespErr(-6, "角色不存在")
		return
	}
	if rolestate != 1 {
		ctx.RespErr(-7, "角色已被禁用")
		return
	}
	if !debug && len(googlesecret) > 0 && !abugo.VerifyGoogleCode(googlesecret, reqdata.VerifyCode) {
		ctx.RespErr(-10, "验证码不正确")
		return
	}
	if rolename == "超级管理员" {
		authdata = AuthData
	}
	if len(token) > 0 {
		http.DelToken(token)
	}
	tokendata := TokenData{}
	tokendata.Account = reqdata.Account
	tokendata.SellerId = sellerid
	tokendata.AuthData = authdata
	token = abugo.GetUuid()
	http.SetToken(token, tokendata)
	sqlstr = "update x_user set Token = ?,LoginCount = LoginCount + 1,LoginTime = now(),LoginIp = ? where id = ?"
	err = admindb.QueryNoResult(sqlstr, token, ctx.GetIp(), userid)
	if err != nil {
		ctx.RespErr(-2, err.Error())
		return
	}
	sqlstr = "insert into x_login_log(UserId,SellerId,Account,Token,LoginIp)values(?,?,?,?,?)"
	err = admindb.QueryNoResult(sqlstr, userid, sellerid, reqdata.Account, token, ctx.GetIp())
	if err != nil {
		ctx.RespErr(-2, err.Error())
		return
	}
	menu := []menu_data_struct{}
	json.Unmarshal([]byte(MenuData), &menu)
	parser := fastjson.Parser{}
	jauthdata, _ := parser.ParseBytes([]byte(authdata))
	//三级菜单
	for i := 0; i < len(menu); i++ {
		for j := 0; j < len(menu[i].Subs); j++ {
			smenu := []menu_data_struct{}
			for k := 0; k < len(menu[i].Subs[j].Subs); k++ {
				open := jauthdata.GetInt(menu[i].Title, menu[i].Subs[j].Title, menu[i].Subs[j].Subs[k].Title, "查")
				if open == 1 {
					smenu = append(smenu, menu[i].Subs[j].Subs[k])
				}
			}
			menu[i].Subs[j].Subs = smenu
		}
	}
	//二级菜单
	for i := 0; i < len(menu); i++ {
		smenu := []menu_data_struct{}
		for j := 0; j < len(menu[i].Subs); j++ {
			open := jauthdata.GetInt(menu[i].Title, menu[i].Subs[j].Title, "查")
			if open == 1 || len(menu[i].Subs[j].Subs) > 0 {
				smenu = append(smenu, menu[i].Subs[j])
			}
		}
		menu[i].Subs = smenu
	}
	//一级菜单
	smenu := []menu_data_struct{}
	for i := 0; i < len(menu); i++ {
		open := jauthdata.GetInt(menu[i].Title, "查")
		if open == 1 || len(menu[i].Subs) > 0 {
			smenu = append(smenu, menu[i])
		}
	}
	menu = smenu
	jauth := make(map[string]interface{})
	json.Unmarshal([]byte(authdata), &jauth)

	ctx.Put("UserId", userid)
	ctx.Put("SellerId", sellerid)
	ctx.Put("Account", reqdata.Account)
	ctx.Put("NickName", nickname)
	ctx.Put("AuthData", jauth)
	ctx.Put("MenuData", menu)
	ctx.Put("Token", token)
	ctx.Put("LoginTime", logintime)
	ctx.Put("Ip", ctx.GetIp())
	ctx.Put("LoginCount", logincount)
	ctx.Put("Version", "1.0.0")
	ctx.RespOK()
}

type login_log_request_struct struct {
	Page     int
	PageSize int
	Account  string
	SellerId int
}

type login_log_data_struct struct {
	Id         int
	UserId     int
	SellerId   int
	Account    string
	LoginIp    string
	CreateTime string
}

func login_log(ctx *abugo.AbuHttpContent) {
	reqdata := login_log_request_struct{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 10
	}
	token := GetToken(ctx)
	if token.SellerId > 0 {
		reqdata.SellerId = token.SellerId
	}
	where := abugo.AbuWhere{}
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	where.AddString("and", "Account", reqdata.Account, "")
	var total int
	admindb.QueryScan(where.CountSql("x_login_log"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := admindb.Conn().Query(where.Sql("x_login_log", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if err != nil {
		logs.Error(err)
		ctx.RespErr(-2, err.Error())
		return
	}
	data := []login_log_data_struct{}
	for dbresult.Next() {
		d := login_log_data_struct{}
		abugo.GetDbResult(dbresult, &d)
		d.CreateTime = abugo.TimeToUtc(d.CreateTime)
		data = append(data, d)
	}
	dbresult.Close()
	ctx.Put("data", data)
	ctx.Put("page", reqdata.Page)
	ctx.Put("pagesize", reqdata.PageSize)
	ctx.Put("total", total)
	ctx.RespOK()
}

type role_request_struct struct {
	Page     int
	PageSize int
	SellerId int
}

type role_data_struct struct {
	Id       int
	RoleName string
	SellerId int
	State    int
	Parent   string
	RoleData string
}

func role_list(ctx *abugo.AbuHttpContent) {
	reqdata := role_request_struct{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	if reqdata.Page == 0 {
		reqdata.Page = 1
	}
	if reqdata.PageSize == 0 {
		reqdata.PageSize = 10
	}
	token := GetToken(ctx)
	if token.SellerId > 0 {
		reqdata.SellerId = token.SellerId
	}
	where := abugo.AbuWhere{}
	where.OrderBy = "ASC"
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	var total int
	admindb.QueryScan(where.CountSql("x_role"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := admindb.Conn().Query(where.Sql("x_role", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if err != nil {
		logs.Error(err)
		ctx.RespErr(-2, err.Error())
		return
	}
	data := []role_data_struct{}
	for dbresult.Next() {
		d := role_data_struct{}
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

func role_listall(ctx *abugo.AbuHttpContent) {
	reqdata := role_request_struct{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if token.SellerId > 0 {
		reqdata.SellerId = token.SellerId
	}
	sql := "select RoleName from x_role where SellerId = ? and State = 1"
	dbresult, err := admindb.Conn().Query(sql, reqdata.SellerId)
	if err != nil {
		ctx.RespErr(1, err.Error())
		return
	}
	names := []string{}
	for dbresult.Next() {
		var RoleName string
		dbresult.Scan(&RoleName)
		names = append(names, RoleName)
	}
	dbresult.Close()
	ctx.RespOK(names)
}

type role_data_request_struct struct {
	SellerId int    `binding:"required"`
	RoleName string `binding:"required"`
}

func role_data(ctx *abugo.AbuHttpContent) {
	reqdata := role_data_request_struct{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if token.SellerId > 0 {
		reqdata.SellerId = token.SellerId
	}
	sql := "select RoleData from x_role where SellerId = ? and State = 1 and RoleName = ?"
	var RoleData string
	admindb.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &RoleData)
	var SuperRoleData string
	sql = "select RoleData from x_role where SellerId = -1 and State = 1 and RoleName = '超级管理员'"
	admindb.QueryScan(sql, []interface{}{}, &SuperRoleData)
	ctx.Put("RoleData", RoleData)
	ctx.Put("SuperRoleData", SuperRoleData)
	ctx.RespOK()
}

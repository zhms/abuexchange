package server

import (
	"encoding/json"
	"fmt"
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
var websocket *abugo.AbuWebsocket
var debug bool = false
var db_seller_tablename string = "ex_seller"

type SellerData struct {
	SellerId   int
	SellerName string
	State      int
	Remark     string
	CreateTime string
}

func seller_list(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
	}
	reqdata := RequestData{}
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
	if token.SellerId != -1 {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if !Auth2(token, "系统管理", "运营商管理", "查") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	where := abugo.AbuWhere{}
	where.OrderBy = "ASC"
	where.OrderKey = "SellerId"
	var total int
	db.QueryScan(where.CountSql(db_seller_tablename), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql(db_seller_tablename, reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if err != nil {
		logs.Error(err)
		ctx.RespErr(-2, err.Error())
		return
	}
	data := []SellerData{}
	for dbresult.Next() {
		d := SellerData{}
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

func seller_add(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerName string `validate:"required"`
		State      int    `validate:"required"`
		Remark     string
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if token.SellerId != -1 {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if !Auth2(token, "系统管理", "运营商管理", "增") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	sql := fmt.Sprintf("insert into %s(SellerName,State,Remark)values(?,?,?)", db_seller_tablename)
	db.QueryNoResult(sql, reqdata.SellerName, reqdata.State, reqdata.Remark)
	WriteAdminLog("添加运营商", ctx, reqdata)
	ctx.RespOK()
}

func seller_modify(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId   int    `validate:"required"`
		SellerName string `validate:"required"`
		State      int    `validate:"required"`
		Remark     string
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if token.SellerId != -1 {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if !Auth2(token, "系统管理", "运营商管理", "改") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	sql := fmt.Sprintf("update %s set SellerName = ?,State = ?,Remark = ? where SellerId = ?", db_seller_tablename)
	db.QueryNoResult(sql, reqdata.SellerName, reqdata.State, reqdata.Remark, reqdata.SellerId)
	WriteAdminLog("修改运营商", ctx, reqdata)
	ctx.RespOK()
}

func seller_delete(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId int `validate:"required"`
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if token.SellerId != -1 {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if !Auth2(token, "系统管理", "运营商管理", "增") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	sql := fmt.Sprintf("delete from  %s where SellerId = ?", db_seller_tablename)
	db.QueryNoResult(sql, reqdata.SellerId)
	WriteAdminLog("删除运营商", ctx, reqdata)
	ctx.RespOK()
}

func Init() {
	abugo.Init()
	debug = viper.GetBool("server.debug")
	http = new(abugo.AbuHttp)
	http.Init("server.http.http.port")
	redis = new(abugo.AbuRedis)
	redis.Init("server.redis")
	db = new(abugo.AbuDb)
	db.Init("server.db")
	SetupDatabase()
	{
		http.PostNoAuth("/admin/user/login", user_login)
		http.Post("/admin/login_log", login_log)
		http.Post("/admin/role/list", role_list)
		http.Post("/admin/role/listall", role_listall)
		http.Post("/admin/role/roledata", role_data)
		http.Post("/admin/role/modify", role_modify)
		http.Post("/admin/role/add", role_add)
		http.Post("/admin/role/delete", role_delete)
		http.Post("/admin/opt_log", opt_log)
		http.Post("/admin/user/list", user_list)
		http.Post("/admin/user/modify", user_modify)
		http.Post("/admin/user/delete", user_delete)
		http.Post("/admin/user/add", user_add)
		http.Post("/admin/user/google", user_google)
		http.Post("seller/name", seller_name)
		http.Post("seller/list", seller_list)
		http.Post("seller/add", seller_add)
		http.Post("seller/delete", seller_delete)
		http.Post("seller/modify", seller_modify)
	}
	sql := "select RoleData from z_admin_role where SellerId = -1 and RoleName = '超级管理员'"
	var dbauthdata string
	db.QueryScan(sql, []interface{}{}, &dbauthdata)
	if dbauthdata != AuthDataStr {
		sql = "select Id,SellerId,RoleName,RoleData from z_admin_role"
		dbresult, _ := db.Conn().Query(sql)
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
			json.Unmarshal([]byte(AuthDataStr), &jnewdata)
			clean_auth(jnewdata)
			jrdata := make(map[string]interface{})
			json.Unmarshal([]byte(roledata), &jrdata)
			for k, v := range jrdata {
				set_auth(k, jnewdata, v.(map[string]interface{}))
			}
			newauthbyte, _ := json.Marshal(&jnewdata)
			sql = "update z_admin_role set RoleData = ? where id = ?"
			db.QueryNoResult(sql, string(newauthbyte), roleid)
		}
		dbresult.Close()
		sql = "update z_admin_role set RoleData = ? where RoleName = '超级管理员'"
		db.QueryNoResult(sql, AuthDataStr)
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
func WriteAdminLog(opt string, ctx *abugo.AbuHttpContent, data interface{}) {
	token := ctx.Token
	strdata, _ := json.Marshal(&data)
	tokendata := GetToken(ctx)
	Ip := ctx.GetIp()
	go func() {
		sql := "insert into z_admin_opt_log(Account,Opt,Token,Data,Ip)values(?,?,?,?,?)"
		db.QueryNoResult(sql, tokendata.Account, opt, token, string(strdata), Ip)
	}()
}
func Auth2(td *TokenData, m string, s string, o string) bool {
	defer recover()
	authdata := make(map[string]interface{})
	json.Unmarshal([]byte(td.AuthData), &authdata)
	im, imok := authdata[m]
	if !imok {
		return false
	}
	is, isok := im.(map[string]interface{})[s]
	if !isok {
		return false
	}
	io, iook := is.(map[string]interface{})[o]
	if !iook {
		return false
	}
	if strings.Index(reflect.TypeOf(io).Name(), "float64") < 0 {
		return false
	}
	if io.(float64) != 1 {
		return false
	}
	return true
}
func user_login(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Account    string `validate:"required"`
		Password   string `validate:"required"`
		VerifyCode string `validate:"required"`
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	type MenuData struct {
		Icon  string     `json:"icon"`
		Index string     `json:"index"`
		Title string     `json:"title"`
		Subs  []MenuData `json:"subs"`
	}
	sqlstr := "select id,SellerId,`Password`,RoleName,Token,State,GoogleSecret,LoginTime,LoginCount from z_admin_user where Account = ?"
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
	qserr, qsresult := db.QueryScan(sqlstr, sqlparam, &userid, &sellerid, &password, &rolename, &token, &userstate, &googlesecret, &logintime, &logincount)
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
	if sellerid != -1 {
		var dbsellerid int
		sqlstr = fmt.Sprintf("select SellerId from %s where SellerId = ?", db_seller_tablename)
		db.QueryScan(sqlstr, []interface{}{sellerid}, &dbsellerid)
		if dbsellerid == 0 {
			ctx.RespErr(-6, "所属运营商已停用")
			return
		}
	}
	var authdata string
	sqlstr = "select RoleData from z_admin_role where RoleName = ? and SellerId = ?"
	sqlparam = []interface{}{rolename, sellerid}
	qserr, qsresult = db.QueryScan(sqlstr, sqlparam, &authdata)
	if qserr != nil {
		ctx.RespErr(-2, qserr.Error())
		return
	}
	if !qsresult {
		ctx.RespErr(-6, "角色不存在")
		return
	}
	if !debug && len(googlesecret) > 0 && !abugo.VerifyGoogleCode(googlesecret, reqdata.VerifyCode) {
		ctx.RespErr(-10, "验证码不正确")
		return
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
	sqlstr = "update z_admin_user set Token = ?,LoginCount = LoginCount + 1,LoginTime = now(),LoginIp = ? where id = ?"
	err = db.QueryNoResult(sqlstr, token, ctx.GetIp(), userid)
	if err != nil {
		ctx.RespErr(-2, err.Error())
		return
	}
	sqlstr = "insert into z_admin_login_log(UserId,SellerId,Account,Token,LoginIp)values(?,?,?,?,?)"
	err = db.QueryNoResult(sqlstr, userid, sellerid, reqdata.Account, token, ctx.GetIp())
	if err != nil {
		ctx.RespErr(-2, err.Error())
		return
	}
	menu := []MenuData{}
	json.Unmarshal([]byte(MenuDataStr), &menu)
	parser := fastjson.Parser{}
	jauthdata, _ := parser.ParseBytes([]byte(authdata))
	//三级菜单
	for i := 0; i < len(menu); i++ {
		for j := 0; j < len(menu[i].Subs); j++ {
			smenu := []MenuData{}
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
		smenu := []MenuData{}
		for j := 0; j < len(menu[i].Subs); j++ {
			open := jauthdata.GetInt(menu[i].Title, menu[i].Subs[j].Title, "查")
			if open == 1 || len(menu[i].Subs[j].Subs) > 0 {
				smenu = append(smenu, menu[i].Subs[j])
			}
		}
		menu[i].Subs = smenu
	}
	//一级菜单
	smenu := []MenuData{}
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
	ctx.Put("AuthData", jauth)
	ctx.Put("MenuData", menu)
	ctx.Put("Token", token)
	ctx.Put("LoginTime", logintime)
	ctx.Put("Ip", ctx.GetIp())
	ctx.Put("LoginCount", logincount)
	ctx.Put("Version", "1.0.0")
	ctx.RespOK()
}
func login_log(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		Account  string
		SellerId int
	}
	reqdata := RequestData{}
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
	if !Auth2(token, "系统管理", "登录日志", "查") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-5, "运营商不正确")
		return
	}
	where := abugo.AbuWhere{}
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	where.AddString("and", "Account", reqdata.Account, "")
	var total int
	db.QueryScan(where.CountSql("z_admin_login_log"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql("z_admin_login_log", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if err != nil {
		logs.Error(err)
		ctx.RespErr(-2, err.Error())
		return
	}
	type ReturnData struct {
		Id         int
		UserId     int
		SellerId   int
		Account    string
		LoginIp    string
		CreateTime string
	}
	data := []ReturnData{}
	for dbresult.Next() {
		d := ReturnData{}
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
func role_list(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		SellerId int
	}
	reqdata := RequestData{}
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
	if !Auth2(token, "系统管理", "角色管理", "查") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-5, "运营商不正确")
		return
	}
	where := abugo.AbuWhere{}
	where.OrderBy = "ASC"
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	var total int
	db.QueryScan(where.CountSql("z_admin_role"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql("z_admin_role", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if err != nil {
		logs.Error(err)
		ctx.RespErr(-2, err.Error())
		return
	}
	type ReturnData struct {
		Id             int
		RoleName       string
		SellerId       int
		ParentSellerId int
		Parent         string
		RoleData       string
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
func role_listall(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		SellerId int
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if !Auth2(token, "系统管理", "角色管理", "查") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-5, "运营商不正确")
		return
	}
	sql := "select RoleName from z_admin_role where SellerId = ?"
	dbresult, err := db.Conn().Query(sql, reqdata.SellerId)
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
func role_data(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId int    `validate:"required"`
		RoleName string `validate:"required"`
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if !Auth2(token, "系统管理", "角色管理", "查") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-5, "运营商不正确")
		return
	}
	sql := "select RoleData from z_admin_role where SellerId = ? and RoleName = ?"
	var RoleData string
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &RoleData)
	var SuperRoleData string
	sql = "select RoleData from z_admin_role where SellerId = -1 and RoleName = '超级管理员'"
	db.QueryScan(sql, []interface{}{}, &SuperRoleData)
	ctx.Put("RoleData", RoleData)
	ctx.Put("SuperRoleData", SuperRoleData)
	ctx.RespOK()
}
func role_check(parent string, parentdata map[string]interface{}, data map[string]interface{}, result *string) {
	defer recover()
	for k, v := range data {
		if strings.Index(reflect.TypeOf(v).Name(), "float") >= 0 {
			if v.(float64) != 1 {
				continue
			}
			path := strings.Split(parent, ".")
			if len(path) == 0 {
				continue
			}
			fk, fok := parentdata[path[0]]
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
				fv := (*pn).(map[string]interface{})[k].(float64)
				if fv != 1 {
					(*result) = "fail"
				}
			} else {
				(*result) = "fail"
			}

		} else {
			role_check(parent+"."+k, parentdata, v.(map[string]interface{}), result)
		}
	}
}
func role_modify(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId int    `validate:"required"`
		RoleName string `validate:"required"`
		RoleData string `validate:"required"`
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if !Auth2(token, "系统管理", "角色管理", "改") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-5, "运营商不正确")
		return
	}
	var ParentSellerId int
	var ParentRoleName string
	sql := "select ParentSellerId,Parent from z_admin_role where SellerId = ? and RoleName = ?"
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &ParentSellerId, &ParentRoleName)
	if len(ParentRoleName) == 0 {
		ctx.RespErr(3, "上级角色不存在")
		return
	}
	var ParentRoleData string
	sql = "select RoleData from z_admin_role where SellerId = ? and RoleName = ?"
	db.QueryScan(sql, []interface{}{ParentSellerId, ParentRoleName}, &ParentRoleData)
	if len(ParentRoleData) == 0 {
		ctx.RespErr(3, "获取上级角色数据失败")
		return
	}
	jparent := make(map[string]interface{})
	err = json.Unmarshal([]byte(ParentRoleData), &jparent)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	jdata := make(map[string]interface{})
	err = json.Unmarshal([]byte(reqdata.RoleData), &jdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	result := ""
	for k, v := range jdata {
		role_check(k, jparent, v.(map[string]interface{}), &result)
	}
	if len(result) > 0 {
		ctx.RespErr(-20, "权限大过上级角色")
		return
	}
	sql = "update z_admin_role set  RoleData = ? where SellerId = ? and RoleName = ?"
	err = db.QueryNoResult(sql, reqdata.RoleData, reqdata.SellerId, reqdata.RoleName)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	WriteAdminLog("修改角色", ctx, reqdata)
	ctx.RespOK()
}
func role_add(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		ParentSellerId int    `validate:"required"`
		Parent         string `validate:"required"`
		SellerId       int    `validate:"required"`
		RoleName       string `validate:"required"`
		RoleData       string `validate:"required"`
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if !Auth2(token, "系统管理", "角色管理", "增") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-5, "运营商不正确")
		return
	}
	if reqdata.SellerId != -1 && reqdata.SellerId != reqdata.ParentSellerId && reqdata.ParentSellerId != -1 {
		ctx.RespErr(-5, "上级角色运营商只能是总后台角色或跟自己所属运营商一致,不可以是别的运营商")
		return
	}
	if reqdata.SellerId == -1 && reqdata.ParentSellerId != -1 {
		ctx.RespErr(-5, "总后台角色上级角色只能是总后台的角色")
		return
	}
	var rid int = 0
	sql := "select id from z_admin_role where SellerId = ? and RoleName = ?"
	db.QueryScan(sql, []interface{}{reqdata.ParentSellerId, reqdata.Parent}, &rid)
	if rid == 0 {
		ctx.RespErr(3, "上级角色不存在")
		return
	}
	rid = 0
	sql = "select id from z_admin_role where SellerId = ? and RoleName = ? "
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &rid)
	if rid > 0 {
		ctx.RespErr(3, "角色已经存在")
		return
	}
	if reqdata.SellerId != -1 {
		sql = fmt.Sprintf("select SellerId from %s where SellerId = ? and state = 1", db_seller_tablename)
		var sellerid int
		db.QueryScan(sql, []interface{}{reqdata.SellerId}, &sellerid)
		if sellerid == 0 {
			ctx.RespErr(3, "运营商不存在")
			return
		}
	}
	sql = "insert into z_admin_role(RoleName,SellerId,ParentSellerId,Parent,RoleData)values(?,?,?,?,?)"
	param := []interface{}{reqdata.RoleName, reqdata.SellerId, reqdata.ParentSellerId, reqdata.Parent, reqdata.RoleData}
	err = db.QueryNoResult(sql, param...)
	if err != nil {
		logs.Error(err)
		ctx.RespErr(-1, err.Error())
	}
	WriteAdminLog("添加角色", ctx, reqdata)
	ctx.RespOK()
}
func role_delete(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		SellerId int    `validate:"required"`
		RoleName string `validate:"required"`
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if !Auth2(token, "系统管理", "角色管理", "删") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-1, "运营商不正确")
		return
	}
	sql := "select id,Parent from z_admin_role where ParentSellerId = ? and Parent = ?"
	var id int
	var parent string
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &id, &parent)
	if id > 0 {
		ctx.RespErr(-2, "该角色有下级角色,不可删除")
		return
	}
	if parent == "god" {
		ctx.RespErr(-5, "该角色不可删除")
		return
	}
	id = 0
	sql = "select id from z_admin_user where RoleSellerId = ? and RoleName = ?"
	db.QueryScan(sql, []interface{}{reqdata.SellerId, reqdata.RoleName}, &id)
	if id > 0 {
		ctx.RespErr(-3, "该角色下存在账号,不可删除")
		return
	}
	sql = "delete from z_admin_role where SellerId = ? and RoleName = ?"
	db.QueryNoResult(sql, reqdata.SellerId, reqdata.RoleName)
	WriteAdminLog("删除角色", ctx, reqdata)
	ctx.RespOK()
}
func opt_log(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		SellerId int
		Account  string
		Opt      string
	}
	reqdata := RequestData{}
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
	if !Auth2(token, "系统管理", "操作日志", "查") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-1, "运营商不正确")
		return
	}
	where := abugo.AbuWhere{}
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	where.AddString("and", "Account", reqdata.Account, "")
	where.AddString("and", "Opt", reqdata.Opt, "")
	var total int
	db.QueryScan(where.CountSql("z_admin_opt_log"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql("z_admin_opt_log", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if err != nil {
		logs.Error(err)
		ctx.RespErr(-2, err.Error())
		return
	}
	type ReturnData struct {
		Id         int
		Account    string
		SellerId   int
		Ip         string
		Opt        string
		Data       string
		CreateTime string
	}
	data := []ReturnData{}
	for dbresult.Next() {
		d := ReturnData{}
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
func user_list(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Page     int
		PageSize int
		Account  string
		SellerId int
	}
	reqdata := RequestData{}
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
	if !Auth2(token, "系统管理", "账号管理", "查") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-1, "运营商不正确")
		return
	}
	where := abugo.AbuWhere{}
	where.OrderBy = "ASC"
	where.AddInt("and", "SellerId", reqdata.SellerId, 0)
	where.AddString("and", "Account", reqdata.Account, "")
	var total int
	db.QueryScan(where.CountSql("z_admin_user"), where.Params, &total)
	if total == 0 {
		ctx.Put("data", []interface{}{})
		ctx.Put("page", reqdata.Page)
		ctx.Put("pagesize", reqdata.PageSize)
		ctx.Put("total", total)
		ctx.RespOK()
		return
	}
	dbresult, err := db.Conn().Query(where.Sql("z_admin_user", reqdata.Page, reqdata.PageSize), where.GetParams()...)
	if err != nil {
		logs.Error(err)
		ctx.RespErr(-2, err.Error())
		return
	}
	type ReturnData struct {
		Id           int
		Account      string
		SellerId     int
		RoleSellerId int
		RoleName     string
		Remark       string
		State        int
		LoginCount   int
		LoginIp      string
		LoginTime    string
		CreateTime   string
	}
	data := []ReturnData{}
	for dbresult.Next() {
		d := ReturnData{}
		abugo.GetDbResult(dbresult, &d)
		d.CreateTime = abugo.TimeToUtc(d.CreateTime)
		d.LoginTime = abugo.TimeToUtc(d.LoginTime)
		data = append(data, d)
	}
	dbresult.Close()
	ctx.Put("data", data)
	ctx.Put("page", reqdata.Page)
	ctx.Put("pagesize", reqdata.PageSize)
	ctx.Put("total", total)
	ctx.RespOK()
}
func user_modify(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Account      string `validate:"required"`
		SellerId     int    `validate:"required"`
		Password     string
		RoleSellerId int    `validate:"required"`
		RoleName     string `validate:"required"`
		State        int    `validate:"required"`
		Remark       string
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if !Auth2(token, "系统管理", "账号管理", "改") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-1, "运营商不正确")
		return
	}
	if reqdata.RoleSellerId != -1 && reqdata.RoleSellerId != reqdata.SellerId {
		ctx.RespErr(-2, "运营商不正确")
		return
	}
	sql := "select id from z_admin_role  where SellerId = ? and RoleName = ?"
	var rid int
	db.QueryScan(sql, []interface{}{reqdata.RoleSellerId, reqdata.RoleName}, &rid)
	if rid == 0 {
		ctx.RespErr(-3, "角色不存在")
		return
	}
	if len(reqdata.Password) > 0 {
		sql = "update z_admin_user set RoleSellerId = ?,RoleName = ?,State = ?,Remark = ?,`Password` = ? where Account = ? and SellerId = ?"
		db.QueryNoResult(sql, reqdata.RoleSellerId, reqdata.RoleName, reqdata.State, reqdata.Remark, reqdata.Password, reqdata.Account, reqdata.SellerId)
	} else {
		sql = "update z_admin_user set RoleSellerId = ?,RoleName = ?,State = ?,Remark = ? where Account = ? and SellerId = ?"
		db.QueryNoResult(sql, reqdata.RoleSellerId, reqdata.RoleName, reqdata.State, reqdata.Remark, reqdata.Account, reqdata.SellerId)
	}
	if reqdata.State != 1 {
		sql = "select Token from z_admin_user where Account = ? and SellerId = ? "
		var tokenstr string
		db.QueryScan(sql, []interface{}{reqdata.Account, reqdata.SellerId}, &tokenstr)
		if len(tokenstr) > 0 {
			http.DelToken(tokenstr)
		}
	}
	WriteAdminLog("修改管理员", ctx, reqdata)
	ctx.RespOK()
}
func user_delete(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Id       int    `validate:"required"`
		Account  string `validate:"required"`
		SellerId int    `validate:"required"`
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if !Auth2(token, "系统管理", "账号管理", "改") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-1, "运营商不正确")
		return
	}
	sql := "delete from z_admin_user where Id = ? and Account = ? and SellerId = ?"
	db.QueryNoResult(sql, reqdata.Id, reqdata.Account, reqdata.SellerId)
	WriteAdminLog("删除管理员", ctx, reqdata)
	ctx.RespOK()
}
func user_add(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Account      string `validate:"required"`
		SellerId     int    `validate:"required"`
		Password     string `validate:"required"`
		RoleSellerId int    `validate:"required"`
		RoleName     string `validate:"required"`
		State        int    `validate:"required"`
		Remark       string
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if !Auth2(token, "系统管理", "账号管理", "增") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-1, "运营商不正确")
		return
	}
	if reqdata.RoleSellerId != -1 && reqdata.RoleSellerId != reqdata.SellerId {
		ctx.RespErr(-2, "运营商不正确")
		return
	}
	sql := "select id from z_admin_role  where SellerId = ? and RoleName = ?"
	var rid int
	db.QueryScan(sql, []interface{}{reqdata.RoleSellerId, reqdata.RoleName}, &rid)
	if rid == 0 {
		ctx.RespErr(-3, "角色不存在")
		return
	}
	sql = "select id from z_admin_user where Account = ? and SellerId = ?"
	var uid int
	db.QueryScan(sql, []interface{}{reqdata.Account, reqdata.SellerId}, &uid)
	if uid > 0 {
		ctx.RespErr(-3, "账号已经存在")
		return
	}
	sql = "insert into z_admin_user(Account,Password,SellerId,RoleSellerId,RoleName,State)values(?,?,?,?,?,?)"
	db.QueryNoResult(sql, reqdata.Account, reqdata.Password, reqdata.SellerId, reqdata.RoleSellerId, reqdata.RoleName, reqdata.State)
	WriteAdminLog("添加管理员", ctx, reqdata)
	ctx.RespOK()
}

func user_google(ctx *abugo.AbuHttpContent) {
	defer recover()
	type RequestData struct {
		Account  string `validate:"required"`
		SellerId int    `validate:"required"`
	}
	reqdata := RequestData{}
	err := ctx.RequestData(&reqdata)
	if err != nil {
		ctx.RespErr(-1, err.Error())
		return
	}
	token := GetToken(ctx)
	if !Auth2(token, "系统管理", "账号管理", "改") {
		ctx.RespErr(-300, "权限不足")
		return
	}
	if token.SellerId > 0 && reqdata.SellerId != token.SellerId {
		ctx.RespErr(-1, "运营商不正确")
		return
	}
	verifykey := abugo.GetGoogleSecret()
	verifyurl := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=abugo", reqdata.Account, verifykey)
	sql := "update z_admin_user set GoogleSecret = ? where Account = ? and SellerId = ?"
	db.QueryNoResult(sql, verifykey, reqdata.Account, reqdata.SellerId)
	ctx.RespOK(verifyurl)
}

func seller_name(ctx *abugo.AbuHttpContent) {
	token := GetToken(ctx)
	if token.SellerId != -1 {
		ctx.Put("data", []interface{}{})
		ctx.RespOK()
		return
	}
	sql := "select * from ex_seller where State = 1"
	dbresult, err := db.Conn().Query(sql)
	if err != nil {
		ctx.RespErr(1, err.Error())
		return
	}
	type ReturnData struct {
		SellerId   int
		SellerName string
	}
	data := []ReturnData{}
	data = append(data, ReturnData{0, "全部"})
	data = append(data, ReturnData{-1, "总后台"})
	for dbresult.Next() {
		d := ReturnData{}
		abugo.GetDbResult(dbresult, &d)
		data = append(data, d)
	}
	dbresult.Close()
	ctx.RespOK(data)
}

package server

var MenuData string = `[
	{
		"icon": "el-icon-lx-home",
		"index": "home",
		"title": "系统首页"
	},
	{
		"icon": "el-icon-setting",
		"index": "7",
		"title": "系统管理",
		"subs":
		[
			{
				"index": "system_account",
				"title": "账号管理"
			},
			{
				"index": "system_role",
				"title": "角色管理"
			},
			{
				"index": "system_log",
				"title": "系统日志"
			},
			{
				"index": "system_login_log",
				"title": "登录日志"
			}
		]
	}
]`

var AuthData = `{
	"系统首页": { "查" : 1 },
	"系统管理": {
		"账号管理": { "查": 1,"增": 1,"删": 1,"改": 1},
		"角色管理": { "查": 1,"增": 1,"删": 1,"改": 1},
		"登录日志": { "查": 1},
		"系统日志": { "查": 1}
	}
}`

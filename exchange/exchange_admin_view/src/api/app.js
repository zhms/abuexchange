export class app {}
import { Message, Loading } from 'element-ui'
import router from '../router'
import axios from 'axios'
var service = axios.create({
	//  baseURL: 'http://192.168.2.46:1221/',
	//	baseURL: 'http://localhost:1221/',
	baseURL: 'http://localhost:1221/',
	timeout: 60000,
})
service.interceptors.request.use(
	(config) => {
		config.headers['x-token'] = sessionStorage.getItem('token')
		return config
	},
	(error) => {
		return Promise.reject(error)
	}
)

app.getInstance = (function() {
	let instance
	return function() {
		if (!instance) {
			instance = new app()
		}
		return instance
	}
})()

app.prototype.clone = function(obj) {
	return JSON.parse(JSON.stringify(obj))
}
//显示加载动画
app.prototype.showLoading = function(show) {
	if (show) {
		if (!this.loading) {
			this.loading = Loading.service({ lock: true, spinner: 'el-icon-loading', background: 'rgba(0, 0, 0, 0.7)' })
		}
	} else {
		if (this.loading) {
			this.loading.close()
			this.loading = null
		}
	}
}
//退回登录界面
app.prototype.showLoginPage = function() {
	router.push('/login')
}
//获取管理员信息
app.prototype.getInfo = function() {
	return this.info
}
//设置管理员信息
app.prototype.setInfo = function(data) {
	if (data) {
		this.info = JSON.parse(data)
	}
}
//get请求
app.prototype.get = function(url, p1, p2) {
	var data = null
	var callback = null
	if (typeof p1 == 'object') {
		data = p1
		if (typeof p2 == 'function') {
			callback = p2
		}
	} else if (typeof p1 == 'function') {
		callback = p1
	}
	if (data) {
		url += '?'
		for (var i in data) {
			url += i
			url += '='
			url += data[i]
			url += '&'
		}
	}
	if (url.charAt(url.length - 1) == '&') {
		url = url.substr(0, url.length - 1)
	}
	service({
		url: url,
		method: 'get',
	})
		.then((result) => {
			if (result.data.errmsg) {
				console.log('get:' + url + ' ' + errmsg)
			} else {
				if (callback) callback(result.data)
			}
		})
		.catch((err) => {
			console.log('get:' + url + ' ' + err)
		})
}
//post请求
app.prototype.post = function(url, data, callback, noloading) {
	noloading = false
	if (!noloading) app.getInstance().showLoading(true)
	service({
		url: url,
		method: 'post',
		data,
	})
		.then((result) => {
			if (!noloading) app.getInstance().showLoading(false)
			if (result.data.errmsg) {
				console.log(result.data.errmsg)
				Message({
					message: result.data.errmsg,
					type: 'error',
					duration: 1000 * 3,
					showClose: true,
					center: true,
				})
			} else {
				if (callback) callback(result.data)
			}
		})
		.catch((err) => {
			if (!noloading) app.getInstance().showLoading(false)
			Message({
				message: err,
				type: 'error',
				duration: 1000 * 3,
				showClose: true,
				center: true,
			})
		})
}

app.prototype.login = function(account, password, verifycode, callback) {
	app.getInstance().post('/login', { account, password, verifycode }, (result) => {
		sessionStorage.setItem('userdata', JSON.stringify(result))
		sessionStorage.setItem('token', result.token)
		callback()
	})
}

app.prototype.auth = (m, s, o) => {
	var info = app.getInstance().getInfo()
	if (!info) {
		return false
	}
	var authm = info.Auth[m]
	if (!authm) {
		return false
	}
	var auths = authm
	if (s) {
		auths = authm[s]
	}
	if (!auths) {
		return false
	}
	var autho = auths[o]
	if (!autho) {
		return false
	}
	return true
}

const config = require('./config.js')
const WebSocket = require('ws')
const uuid = require('uuid')
const wserver = new WebSocket.Server({ port: config.port })
const sredis = require('redis').createClient({ url: `redis://:${config.redis.password}@${config.redis.host}:${config.redis.port}` })
const predis = require('redis').createClient({ url: `redis://:${config.redis.password}@${config.redis.host}:${config.redis.port}` })
const redis = require('redis').createClient({ url: `redis://:${config.redis.password}@${config.redis.host}:${config.redis.port}` })
process.on('uncaughtException', function (err) {
	console.log(err)
})
let intervals = ['1m', '3m', '5m', '15m', '30m', '2h', '4h', '6h', '8h', '12h', '1d', '3d', '1w', '1M']
let intervalseconds = {}
intervalseconds['1m'] = 60
intervalseconds['3m'] = 60 * 3
intervalseconds['5m'] = 60 * 5
intervalseconds['15m'] = 60 * 15
intervalseconds['30m'] = 60 * 30
intervalseconds['2h'] = 60 * 60 * 2
intervalseconds['4h'] = 60 * 60 * 4
intervalseconds['6h'] = 60 * 60 * 6
intervalseconds['8h'] = 60 * 60 * 8
intervalseconds['12h'] = 60 * 60 * 12
intervalseconds['1d'] = 60 * 60 * 24
intervalseconds['3d'] = 60 * 60 * 24 * 3
intervalseconds['1w'] = 60 * 60 * 24 * 7
let id = 1
let server
let connections = {}
let userid2conn = {}
let topics = {}
let last_topics_data = {}
let market_price = {}
let futures_price = {}
let publish
sredis.on('error', (err) => {})
predis.on('error', (err) => {})
redis.on('error', (err) => {})
sredis.connect().then(function () {
	predis.connect().then(function () {
		redis.connect().then(function () {
			server()
		})
	})
})
wserver.on('connection', (conn) => {
	conn.uniqueid = uuid.v4()
	connections[conn.uniqueid] = conn
	conn.send(JSON.stringify({ msgid: 'connected' }))
	conn.on('close', () => {
		if (conn.topics) {
			conn.topics.forEach((element) => {
				delete topics[element].delete(conn.uniqueid)
			})
		}
		if (conn.UserId) delete userid2conn[conn.UserId]
		if (conn.uniqueid) delete connections[conn.uniqueid]
	})
	conn.on('message', (message) => {
		try {
			message = JSON.parse(message.toString('utf-8'))
		} catch (e) {
			return conn.close()
		}
		if (!message.msgid) return conn.close()
		if (message.msgid == 'attach') {
			if (conn.UserId) delete userid2conn[conn.UserId]
			conn.UserId = message.data
			userid2conn[conn.UserId] = conn
			conn.send(
				JSON.stringify({
					msgid: 'attach',
					result: true,
				})
			)
		} else if (message.msgid == 'subscribe') {
			topics[message.data] = topics[message.data] || new Set()
			topics[message.data].add(conn.uniqueid)
			conn.topics = conn.topics || new Set()
			conn.topics.add(message.data)
			let lastdata = last_topics_data[message.data]
			if (lastdata) {
				conn.send(
					JSON.stringify({
						msgid: message.data,
						data: lastdata,
					})
				)
			}
		} else if (message.msgid == 'unsubscribe') {
			topics[message.data] = topics[message.data] || new Set()
			conn.topics.delete(message.data)
			topics[message.data].delete(conn.uniqueid)
		} else if (message.msgid == 'get') {
			let data = message.data.split('-')
			if (data.length < 3) return
			let KlineType = data[0]
			let symbol = data[1]
			let interval = data[2]
			if (interval == '1M') {
				if (KlineType == 'MarketKlineHistory') {
					redis.HGETALL(`reptile:market:kline:${symbol.replace('/', '')}:${interval}`).then((kdata) => {
						var keys = []
						for (var i in kdata) {
							keys.push(i)
						}
						keys.sort((a, b) => {
							return parseInt(b) - parseInt(a)
						})
						var senddata = []
						for (var i = 0; i < keys.length; i++) {
							senddata.push(JSON.parse(kdata[keys[i]]))
						}
						conn.send(
							JSON.stringify({
								msgid: `${KlineType}-${symbol}-${interval}`,
								data: senddata,
							})
						)
					})
				} else if (KlineType == 'FuturesKlineHistory') {
					redis.HGETALL(`reptile:futures:kline:${symbol.replace('/', '')}:${interval}`).then((kdata) => {
						var keys = []
						for (var i in kdata) {
							keys.push(i)
						}
						keys.sort((a, b) => {
							return parseInt(b) - parseInt(a)
						})
						var senddata = []
						for (var i = 0; i < keys.length; i++) {
							senddata.push(JSON.parse(kdata[keys[i]]))
						}
						conn.send(
							JSON.stringify({
								msgid: `${KlineType}-${symbol}-${interval}`,
								data: senddata,
							})
						)
					})
				}
				return
			}
			let startid = data[3]
			let seconds = intervalseconds[interval]
			if (!seconds) return
			if (KlineType == 'MarketKlineHistory') {
				let lastkey = `reptile:market:kline:${symbol.replace('/', '')}:${interval}:lastid`
				redis.get(lastkey).then((lastid) => {
					if (!lastid) return
					if (!startid) startid = lastid
					let ids = []
					for (let i = 0; i < 200; i++) {
						ids.push(`${startid - seconds * (i + 1)}`)
					}
					redis.HMGET(`reptile:market:kline:${symbol.replace('/', '')}:${interval}`, ids).then((kdata) => {
						let senddata = []
						for (let i = 0; i < kdata.length; i++) {
							if (kdata[i]) {
								senddata.push(JSON.parse(kdata[i]))
							}
						}
						conn.send(
							JSON.stringify({
								msgid: `${KlineType}-${symbol}-${interval}`,
								data: senddata,
							})
						)
					})
				})
			} else if (KlineType == 'FuturesKlineHistory') {
				let lastkey = `reptile:futures:kline:${symbol.replace('/', '')}:${interval}:lastid`
				redis.get(lastkey).then((lastid) => {
					if (!lastid) return
					if (!startid) startid = lastid
					let ids = []
					for (let i = 0; i < 200; i++) {
						ids.push(`${startid - seconds * (i + 1)}`)
					}
					redis.HMGET(`reptile:futures:kline:${symbol.replace('/', '')}:${interval}`, ids).then((kdata) => {
						let senddata = []
						for (let i = 0; i < kdata.length; i++) {
							if (kdata[i]) {
								senddata.push(JSON.parse(kdata[i]))
							}
						}
						conn.send(
							JSON.stringify({
								msgid: `${KlineType}-${symbol}-${interval}`,
								data: senddata,
							})
						)
					})
				})
			}
		}
	})
})
publish = function (topic, data) {
	last_topics_data[topic] = data
	let subscribes = topics[topic]
	if (!subscribes) return
	let uniid = id
	id++
	subscribes.forEach((element) => {
		let conn = connections[element]
		if (conn) {
			conn.send(
				JSON.stringify({
					msgid: topic,
					data: data,
					id: uniid,
				})
			)
		}
	})
}
function send_market_depth() {
	redis.hGetAll('reptile:config:depth:market').then((data) => {
		for (let symbol in data) {
			let levels = data[symbol].split('@')
			for (let level in levels) {
				level = levels[level]
				let topic = `MarketDepth-${symbol}-${parseFloat(level)}`
				let rediskey = `reptile:market:depth:${symbol}:${level}`
				redis.get(rediskey).then((rdata) => {
					publish(topic, JSON.parse(rdata))
				})
			}
		}
	})
}
function send_futures_depth() {
	redis.hGetAll('reptile:config:depth:futures').then((data) => {
		for (let symbol in data) {
			let levels = data[symbol].split('@')
			for (let level in levels) {
				level = levels[level]
				let topic = `FuturesDepth-${symbol}-${parseFloat(level)}`
				let rediskey = `reptile:futures:depth:${symbol}:${level}`
				redis.get(rediskey).then((rdata) => {
					publish(topic, JSON.parse(rdata))
				})
			}
		}
	})
}
server = function () {
	sredis.subscribe('msg_to_client', (msg) => {
		try {
			msg = JSON.parse(msg)
		} catch (e) {
			return
		}
		let uniid = id
		id++
		let conn = userid2conn[msg.UserId]
		if (conn) {
			conn.send(
				JSON.stringify({
					msgid: msg.msgid,
					data: msg.data,
					id: uniid,
				})
			)
		}
	})
	sredis.subscribe('reptile_market_kline', (message) => {
		let data = message.split('@')
		let symbol = data[0]
		let interval = data[1]
		let info = JSON.parse(data[2])
		let topic = `MarketKline-${symbol.toLowerCase()}-${interval}`
		publish(topic, info)
	})
	sredis.subscribe('reptile_market_ticker', (message) => {
		let data = message.split('@')
		let symbol = data[0]
		let info = JSON.parse(data[1])
		let topic = `MarketTicker-${symbol.toLowerCase()}`
		publish(topic, info)
	})
	sredis.subscribe('reptile_market_trade', (message) => {
		let data = message.split('@')
		let symbol = data[0]
		let info = JSON.parse(data[1])
		let topic = `MarketTrade-${symbol.toLowerCase()}`
		publish(topic, info)
	})
	sredis.subscribe('reptile_market_price', (message) => {
		let data = message.split('@')
		let symbol = data[0]
		let price = parseFloat(data[1])
		if (market_price[symbol] == price) return
		market_price[symbol] = price
		let topic = `MarketPrice-${symbol.toLowerCase()}`
		publish(topic, { price })
	})
	setInterval(send_market_depth, 1000)
	/////////////////////////////////////////////////////////////////////////////////////////////////
	sredis.subscribe('reptile_futures_kline', (message) => {
		let data = message.split('@')
		let symbol = data[0]
		let interval = data[1]
		let info = JSON.parse(data[2])
		let topic = `FuturesKline-${symbol.toLowerCase()}-${interval}`
		publish(topic, info)
	})
	sredis.subscribe('reptile_futures_ticker', (message) => {
		let data = message.split('@')
		let symbol = data[0]
		let info = JSON.parse(data[1])
		let topic = `FuturesTicker-${symbol.toLowerCase()}`
		publish(topic, info)
	})
	sredis.subscribe('reptile_futures_trade', (message) => {
		let data = message.split('@')
		let symbol = data[0]
		let info = JSON.parse(data[1])
		let topic = `FuturesTrade-${symbol.toLowerCase()}`
		publish(topic, info)
	})
	sredis.subscribe('reptile_futures_price', (message) => {
		let data = message.split('@')
		let symbol = data[0]
		let price = parseFloat(data[1])
		if (futures_price[symbol] == price) return
		futures_price[symbol] = price
		let topic = `FuturesPrice-${symbol.toLowerCase()}`
		publish(topic, { price })
	})
	sredis.subscribe('reptile_futures_info', (message) => {
		let data = message.split('@')
		let symbol = data[0]
		let info = JSON.parse(data[1])
		let topic = `FuturesInfo-${symbol.toLowerCase()}`
		publish(topic, info)
	})
	setInterval(send_futures_depth, 1000)
}
console.log('*********************start*********************')

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
	setInterval(() => {
		for (let i = 0; i < config.symbols.length; i++) {
			let symbol = config.symbols[i].replace('/', '')
			redis.HGETALL(`reptile:market:depth:${symbol}:asks`).then((asks) => {
				redis.HGETALL(`reptile:market:depth:${symbol}:bids`).then((bids) => {
					let arrasks = []
					for (let ai in asks) {
						arrasks.push([parseFloat(ai), parseFloat(asks[ai])])
					}
					let arrbids = []
					for (let bi in bids) {
						arrbids.push([parseFloat(bi), parseFloat(bids[bi])])
					}
					arrasks.sort((a, b) => {
						return a[0] - b[0]
					})
					arrbids.sort((a, b) => {
						return b[0] - a[0]
					})
					if (config.depths[config.symbols[i]]) {
						for (let k = 0; k < config.depths[config.symbols[i]].length; k++) {
							let d = config.depths[config.symbols[i]][k]
							let sd = `${d}`.split('.')
							let decimal = 0
							if (sd.length > 1) decimal = sd[1].length
							d = d * Math.pow(10, decimal)
							let mapasks = {}
							for (let ad in arrasks) {
								let price = arrasks[ad][0]
								let amount = arrasks[ad][1]
								price = Math.floor(price * Math.pow(10, decimal))
								price = Math.floor(Math.floor(price / d) * d)
								price = price / Math.pow(10, decimal)
								mapasks[price] = mapasks[price] || 0
								mapasks[price] += amount
							}
							let mapbids = {}
							for (let ad in arrbids) {
								let price = arrbids[ad][0]
								let amount = arrbids[ad][1]
								price = Math.floor(price * Math.pow(10, decimal))
								price = Math.floor(Math.floor(price / d) * d)
								price = price / Math.pow(10, decimal)
								mapbids[price] = mapbids[price] || 0
								mapbids[price] += amount
							}
							let arrmapasks = []
							for (let ai in mapasks) {
								arrmapasks.push([parseFloat(ai), parseFloat(mapasks[ai])])
							}
							let arrmapbids = []
							for (let bi in mapbids) {
								arrmapbids.push([parseFloat(bi), parseFloat(mapbids[bi])])
							}
							arrmapasks.sort((a, b) => {
								return a[0] - b[0]
							})
							arrmapbids.sort((a, b) => {
								return b[0] - a[0]
							})
							let publishdata = { asks: [], bids: [] }
							for (let j = 0; j < arrmapasks.length; j++) {
								publishdata.asks.push(arrmapasks[j])
								if (j >= 30) break
							}
							for (let j = 0; j < arrmapbids.length; j++) {
								publishdata.bids.push(arrmapbids[j])
								if (j >= 30) break
							}
							let topic = `MarketDepth-${symbol}-${parseFloat(d / Math.pow(10, decimal))}`
							publish(topic, publishdata)
						}
					}
				})
			})
		}
	}, 1000)
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
}

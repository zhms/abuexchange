const config = require('./config.js')
const redis = require('redis').createClient({ url: `redis://:${config.redis.password}@${config.redis.host}:${config.redis.port}` })
redis.on('error', (err) => {})
let server
process.on('uncaughtException', function (err) {
	console.log(err)
})
redis.connect().then(function () {
	server()
})

server = function () {
	setInterval(() => {
		for (let i = 0; i < config.symbols.length; i++) {
			let symbol = config.symbols[i].replace('/', '')
			redis.HGETALL(`reptile:futures:depth:${symbol}:asks`).then((asks) => {
				redis.HGETALL(`reptile:futures:depth:${symbol}:bids`).then((bids) => {
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
							let rediskey = `reptile:futures:depth:${symbol}:${parseFloat(d / Math.pow(10, decimal))}`
							redis.set(rediskey, JSON.stringify(publishdata))
						}
					}
				})
			})
		}
	}, 1000)
}

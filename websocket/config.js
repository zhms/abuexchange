module.exports = {
	port: 4754,
	redis: {
		host: '127.0.0.1',
		port: 6379,
		password: 'Ho9mpyeqaILEOfjM',
	},
	symbols: ['btc/usdt', 'eth/usdt'],
	depths: {
		'btc/usdt': [0.01, 0.1, 1, 10, 50, 100],
		'eth/usdt': [0.01, 0.1, 1, 10, 50, 100],
	},
}

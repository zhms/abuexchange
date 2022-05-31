<template>
	<p>房间统计</p>
</template>
<script>
import { app } from '@/api/app.js'
import '@/assets/css/k.css'
export default {
	data() {
		return {}
	},
	components: {},
	computed: {},
	created() {},
	methods: {
		auth(o) {
			return app.auth2('玩家管理', '账号管理', o)
		},
	},
}
</script>

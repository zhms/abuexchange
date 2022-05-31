<template>
	<div>
		<el-dialog :title="config_title" :visible.sync="dialog" width="400px" center>
			<el-form :inline="true" label-width="100px">
				<el-form-item label="玩家胜率:">
					<el-input v-model="config.UserWinRate"></el-input>
				</el-form-item>
				<el-form-item label="下注筹码:">
					<el-input v-model="chip" placeholder="逗号分开,例如:100,200,300"></el-input>
				</el-form-item>
			</el-form>
			<span slot="footer" class="dialog-footer">
				<el-button type="primary" @click="handleConfirm">确 定</el-button>
			</span>
		</el-dialog>
	</div>
</template>
<script>
import { app } from '@/api/app.js'
export default {
	data() {
		return {
			dialog: false,
			dialog_data: null,
			chip: '',
			config: {
				chip: [],
				UserWinRate: null,
			},
		}
	},
	props: ['config_title', 'show', 'confirm'],
	components: {},
	computed: {},
	created() {
		this.show('long_hu_da_zhan', (data) => {
			this.dialog_data = data
			var cfgstr = this.dialog_data.config
			if (!cfgstr || cfgstr == '') cfgstr = '{}'
			var cfg = {}
			try {
				cfg = JSON.parse(cfgstr)
			} catch (e) {}
			this.chip = ''
			if (cfg.chip) {
				for (var i in cfg.chip) {
					if (this.chip.length > 0) {
						this.chip += ','
					}
					this.chip += cfg.chip[i]
				}
			}
			this.config.UserWinRate = cfg.UserWinRate
			this.dialog = true
		})
	},
	methods: {
		setValue(dest, src, filed) {
			if (src[filed]) {
				dest[filed] = src[filed]
			}
		},
		handleConfirm() {
			var arr = this.chip.split(',')
			this.config.chip = []
			for (var i in arr) {
				this.config.chip.push(Number(arr[i]))
			}
			var data = {
				gameid: this.dialog_data.gameid,
				level: this.dialog_data.level,
				config: JSON.stringify(this.config),
			}
			app.post('/game/room/config', data, () => {
				this.confirm(data.config)
				this.$message.success('操作成功')
				this.dialog = false
			})
		},
	},
}
</script>

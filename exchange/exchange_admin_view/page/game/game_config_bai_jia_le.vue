<template>
	<div>
		<el-dialog :title="config_title" :visible.sync="dialog" width="380px" center>
			<el-form :inline="true" label-width="100px">
				<el-form-item label="空闲时间:">
					<el-input v-model="GameConfig.FreeDelay"></el-input>
				</el-form-item>
				<el-form-item label="下注时间:">
					<el-input v-model="GameConfig.BetDelay"></el-input>
				</el-form-item>
				<el-form-item label="开牌时间:">
					<el-input v-model="GameConfig.OpenDelay"></el-input>
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
			GameConfig: {
				FreeDelay: null,
				BetDelay: null,
				OpenDelay: null,
			},
		}
	},
	props: ['config_title', 'show', 'confirm'],
	components: {},
	computed: {},
	created() {
		this.show('BaiJiaLe', (data) => {
			this.dialog_data = data
			var cfgstr = this.dialog_data.GameConfig
			if (!cfgstr || cfgstr == '') cfgstr = '{}'
			var cfg = {}
			try {
				cfg = JSON.parse(cfgstr)
			} catch (e) {}
			this.setValue(this.GameConfig, cfg, 'FreeDelay')
			this.setValue(this.GameConfig, cfg, 'BetDelay')
			this.setValue(this.GameConfig, cfg, 'OpenDelay')
			this.dialog = true
			console.log(this.config_title)
		})
	},
	methods: {
		setValue(dest, src, filed) {
			if (src[filed]) {
				dest[filed] = src[filed]
			}
		},
		handleConfirm() {
			this.dialog_data.GameConfig = JSON.stringify(this.GameConfig)
			app.post('/game/game/modify', this.dialog_data, () => {
				this.confirm(this.dialog_data.GameConfig)
				this.$message.success('操作成功')
				this.dialog = false
			})
		},
	},
}
</script>

<template>
	<div>
		<el-dialog :title="control_title" :visible.sync="dialog" width="380px" center>
			<el-form :inline="true" label-width="100px">
				<el-form-item label="玩家胜率:">
					<el-input v-model="GameControl.UserWinRate"></el-input>
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
			GameControl: {
				UserWinRate: null,
			},
		}
	},
	props: ['control_title', 'show', 'confirm'],
	components: {},
	computed: {},
	created() {
		this.show('BaiJiaLe', (data) => {
			this.dialog_data = data
			var cfgstr = this.dialog_data.GameControl
			if (!cfgstr || cfgstr == '') cfgstr = '{}'
			var cfg = {}
			try {
				cfg = JSON.parse(cfgstr)
			} catch (e) {}
			this.setValue(this.GameControl, cfg, 'UserWinRate')
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
			this.dialog_data.GameControl = JSON.stringify(this.GameControl)
			app.post('/game/game/modify', this.dialog_data, () => {
				this.confirm(this.dialog_data.GameControl)
				this.$message.success('操作成功')
				this.dialog = false
			})
		},
	},
}
</script>

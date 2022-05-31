<template>
	<div>
		<el-dialog :title="control_title" :visible.sync="dialog" width="380px" center>
			<el-form :inline="true" label-width="100px">
				<el-form-item label="样本编号:">
					<el-input v-model="GameControl.no" placeholder="请输入内容"></el-input>
				</el-form-item>
				<el-form-item label="收益率:">
					<el-select v-model="GameControl.rtp" placeholder="请选择">
						<el-option v-for="item in options" :key="item.label" :label="item.label" :value="item.value"> </el-option>
					</el-select>
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
			options: [
				{ label: '110%', value: 110 },
				{ label: '105%', value: 105 },
				{ label: '95%', value: 95 },
				{ label: '90%', value: 90 },
				{ label: '85%', value: 85 },
				{ label: '80%', value: 80 },
			],
			GameControl: {
				no: null,
				rtp: null,
			},
		}
	},
	props: ['control_title', 'show', 'confirm'],
	components: {},
	computed: {},
	created() {
		this.show('JiuXianLaBa', (data) => {
			this.dialog_data = data
			var cfgstr = this.dialog_data.GameControl
			if (!cfgstr || cfgstr == '') cfgstr = '{}'
			var cfg = {}
			try {
				cfg = JSON.parse(cfgstr)
			} catch (e) {}
			this.setValue(this.GameControl, cfg, 'no')
			this.setValue(this.GameControl, cfg, 'rtp')
			this.dialog = true
		})
	},
	methods: {
		setValue(dest, src, filed) {
			if (src[filed] != null && src[filed] != undefined) {
				dest[filed] = src[filed]
			}
		},
		handleConfirm() {
			this.GameControl.no = Number(this.GameControl.no)
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

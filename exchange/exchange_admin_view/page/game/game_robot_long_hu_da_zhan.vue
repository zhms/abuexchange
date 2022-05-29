<template>
	<div>
		<el-dialog :title="robot_title" :visible.sync="dialog" width="450px" center>
			<el-form :inline="true" label-width="150px">
				<el-form-item label="每局最大下注次数:">
					<el-input v-model="GameRobot.MaxBetCount"></el-input>
				</el-form-item>
				<el-form-item label="下注筹码权重:">
					<el-input v-model="BetChipWeight"></el-input>
				</el-form-item>
				<el-form-item label="下注区域权重:">
					<el-input v-model="BetAreaWeight" placeholder="龙,虎,和"></el-input>
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
			BetChipWeight: '',
			BetAreaWeight: '',
			GameRobot: {
				MaxBetCount: null,
			},
		}
	},
	props: ['robot_title', 'show', 'confirm'],
	components: {},
	computed: {},
	created() {
		this.show('LongHuDaZhan', (data) => {
			this.dialog_data = data
			var robotstr = this.dialog_data.GameRobot
			if (!robotstr || robotstr == '') robotstr = '{}'
			var cfg = {}
			try {
				cfg = JSON.parse(robotstr)
			} catch (e) {}
			this.GameRobot.MaxBetCount = cfg.MaxBetCount
			this.BetChipWeight = ''
			this.BetAreaWeight = ''
			if (cfg.BetChipWeight) {
				for (var i in cfg.BetChipWeight) {
					if (this.BetChipWeight.length > 0) {
						this.BetChipWeight += ','
					}
					this.BetChipWeight += cfg.BetChipWeight[i]
				}
			}
			if (cfg.BetAreaWeight) {
				for (var i in cfg.BetAreaWeight) {
					if (this.BetAreaWeight.length > 0) {
						this.BetAreaWeight += ','
					}
					this.BetAreaWeight += cfg.BetAreaWeight[i]
				}
			}
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
			var arr = this.BetChipWeight.split(',')
			this.GameRobot.BetChipWeight = []
			for (var i in arr) {
				this.GameRobot.BetChipWeight.push(Number(arr[i]))
			}
			arr = this.BetAreaWeight.split(',')
			this.GameRobot.BetAreaWeight = []
			for (var i in arr) {
				this.GameRobot.BetAreaWeight.push(Number(arr[i]))
			}
			this.dialog_data.GameRobot = JSON.stringify(this.GameRobot)
			app.getInstance().post('/game/game/modify', this.dialog_data, () => {
				this.confirm(this.dialog_data.GameRobot)
				this.$message.success('操作成功')
				this.dialog = false
			})
		},
	},
}
</script>
